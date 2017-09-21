package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ = runtime.Object(&WeblogicDomain{})

const (
	defaultDomainVersion  = "12.2.1.2"
	defaultDomainReplicas = 1
)

// DomainCRDResourcePlural defines the custom resource name for weblogicdomain
const DomainCRDResourcePlural = "weblogicdomains"

var validDomainVersions = []string{
	defaultDomainVersion,
}

// WeblogicServerSpec defines the attributes a user can specify when creating a server
type WeblogicDomainSpec struct {
	// Version defines the Weblogic Docker image version
	Version string `json:"version"`
	// Replicas defines the number of running Weblogic server instances
	Replicas int32 `json:"replicas,omitempty"`
	// NodeSelector is a selector which must be true for the pod to fit on a node.
	// Selector which must match a node's labels for the pod to be scheduled on that node.
	// More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
}

// WeblogicServerPhase describes the state of the server.
type WeblogicDomainPhase string

const (
	// WeblogicDomainPending means the domain has been accepted by the system,
	// This includes time before being bound to a node, as well as time spent
	// pulling images onto the host.
	WeblogicDomainPending WeblogicDomainPhase = "Pending"
	/* WeblogicDomainRunning means the domain has been created, all of it's
	 required components are present, and there is at least one endpoint that
	 weblogic client can connect to.*/
	WeblogicDomainCreated WeblogicDomainPhase = "Created"
	// WeblogicDomainFailed means that all containers in the pod have terminated,
	// and at least one container has terminated in a failure (exited with a
	// non-zero exit code or was stopped by the system).
	WeblogicDomainFailed WeblogicDomainPhase = "Failed"
	// WeblogicDomainUnknown means that for some reason the state of the Domain
	// could not be obtained, typically due to an error in communicating with
	// the host of the pod.
	WeblogicDomainUnknown WeblogicDomainPhase = ""
)

var WeblogicDomainValidPhases = []WeblogicDomainPhase{WeblogicDomainPending,
	WeblogicDomainCreated,
	WeblogicDomainFailed,
	WeblogicDomainUnknown}

type WeblogicDomainStatus struct {
	Phase  WeblogicDomainPhase `json:"phase"`
	Errors []string            `json:"errors"`
}

// WeblogicDomain represents a doamin spec and associated metadata
type WeblogicDomain struct {
	metav1.TypeMeta             `json:",inline"`
	metav1.ObjectMeta           `json:"metadata"`
	Spec   WeblogicDomainSpec   `json:"spec"`
	Status WeblogicDomainStatus `json:"status"`
}

type WeblogicDomainList struct {
	metav1.TypeMeta        `json:",inline"`
	metav1.ListMeta        `json:"metadata"`
	Items []WeblogicServer `json:"items"`
}

// Validate returns an error if a server is invalid
//func (c *WeblogicDomain) Validate() error {
//	return validateServer(c).ToAggregate()
//}

// EnsureDefaults will ensure that if a user omits and fields in the
// spec that are required, we set some sensible defaults.
// For example a user can choose to omit the version
// and number of replicas
func (c *WeblogicDomain) EnsureDefaults() *WeblogicDomain {
	if c.Spec.Replicas == 0 {
		c.Spec.Replicas = defaultDomainReplicas
	}

	if c.Spec.Version == "" {
		c.Spec.Version = defaultDomainVersion
	}

	return c
}

func (c *WeblogicDomain) GetObjectKind() schema.ObjectKind {
	return &c.TypeMeta
}

func (c *WeblogicDomainList) GetObjectKind() schema.ObjectKind {
	return &c.TypeMeta
}
