package domain

import (
	"time"

	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"github.com/golang/glog"
	"weblogic-operator/pkg/types"
	"weblogic-operator/pkg/constants"
)

type StoreToWebLogicDomainLister struct {
	cache.Store
}

type StoreToWebLogicDomainReplicaSetLister struct {
	cache.Store
}

// The WebLogicDomainController watches the Kubernetes API for changes to WebLogicDomain resources
type WebLogicDomainController struct {
	client                        kubernetes.Interface
	restClient                    *rest.RESTClient
	startTime                     time.Time
	shutdown                      bool
	weblogicDomainController      cache.Controller
	weblogicDomainStore           StoreToWebLogicDomainLister
	weblogicDomainReplicaSet      cache.Controller
	weblogicDomainReplicaSetStore StoreToWebLogicDomainReplicaSetLister
}

// NewController creates a new WebLogicDomainController.
func NewController(kubeClient kubernetes.Interface, restClient *rest.RESTClient, resyncPeriod time.Duration, namespace string) (*WebLogicDomainController, error) {
	m := WebLogicDomainController{
		client:     kubeClient,
		restClient: restClient,
		startTime:  time.Now(),
	}

	weblogicDomainHandlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    m.onAdd,
		DeleteFunc: m.onDelete,
		UpdateFunc: m.onUpdate,
	}

	watcher := cache.NewListWatchFromClient(restClient, constants.WebLogicDomainResourceKindPlural, namespace, fields.Everything())

	m.weblogicDomainStore.Store, m.weblogicDomainController = cache.NewInformer(
		watcher,
		&types.WebLogicDomain{},
		resyncPeriod,
		weblogicDomainHandlers)

	replicaSetHandler := cache.ResourceEventHandlerFuncs{
		AddFunc:    m.onReplicaSetAdd,
		DeleteFunc: m.onReplicaSetDelete,
		UpdateFunc: m.onReplicaSetUpdate,
	}

	m.weblogicDomainReplicaSetStore.Store, m.weblogicDomainReplicaSet = cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				options.LabelSelector = constants.WebLogicDomainLabel
				return kubeClient.ExtensionsV1beta1().ReplicaSets(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				options.LabelSelector = constants.WebLogicDomainLabel
				return kubeClient.ExtensionsV1beta1().ReplicaSets(namespace).Watch(options)
			},
		},
		&v1beta1.ReplicaSet{},
		resyncPeriod,
		replicaSetHandler,
	)

	return &m, nil
}

func (m *WebLogicDomainController) onAdd(obj interface{}) {
	glog.V(4).Info("WebLogicDomainController.onAdd() called")

	weblogicDomain := obj.(*types.WebLogicDomain)
	err := createWebLogicDomain(weblogicDomain, m.client, m.restClient)
	if err != nil {
		glog.Errorf("Failed to create weblogicDomain: %s", err)
	}
}

func (m *WebLogicDomainController) onDelete(obj interface{}) {
	glog.V(4).Info("WebLogicDomainController.onDelete() called")

	weblogicDomain := obj.(*types.WebLogicDomain)
	err := deleteWebLogicDomain(weblogicDomain, m.client, m.restClient)
	if err != nil {
		glog.Errorf("Failed to delete weblogicDomain: %s", err)
	}
}

func (m *WebLogicDomainController) onUpdate(old, cur interface{}) {
	glog.V(4).Info("WebLogicDomainController.onUpdate() called")
	curDomain := cur.(*types.WebLogicDomain)
	oldDomain := old.(*types.WebLogicDomain)
	if curDomain.ResourceVersion == oldDomain.ResourceVersion {
		return
	}

	err := createWebLogicDomain(curDomain, m.client, m.restClient)
	if err != nil {
		glog.Errorf("Failed to update domain: %s", err)
	}
}

func (m *WebLogicDomainController) onReplicaSetAdd(obj interface{}) {
	glog.V(4).Info("WebLogicDomainController.onReplicaSetAdd() called")

	replicaSet := obj.(*v1beta1.ReplicaSet)

	weblogicDomain, err := GetDomainForReplicaSet(replicaSet, m.restClient)
	if err != nil {
		// FIXME: Should we delete the replica set here???
		// it has no domain but it has the label.
		glog.Errorf("Failed to find domain for replica set: %s(%s):%#v", replicaSet.Name, err.Error(), replicaSet.Labels)
		return
	}
	err = updateDomainWithReplicaSet(weblogicDomain, replicaSet, m.client, m.restClient)
	if err != nil {
		glog.Errorf("Failed to create Domain: %s", err)
	}
}

//TODO Fix hanldings here. Need to call onStatefulSetAdd ???
func (m *WebLogicDomainController) onReplicaSetDelete(obj interface{}) {
	glog.V(4).Info("WebLogicDomainController.onReplicaSetDelete() called")
	m.onReplicaSetAdd(obj)
}

func (m *WebLogicDomainController) onReplicaSetUpdate(old, new interface{}) {
	glog.V(4).Info("WebLogicDomainController.onReplicaSetUpdate() called")
	m.onReplicaSetAdd(new)
}

// Run the WebLogic controller
func (m *WebLogicDomainController) Run(stopChan <-chan struct{}) {
	glog.Infof("Starting WebLogic Domain controller")
	go m.weblogicDomainController.Run(stopChan)
	//go m.weblogicStatefulSetController.Run(stopChan)
	go m.weblogicDomainReplicaSet.Run(stopChan)
	<-stopChan
	glog.Infof("Shutting down WebLogic Domain controller")
}
