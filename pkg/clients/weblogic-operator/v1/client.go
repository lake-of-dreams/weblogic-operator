package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"

	"weblogic-operator/pkg/types"
)

type MySQLV1Interface interface {
	RESTClient() rest.Interface
	MySQLClustersGetter
	MySQLBackupsGetter
	MySQLRestoresGetter
}

// MySQLV1Client is used to interact with features provided by the group.
type MySQLV1Client struct {
	restClient rest.Interface
}

// New creates a new MySQLV1Client for the given RESTClient.
func New(c rest.Interface) *MySQLV1Client {
	return &MySQLV1Client{c}
}

// NewForConfig creates a new MySQLV1Client for the given config.
func NewForConfig(c *rest.Config) (*MySQLV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &MySQLV1Client{client}, nil
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

// MySQLClusters returns a MySQLClusterInterface used to interact with
// MySQLCluster custom resources.
func (c *MySQLV1Client) MySQLClusters(namespace string) MySQLClusterInterface {
	return newMySQLClusters(c, namespace)
}

// MySQLBackups returns a MySQLBackupInterface used to interact with
// MySQLBackup custom resources.
func (c *MySQLV1Client) MySQLBackups(namespace string) MySQLBackupInterface {
	return newMySQLBackups(c, namespace)
}

// MySQLRestores returns a MySQLRestoreInterface used to interact with
// MySQLRestore custom resources.
func (c *MySQLV1Client) MySQLRestores(namespace string) MySQLRestoreInterface {
	return newMySQLRestores(c, namespace)
}

// RESTClient returns a RESTClient that is used to communicate with the API
// server used by this client implementation.
func (c *MySQLV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
