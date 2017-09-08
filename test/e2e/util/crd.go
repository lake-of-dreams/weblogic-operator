package util

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"weblogic-operator/pkg/clients/weblogic-operator"
	"weblogic-operator/pkg/types"
)

func CreateMySQLCluster(t *testing.T, mysqlopClient weblogic_operator.Interface, ns string, cluster *types.MySQLCluster) (*types.MySQLCluster, error) {
	cluster.Namespace = ns
	res, err := mysqlopClient.MySQLV1().MySQLClusters(ns).Create(cluster)
	if err != nil {
		return nil, err
	}
	t.Logf("Creating mysql cluster: %s", res.Name)
	return res, nil
}

// TODO(apryde): Wait for deletion of underlying resources.
func DeleteMySQLCluster(t *testing.T, mysqlopClient weblogic_operator.Interface, cluster *types.MySQLCluster) error {
	t.Logf("Deleting mysql cluster: %s", cluster.Name)
	err := mysqlopClient.MySQLV1().MySQLClusters(cluster.Namespace).Delete(cluster.Name, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
