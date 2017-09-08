package util

import (
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	"weblogic-operator/pkg/clients/weblogic-operator"
	"weblogic-operator/pkg/types"
	"weblogic-operator/pkg/util/retry"
)

// DefaultRetry is the default backoff for e2e tests.
var DefaultRetry = wait.Backoff{
	Steps:    3,
	Duration: 10 * time.Second,
	Factor:   1.0,
	Jitter:   0.1,
}

// NewDefaultRetyWithDuration creates a customized backoff for e2e tests.
func NewDefaultRetyWithDuration(seconds time.Duration) wait.Backoff {
	return wait.Backoff{
		Steps:    3,
		Duration: seconds * time.Second,
		Factor:   1.0,
		Jitter:   0.1,
	}
}

// WaitForClusterPhase retries until a cluster reaches a given phase or a
// timeout is reached.
func WaitForClusterPhase(
	t *testing.T,
	cluster *types.MySQLCluster,
	phase types.MySQLClusterPhase,
	backoff wait.Backoff,
	mySQLOpClient weblogic_operator.Interface,
) (*types.MySQLCluster, error) {
	t.Logf("Waiting for cluster %s to reach phase %s...", cluster.Name, phase)
	var cl *types.MySQLCluster
	var err error
	err = retry.Retry(backoff, func() (bool, error) {
		cl, err = mySQLOpClient.MySQLV1().MySQLClusters(cluster.Namespace).Get(cluster.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		t.Logf("%s in phase %s, want %s", cluster.Name, cl.Status.Phase, phase)
		return cl.Status.Phase == phase, err
	})
	if err != nil {
		return nil, err
	}
	return cl, nil
}

// WaitForBackupPhase retries until a backup completes or timeout is reached.
func WaitForBackupPhase(
	t *testing.T,
	backup *types.MySQLBackup,
	phase types.MySQLBackupPhase,
	backoff wait.Backoff,
	mySQLOpClient weblogic_operator.Interface,
) (*types.MySQLBackup, error) {
	t.Logf("Waiting for backup %s to reach phase %s...", backup.Name, phase)
	var latest *types.MySQLBackup
	var err error
	err = retry.Retry(backoff, func() (bool, error) {
		latest, err = mySQLOpClient.MySQLV1().MySQLBackups(backup.Namespace).Get(backup.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		t.Logf("%s in phase %s, want %s", backup.Name, latest.Status.Phase, phase)
		return latest.Status.Phase == phase, err
	})
	if err != nil {
		return nil, err
	}
	return latest, nil
}
