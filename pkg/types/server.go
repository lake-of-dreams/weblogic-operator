package types

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ = runtime.Object(&WebLogicManagedServer{})

const (
	defaultVersion  = "12.2.1.2"
	defaultReplicas = 0
)

// WebLogicManagedServerSpec defines the attributes a user can specify when creating a server
type WebLogicManagedServerSpec struct {
	// Version defines the Weblogic Docker image version
	Version string `json:"version"`
	// Replicas defines the number of running Weblogic server instances
	Replicas int32 `json:"replicas,omitempty"`
	// NodeSelector is a selector which must be true for the pod to fit on a node.
	// Selector which must match a node's labels for the pod to be scheduled on that node.
	// More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// Compute Resources required by this container.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#resources
	// +optional
	Resources  v1.ResourceRequirements `json:"resources,omitempty" protobuf:"bytes,8,opt,name=resources"`
	DomainName string                  `json:"domainName"`
}

// WebLogicManagedServer represents a server spec and associated metadata
type WebLogicManagedServer struct {
	metav1.TypeMeta                `json:",inline"`
	metav1.ObjectMeta              `json:"metadata"`
	Spec WebLogicManagedServerSpec `json:"spec"`
}

type WebLogicManagedServerList struct {
	metav1.TypeMeta               `json:",inline"`
	metav1.ListMeta               `json:"metadata"`
	Items []WebLogicManagedServer `json:"items"`
}

// EnsureDefaults will ensure that if a user omits and fields in the
// spec that are required, we set some sensible defaults.
// For example a user can choose to omit the version
// and number of replicas
func (c *WebLogicManagedServer) EnsureDefaults() *WebLogicManagedServer {
	if c.Spec.Replicas == 0 {
		c.Spec.Replicas = defaultReplicas
	}

	if c.Spec.Version == "" {
		c.Spec.Version = defaultVersion
	}

	return c
}

func (c *WebLogicManagedServer) GetObjectKind() schema.ObjectKind {
	return &c.TypeMeta
}

func (c *WebLogicManagedServerList) GetObjectKind() schema.ObjectKind {
	return &c.TypeMeta
}
