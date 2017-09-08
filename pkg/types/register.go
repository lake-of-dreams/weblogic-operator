package types

// This package will auto register types with the Kubernetes API

import (
	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
)

const (
	groupName                   = "weblogic.oracle.com"
	schemeVersion               = "v1"
	WeblogicServerCRDResourceKind = "WeblogicServer"
)

var (
	schemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = schemeBuilder.AddToScheme
	SchemeGroupVersion = schema.GroupVersion{Group: groupName, Version: schemeVersion}
)

// addKnownTypes adds the set of types defined in this package to the supplied
// scheme.
func addKnownTypes(s *runtime.Scheme) error {
	s.AddKnownTypes(SchemeGroupVersion,
		&WeblogicServer{},
		&WeblogicServerList{})
	metav1.AddToGroupVersion(s, SchemeGroupVersion)
	return nil
}

func registerDefaults(scheme *runtime.Scheme) error {
	scheme.AddTypeDefaultingFunc(&WeblogicServer{}, defaultWeblogicServer)
	scheme.AddTypeDefaultingFunc(&WeblogicServerList{}, defaultWeblogicServerList)
	return nil
}

// TODO currently unused

func defaultWeblogicServerList(obj interface{}) {
	clusterList := obj.(*WeblogicServerList)
	for _, cluster := range clusterList.Items {
		defaultWeblogicServer(cluster)
	}
}

func defaultWeblogicServer(obj interface{}) {
	cluster := obj.(*WeblogicServer)
	cluster.Spec.Replicas = defaultReplicas
	cluster.Spec.Version = defaultVersion
	defaultWeblogicServerStatus(cluster.Status)
}

func defaultWeblogicServerStatus(obj interface{}) {
	clusterStatus := obj.(*WeblogicServerStatus)
	clusterStatus.Phase = WeblogicServerUnknown
	clusterStatus.Errors = []string{}
}

func init() {
	glog.Info("Registering Types")
	addKnownTypes(scheme.Scheme)
	registerDefaults(scheme.Scheme)
	glog.V(4).Infof("All types: %#v", scheme.Scheme.AllKnownTypes())
}
