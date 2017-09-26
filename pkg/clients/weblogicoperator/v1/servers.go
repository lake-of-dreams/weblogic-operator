package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"weblogic-operator/pkg/types"
	"weblogic-operator/pkg/constants"
)

// WebLogicManagedServersGetter has a method to return a WebLogicManagedServerInterface.
type WebLogicManagedServersGetter interface {
	WebLogicManagedServers(namespace string) WebLogicManagedServerInterface
}

// WebLogicManagedServerInterface has methods to work with Weblogic Server custom
// resources.
type WebLogicManagedServerInterface interface {
	Create(*types.WebLogicManagedServer) (*types.WebLogicManagedServer, error)
	Update(*types.WebLogicManagedServer) (*types.WebLogicManagedServer, error)
	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*types.WebLogicManagedServer, error)
	List(opts metav1.ListOptions) (*types.WebLogicManagedServerList, error)
}

// weblogicServers implements the WebLogicManagedServerInterface
type weblogicServers struct {
	client rest.Interface
	ns     string
}

// newWebLogicManagedServers returns a weblogicServers.
func newWebLogicManagedServers(c *WeblogicV1Client, namespace string) *weblogicServers {
	return &weblogicServers{client: c.RESTClient(), ns: namespace}
}

// Create takes the representation of a WebLogicManagedServer and creates it. Returns
// the server's representation of the WebLogicManagedServer, and an error, if there is
// any.
func (c *weblogicServers) Create(server *types.WebLogicManagedServer) (result *types.WebLogicManagedServer, err error) {
	result = &types.WebLogicManagedServer{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource(constants.WebLogicManagedServerResourceKindPlural).
		Body(server).
		Do().
		Into(result)
	return
}

// Update takes the representation of a WebLogicManagedServer and updates it. Returns
// the server's representation of the WebLogicManagedServer, and an error, if there is
// any.
func (c *weblogicServers) Update(server *types.WebLogicManagedServer) (result *types.WebLogicManagedServer, err error) {
	result = &types.WebLogicManagedServer{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource(constants.WebLogicManagedServerResourceKindPlural).
		Name(server.Name).
		Body(server).
		Do().
		Into(result)
	return
}

// Delete takes name of the server and deletes it. Returns an error if one
// occurs.
func (c *weblogicServers) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource(constants.WebLogicManagedServerResourceKindPlural).
		Name(name).
		Body(options).
		Do().
		Error()
}

// Get takes name of the WebLogicManagedServer, and returns the corresponding
// WebLogicManagedServer object, and an error if there is any.
func (c *weblogicServers) Get(name string, options metav1.GetOptions) (result *types.WebLogicManagedServer, err error) {
	result = &types.WebLogicManagedServer{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource(constants.WebLogicManagedServerResourceKindPlural).
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of
// WebLogicManagedServerList that match those selectors.
func (c *weblogicServers) List(opts metav1.ListOptions) (result *types.WebLogicManagedServerList, err error) {
	result = &types.WebLogicManagedServerList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource(constants.WebLogicManagedServerResourceKindPlural).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}
