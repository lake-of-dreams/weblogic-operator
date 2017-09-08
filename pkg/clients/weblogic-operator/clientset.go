package weblogic_operator

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/flowcontrol"

	weblogicopv1 "weblogic-operator/pkg/clients/weblogic-operator/v1"
)

// Interface for the weblogic operator client.
type Interface interface {
	WeblogicV1() weblogicopv1.WeblogicV1Interface
}

// Clientset contains the clients for the Weblogic operator API groups.
type Clientset struct {
	*weblogicopv1.WeblogicV1Client
}

// WeblogicV1 retrieves the WeblogicV1Client
func (c *Clientset) WeblogicV1() weblogicopv1.WeblogicV1Interface {
	if c == nil {
		return nil
	}
	return c.WeblogicV1Client
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.WeblogicV1Client = weblogicopv1.New(c)
	return &cs
}

// NewForConfig creates a new Clientset for the given config.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs Clientset
	var err error
	cs.WeblogicV1Client, err = weblogicopv1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}
