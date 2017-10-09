package server

import (
	"time"

	"github.com/golang/glog"
	"k8s.io/api/autoscaling/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"weblogic-operator/pkg/constants"
	"weblogic-operator/pkg/types"
)

type StoreToWebLogicManagedServerLister struct {
	cache.Store
}

type StoreToWebLogicManagedServerReplicaSetLister struct {
	cache.Store
}

type StoreToWebLogicManagedServerHorizontalPodAutoscalingLister struct {
	cache.Store
}

// The WebLogicManagedServerController watches the Kubernetes API for changes to WebLogicManagedServer resources
type WebLogicManagedServerController struct {
	client                                             kubernetes.Interface
	restClient                                         *rest.RESTClient
	startTime                                          time.Time
	shutdown                                           bool
	weblogicManagedServerController                    cache.Controller
	weblogicManagedServerStore                         StoreToWebLogicManagedServerLister
	weblogicManagedServerReplicaSet                    cache.Controller
	weblogicManagedServerReplicaSetStore               StoreToWebLogicManagedServerReplicaSetLister
	weblogicManagedServerHorizontalPodAutoscaling      cache.Controller
	weblogicManagedServerHorizontalPodAutoscalingStore StoreToWebLogicManagedServerHorizontalPodAutoscalingLister
}

// NewController creates a new WebLogicManagedServerController.
func NewController(kubeClient kubernetes.Interface, restClient *rest.RESTClient, resyncPeriod time.Duration, namespace string) (*WebLogicManagedServerController, error) {
	m := WebLogicManagedServerController{
		client:     kubeClient,
		restClient: restClient,
		startTime:  time.Now(),
	}

	weblogicManagedServerHandlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    m.onAdd,
		DeleteFunc: m.onDelete,
		UpdateFunc: m.onUpdate,
	}

	watcher := cache.NewListWatchFromClient(restClient, constants.WebLogicManagedServerResourceKindPlural, namespace, fields.Everything())

	m.weblogicManagedServerStore.Store, m.weblogicManagedServerController = cache.NewInformer(
		watcher,
		&types.WebLogicManagedServer{},
		resyncPeriod,
		weblogicManagedServerHandlers)

	replicaSetHandler := cache.ResourceEventHandlerFuncs{
		AddFunc:    m.onReplicaSetAdd,
		DeleteFunc: m.onReplicaSetDelete,
		UpdateFunc: m.onReplicaSetUpdate,
	}

	m.weblogicManagedServerReplicaSetStore.Store, m.weblogicManagedServerReplicaSet = cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				options.LabelSelector = constants.WebLogicManagedServerLabel
				return kubeClient.ExtensionsV1beta1().ReplicaSets(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				options.LabelSelector = constants.WebLogicManagedServerLabel
				return kubeClient.ExtensionsV1beta1().ReplicaSets(namespace).Watch(options)
			},
		},
		&v1beta1.ReplicaSet{},
		resyncPeriod,
		replicaSetHandler)

	horizontalPodAutoscalerHandler := cache.ResourceEventHandlerFuncs{
		AddFunc:    m.onHorizontalPodAutoscalerAdd,
		DeleteFunc: m.onHorizontalPodAutoscalerDelete,
		UpdateFunc: m.onHorizontalPodAutoscalerUpdate,
	}

	m.weblogicManagedServerHorizontalPodAutoscalingStore.Store, m.weblogicManagedServerHorizontalPodAutoscaling = cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				options.LabelSelector = constants.WebLogicManagedServerLabel
				return kubeClient.AutoscalingV1().HorizontalPodAutoscalers(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				options.LabelSelector = constants.WebLogicManagedServerLabel
				return kubeClient.AutoscalingV1().HorizontalPodAutoscalers(namespace).Watch(options)
			},
		},
		&v1.HorizontalPodAutoscaler{},
		resyncPeriod,
		horizontalPodAutoscalerHandler)

	return &m, nil
}

func (m *WebLogicManagedServerController) onAdd(obj interface{}) {
	glog.V(4).Info("WebLogicManagedServerController.onAdd() called")

	weblogicManagedServer := obj.(*types.WebLogicManagedServer)
	err := createWebLogicManagedServer(weblogicManagedServer, m.client, m.restClient)
	if err != nil {
		glog.Errorf("Failed to create weblogicManagedServer: %s", err)
	}
}

func (m *WebLogicManagedServerController) onDelete(obj interface{}) {
	glog.V(4).Info("WebLogicManagedServerController.onDelete() called")

	weblogicManagedServer := obj.(*types.WebLogicManagedServer)
	err := deleteWebLogicManagedServer(weblogicManagedServer, m.client, m.restClient)
	if err != nil {
		glog.Errorf("Failed to delete weblogicManagedServer: %s", err)
	}
}

func (m *WebLogicManagedServerController) onUpdate(old, cur interface{}) {
	glog.V(4).Info("WebLogicManagedServerController.onUpdate() called")
	curServer := cur.(*types.WebLogicManagedServer)
	oldServer := old.(*types.WebLogicManagedServer)
	if curServer.ResourceVersion == oldServer.ResourceVersion {
		// Periodic resync will send update events for all known servers.
		// Two different versions of the same server will always have
		// different RVs.
		return
	}

	err := updateWebLogicManagedServer(curServer, m.client, m.restClient)
	if err != nil {
		glog.Errorf("Failed to update server: %s", err)
	}
}

func (m *WebLogicManagedServerController) onReplicaSetAdd(obj interface{}) {
	glog.V(4).Info("WebLogicManagedServerController.onReplicaSetAdd() called")

	replicaSet := obj.(*v1beta1.ReplicaSet)

	weblogicManagedServer, err := GetServerForReplicaSet(replicaSet, m.restClient)
	if err != nil {
		// FIXME: Should we delete the replica set here???
		// it has no server but it has the label.
		glog.Errorf("Failed to find server for replica set: %s(%s):%#v", replicaSet.Name, err.Error(), replicaSet.Labels)
		return
	}
	err = updateServerWithReplicaSet(weblogicManagedServer, replicaSet, m.client, m.restClient)
	if err != nil {
		glog.Errorf("Failed to create update Server: %s", err)
	}
}

//TODO Fix hanldings here. Need to call onStatefulSetAdd ???
func (m *WebLogicManagedServerController) onReplicaSetDelete(obj interface{}) {
	glog.V(4).Info("WebLogicManagedServerController.onReplicaSetDelete() called")
	m.onReplicaSetAdd(obj)
}

func (m *WebLogicManagedServerController) onReplicaSetUpdate(old, new interface{}) {
	glog.V(4).Info("WebLogicManagedServerController.onReplicaSetUpdate() called")
	m.onReplicaSetAdd(new)
}

func (m *WebLogicManagedServerController) onHorizontalPodAutoscalerAdd(obj interface{}) {
	glog.V(4).Info("WebLogicManagedServerController.onHorizontalPodAutoscalerAdd() called")

	horizontalPodAutoscaler := obj.(*v1.HorizontalPodAutoscaler)

	weblogicManagedServer, err := GetServerForHorizontalPodAutoscaler(horizontalPodAutoscaler, m.restClient)
	if err != nil {
		// FIXME: Should we delete the replica set here???
		// it has no server but it has the label.
		glog.Errorf("Failed to find server for horizontal pod autoscaler: %s(%s):%#v", horizontalPodAutoscaler.Name, err.Error(), horizontalPodAutoscaler.Labels)
		return
	}
	err = updateServerWithHorizontalPodAutoscaler(weblogicManagedServer, horizontalPodAutoscaler, m.client, m.restClient)
	if err != nil {
		glog.Errorf("Failed to create update Server: %s", err)
	}
}

//TODO Fix hanldings here. Need to call onStatefulSetAdd ???
func (m *WebLogicManagedServerController) onHorizontalPodAutoscalerDelete(obj interface{}) {
	glog.V(4).Info("WebLogicManagedServerController.onHorizontalPodAutoscalerDelete() called")
	m.onHorizontalPodAutoscalerAdd(obj)
}

func (m *WebLogicManagedServerController) onHorizontalPodAutoscalerUpdate(old, new interface{}) {
	glog.V(4).Info("WebLogicManagedServerController.onHorizontalPodAutoscalerUpdate() called")
	m.onHorizontalPodAutoscalerAdd(new)
}

// Run the WebLogic controller
func (m *WebLogicManagedServerController) Run(stopChan <-chan struct{}) {
	glog.Infof("Starting WebLogic controller")
	go m.weblogicManagedServerController.Run(stopChan)
	go m.weblogicManagedServerReplicaSet.Run(stopChan)
	go m.weblogicManagedServerHorizontalPodAutoscaling.Run(stopChan)
	<-stopChan
	glog.Infof("Shutting down WebLogic controller")
}
