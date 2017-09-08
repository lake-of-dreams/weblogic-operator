package util

import (
	"k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"weblogic-operator/pkg/types"
)

func NewMySQLCluster(genName string, replicas int32) *types.MySQLCluster {
	return &types.MySQLCluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       types.MySQLClusterCRDResourceKind,
			APIVersion: types.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: genName,
		},
		Spec: types.MySQLClusterSpec{
			Replicas: replicas,
		},
	}
}

func NewMySQLBackup(clusterName string, backupName string, ossCredsSecretRef string) *types.MySQLBackup {
	return &types.MySQLBackup{
		TypeMeta: metav1.TypeMeta{
			Kind:       types.MySQLBackupCRDResourceKind,
			APIVersion: types.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: backupName,
		},
		Cluster: &v1.LocalObjectReference{
			Name: clusterName,
		},
		SecretRef: &v1.LocalObjectReference{
			Name: ossCredsSecretRef,
		},
	}
}
