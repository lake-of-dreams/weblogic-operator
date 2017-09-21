package domain

import (
	"time"

	"k8s.io/api/apps/v1beta1"
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

type StoreToWeblogicDomainLister struct {
	cache.Store
}

type StoreToWeblogicDomainStatefulSetLister struct {
	cache.Store
}

// The WeblogicController watches the Kubernetes API for changes to Weblogic resources
type WeblogicDomainController struct {
	client                              kubernetes.Interface
	restClient                          *rest.RESTClient
	startTime                           time.Time
	shutdown                            bool
	weblogicDomainController            cache.Controller
	weblogicDomainStore                 StoreToWeblogicDomainLister
	weblogicDomainStatefulSetController cache.Controller
	weblogicDomainStatefulSetStore      StoreToWeblogicDomainStatefulSetLister
}

// NewController creates a new WeblogicController.
func NewController(kubeClient kubernetes.Interface, restClient *rest.RESTClient, resyncPeriod time.Duration, namespace string) (*WeblogicDomainController, error) {
	m := WeblogicDomainController{
		client:     kubeClient,
		restClient: restClient,
		startTime:  time.Now(),
	}

	weblogicDomainHandlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    m.onAdd,
		DeleteFunc: m.onDelete,
		UpdateFunc: m.onUpdate,
	}

	watcher := cache.NewListWatchFromClient(restClient, types.DomainCRDResourcePlural, namespace, fields.Everything())
	m.weblogicDomainStore.Store, m.weblogicDomainController = cache.NewInformer(
		watcher,
		&types.WeblogicDomain{},
		resyncPeriod,
		weblogicDomainHandlers)

	statefulSetHandler := cache.ResourceEventHandlerFuncs{
		AddFunc:    m.onStatefulSetAdd,
		DeleteFunc: m.onStatefulSetDelete,
		UpdateFunc: m.onStatefulSetUpdate,
	}

	m.weblogicDomainStatefulSetStore.Store, m.weblogicDomainStatefulSetController = cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				options.LabelSelector = constants.WebLogicDomainLabel
				return kubeClient.AppsV1beta1().StatefulSets(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				options.LabelSelector = constants.WebLogicDomainLabel
				return kubeClient.AppsV1beta1().StatefulSets(namespace).Watch(options)
			},
		},
		&v1beta1.StatefulSet{},
		resyncPeriod,
		statefulSetHandler)

	return &m, nil
}

func (m *WeblogicDomainController) onAdd(obj interface{}) {
	glog.V(4).Info("WeblogicController.onAdd() called")

	weblogicDomain := obj.(*types.WeblogicDomain)
	err := createWeblogicDomain(weblogicDomain, m.client, m.restClient)
	if err != nil {
		glog.Errorf("Failed to create weblogicDomain: %s", err)
		err = setWeblogicDomainState(weblogicDomain, m.restClient, types.WeblogicDomainFailed, err)
	}
}

func (m *WeblogicDomainController) onDelete(obj interface{}) {
	glog.V(4).Info("WeblogicController.onDelete() called")

	weblogicDomain := obj.(*types.WeblogicDomain)
	err := deleteWeblogicDomain(weblogicDomain, m.client, m.restClient)
	if err != nil {
		glog.Errorf("Failed to delete weblogicDomain: %s", err)
		err = setWeblogicDomainState(weblogicDomain, m.restClient, types.WeblogicDomainFailed, err)
	}
}

func (m *WeblogicDomainController) onUpdate(old, cur interface{}) {
	glog.V(4).Info("WeblogicController.onUpdate() called")
	curDomain := cur.(*types.WeblogicDomain)
	oldDomain := old.(*types.WeblogicDomain)
	if curDomain.ResourceVersion == oldDomain.ResourceVersion {
		// Periodic resync will send update events for all known servers.
		// Two different versions of the same server will always have
		// different RVs.
		return
	}

	err := createWeblogicDomain(curDomain, m.client, m.restClient)
	if err != nil {
		glog.Errorf("Failed to update domain: %s", err)
		err = setWeblogicDomainState(curDomain, m.restClient, types.WeblogicDomainFailed, err)
	}
}

func (m *WeblogicDomainController) onStatefulSetAdd(obj interface{}) {
	glog.V(4).Info("WeblogicController.onStatefulSetAdd() called")

	statefulSet := obj.(*v1beta1.StatefulSet)

	weblogicDomain, err := GetDomainForStatefulSet(statefulSet, m.restClient)
	if err != nil {
		// FIXME: Should we delete the stateful set here???
		// it has no server but it has the label.
		glog.Errorf("Failed to find server for stateful set: %s(%s):%#v", statefulSet.Name, err.Error(), statefulSet.Labels)
		return
	}
	err = updateDomainWithStatefulSet(weblogicDomain, statefulSet, m.client, m.restClient)
	if err != nil {
		glog.Errorf("Failed to create update Server: %s", err)
	}
}

//TODO Fix hanldings here. Need to call onStatefulSetAdd ???
func (m *WeblogicDomainController) onStatefulSetDelete(obj interface{}) {
	glog.V(4).Info("WeblogicController.onStatefulSetDelete() called")
	m.onStatefulSetAdd(obj)
}

func (m *WeblogicDomainController) onStatefulSetUpdate(old, new interface{}) {
	glog.V(4).Info("WeblogicController.onStatefulSetUpdate() called")
	m.onStatefulSetAdd(new)
}

// Run the Weblogic controller
func (m *WeblogicDomainController) Run(stopChan <-chan struct{}) {
	glog.Infof("Starting Weblogic controller")
	go m.weblogicDomainController.Run(stopChan)
	go m.weblogicDomainStatefulSetController.Run(stopChan)
	<-stopChan
	glog.Infof("Shutting down Weblogic controller")
}
