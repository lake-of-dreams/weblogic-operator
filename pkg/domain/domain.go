package domain

import (
	"fmt"

	"k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/golang/glog"

	"weblogic-operator/pkg/constants"
	"weblogic-operator/pkg/resources/services"
	"weblogic-operator/pkg/resources/statefulsets"
	"weblogic-operator/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
)

// HasServerNameLabel returns true if the given labels map matches the given
// server name.
func HasDomainNameLabel(labels map[string]string, domainname string) bool {
	for label, value := range labels {
		if label == constants.WebLogicDomainLabel {
			if value == domainname {
				return true
			}
		}
	}
	return false
}

// Return a label that uniquely identifies a Weblogic server
func getLabelSelectorForDomain(domain *types.WebLogicDomain) string {
	return fmt.Sprintf("%s=%s", constants.WebLogicDomainLabel, domain.Name)
}

// GetStatefulSetForWebLogicManagedServer finds the associated StatefulSet for a Weblogic server
func GetStatefulSetForWebLogicDomain(domain *types.WebLogicDomain, kubeClient kubernetes.Interface) (*v1beta1.StatefulSet, error) {
	opts := metav1.ListOptions{LabelSelector: getLabelSelectorForDomain(domain)}
	statefulsets, err := kubeClient.AppsV1beta1().StatefulSets(domain.Namespace).List(opts)
	if err != nil {
		glog.Errorf("Unable to list stateful sets for %s: %s", domain.Name, err)
		return nil, err
	}

	for _, ss := range statefulsets.Items {
		if HasDomainNameLabel(ss.Labels, domain.Name) {
			return &ss, nil
		}
	}
	return nil, nil
}

// CreateStatefulSetForWebLogicManagedServer will create a new Kubernetes StatefulSet based on a predefined template
func CreateStatefulSetForWebLogicDomain(clientset kubernetes.Interface, domain *types.WebLogicDomain, service *v1.Service) (*v1beta1.StatefulSet, error) {
	// Find StatefulSet and if it does not exist create it
	existingStatefulSet, err := GetStatefulSetForWebLogicDomain(domain, clientset)
	if err != nil {
		glog.Errorf("Error finding stateful set for domain: %v", err)
		return nil, err
	}

	if existingStatefulSet != nil {
		glog.V(2).Infof("Stateful set with label %s already exists", getLabelSelectorForDomain(domain))
		return existingStatefulSet, nil
	}

	glog.V(4).Infof("Creating a new stateful set for domain %s", domain.Name)
	ss := statefulsets.NewServiceForDomain(domain, service.Name)

	glog.V(4).Infof("Creating domain %+v", ss)
	return clientset.AppsV1beta1().StatefulSets(domain.Namespace).Create(ss)
}

// DeleteStatefulSetForWebLogicManagedServer will delete a stateful set by name
func DeleteStatefulSetForWebLogicDomain(clientset kubernetes.Interface, domain *types.WebLogicDomain) error {
	statefulSet, err := GetStatefulSetForWebLogicDomain(domain, clientset)
	if err != nil || statefulSet == nil {
		glog.Errorf("Could not delete stateful set: %s", err)
		return err
	}

	glog.V(4).Infof("Deleting stateful set %s", statefulSet.Name)
	var policy = metav1.DeletePropagationBackground
	return clientset.AppsV1beta1().
		StatefulSets(domain.Namespace).
		Delete(statefulSet.Name, &metav1.DeleteOptions{PropagationPolicy: &policy})
}

func createWebLogicDomain(domain *types.WebLogicDomain, kubeClient kubernetes.Interface, restClient *rest.RESTClient) error {
	domain.EnsureDefaults()

	//err := domain.Validate()
	//if err != nil {
	//	return err
	//}

	// Validate that a label is set on the server
	if !HasDomainNameLabel(domain.Labels, domain.Name) {
		glog.V(4).Infof("Setting label on domain %s", getLabelSelectorForDomain(domain))
		if domain.Labels == nil {
			domain.Labels = make(map[string]string)
		}
		domain.Labels[constants.WebLogicDomainLabel] = domain.Name
		return updateWebLogicDomain(domain, restClient)
	}

	domainService, err := CreateServiceForWebLogicDomain(kubeClient, domain)
	if err != nil {
		return err
	}

	_, err = CreateStatefulSetForWebLogicDomain(kubeClient, domain, domainService)
	if err != nil {
		return err
	}

	return nil
}

func updateWebLogicDomain(domain *types.WebLogicDomain, restClient *rest.RESTClient) error {
	result := restClient.Put().
		Resource(types.DomainCRDResourcePlural).
		Namespace(domain.Namespace).
		Name(domain.Name).
		Body(domain).
		Do()
	return result.Error()
}

// When delete server is called we will delete the stateful set (which also deletes the associated service)
//TODO handling to call stopWeblogic.sh needs to be done here
func deleteWebLogicDomain(domain *types.WebLogicDomain, kubeClient kubernetes.Interface, restClient *rest.RESTClient) error {
	//err := domain.Validate()
	//if err != nil {
	//	return err
	//}

	//err = RunStopForWebLogicManagedServer(kubeClient, restClient, server)
	//if err != nil {
	//	return err
	//}

	//err = DeleteStatefulSetForWebLogicDomain(kubeClient, domain)
	//if err != nil {
	//	return err
	//}
	//
	//err = DeleteServiceForWebLogicDomain(kubeClient, domain)
	//if err != nil {
	//	return err
	//}

	return nil
}

// GetServiceForWebLogicManagedServer returns the associated service for a given server
func GetServiceForWebLogicDomain(domain *types.WebLogicDomain, clientset kubernetes.Interface) (*v1.Service, error) {
	opts := metav1.ListOptions{LabelSelector: getLabelSelectorForDomain(domain)}
	services, err := clientset.CoreV1().Services(domain.Namespace).List(opts)
	if err != nil {
		glog.Errorf("Unable to list services for %s: %s", domain.Name, err)
		return nil, err
	}

	for _, svc := range services.Items {
		if HasDomainNameLabel(svc.Labels, domain.Name) {
			return &svc, nil
		}
	}
	return nil, nil
}

// CreateServiceForWebLogicManagedServer will create a new Kubernetes Service based on a predefined template
func CreateServiceForWebLogicDomain(clientset kubernetes.Interface, domain *types.WebLogicDomain) (*v1.Service, error) {
	// Find Service and if it does not exist create it
	existingService, err := GetServiceForWebLogicDomain(domain, clientset)
	if err != nil {
		glog.Errorf("Error finding service for domain: %s", err)
		return nil, err
	}

	if existingService != nil {
		glog.V(2).Infof("Service with label %s already exists", getLabelSelectorForDomain(domain))
		return existingService, nil
	}

	glog.V(4).Infof("Creating a new service for domain %s", domain.Name)

	svc := services.NewServiceForDomain(domain)
	return clientset.CoreV1().Services(domain.Namespace).Create(svc)
}

// DeleteServiceForWebLogicManagedServer deletes the Service associated with a Weblogic server.
func DeleteServiceForWebLogicDomain(clientset kubernetes.Interface, domain *types.WebLogicDomain) error {
	service, err := GetServiceForWebLogicDomain(domain, clientset)
	if err != nil || service == nil {
		glog.Errorf("Could not delete service: %s", err)
		return err
	}
	glog.V(4).Infof("Deleting service %s", service.Name)
	return clientset.CoreV1().Services(domain.Namespace).Delete(service.Name, nil)
}

func GetDomainForStatefulSet(statefulSet *v1beta1.StatefulSet, restClient *rest.RESTClient) (domain *types.WebLogicDomain, err error) {
	if weblogicDomainName, ok := statefulSet.Labels[constants.WebLogicDomainLabel]; ok {
		domain = &types.WebLogicDomain{}
		result := restClient.Get().
			Resource(types.DomainCRDResourcePlural).
			Namespace(statefulSet.Namespace).
			Name(weblogicDomainName).
			Do().
			Into(domain)
		return domain, result
	}
	return nil, fmt.Errorf("unable to get Label %s from statefulset. Not part of domain", constants.WebLogicDomainLabel)
}

func setWebLogicDomainState(domain *types.WebLogicDomain, restClient *rest.RESTClient, phase types.WebLogicDomainPhase, err error) error {
	modified := false
	if domain.Status.Phase != phase {
		domain.Status.Phase = phase
		modified = true
	}

	l := len(domain.Status.Errors)
	if err != nil && (l < 1 || domain.Status.Errors[l-1] != err.Error()) {
		domain.Status.Errors = append(domain.Status.Errors, err.Error())
		modified = true
	} else if l == 0 {
		domain.Status.Errors = []string{}
		modified = true
	}

	if modified {
		result := restClient.Put().
			Resource(types.DomainCRDResourcePlural).
			Namespace(domain.Namespace).
			Name(domain.Name).
			Body(domain).
			Do()

		return result.Error()
	}

	return nil
}

func updateDomainWithStatefulSet(domain *types.WebLogicDomain, statefulSet *v1beta1.StatefulSet, kubeClient kubernetes.Interface, restClient *rest.RESTClient) (err error) {
	// Some simple logic for the time being.
	// To add
	// connection to the server
	// validate each pod?
	// Check how a rolling upgrade effects this
	// check version of each pod

	if statefulSet.Status.ReadyReplicas < statefulSet.Status.Replicas {
		setWebLogicDomainState(domain, restClient, types.WebLogicDomainPending, nil)
	} else if statefulSet.Status.ReadyReplicas == statefulSet.Status.Replicas {
		setWebLogicDomainState(domain, restClient, types.WebLogicDomainPending, nil)
	}
	return err
}

// GetPodForWebLogicManagedServer finds the associated pod for a Weblogic server
func GetPodForWebLogicDomain(domain *types.WebLogicDomain, clientset kubernetes.Interface) (*v1.Pod, error) {
	opts := metav1.ListOptions{LabelSelector: getLabelSelectorForDomain(domain)}
	pods, err := clientset.CoreV1().Pods(domain.Namespace).List(opts)
	if err != nil {
		glog.Errorf("Unable to list pods for %s: %s", domain.Name, err)
		return nil, err
	}

	for _, pod := range pods.Items {
		if HasDomainNameLabel(pod.Labels, domain.Name) {
			return &pod, nil
		}
	}
	return nil, nil
}

// GetPodForWebLogicManagedServer finds the associated pod for a Weblogic server
func GetContainerForPod(domain *types.WebLogicDomain, pod *v1.Pod) (*v1.Container, error) {
	containers := pod.Spec.Containers

	for _, container := range containers {
		if container.Name == domain.Name {
			return &container, nil
		}
	}

	return nil, nil
}

// RunStopForWebLogicManagedServer will run stopWebLogic to stop a Weblogic server container in a pod
func CreateWebLogicDomain(clientset kubernetes.Interface, restClient *rest.RESTClient, domain *types.WebLogicDomain) error {
	pod, err := GetPodForWebLogicDomain(domain, clientset)
	if err != nil || pod == nil {
		glog.Errorf("Could not find pod: %s", err)
		return err
	}

	container, err := GetContainerForPod(domain, pod)
	if err != nil || container == nil {
		glog.Errorf("Could not find container %s in pod %s: %s", domain.Name, pod.Name, err)
		return err
	}

	glog.V(4).Infof("Running domainSetup.sh for container %s in pod %s", domain.Name, pod.Name)
	command := []string{"domainSetup.sh /u01/oracle/user_projects/domains/base_domain weblogic welcome1 localhost 7001 2 5556"}
	cmdErr := ExecuteCommandInContainer(restClient, pod, container, command)
	if cmdErr != nil {
		glog.Errorf("Error executing command : %s", cmdErr)
		return cmdErr
	}

	return nil
}

// ExecuteCommandInContainer will run a command in a container in a pod
func ExecuteCommandInContainer(restClient *rest.RESTClient, pod *v1.Pod, container *v1.Container, command []string) error {
	result := restClient.Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec").
		Param("container", container.Name).
		VersionedParams(&v1.PodExecOptions{
		Container: container.Name,
		Command:   command,
	}, scheme.ParameterCodec).
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
