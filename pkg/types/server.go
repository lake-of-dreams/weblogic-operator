package types

import (
	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"weblogic-operator/pkg/constants"
)

var _ = runtime.Object(&WebLogicManagedServer{})
var ServerRESTClient *rest.RESTClient

const (
	defaultServersToRun = 0
)

// WebLogicManagedServerSpec defines the attributes a user can specify when creating a server
type WebLogicManagedServerSpec struct {
	DomainName   string `json:"domainName"`
	ServersToRun int32  `json:"serversToRun,omitempty"`
	Domain       WebLogicDomain
	// NodeSelector is a selector which must be true for the pod to fit on a node.
	// Selector which must match a node's labels for the pod to be scheduled on that node.
	// More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// Compute Resources required by this container.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#resources
	// +optional
	Resources v1.ResourceRequirements `json:"resources,omitempty" protobuf:"bytes,8,opt,name=resources"`
}

// WebLogicManagedServer represents a server spec and associated metadata
type WebLogicManagedServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              WebLogicManagedServerSpec `json:"spec"`
}

type WebLogicManagedServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []WebLogicManagedServer `json:"items"`
}

// EnsureDefaults will ensure that if a user omits and fields in the
// spec that are required, we set some sensible defaults.
// For example a user can choose to omit the version
// and number of replicas
func (c *WebLogicManagedServer) EnsureDefaults() *WebLogicManagedServer {
	if c.Spec.ServersToRun == 0 {
		c.Spec.ServersToRun = defaultServersToRun
	}

	return c
}

func (c *WebLogicManagedServer) PopulateDomain() *WebLogicManagedServer {
	domain := &WebLogicDomain{}
	result := DomainRESTClient.Get().
		Resource(constants.WebLogicDomainResourceKindPlural).
		Namespace(c.Namespace).
		Name(c.Spec.DomainName).
		Do().
		Into(domain)

	if result != nil {
		glog.V(4).Info("Extracted domain %s for server %s", domain.Name, c.Name)
	}

	c.Spec.Domain = *domain
	return c
}

func (c *WebLogicManagedServer) GetObjectKind() schema.ObjectKind {
	return &c.TypeMeta
}

func (c *WebLogicManagedServerList) GetObjectKind() schema.ObjectKind {
	return &c.TypeMeta
}

func NewManagedServerRESTClient(config *rest.Config) (*rest.RESTClient, error) {
	//if err := types.AddToScheme(scheme.Scheme); err != nil {
	//	return nil, err
	//}
	config.GroupVersion = &WeblogicManagedServerSchemeGroupVersion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: serializer.NewCodecFactory(scheme.Scheme)}

	ServerRESTClient, _ = rest.RESTClientFor(config)
	return rest.RESTClientFor(config)
}
