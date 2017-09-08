package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"weblogic-operator/pkg/types"
)

// MySQLBackupsGetter has a method to return a MySQLBackupInterface.
type MySQLBackupsGetter interface {
	MySQLBackups(namespace string) MySQLBackupInterface
}

// MySQLBackupInterface has methods to work with MySQL backup custom resources.
type MySQLBackupInterface interface {
	Create(*types.MySQLBackup) (*types.MySQLBackup, error)
	Update(*types.MySQLBackup) (*types.MySQLBackup, error)
	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*types.MySQLBackup, error)
	List(opts metav1.ListOptions) (*types.MySQLBackupList, error)
}

// MySQLBackups implements the MySQLBackupInterface
type mySQLBackups struct {
	client rest.Interface
	ns     string
}

// newMySQLBackups returns a mySQLBackups.
func newMySQLBackups(c *MySQLV1Client, namespace string) *mySQLBackups {
	return &mySQLBackups{client: c.RESTClient(), ns: namespace}
}

// Create takes the representation of a MySQLBackup and creates it. Returns
// the server's representation of the MySQLBackup, and an error, if there is
// any.
func (c *mySQLBackups) Create(backup *types.MySQLBackup) (result *types.MySQLBackup, err error) {
	result = &types.MySQLBackup{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource(types.BackupCRDResourcePlural).
		Body(backup).
		Do().
		Into(result)
	return
}

// Update takes the representation of a MySQLBackup and updates it. Returns
// the server's representation of the MySQLBackup, and an error, if there is
// any.
func (c *mySQLBackups) Update(backup *types.MySQLBackup) (result *types.MySQLBackup, err error) {
	result = &types.MySQLBackup{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource(types.BackupCRDResourcePlural).
		Name(backup.Name).
		Body(backup).
		Do().
		Into(result)
	return
}

// Delete takes name of the backup and deletes it. Returns an error if one
// occurs.
func (c *mySQLBackups) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource(types.BackupCRDResourcePlural).
		Name(name).
		Body(options).
		Do().
		Error()
}

// Get takes name of the MySQLBackup, and returns the corresponding
// MySQLBackup object, and an error if there is any.
func (c *mySQLBackups) Get(name string, options metav1.GetOptions) (result *types.MySQLBackup, err error) {
	result = &types.MySQLBackup{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource(types.BackupCRDResourcePlural).
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of
// MySQLBackupList that match those selectors.
func (c *mySQLBackups) List(opts metav1.ListOptions) (result *types.MySQLBackupList, err error) {
	result = &types.MySQLBackupList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource(types.BackupCRDResourcePlural).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}
