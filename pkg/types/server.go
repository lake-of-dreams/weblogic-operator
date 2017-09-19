package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ = runtime.Object(&WeblogicServer{})

const (
	defaultVersion  = "12.2.1.2"
	defaultReplicas = 1
)

// ServerCRDResourcePlural defines the custom resource name for weblogicservers
const ServerCRDResourcePlural = "weblogicservers"

var validVersions = []string{
	defaultVersion,
}

// WeblogicServerSpec defines the attributes a user can specify when creating a server
type WeblogicServerSpec struct {
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
type WeblogicServerPhase string

const (
	// WeblogicServerPending means the server has been accepted by the system,
	// but one or more of the services or statefulsets has not been started.
	// This includes time before being bound to a node, as well as time spent
	// pulling images onto the host.
	WeblogicServerPending WeblogicServerPhase = "Pending"
	// WeblogicServerRunning means the server has been created, all of it's
	// required components are present, and there is at least one endpoint that
	// weblogic client can connect to.
	WeblogicServerRunning WeblogicServerPhase = "Running"
	// WeblogicServerStopped means that all containers in the pod have
	// voluntarily terminated with a container exit code of 0, and the system
	// is not going to restart any of these containers.
	WeblogicServerStopped WeblogicServerPhase = "Stopped"
	// WeblogicServerFailed means that all containers in the pod have terminated,
	// and at least one container has terminated in a failure (exited with a
	// non-zero exit code or was stopped by the system).
	WeblogicServerFailed WeblogicServerPhase = "Failed"
	// WeblogicServerUnknown means that for some reason the state of the server
	// could not be obtained, typically due to an error in communicating with
	// the host of the pod.
	WeblogicServerUnknown WeblogicServerPhase = ""
)

var WeblogicServerValidPhases = []WeblogicServerPhase{WeblogicServerPending,
	WeblogicServerRunning,
	WeblogicServerStopped,
	WeblogicServerFailed,
	WeblogicServerUnknown}

type WeblogicServerStatus struct {
	Phase  WeblogicServerPhase `json:"phase"`
	Errors []string            `json:"errors"`
}

// WeblogicServer represents a server spec and associated metadata
type WeblogicServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              WeblogicServerSpec   `json:"spec"`
	Status            WeblogicServerStatus `json:"status"`
}

type WeblogicServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []WeblogicServer `json:"items"`
}

// Validate returns an error if a server is invalid
func (c *WeblogicServer) Validate() error {
	return validateServer(c).ToAggregate()
}

// EnsureDefaults will ensure that if a user omits and fields in the
// spec that are required, we set some sensible defaults.
// For example a user can choose to omit the version
// and number of replicas
func (c *WeblogicServer) EnsureDefaults() *WeblogicServer {
	if c.Spec.Replicas == 0 {
		c.Spec.Replicas = defaultReplicas
	}

	if c.Spec.Version == "" {
		c.Spec.Version = defaultVersion
	}

	return c
}

func (c *WeblogicServer) GetObjectKind() schema.ObjectKind {
	return &c.TypeMeta
}

func (c *WeblogicServerList) GetObjectKind() schema.ObjectKind {
	return &c.TypeMeta
}
