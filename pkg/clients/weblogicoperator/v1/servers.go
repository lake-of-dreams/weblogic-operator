package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"weblogic-operator/pkg/types"
	"weblogic-operator/pkg/constants"
)

// WeblogicServersGetter has a method to return a WeblogicServerInterface.
type WeblogicServersGetter interface {
	WeblogicServers(namespace string) WeblogicServerInterface
}

// WeblogicServerInterface has methods to work with Weblogic Server custom
// resources.
type WeblogicServerInterface interface {
	Create(*types.WeblogicServer) (*types.WeblogicServer, error)
	Update(*types.WeblogicServer) (*types.WeblogicServer, error)
	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*types.WeblogicServer, error)
	List(opts metav1.ListOptions) (*types.WeblogicServerList, error)
}

// weblogicServers implements the WeblogicServerInterface
type weblogicServers struct {
	client rest.Interface
	ns     string
}

// newWeblogicServers returns a weblogicServers.
func newWeblogicServers(c *WeblogicV1Client, namespace string) *weblogicServers {
	return &weblogicServers{client: c.RESTClient(), ns: namespace}
}

// Create takes the representation of a WeblogicServer and creates it. Returns
// the server's representation of the WeblogicServer, and an error, if there is
// any.
func (c *weblogicServers) Create(server *types.WeblogicServer) (result *types.WeblogicServer, err error) {
	result = &types.WeblogicServer{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource(constants.WeblogicServerResourceKindPlural).
		Body(server).
		Do().
		Into(result)
	return
}

// Update takes the representation of a WeblogicServer and updates it. Returns
// the server's representation of the WeblogicServer, and an error, if there is
// any.
func (c *weblogicServers) Update(server *types.WeblogicServer) (result *types.WeblogicServer, err error) {
	result = &types.WeblogicServer{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource(constants.WeblogicServerResourceKindPlural).
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
		Resource(constants.WeblogicServerResourceKindPlural).
		Name(name).
		Body(options).
		Do().
		Error()
}

// Get takes name of the WeblogicServer, and returns the corresponding
// WeblogicServer object, and an error if there is any.
func (c *weblogicServers) Get(name string, options metav1.GetOptions) (result *types.WeblogicServer, err error) {
	result = &types.WeblogicServer{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource(constants.WeblogicServerResourceKindPlural).
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of
// WeblogicServerList that match those selectors.
func (c *weblogicServers) List(opts metav1.ListOptions) (result *types.WeblogicServerList, err error) {
	result = &types.WeblogicServerList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource(constants.WeblogicServerResourceKindPlural).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}
