package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"weblogic-operator/pkg/types"
)

// MySQLClustersGetter has a method to return a MySQLClusterInterface.
type MySQLClustersGetter interface {
	MySQLClusters(namespace string) MySQLClusterInterface
}

// MySQLClusterInterface has methods to work with MySQL Cluster custom
// resources.
type MySQLClusterInterface interface {
	Create(*types.MySQLCluster) (*types.MySQLCluster, error)
	Update(*types.MySQLCluster) (*types.MySQLCluster, error)
	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*types.MySQLCluster, error)
	List(opts metav1.ListOptions) (*types.MySQLClusterList, error)
}

// mySQLClusters implements the MySQLClusterInterface
type mySQLClusters struct {
	client rest.Interface
	ns     string
}

// newMySQLClusters returns a mySQLClusters.
func newMySQLClusters(c *MySQLV1Client, namespace string) *mySQLClusters {
	return &mySQLClusters{client: c.RESTClient(), ns: namespace}
}

// Create takes the representation of a MySQLCluster and creates it. Returns
// the server's representation of the MySQLCluster, and an error, if there is
// any.
func (c *mySQLClusters) Create(cluster *types.MySQLCluster) (result *types.MySQLCluster, err error) {
	result = &types.MySQLCluster{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource(types.ClusterCRDResourcePlural).
		Body(cluster).
		Do().
		Into(result)
	return
}

// Update takes the representation of a MySQLCluster and updates it. Returns
// the server's representation of the MySQLCluster, and an error, if there is
// any.
func (c *mySQLClusters) Update(cluster *types.MySQLCluster) (result *types.MySQLCluster, err error) {
	result = &types.MySQLCluster{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource(types.ClusterCRDResourcePlural).
		Name(cluster.Name).
		Body(cluster).
		Do().
		Into(result)
	return
}

// Delete takes name of the cluster and deletes it. Returns an error if one
// occurs.
func (c *mySQLClusters) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource(types.ClusterCRDResourcePlural).
		Name(name).
		Body(options).
		Do().
		Error()
}

// Get takes name of the MySQLCluster, and returns the corresponding
// MySQLCluster object, and an error if there is any.
func (c *mySQLClusters) Get(name string, options metav1.GetOptions) (result *types.MySQLCluster, err error) {
	result = &types.MySQLCluster{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource(types.ClusterCRDResourcePlural).
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of
// MySQLClusterList that match those selectors.
func (c *mySQLClusters) List(opts metav1.ListOptions) (result *types.MySQLClusterList, err error) {
	result = &types.MySQLClusterList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource(types.ClusterCRDResourcePlural).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}
