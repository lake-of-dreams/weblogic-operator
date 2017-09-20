package types

// This package will auto register types with the Kubernetes API

import (
	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"weblogic-operator/pkg/constants"
)

var (
	schemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = schemeBuilder.AddToScheme
	SchemeGroupVersion = schema.GroupVersion{Group: constants.WeblogicServerGroupName, Version: constants.WeblogicServerSchemeVersion}
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
	serverList := obj.(*WeblogicServerList)
	for _, server := range serverList.Items {
		defaultWeblogicServer(server)
	}
}

func defaultWeblogicServer(obj interface{}) {
	server := obj.(*WeblogicServer)
	server.Spec.Replicas = defaultReplicas
	server.Spec.Version = defaultVersion
	defaultWeblogicServerStatus(server.Status)
}

func defaultWeblogicServerStatus(obj interface{}) {
	serverStatus := obj.(*WeblogicServerStatus)
	serverStatus.Phase = WeblogicServerUnknown
	serverStatus.Errors = []string{}
}

func init() {
	glog.Info("Registering Types")
	addKnownTypes(scheme.Scheme)
	registerDefaults(scheme.Scheme)
	glog.V(4).Infof("All types: %#v", scheme.Scheme.AllKnownTypes())
}
