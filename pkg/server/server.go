// Copyright 2017 The mysql-operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"fmt"

	"k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/golang/glog"

	"gitlab-odx.oracle.com/odx/mysql-operator/pkg/constants"
	"gitlab-odx.oracle.com/odx/mysql-operator/pkg/resources/secrets"
	"gitlab-odx.oracle.com/odx/mysql-operator/pkg/resources/services"
	"gitlab-odx.oracle.com/odx/mysql-operator/pkg/resources/statefulsets"
	"gitlab-odx.oracle.com/odx/mysql-operator/pkg/types"
)

// HasClusterNameLabel returns true if the given labels map matches the given
// cluster name.
func HasClusterNameLabel(labels map[string]string, clustername string) bool {
	for label, value := range labels {
		if label == constants.MySQLClusterLabel {
			if value == clustername {
				return true
			}
		}
	}
	return false
}

// Return a label that uniquely identifies a MySQL cluster
func getLabelSelectorForCluster(cluster *types.MySQLCluster) string {
	return fmt.Sprintf("%s=%s", constants.MySQLClusterLabel, cluster.Name)
}

// GetStatefulSetForCluster finds the associated StatefulSet for a MySQL cluster
func GetStatefulSetForCluster(cluster *types.MySQLCluster, kubeClient kubernetes.Interface) (*v1beta1.StatefulSet, error) {
	opts := metav1.ListOptions{LabelSelector: getLabelSelectorForCluster(cluster)}
	statefulsets, err := kubeClient.AppsV1beta1().StatefulSets(cluster.Namespace).List(opts)
	if err != nil {
		glog.Errorf("Unable to list stateful sets for %s: %s", cluster.Name, err)
		return nil, err
	}

	for _, ss := range statefulsets.Items {
		if HasClusterNameLabel(ss.Labels, cluster.Name) {
			return &ss, nil
		}
	}
	return nil, nil
}

// GetServiceForMySQLCluster returns the associated service for a given cluster
func GetServiceForMySQLCluster(cluster *types.MySQLCluster, clientset kubernetes.Interface) (*v1.Service, error) {
	opts := metav1.ListOptions{LabelSelector: getLabelSelectorForCluster(cluster)}
	services, err := clientset.CoreV1().Services(cluster.Namespace).List(opts)
	if err != nil {
		glog.Errorf("Unable to list services for %s: %s", cluster.Name, err)
		return nil, err
	}

	for _, svc := range services.Items {
		if HasClusterNameLabel(svc.Labels, cluster.Name) {
			return &svc, nil
		}
	}
	return nil, nil
}

// GetSecretForMySQLCluster returns the root password secret for a given MySQL
// cluster.
func GetSecretForMySQLCluster(cluster *types.MySQLCluster, clientset kubernetes.Interface) (*v1.Secret, error) {
	opts := metav1.ListOptions{LabelSelector: getLabelSelectorForCluster(cluster)}
	r, err := clientset.CoreV1().Secrets(cluster.Namespace).List(opts)
	if err != nil {
		glog.Errorf("Unable to list secrets for %s: %s", cluster.Name, err)
		return nil, err
	}

	for _, secret := range r.Items {
		if HasClusterNameLabel(secret.Labels, cluster.Name) {
			return &secret, nil
		}
	}
	return nil, nil
}

func updateCluster(cluster *types.MySQLCluster, restClient *rest.RESTClient) error {
	// TODO(apryde): Use retry.RetryOnConflict()?
	result := restClient.Put().
		Resource(types.ClusterCRDResourcePlural).
		Namespace(cluster.Namespace).
		Name(cluster.Name).
		Body(cluster).
		Do()
	return result.Error()
}

func setMySQLClusterState(cluster *types.MySQLCluster, restClient *rest.RESTClient, phase types.MySQLClusterPhase, err error) error {
	modified := false
	if cluster.Status.Phase != phase {
		cluster.Status.Phase = phase
		modified = true
	}

	l := len(cluster.Status.Errors)
	if err != nil && (l < 1 || cluster.Status.Errors[l-1] != err.Error()) {
		cluster.Status.Errors = append(cluster.Status.Errors, err.Error())
		modified = true
	} else if l == 0 {
		cluster.Status.Errors = []string{}
		modified = true
	}

	// TODO(apryde): Use retry.RetryOnConflict()?
	if modified {
		result := restClient.Put().
			Resource(types.ClusterCRDResourcePlural).
			Namespace(cluster.Namespace).
			Name(cluster.Name).
			Body(cluster).
			Do()

		return result.Error()
	}

	return nil
}

func createCluster(cluster *types.MySQLCluster, kubeClient kubernetes.Interface, restClient *rest.RESTClient) error {
	cluster.EnsureDefaults()

	err := cluster.Validate()
	if err != nil {
		return err
	}

	// Validate that a label is set on the cluster
	if !HasClusterNameLabel(cluster.Labels, cluster.Name) {
		glog.V(4).Infof("Setting label on cluster %s", getLabelSelectorForCluster(cluster))
		if cluster.Labels == nil {
			cluster.Labels = make(map[string]string)
		}
		cluster.Labels[constants.MySQLClusterLabel] = cluster.Name
		return updateCluster(cluster, restClient)
	}

	if cluster.RequiresSecret() {
		_, err = CreateSecret(kubeClient, cluster)
		if err != nil {
			return err
		}
	}

	clusterService, err := CreateServiceForMySQLCluster(kubeClient, cluster)
	if err != nil {
		return err
	}

	_, err = CreateStatefulSet(kubeClient, cluster, clusterService)
	if err != nil {
		return err
	}

	return nil
}

// When delete cluster is called we will delete the stateful set (which also deletes the associated service) and
// delete any secrets associated with the cluster
func deleteCluster(cluster *types.MySQLCluster, kubeClient kubernetes.Interface, restClient *rest.RESTClient) error {
	err := cluster.Validate()
	if err != nil {
		return err
	}

	err = DeleteStatefulSet(kubeClient, cluster)
	if err != nil {
		return err
	}

	err = DeleteService(kubeClient, cluster)
	if err != nil {
		return err
	}

	if cluster.RequiresSecret() {
		err = DeleteSecret(kubeClient, cluster)
		if err != nil {
			return err
		}
	}

	return nil
}

func updateClusterWithStatefulSet(cluster *types.MySQLCluster,
	statefulSet *v1beta1.StatefulSet,
	kubeClient kubernetes.Interface,
	restClient *rest.RESTClient) (err error) {
	// Some simple logic for the time being.
	// To add
	// connection to the cluster
	// validate each pod?
	// Check how a rolling upgrade effects this
	// check version of each pod

	if statefulSet.Status.ReadyReplicas < statefulSet.Status.Replicas {
		setMySQLClusterState(cluster, restClient, types.MySQLClusterPending, nil)
	} else if statefulSet.Status.ReadyReplicas == statefulSet.Status.Replicas {
		setMySQLClusterState(cluster, restClient, types.MySQLClusterRunning, nil)
	}
	return err
}

// CreateStatefulSet will create a new Kubernetes StatefulSet based on a predefined template
func CreateStatefulSet(clientset kubernetes.Interface, cluster *types.MySQLCluster, service *v1.Service) (*v1beta1.StatefulSet, error) {
	// Find StatefulSet and if it does not exist create it
	existingStatefulSet, err := GetStatefulSetForCluster(cluster, clientset)
	if err != nil {
		glog.Errorf("Error finding stateful set for cluster: %v", err)
		return nil, err
	}

	if existingStatefulSet != nil {
		glog.V(2).Infof("Stateful set with label %s already exists", getLabelSelectorForCluster(cluster))
		return existingStatefulSet, nil
	}

	glog.V(4).Infof("Creating a new stateful set for cluster %s", cluster.Name)
	ss := statefulsets.NewForCluster(cluster, service.Name)

	glog.V(4).Infof("Creating cluster %+v", ss)
	return clientset.AppsV1beta1().StatefulSets(cluster.Namespace).Create(ss)
}

func GetClusterForStatefulSet(statefulSet *v1beta1.StatefulSet, restClient *rest.RESTClient) (cluster *types.MySQLCluster, err error) {
	if mySQLClusterName, ok := statefulSet.Labels[constants.MySQLClusterLabel]; ok {
		cluster = &types.MySQLCluster{}
		result := restClient.Get().
			Resource(types.ClusterCRDResourcePlural).
			Namespace(statefulSet.Namespace).
			Name(mySQLClusterName).
			Do().
			Into(cluster)
		return cluster, result
	}
	return nil, fmt.Errorf("unable to get Label %s from statefulset. Not part of cluster", constants.MySQLClusterLabel)
}

// DeleteStatefulSet will delete a stateful set by name
func DeleteStatefulSet(clientset kubernetes.Interface, cluster *types.MySQLCluster) error {
	statefulSet, err := GetStatefulSetForCluster(cluster, clientset)
	if err != nil || statefulSet == nil {
		glog.Errorf("Could not delete stateful set: %s", err)
		return err
	}

	glog.V(4).Infof("Deleting stateful set %s", statefulSet.Name)
	var policy = metav1.DeletePropagationBackground
	return clientset.AppsV1beta1().
		StatefulSets(cluster.Namespace).
		Delete(statefulSet.Name, &metav1.DeleteOptions{PropagationPolicy: &policy})
}

// CreateServiceForMySQLCluster will create a new Kubernetes Service based on a predefined template
func CreateServiceForMySQLCluster(clientset kubernetes.Interface, cluster *types.MySQLCluster) (*v1.Service, error) {
	// Find Service and if it does not exist create it
	existingService, err := GetServiceForMySQLCluster(cluster, clientset)
	if err != nil {
		glog.Errorf("Error finding service for cluster: %s", err)
		return nil, err
	}

	if existingService != nil {
		glog.V(2).Infof("Service with label %s already exists", getLabelSelectorForCluster(cluster))
		return existingService, nil
	}

	glog.V(4).Infof("Creating a new service for cluster %s", cluster.Name)

	svc := services.NewForCluster(cluster)
	return clientset.CoreV1().Services(cluster.Namespace).Create(svc)
}

// DeleteService deletes the Service associated with a MySQL cluster.
func DeleteService(clientset kubernetes.Interface, cluster *types.MySQLCluster) error {
	service, err := GetServiceForMySQLCluster(cluster, clientset)
	if err != nil || service == nil {
		glog.Errorf("Could not delete service: %s", err)
		return err
	}
	glog.V(4).Infof("Deleting service %s", service.Name)
	return clientset.CoreV1().Services(cluster.Namespace).Delete(service.Name, nil)
}

// CreateSecret creates the Secret associated with a MySQL cluster.
func CreateSecret(clientset kubernetes.Interface, cluster *types.MySQLCluster) (*v1.Secret, error) {
	existingSecret, err := GetSecretForMySQLCluster(cluster, clientset)
	if err != nil {
		return nil, err
	}

	if existingSecret != nil {
		glog.V(2).Infof("Secret with label %s already exists", getLabelSelectorForCluster(cluster))
		return existingSecret, nil
	}
	glog.V(4).Infof("Creating a new secret for cluster %s", cluster.Name)
	secret := secrets.NewMysqlRootPassword(cluster)
	return clientset.CoreV1().Secrets(cluster.Namespace).Create(secret)
}

// DeleteSecret will delete the MySQL secret for a given cluster if it exists
func DeleteSecret(clientset kubernetes.Interface, cluster *types.MySQLCluster) error {
	secret, err := GetSecretForMySQLCluster(cluster, clientset)
	if err != nil || secret == nil {
		glog.Errorf("Unable to find secret %s for deletion: %s", cluster.Name, err)
		return err
	}
	glog.V(4).Infof("Deleting secret %s", secret.Name)
	return clientset.CoreV1().Secrets(cluster.Namespace).Delete(secret.Name, nil)
}
