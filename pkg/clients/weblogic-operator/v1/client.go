package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"

	"weblogic-operator/pkg/types"
)

type WeblogicV1Interface interface {
	RESTClient() rest.Interface
	WeblogicServersGetter
}

// WeblogicV1Client is used to interact with features provided by the group.
type WeblogicV1Client struct {
	restClient rest.Interface
}

// New creates a new WeblogicV1Client for the given RESTClient.
func New(c rest.Interface) *WeblogicV1Client {
	return &WeblogicV1Client{c}
}

// NewForConfig creates a new WeblogicV1Client for the given config.
func NewForConfig(c *rest.Config) (*WeblogicV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &WeblogicV1Client{client}, nil
}

func setConfigDefaults(config *rest.Config) error {
	crScheme := runtime.NewScheme()
	// TODO(apryde): Is this necessary? Can we do this in one place?
	if err := types.AddToScheme(crScheme); err != nil {
		return err
	}

	gv := types.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{
		CodecFactory: serializer.NewCodecFactory(crScheme),
	}

	return nil
}

// WeblogicServers returns a WeblogicServerInterface used to interact with
// WeblogicServer custom resources.
func (c *WeblogicV1Client) WeblogicServers(namespace string) WeblogicServerInterface {
	return newWeblogicServers(c, namespace)
}

// RESTClient returns a RESTClient that is used to communicate with the API
// server used by this client implementation.
func (c *WeblogicV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
