package weblogic_operator

import (
	"k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"

	mysqlopv1 "weblogic-operator/pkg/clients/weblogic-operator/v1"
)

// Interface for the mysql operator client.
type Interface interface {
	MySQLV1() mysqlopv1.MySQLV1Interface
}

// Clientset contains the clients for the MySQL operator API groups.
type Clientset struct {
	*mysqlopv1.MySQLV1Client
}

// MySQLV1 retrieves the MySQLV1Client
func (c *Clientset) MySQLV1() mysqlopv1.MySQLV1Interface {
	if c == nil {
		return nil
	}
	return c.MySQLV1Client
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.MySQLV1Client = mysqlopv1.New(c)
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
	cs.MySQLV1Client, err = mysqlopv1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}
