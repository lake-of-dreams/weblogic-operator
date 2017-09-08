package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"weblogic-operator/pkg/types"
)

// MySQLRestoresGetter has a method to return a MySQLRestoreInterface.
type MySQLRestoresGetter interface {
	MySQLRestores(namespace string) MySQLRestoreInterface
}

// MySQLRestoreInterface has methods to work with MySQL restore custom resources.
type MySQLRestoreInterface interface {
	Create(*types.MySQLRestore) (*types.MySQLRestore, error)
	Update(*types.MySQLRestore) (*types.MySQLRestore, error)
	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*types.MySQLRestore, error)
	List(opts metav1.ListOptions) (*types.MySQLRestoreList, error)
}

// MySQLRestores implements the MySQLRestoreInterface
type mySQLRestores struct {
	client rest.Interface
	ns     string
}

// newMySQLRestores returns a mySQLRestores.
func newMySQLRestores(c *MySQLV1Client, namespace string) *mySQLRestores {
	return &mySQLRestores{client: c.RESTClient(), ns: namespace}
}

// Create takes the representation of a MySQLRestore and creates it. Returns
// the server's representation of the MySQLRestore, and an error, if there is
// any.
func (c *mySQLRestores) Create(restore *types.MySQLRestore) (result *types.MySQLRestore, err error) {
	result = &types.MySQLRestore{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource(types.RestoreCRDResourcePlural).
		Body(restore).
		Do().
		Into(result)
	return
}

// Update takes the representation of a MySQLRestore and updates it. Returns
// the server's representation of the MySQLRestore, and an error, if there is
// any.
func (c *mySQLRestores) Update(restore *types.MySQLRestore) (result *types.MySQLRestore, err error) {
	result = &types.MySQLRestore{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource(types.RestoreCRDResourcePlural).
		Name(restore.Name).
		Body(restore).
		Do().
		Into(result)
	return
}

// Delete takes name of the restore and deletes it. Returns an error if one
// occurs.
func (c *mySQLRestores) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource(types.RestoreCRDResourcePlural).
		Name(name).
		Body(options).
		Do().
		Error()
}

// Get takes name of the MySQLRestore, and returns the corresponding
// MySQLRestore object, and an error if there is any.
func (c *mySQLRestores) Get(name string, options metav1.GetOptions) (result *types.MySQLRestore, err error) {
	result = &types.MySQLRestore{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource(types.RestoreCRDResourcePlural).
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of
// MySQLRestoreList that match those selectors.
func (c *mySQLRestores) List(opts metav1.ListOptions) (result *types.MySQLRestoreList, err error) {
	result = &types.MySQLRestoreList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource(types.RestoreCRDResourcePlural).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}
