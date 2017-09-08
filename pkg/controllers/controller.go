package controllers

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

// Controller provides an interface for controller executors.
type Controller interface {
	// Run executes the controller blocking until it recieves on the
	// stopChan.
	Run(stopChan <-chan struct{})
}

// StoreToMySQLClusterLister TODO add doc strings explaining what this is for
type StoreToMySQLClusterLister struct {
	cache.Store
}

// StoreToMySQLStatefulSetLister TODO add doc strings explaining what this is for
type StoreToMySQLStatefulSetLister struct {
	cache.Store
}

// The MySQLController watches the Kubernetes API for changes to MySQL resources
type MySQLController struct {
	client                     kubernetes.Interface
	restClient                 *rest.RESTClient
	startTime                  time.Time
	shutdown                   bool
	mySQLClusterController     cache.Controller
	mySQLClusterStore          StoreToMySQLClusterLister
	mySQLStatefulSetController cache.Controller
	mySQLStatefulSetStore      StoreToMySQLStatefulSetLister
}

// NewController creates a new MySQLController.
func NewController(kubeClient kubernetes.Interface, restClient *rest.RESTClient, resyncPeriod time.Duration, namespace string) (*MySQLController, error) {
	m := MySQLController{
		client:     kubeClient,
		restClient: restClient,
		startTime:  time.Now(),
	}

	mySQLClusterHandlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    m.onAdd,
		DeleteFunc: m.onDelete,
		UpdateFunc: m.onUpdate,
	}

	watcher := cache.NewListWatchFromClient(restClient, types.ClusterCRDResourcePlural, namespace, fields.Everything())
	m.mySQLClusterStore.Store, m.mySQLClusterController = cache.NewInformer(
		watcher,
		&types.MySQLCluster{},
		resyncPeriod,
		mySQLClusterHandlers)

	statefulSetHandler := cache.ResourceEventHandlerFuncs{
		AddFunc:    m.onStatefulSetAdd,
		DeleteFunc: m.onStatefulSetDelete,
		UpdateFunc: m.onStatefulSetUpdate,
	}

	m.mySQLStatefulSetStore.Store, m.mySQLStatefulSetController = cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				options.LabelSelector = constants.MySQLClusterLabel
				return kubeClient.AppsV1beta1().StatefulSets(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				options.LabelSelector = constants.MySQLClusterLabel
				return kubeClient.AppsV1beta1().StatefulSets(namespace).Watch(options)
			},
		},
		&v1beta1.StatefulSet{},
		resyncPeriod,
		statefulSetHandler)

	return &m, nil
}

func (m *MySQLController) onAdd(obj interface{}) {
	glog.V(4).Info("MySQLController.onAdd() called")

	mySQLCluster := obj.(*types.MySQLCluster)
	err := createCluster(mySQLCluster, m.client, m.restClient)
	if err != nil {
		glog.Errorf("Failed to create mySQLCluster: %s", err)
		err = setMySQLClusterState(mySQLCluster, m.restClient, types.MySQLClusterFailed, err)
	}
}

func (m *MySQLController) onDelete(obj interface{}) {
	glog.V(4).Info("MySQLController.onDelete() called")

	mySQLCluster := obj.(*types.MySQLCluster)
	err := deleteCluster(mySQLCluster, m.client, m.restClient)
	if err != nil {
		glog.Errorf("Failed to delete mySQLCluster: %s", err)
		err = setMySQLClusterState(mySQLCluster, m.restClient, types.MySQLClusterFailed, err)
	}
}

func (m *MySQLController) onUpdate(old, cur interface{}) {
	glog.V(4).Info("MySQLController.onUpdate() called")
	curCluster := cur.(*types.MySQLCluster)
	oldCluster := old.(*types.MySQLCluster)
	if curCluster.ResourceVersion == oldCluster.ResourceVersion {
		// Periodic resync will send update events for all known clusters.
		// Two different versions of the same cluster will always have
		// different RVs.
		return
	}

	err := createCluster(curCluster, m.client, m.restClient)
	if err != nil {
		glog.Errorf("Failed to update cluster: %s", err)
		err = setMySQLClusterState(curCluster, m.restClient, types.MySQLClusterFailed, err)
	}
}

func (m *MySQLController) onStatefulSetAdd(obj interface{}) {
	glog.V(4).Info("MySQLController.onStatefulSetAdd() called")

	statefulSet := obj.(*v1beta1.StatefulSet)

	mySQLCluster, err := GetClusterForStatefulSet(statefulSet, m.restClient)
	if err != nil {
		// FIXME: Should we delete the stateful set here???
		// it has no cluster but it has the label.
		glog.Errorf("Failed to find cluster for stateful set: %s(%s):%#v", statefulSet.Name, err.Error(), statefulSet.Labels)
		return
	}
	err = updateClusterWithStatefulSet(mySQLCluster, statefulSet, m.client, m.restClient)
	if err != nil {
		glog.Errorf("Failed to create update Cluster: %s", err)
	}
}

func (m *MySQLController) onStatefulSetDelete(obj interface{}) {
	glog.V(4).Info("MySQLController.onStatefulSetDelete() called")
	m.onStatefulSetAdd(obj)
}

func (m *MySQLController) onStatefulSetUpdate(old, new interface{}) {
	glog.V(4).Info("MySQLController.onStatefulSetUpdate() called")
	m.onStatefulSetAdd(new)
}

// Run the MySQL controller
func (m *MySQLController) Run(stopChan <-chan struct{}) {
	glog.Infof("Starting MySQL controller")
	go m.mySQLClusterController.Run(stopChan)
	go m.mySQLStatefulSetController.Run(stopChan)
	<-stopChan
	glog.Infof("Shutting down MySQL controller")
}
