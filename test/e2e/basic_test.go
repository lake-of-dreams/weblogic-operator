package e2e

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"weblogic-operator/pkg/server"
	"weblogic-operator/pkg/types"
	"weblogic-operator/test/e2e/framework"
	e2eutil "weblogic-operator/test/e2e/util"
)

func TestCreateCluster(t *testing.T) {
	f := framework.Global
	replicas := int32(1)

	res, err := f.MySQLOpClient.MySQLV1().
		MySQLClusters(f.Namespace).
		Create(e2eutil.NewMySQLCluster("test-create-cluster-", replicas))
	if err != nil {
		t.Fatalf("Failed to create cluster: %v", err)
	}
	defer func() {
		err = f.MySQLOpClient.MySQLV1().
			MySQLClusters(f.Namespace).
			Delete(res.Name, &metav1.DeleteOptions{})
		if err != nil {
			t.Fatalf("Failed clean up cluster: %v", err)
		}
	}()

	cl, err := e2eutil.WaitForClusterPhase(t, res, types.MySQLClusterRunning, e2eutil.DefaultRetry, f.MySQLOpClient)
	if err != nil {
		t.Fatalf("Cluster failed to reach phase %q: %v", types.MySQLClusterRunning, err)
	}

	if cl.Spec.Replicas != replicas {
		t.Errorf("Got cluster with %d replica(s), want %d", cl.Spec.Replicas, replicas)
	}

	// Do we have a statefulset?
	ss, err := server.GetStatefulSetForCluster(cl, f.KubeClient)
	if err != nil {
		t.Errorf("Error getting statefulset for cluster %s: %v", cl.Name, err)
	} else {
		if ss.Status.ReadyReplicas != replicas {
			t.Logf("%#v", ss.Status)
			t.Errorf("Got statefulset with %d ready replica(s), want %d", ss.Status.ReadyReplicas, replicas)
		}
	}

	// Do we have a service?
	_, err = server.GetServiceForMySQLCluster(cl, f.KubeClient)
	if err != nil {
		t.Errorf("Error getting service for cluster %s: %v", cl.Name, err)
	}

	// Do we have a root password secret?
	_, err = server.GetSecretForMySQLCluster(cl, f.KubeClient)
	if err != nil {
		t.Errorf("Error getting root password secret for cluster %s: %v", cl.Name, err)
	}
}
