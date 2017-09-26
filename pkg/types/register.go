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
	schemeBuilder                           = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme                             = schemeBuilder.AddToScheme
	WeblogicManagedServerSchemeGroupVersion = schema.GroupVersion{Group: constants.WebLogicGroupName, Version: constants.WebLogicManagedServerSchemeVersion}
	WebLogicDomainSchemeGroupVersion        = schema.GroupVersion{Group: constants.WebLogicGroupName, Version: constants.WebLogicDomainSchemeVersion}
)

// addKnownTypes adds the set of types defined in this package to the supplied
// scheme.
func addKnownTypes(s *runtime.Scheme) error {
	s.AddKnownTypes(WeblogicManagedServerSchemeGroupVersion,
		&WebLogicManagedServer{},
		&WebLogicManagedServerList{})
	metav1.AddToGroupVersion(s, WeblogicManagedServerSchemeGroupVersion)

	s.AddKnownTypes(WebLogicDomainSchemeGroupVersion,
		&WebLogicDomain{},
		&WebLogicDomainList{})
	metav1.AddToGroupVersion(s, WebLogicDomainSchemeGroupVersion)
	return nil
}

func registerDefaults(scheme *runtime.Scheme) error {
	scheme.AddTypeDefaultingFunc(&WebLogicManagedServer{}, defaultWebLogicManagedServer)
	scheme.AddTypeDefaultingFunc(&WebLogicManagedServerList{}, defaultWebLogicManagedServerList)
	scheme.AddTypeDefaultingFunc(&WebLogicDomain{}, defaultWebLogicDomain)
	scheme.AddTypeDefaultingFunc(&WebLogicManagedServerList{}, defaultWebLogicDomainList)
	return nil
}

// TODO currently unused

func defaultWebLogicManagedServerList(obj interface{}) {
	serverList := obj.(*WebLogicManagedServerList)
	for _, server := range serverList.Items {
		defaultWebLogicManagedServer(server)
	}
}

func defaultWebLogicManagedServer(obj interface{}) {
	server := obj.(*WebLogicManagedServer)
	server.Spec.Replicas = defaultReplicas
	server.Spec.Version = defaultVersion
}

func defaultWebLogicDomainList(obj interface{}) {
	domainList := obj.(*WebLogicDomainList)
	for _, domain := range domainList.Items {
		defaultWebLogicDomain(domain)
	}
}

func defaultWebLogicDomain(obj interface{}) {
	domain := obj.(*WebLogicDomain)
	domain.Spec.ManagedServerCount = defaultDomainManagedServerCount
	domain.Spec.Version = defaultDomainVersion
}

func init() {
	glog.Info("Registering Types")
	addKnownTypes(scheme.Scheme)
	registerDefaults(scheme.Scheme)
	glog.V(4).Infof("All types: %#v", scheme.Scheme.AllKnownTypes())
}
