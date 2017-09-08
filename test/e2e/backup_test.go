package e2e

import (
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"

	"weblogic-operator/pkg/server"
	"weblogic-operator/pkg/types"
	"weblogic-operator/test/e2e/framework"
	e2eutil "weblogic-operator/test/e2e/util"
)

func TestBackUpRestore(t *testing.T) {
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
	} else {
		// TODO: Needs to be wrapped in a more efficient 'WaitForXXX' handler.
		// wait for mysqld to start
		time.Sleep(time.Second * 20)
	}
	// Do we have a root password secret?
	secret, err := server.GetSecretForMySQLCluster(cl, f.KubeClient)
	if err != nil {
		t.Errorf("Error getting root password secret for cluster %s: %v", cl.Name, err)
	}
	// Create some simple test data to backup and restore.
	clustername := res.Name
	username := "root"
	podname := string(clustername + "-0")
	password := string(secret.Data["password"])
	executor := e2eutil.NewKubectlSimpleSQLExecutor(t, podname, username, password)
	dbHelper := e2eutil.NewMySQLDBTestHelper(t, executor)
	db1, table1, column1, value1 := createDbTableValue(t, dbHelper, "1")
	t.Logf("should backup db value: %v.%v.%v = %v", db1, table1, column1, value1)
	db2, table2, column2, value2 := createDbTableValue(t, dbHelper, "2")
	t.Logf("should backup db value: %v.%v.%v = %v", db2, table2, column2, value2)
	// Create a backup resource to perform the backup.
	backupName := "e2e-test-snapshot-backup-"
	ossCredsSecretRef := "bmcs-upload-credentials"
	backupSpec := e2eutil.NewMySQLBackup(clustername, backupName, ossCredsSecretRef)
	backup, err := f.MySQLOpClient.MySQLV1().MySQLBackups(f.Namespace).Create(backupSpec)
	if err != nil {
		t.Fatalf("Failed to create backup: %v", err)
	}
	backup, err = e2eutil.WaitForBackupPhase(
		t, backup, types.MySQLBackupPhaseComplete, e2eutil.NewDefaultRetyWithDuration(10), f.MySQLOpClient)
	if err != nil {
		t.Fatalf("Backup failed to reach phase %q: %v", types.MySQLClusterRunning, err)
	}
	// Check upload exists in OSS
	// TODO
	// example-mysql-cluster-with-volume-example-snapshot-backup.20170825165852.snapshot.mbi
}

func createDbTableValue(
	t *testing.T,
	dbh *e2eutil.MySQLDBTestHelper,
	ident string,
) (db string, table string, column string, value string) {
	db = "test_db_" + ident
	table = "test_table_" + ident
	column = "test_table_column_" + ident
	value = string(uuid.NewUUID())
	dbh.EnsureDBTableValue(db, table, column, value)
	if dbh.HasDBTableValue(db, table, column, value) {
		t.Logf("created db value: %v.%v.%v = %v", db, table, column, value)
	}
	return db, table, column, value
}
