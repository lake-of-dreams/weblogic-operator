package server

import (
	"fmt"

	autoscalingv1 "k8s.io/api/autoscaling/v1"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/golang/glog"

	"strings"
	"weblogic-operator/pkg/constants"
	"weblogic-operator/pkg/resources/horizontalpodautoscalers"
	"weblogic-operator/pkg/resources/replicasets"
	"weblogic-operator/pkg/resources/services"
	"weblogic-operator/pkg/types"
)

// HasServerNameLabel returns true if the given labels map matches the given
// server name.
func HasServerNameLabel(labels map[string]string, servername string) bool {
	for label, value := range labels {
		if label == constants.WebLogicManagedServerLabel {
			if value == servername {
				return true
			}
		}
	}
	return false
}

// Return a label that uniquely identifies a Weblogic server
func getLabelSelectorForServer(server *types.WebLogicManagedServer) string {
	return fmt.Sprintf("%s=%s", constants.WebLogicManagedServerLabel, server.Name)
}

// GetReplicaSetForWebLogicManagedServer finds the associated ReplicaSet for a Weblogic server
func GetReplicaSetForWebLogicManagedServer(server *types.WebLogicManagedServer, kubeClient kubernetes.Interface) (*v1beta1.ReplicaSet, error) {
	opts := metav1.ListOptions{LabelSelector: getLabelSelectorForServer(server)}
	replicasets, err := kubeClient.ExtensionsV1beta1().ReplicaSets(server.Namespace).List(opts)
	if err != nil {
		glog.Errorf("Unable to list replica sets for %s: %s", server.Name, err)
		return nil, err
	}

	for _, rc := range replicasets.Items {
		if HasServerNameLabel(rc.Labels, server.Name) {
			return &rc, nil
		}
	}
	return nil, nil
}

// CreateReplicaSetForWebLogicManagedServer will create a new Kubernetes ReplicaSet based on a predefined template
func CreateReplicaSetForWebLogicManagedServer(clientset kubernetes.Interface, server *types.WebLogicManagedServer, service *v1.Service) (controller *v1beta1.ReplicaSet, err error) {
	// Find ReplicaSet and if it does not exist create it
	existingReplicaSet, err := GetReplicaSetForWebLogicManagedServer(server, clientset)
	if err != nil {
		glog.Errorf("Error finding replica set for server: %v", err)
		return nil, err
	}

	if existingReplicaSet != nil {
		glog.V(2).Infof("Replica set with label %s already exists", getLabelSelectorForServer(server))
		return existingReplicaSet, nil
	}

	glog.V(4).Infof("Creating a new replica set for server %s", server.Name)
	rs := replicasets.NewForServer(server, service.Name)

	glog.V(4).Infof("Creating server %+v", rs)
	return clientset.ExtensionsV1beta1().ReplicaSets(server.Namespace).Create(rs)
}

func UpdateReplicaSetForWebLogicManagedServer(clientset kubernetes.Interface, server *types.WebLogicManagedServer, service *v1.Service) (controller *v1beta1.ReplicaSet, err error) {
	// Find ReplicaSet and if it does not exist create it
	existingReplicaSet, err := GetReplicaSetForWebLogicManagedServer(server, clientset)
	if err != nil {
		glog.Errorf("Error finding replica set for server: %v", err)
		return nil, err
	}

	if existingReplicaSet != nil {
		glog.V(2).Infof("Updating existing Replica set with label %s", getLabelSelectorForServer(server))

		glog.V(4).Infof("Creating updated replica set for server %s", server.Name)
		rs := replicasets.NewForServer(server, service.Name)

		glog.V(4).Infof("Creating server %+v", rs)
		return clientset.ExtensionsV1beta1().ReplicaSets(server.Namespace).Update(rs)
	}

	return nil, nil
}

// DeleteReplicaSetForWebLogicManagedServer will delete a replica set by name
func DeleteReplicaSetForWebLogicManagedServer(clientset kubernetes.Interface, server *types.WebLogicManagedServer) error {

	if strings.EqualFold(server.Name, constants.HorizontalPodAutoscalerTargetLabel) {
		glog.V(4).Infof("Deleting HPA for managed server!!!!!!!!!!!!!")
		err := DeleteHorizontalPodAutoscalerForWebLogicManagedServer(clientset, server)
		if err != nil {
			return err
		}
	}

	replicaSet, err := GetReplicaSetForWebLogicManagedServer(server, clientset)
	if err != nil || replicaSet == nil {
		glog.Errorf("Could not delete replica set: %s", err)
		return err
	}

	glog.V(4).Infof("Deleting replica set %s", replicaSet.Name)
	var policy = metav1.DeletePropagationBackground
	return clientset.ExtensionsV1beta1().
		ReplicaSets(server.Namespace).
		Delete(replicaSet.Name, &metav1.DeleteOptions{PropagationPolicy: &policy})
}

// GetHorizontalPodAutoscalerForWebLogicManagedServer finds the associated ReplicaSet for a Weblogic server
func GetHorizontalPodAutoscalerForWebLogicManagedServer(server *types.WebLogicManagedServer, kubeClient kubernetes.Interface) (*autoscalingv1.HorizontalPodAutoscaler, error) {
	opts := metav1.ListOptions{LabelSelector: getLabelSelectorForServer(server)}
	horizontalpodautoscalers, err := kubeClient.AutoscalingV1().HorizontalPodAutoscalers(server.Namespace).List(opts)
	if err != nil {
		glog.Errorf("Unable to list horizontal pod autoscalers for %s: %s", server.Name, err)
		return nil, err
	}

	for _, rc := range horizontalpodautoscalers.Items {
		if HasServerNameLabel(rc.Labels, server.Name) {
			glog.V(4).Infof("Value of RC!!!")
			return &rc, nil
		}
	}
	return nil, nil
}

// CreateHorizontalPodAutoscalerForWebLogicManagedServer will create a new Kubernetes HorizontalPodAutoscaler based on a predefined template
func CreateHorizontalPodAutoscalerForWebLogicManagedServer(clientset kubernetes.Interface, server *types.WebLogicManagedServer, service *v1.Service) (controller *autoscalingv1.HorizontalPodAutoscaler, err error) {
	// Find ReplicaSet and if it does not exist create it
	existingHorizontalPodAutoscaler, err := GetHorizontalPodAutoscalerForWebLogicManagedServer(server, clientset)
	if err != nil {
		glog.Errorf("Error finding Horizontal Pod Autoscaler for server: %v", err)
		return nil, err
	}

	if existingHorizontalPodAutoscaler != nil {
		glog.V(2).Infof("Replica set with label %s already exists", getLabelSelectorForServer(server))
		return existingHorizontalPodAutoscaler, nil
	}

	glog.V(4).Infof("Creating a new Horizontal Pod Autoscalers for server %s", server.Name)
	rs := horizontalpodautoscalers.NewForHorizontalPodAutoscaling(server, service.Name)

	glog.V(4).Infof("Creating server %+v", rs)
	return clientset.AutoscalingV1().HorizontalPodAutoscalers(server.Namespace).Create(rs)
}

// DeleteHorizontalPodAutoscalerForWebLogicManagedServer will delete a replica set by name
func DeleteHorizontalPodAutoscalerForWebLogicManagedServer(clientset kubernetes.Interface, server *types.WebLogicManagedServer) error {
	horizontalPodAutoscaler, err := GetHorizontalPodAutoscalerForWebLogicManagedServer(server, clientset)
	glog.V(4).Infof("value of horizontalPodAutoscaler")
	if err != nil && horizontalPodAutoscaler == nil {
		glog.Errorf("Could not delete Horizontal Pod Autoscaler: %s", err)
		return err
	}

	glog.V(4).Infof("Deleting Horizontal Pod Autoscaler %s", horizontalPodAutoscaler.Name)
	var policy = metav1.DeletePropagationBackground
	return clientset.AutoscalingV1().HorizontalPodAutoscalers(server.Namespace).
		Delete(horizontalPodAutoscaler.Name, &metav1.DeleteOptions{PropagationPolicy: &policy})
}

func createWebLogicManagedServer(server *types.WebLogicManagedServer, kubeClient kubernetes.Interface, restClient *rest.RESTClient) error {
	server.EnsureDefaults()
	server.PopulateDomain()

	// Validate that a label is set on the server
	if !HasServerNameLabel(server.Labels, server.Name) {
		glog.V(4).Infof("Setting label on server %s", getLabelSelectorForServer(server))
		if server.Labels == nil {
			server.Labels = make(map[string]string)
		}
		server.Labels[constants.WebLogicManagedServerLabel] = server.Name
		server.Labels[server.Spec.Domain.Name] = "managedserver"
		return updateWebLogicManagedServerLabel(server, restClient)
	}

	serverService, err := CreateServiceForWebLogicManagedServer(kubeClient, server)
	if err != nil {
		return err
	}

	_, err = CreateReplicaSetForWebLogicManagedServer(kubeClient, server, serverService)
	if err != nil {
		return err
	}

	_, err = CreateHorizontalPodAutoscalerForWebLogicManagedServer(kubeClient, server, serverService)
	if err != nil {
		return err
	}

	return nil
}

//TODO update the replica set
func updateWebLogicManagedServer(server *types.WebLogicManagedServer, kubeClient kubernetes.Interface, restClient *rest.RESTClient) error {
	// Find Service and if it does not exist create it
	existingService, err := GetServiceForWebLogicManagedServer(server, kubeClient)
	if err != nil {
		glog.Errorf("Error finding service for server: %s", err)
		return err
	}
	_, err = UpdateReplicaSetForWebLogicManagedServer(kubeClient, server, existingService)
	if err != nil {
		return err
	}

	return nil
}

func updateWebLogicManagedServerLabel(server *types.WebLogicManagedServer, restClient *rest.RESTClient) error {
	result := restClient.Put().
		Resource(constants.WebLogicManagedServerResourceKindPlural).
		Namespace(server.Namespace).
		Name(server.Name).
		Body(server).
		Do()
	return result.Error()
}

// When delete server is called we will delete the stateful set (which also deletes the associated service)
//TODO handling to call stopWeblogic.sh needs to be done here
func deleteWebLogicManagedServer(server *types.WebLogicManagedServer, kubeClient kubernetes.Interface, restClient *rest.RESTClient) error {
	//err = RunStopForWebLogicManagedServer(kubeClient, restClient, server)
	//if err != nil {
	//	return err
	//}

	err := DeleteReplicaSetForWebLogicManagedServer(kubeClient, server)
	if err != nil {
		return err
	}

	err = DeleteServiceForWebLogicManagedServer(kubeClient, server)
	if err != nil {
		return err
	}

	return nil
}

// GetServiceForWebLogicManagedServer returns the associated service for a given server
func GetServiceForWebLogicManagedServer(server *types.WebLogicManagedServer, clientset kubernetes.Interface) (*v1.Service, error) {
	opts := metav1.ListOptions{LabelSelector: getLabelSelectorForServer(server)}
	services, err := clientset.CoreV1().Services(server.Namespace).List(opts)
	if err != nil {
		glog.Errorf("Unable to list services for %s: %s", server.Name, err)
		return nil, err
	}

	for _, svc := range services.Items {
		if HasServerNameLabel(svc.Labels, server.Name) {
			return &svc, nil
		}
	}
	return nil, nil
}

// CreateServiceForWebLogicManagedServer will create a new Kubernetes Service based on a predefined template
func CreateServiceForWebLogicManagedServer(clientset kubernetes.Interface, server *types.WebLogicManagedServer) (*v1.Service, error) {
	// Find Service and if it does not exist create it
	existingService, err := GetServiceForWebLogicManagedServer(server, clientset)
	if err != nil {
		glog.Errorf("Error finding service for server: %s", err)
		return nil, err
	}

	if existingService != nil {
		glog.V(2).Infof("Service with label %s already exists", getLabelSelectorForServer(server))
		return existingService, nil
	}

	glog.V(4).Infof("Creating a new service for server %s", server.Name)

	svc := services.NewServiceForServer(server)
	return clientset.CoreV1().Services(server.Namespace).Create(svc)
}

// DeleteServiceForWebLogicManagedServer deletes the Service associated with a Weblogic server.
func DeleteServiceForWebLogicManagedServer(clientset kubernetes.Interface, server *types.WebLogicManagedServer) error {
	service, err := GetServiceForWebLogicManagedServer(server, clientset)
	if err != nil || service == nil {
		glog.Errorf("Could not delete service: %s", err)
		return err
	}
	glog.V(4).Infof("Deleting service %s", service.Name)
	return clientset.CoreV1().Services(server.Namespace).Delete(service.Name, nil)
}

func GetServerForReplicaSet(replicaset *v1beta1.ReplicaSet, restClient *rest.RESTClient) (server *types.WebLogicManagedServer, err error) {
	if weblogicServerName, ok := replicaset.Labels[constants.WebLogicManagedServerLabel]; ok {
		server = &types.WebLogicManagedServer{}
		result := restClient.Get().
			Resource(constants.WebLogicManagedServerResourceKindPlural).
			Namespace(replicaset.Namespace).
			Name(weblogicServerName).
			Do().
			Into(server)
		return server, result
	}
	return nil, fmt.Errorf("unable to get Label %s from replicaset. Not part of server", constants.WebLogicManagedServerLabel)
}

func updateServerWithReplicaSet(server *types.WebLogicManagedServer, replicaSet *v1beta1.ReplicaSet, kubeClient kubernetes.Interface, restClient *rest.RESTClient) (err error) {
	// Some simple logic for the time being.
	// To add
	// connection to the server
	// validate each pod?
	// Check how a rolling upgrade effects this
	// check version of each pod

	return nil
}

func GetServerForHorizontalPodAutoscaler(horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler, restClient *rest.RESTClient) (server *types.WebLogicManagedServer, err error) {
	if weblogicServerName, ok := horizontalPodAutoscaler.Labels[constants.HorizontalPodAutoscalerTargetLabel]; ok {
		server = &types.WebLogicManagedServer{}
		result := restClient.Get().
			Resource(constants.WebLogicManagedServerResourceKindPlural).
			Namespace(horizontalPodAutoscaler.Namespace).
			Name(weblogicServerName).
			Do().
			Into(server)
		return server, result
	}
	return nil, fmt.Errorf("unable to get Label %s from horizontalPodAutoscaler. Not part of server", constants.HorizontalPodAutoscalerTargetLabel)
}

func updateServerWithHorizontalPodAutoscaler(server *types.WebLogicManagedServer, horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler, kubeClient kubernetes.Interface, restClient *rest.RESTClient) (err error) {
	// Some simple logic for the time being.
	// To add
	// connection to the server
	// validate each pod?
	// Check how a rolling upgrade effects this
	// check version of each pod

	return nil
}

// GetPodForWebLogicManagedServer finds the associated pod for a Weblogic server TODO
func GetPodForWebLogicManagedServer(server *types.WebLogicManagedServer, clientset kubernetes.Interface) (*v1.Pod, error) {
	opts := metav1.ListOptions{LabelSelector: getLabelSelectorForServer(server)}
	pods, err := clientset.CoreV1().Pods(server.Namespace).List(opts)
	if err != nil {
		glog.Errorf("Unable to list pods for %s: %s", server.Name, err)
		return nil, err
	}

	for _, pod := range pods.Items {
		if HasServerNameLabel(pod.Labels, server.Name) {
			return &pod, nil
		}
	}
	return nil, nil
}

// GetPodForWebLogicManagedServer finds the associated pod for a Weblogic server TODO
func GetContainerForPod(server *types.WebLogicManagedServer, pod *v1.Pod) (*v1.Container, error) {
	containers := pod.Spec.Containers

	for _, container := range containers {
		if container.Name == server.Name {
			return &container, nil
		}
	}

	return nil, nil
}

// RunStopForWebLogicManagedServer will run stopWebLogic to stop a Weblogic server container in a pod TODO
func RunStopForWebLogicManagedServer(clientset kubernetes.Interface, restClient *rest.RESTClient, server *types.WebLogicManagedServer) error {
	pod, err := GetPodForWebLogicManagedServer(server, clientset)
	if err != nil || pod == nil {
		glog.Errorf("Could not find pod: %s", err)
		return err
	}

	container, err := GetContainerForPod(server, pod)
	if err != nil || container == nil {
		glog.Errorf("Could not find container %s in pod %s: %s", server.Name, pod.Name, err)
		return err
	}

	glog.V(4).Infof("Running stopWeblogic.sh for container %s in pod %s", server.Name, pod.Name)
	command := []string{"/u01/oracle/user_projects/domains/base_domain/bin/stopWebLogic.sh"}
	cmdErr := ExecuteCommandInContainer(restClient, pod, container, command)
	if cmdErr != nil {
		glog.Errorf("Error executing command : %s", cmdErr)
		return cmdErr
	}

	return nil
}

// ExecuteCommandInContainer will run a command in a container in a pod TODO
func ExecuteCommandInContainer(restClient *rest.RESTClient, pod *v1.Pod, container *v1.Container, command []string) error {
	//TODO the restClient to be used should be the k8s one and not weblogic one-------------------
	result :=
		restClient.Post().
			Namespace(pod.Namespace).
			Resource("pods").
			Name(pod.Name).
			SubResource("exec").
			Param("container", container.Name).
			Param("command", strings.Join(command, " ")).
			Do()

	if result.Error() != nil {
		glog.Infof("Result of executing command is not nil")
		//if metav1.Status(result.Error).Status != metav1.StatusSuccess {
		glog.Errorf("Error executing command: %s", result.Error())
		return result.Error()
		//}
	}

	return nil
}
