package server

import (
	"fmt"

	"k8s.io/api/extensions/v1beta1"
	"k8s.io/api/core/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/golang/glog"

	"weblogic-operator/pkg/constants"
	"weblogic-operator/pkg/resources/services"
	"weblogic-operator/pkg/types"
	"weblogic-operator/pkg/resources/replicasets"
	"weblogic-operator/pkg/resources/horizontalpodautoscalers"
	"strings"
)

// HasServerNameLabel returns true if the given labels map matches the given
// server name.
func HasServerNameLabel(labels map[string]string, servername string) bool {
	for label, value := range labels {
		if label == constants.WeblogicServerLabel {
			if value == servername {
				return true
			}
		}
	}
	return false
}

// Return a label that uniquely identifies a Weblogic server
func getLabelSelectorForServer(server *types.WeblogicServer) string {
	return fmt.Sprintf("%s=%s", constants.WeblogicServerLabel, server.Name)
}

//// GetStatefulSetForWeblogicServer finds the associated StatefulSet for a Weblogic server
//func GetStatefulSetForWeblogicServer(server *types.WeblogicServer, kubeClient kubernetes.Interface) (*v1beta1.StatefulSet, error) {
//	opts := metav1.ListOptions{LabelSelector: getLabelSelectorForServer(server)}
//	statefulsets, err := kubeClient.AppsV1beta1().StatefulSets(server.Namespace).List(opts)
//	if err != nil {
//		glog.Errorf("Unable to list stateful sets for %s: %s", server.Name, err)
//		return nil, err
//	}
//
//	for _, ss := range statefulsets.Items {
//		if HasServerNameLabel(ss.Labels, server.Name) {
//			return &ss, nil
//		}
//	}
//	return nil, nil
//}
//
//// CreateStatefulSetForWeblogicServer will create a new Kubernetes StatefulSet based on a predefined template
//func CreateStatefulSetForWeblogicServer(clientset kubernetes.Interface, server *types.WeblogicServer, service *v1.Service) (*v1beta1.StatefulSet, error) {
//	// Find StatefulSet and if it does not exist create it
//	existingStatefulSet, err := GetStatefulSetForWeblogicServer(server, clientset)
//	if err != nil {
//		glog.Errorf("Error finding stateful set for server: %v", err)
//		return nil, err
//	}
//
//	if existingStatefulSet != nil {
//		glog.V(2).Infof("Stateful set with label %s already exists", getLabelSelectorForServer(server))
//		return existingStatefulSet, nil
//	}
//
//	glog.V(4).Infof("Creating a new stateful set for server %s", server.Name)
//	ss := statefulsets.NewForServer(server, service.Name)
//
//	glog.V(4).Infof("Creating server %+v", ss)
//	return clientset.AppsV1beta1().StatefulSets(server.Namespace).Create(ss)
//}
//
//// DeleteStatefulSetForWeblogicServer will delete a stateful set by name
//func DeleteStatefulSetForWeblogicServer(clientset kubernetes.Interface, server *types.WeblogicServer) error {
//	statefulSet, err := GetStatefulSetForWeblogicServer(server, clientset)
//	if err != nil || statefulSet == nil {
//		glog.Errorf("Could not delete stateful set: %s", err)
//		return err
//	}
//
//	glog.V(4).Infof("Deleting stateful set %s", statefulSet.Name)
//	var policy = metav1.DeletePropagationBackground
//	return clientset.AppsV1beta1().
//		StatefulSets(server.Namespace).
//		Delete(statefulSet.Name, &metav1.DeleteOptions{PropagationPolicy: &policy})
//}

// GetReplicaSetForWeblogicServer finds the associated ReplicaSet for a Weblogic server
func GetReplicaSetForWeblogicServer(server *types.WeblogicServer, kubeClient kubernetes.Interface) (*v1beta1.ReplicaSet, error) {
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

// CreateReplicaSetForWeblogicServer will create a new Kubernetes ReplicaSet based on a predefined template
func CreateReplicaSetForWeblogicServer(clientset kubernetes.Interface, server *types.WeblogicServer, service *v1.Service) (controller *v1beta1.ReplicaSet, err error) {
	// Find ReplicaSet and if it does not exist create it
	existingReplicaSet, err := GetReplicaSetForWeblogicServer(server, clientset)
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

// DeleteReplicaSetForWeblogicServer will delete a replica set by name
func DeleteReplicaSetForWeblogicServer(clientset kubernetes.Interface, server *types.WeblogicServer) error {

	if strings.EqualFold(server.Name, constants.HorizontalPodAutoscalerTargetLabel) {
		glog.V(4).Infof("Deleting HPA for managed server!!!!!!!!!!!!!")
		err := DeleteHorizontalPodAutoscalerForWeblogicServer(clientset, server)
		if err != nil {
			return err
		}
	}

	replicaSet, err := GetReplicaSetForWeblogicServer(server, clientset)
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

// GetHorizontalPodAutoscalerForWeblogicServer finds the associated ReplicaSet for a Weblogic server
func GetHorizontalPodAutoscalerForWeblogicServer(server *types.WeblogicServer, kubeClient kubernetes.Interface) (*autoscalingv1.HorizontalPodAutoscaler, error) {
	opts := metav1.ListOptions{LabelSelector: getLabelSelectorForServer(server)}
	horizontalpodautoscalers, err := kubeClient.AutoscalingV1().HorizontalPodAutoscalers(server.Namespace).List(opts)
	glog.V(4).Infof("5555555555555555555")
	if err != nil {
		glog.V(4).Infof("22222222222222222222")
		glog.Errorf("Unable to list horizontal pod autoscalers for %s: %s", server.Name, err)
		return nil, err
	}

	for _, rc := range horizontalpodautoscalers.Items {
		glog.V(4).Infof("1111111111111111111111111111111111111")
		if HasServerNameLabel(rc.Labels, server.Name) {
			glog.V(4).Infof("Value of RC!!!")
			return &rc, nil
		}
	}
	return nil, nil
}

// CreateHorizontalPodAutoscalerForWeblogicServer will create a new Kubernetes HorizontalPodAutoscaler based on a predefined template
func CreateHorizontalPodAutoscalerForWeblogicServer(clientset kubernetes.Interface, server *types.WeblogicServer, service *v1.Service) (controller *autoscalingv1.HorizontalPodAutoscaler, err error) {
	// Find ReplicaSet and if it does not exist create it
	existingHorizontalPodAutoscaler, err := GetHorizontalPodAutoscalerForWeblogicServer(server, clientset)
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

// DeleteHorizontalPodAutoscalerForWeblogicServer will delete a replica set by name
func DeleteHorizontalPodAutoscalerForWeblogicServer(clientset kubernetes.Interface, server *types.WeblogicServer) error {
	horizontalPodAutoscaler, err := GetHorizontalPodAutoscalerForWeblogicServer(server, clientset)
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

func createWeblogicServer(server *types.WeblogicServer, kubeClient kubernetes.Interface, restClient *rest.RESTClient) error {
	server.EnsureDefaults()

	err := server.Validate()
	if err != nil {
		return err
	}

	// Validate that a label is set on the server
	if !HasServerNameLabel(server.Labels, server.Name) {
		glog.V(4).Infof("Setting label on server %s", getLabelSelectorForServer(server))
		if server.Labels == nil {
			server.Labels = make(map[string]string)
		}
		server.Labels[constants.WeblogicServerLabel] = server.Name
		return updateWeblogicServer(server, restClient)
	}

	serverService, err := CreateServiceForWeblogicServer(kubeClient, server)
	if err != nil {
		return err
	}

	//_, err = CreateStatefulSetForWeblogicServer(kubeClient, server, serverService)
	//if err != nil {
	//	return err
	//}

	_, err = CreateReplicaSetForWeblogicServer(kubeClient, server, serverService)
	if err != nil {
		return err
	}

	_, err = CreateHorizontalPodAutoscalerForWeblogicServer(kubeClient, server, serverService)
	if err != nil {
		return err
	}

	return nil
}

func updateWeblogicServer(server *types.WeblogicServer, restClient *rest.RESTClient) error {
	result := restClient.Put().
		Resource(constants.WeblogicServerResourceKindPlural).
		Namespace(server.Namespace).
		Name(server.Name).
		Body(server).
		Do()
	return result.Error()
}

// When delete server is called we will delete the stateful set (which also deletes the associated service)
//TODO handling to call stopWeblogic.sh needs to be done here
func deleteWeblogicServer(server *types.WeblogicServer, kubeClient kubernetes.Interface, restClient *rest.RESTClient) error {
	err := server.Validate()
	if err != nil {
		return err
	}

	err = DeleteReplicaSetForWeblogicServer(kubeClient, server)
	if err != nil {
		return err
	}

	err = RunStopForWeblogicServer(kubeClient, restClient, server)
	if err != nil {
		return err
	}

	//err = DeleteStatefulSetForWeblogicServer(kubeClient, server)
	//if err != nil {
	//	return err
	//}
	//
	//err = DeleteServiceForWeblogicServer(kubeClient, server)
	//if err != nil {
	//	return err
	//}

	return nil
}

// GetServiceForWeblogicServer returns the associated service for a given server
func GetServiceForWeblogicServer(server *types.WeblogicServer, clientset kubernetes.Interface) (*v1.Service, error) {
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

// CreateServiceForWeblogicServer will create a new Kubernetes Service based on a predefined template
func CreateServiceForWeblogicServer(clientset kubernetes.Interface, server *types.WeblogicServer) (*v1.Service, error) {
	// Find Service and if it does not exist create it
	existingService, err := GetServiceForWeblogicServer(server, clientset)
	if err != nil {
		glog.Errorf("Error finding service for server: %s", err)
		return nil, err
	}

	if existingService != nil {
		glog.V(2).Infof("Service with label %s already exists", getLabelSelectorForServer(server))
		return existingService, nil
	}

	glog.V(4).Infof("Creating a new service for server %s", server.Name)

	svc := services.NewForServer(server)
	return clientset.CoreV1().Services(server.Namespace).Create(svc)
}

// DeleteServiceForWeblogicServer deletes the Service associated with a Weblogic server.
func DeleteServiceForWeblogicServer(clientset kubernetes.Interface, server *types.WeblogicServer) error {
	service, err := GetServiceForWeblogicServer(server, clientset)
	if err != nil || service == nil {
		glog.Errorf("Could not delete service: %s", err)
		return err
	}
	glog.V(4).Infof("Deleting service %s", service.Name)
	return clientset.CoreV1().Services(server.Namespace).Delete(service.Name, nil)
}

//func GetServerForStatefulSet(statefulSet *v1beta1.StatefulSet, restClient *rest.RESTClient) (server *types.WeblogicServer, err error) {
//	if weblogicServerName, ok := statefulSet.Labels[constants.WeblogicServerLabel]; ok {
//		server = &types.WeblogicServer{}
//		result := restClient.Get().
//			Resource(constants.WeblogicServerResourceKindPlural).
//			Namespace(statefulSet.Namespace).
//			Name(weblogicServerName).
//			Do().
//			Into(server)
//		return server, result
//	}
//	return nil, fmt.Errorf("unable to get Label %s from statefulset. Not part of server", constants.WeblogicServerLabel)
//}

func setWeblogicServerState(server *types.WeblogicServer, restClient *rest.RESTClient, phase types.WeblogicServerPhase, err error) error {
	modified := false
	if server.Status.Phase != phase {
		server.Status.Phase = phase
		modified = true
	}

	l := len(server.Status.Errors)
	if err != nil && (l < 1 || server.Status.Errors[l-1] != err.Error()) {
		server.Status.Errors = append(server.Status.Errors, err.Error())
		modified = true
	} else if l == 0 {
		server.Status.Errors = []string{}
		modified = true
	}

	if modified {
		result := restClient.Put().
			Resource(constants.WeblogicServerResourceKindPlural).
			Namespace(server.Namespace).
			Name(server.Name).
			Body(server).
			Do()

		return result.Error()
	}

	return nil
}

//func updateServerWithStatefulSet(server *types.WeblogicServer, statefulSet *v1beta1.StatefulSet, kubeClient kubernetes.Interface, restClient *rest.RESTClient) (err error) {
//	// Some simple logic for the time being.
//	// To add
//	// connection to the server
//	// validate each pod?
//	// Check how a rolling upgrade effects this
//	// check version of each pod
//
//	if statefulSet.Status.ReadyReplicas < statefulSet.Status.Replicas {
//		setWeblogicServerState(server, restClient, types.WeblogicServerPending, nil)
//	} else if statefulSet.Status.ReadyReplicas == statefulSet.Status.Replicas {
//		setWeblogicServerState(server, restClient, types.WeblogicServerRunning, nil)
//	}
//	return err
//}

func GetServerForReplicaSet(replicaset *v1beta1.ReplicaSet, restClient *rest.RESTClient) (server *types.WeblogicServer, err error) {
	if weblogicServerName, ok := replicaset.Labels[constants.WeblogicServerLabel]; ok {
		server = &types.WeblogicServer{}
		result := restClient.Get().
			Resource(constants.WeblogicServerResourceKindPlural).
			Namespace(replicaset.Namespace).
			Name(weblogicServerName).
			Do().
			Into(server)
		return server, result
	}
	return nil, fmt.Errorf("unable to get Label %s from replicaset. Not part of server", constants.WeblogicServerLabel)
}

func updateServerWithReplicaSet(server *types.WeblogicServer, replicaSet *v1beta1.ReplicaSet, kubeClient kubernetes.Interface, restClient *rest.RESTClient) (err error) {
	// Some simple logic for the time being.
	// To add
	// connection to the server
	// validate each pod?
	// Check how a rolling upgrade effects this
	// check version of each pod

	if replicaSet.Status.ReadyReplicas < replicaSet.Status.Replicas {
		setWeblogicServerState(server, restClient, types.WeblogicServerPending, nil)
	} else if replicaSet.Status.ReadyReplicas == replicaSet.Status.Replicas {
		setWeblogicServerState(server, restClient, types.WeblogicServerRunning, nil)
	}
	return err
}

func GetServerForHorizontalPodAutoscaler(horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler, restClient *rest.RESTClient) (server *types.WeblogicServer, err error) {
	if weblogicServerName, ok := horizontalPodAutoscaler.Labels[constants.HorizontalPodAutoscalerTargetLabel]; ok {
		server = &types.WeblogicServer{}
		result := restClient.Get().
			Resource(constants.WeblogicServerResourceKindPlural).
			Namespace(horizontalPodAutoscaler.Namespace).
			Name(weblogicServerName).
			Do().
			Into(server)
		return server, result
	}
	return nil, fmt.Errorf("unable to get Label %s from horizontalPodAutoscaler. Not part of server", constants.HorizontalPodAutoscalerTargetLabel)
}

func updateServerWithHorizontalPodAutoscaler(server *types.WeblogicServer, horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler, kubeClient kubernetes.Interface, restClient *rest.RESTClient) (err error) {
	// Some simple logic for the time being.
	// To add
	// connection to the server
	// validate each pod?
	// Check how a rolling upgrade effects this
	// check version of each pod

	if horizontalPodAutoscaler.Status.CurrentReplicas < horizontalPodAutoscaler.Status.DesiredReplicas {
		setWeblogicServerState(server, restClient, types.WeblogicServerPending, nil)
	} else if horizontalPodAutoscaler.Status.CurrentReplicas == horizontalPodAutoscaler.Status.DesiredReplicas {
		setWeblogicServerState(server, restClient, types.WeblogicServerRunning, nil)
	}
	return err
}

// GetPodForWeblogicServer finds the associated pod for a Weblogic server
func GetPodForWeblogicServer(server *types.WeblogicServer, clientset kubernetes.Interface) (*v1.Pod, error) {
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

// GetPodForWeblogicServer finds the associated pod for a Weblogic server
func GetContainerForPod(server *types.WeblogicServer, pod *v1.Pod) (*v1.Container, error) {
	containers := pod.Spec.Containers

	for _, container := range containers {
		if container.Name == server.Name {
			return &container, nil
		}
	}

	return nil, nil
}

// RunStopForWeblogicServer will run stopWebLogic to stop a Weblogic server container in a pod
func RunStopForWeblogicServer(clientset kubernetes.Interface, restClient *rest.RESTClient, server *types.WeblogicServer) error {
	pod, err := GetPodForWeblogicServer(server, clientset)
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

// ExecuteCommandInContainer will run a command in a container in a pod
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
