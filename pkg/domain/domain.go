package domain

import (
	"fmt"

	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/golang/glog"

	"weblogic-operator/pkg/constants"
	"weblogic-operator/pkg/resources/replicasets"
	"weblogic-operator/pkg/resources/services"
	"weblogic-operator/pkg/types"
	"io/ioutil"
	"encoding/json"
	"reflect"
)

// HasDomainNameLabel returns true if the given labels map matches the given
// domain name.
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

// Return a label that uniquely identifies a Weblogic domain
func getLabelSelectorForDomain(domain *types.WebLogicDomain) string {
	return fmt.Sprintf("%s=%s", constants.WebLogicDomainLabel, domain.Name)
}

// GetReplicaSetForWebLogicDomain finds the associated ReplicaSet for a Weblogic domain
func GetReplicaSetForWebLogicDomain(domain *types.WebLogicDomain, kubeClient kubernetes.Interface) (*v1beta1.ReplicaSet, error) {
	opts := metav1.ListOptions{LabelSelector: getLabelSelectorForDomain(domain)}
	replicaSets, err := kubeClient.ExtensionsV1beta1().ReplicaSets(domain.Namespace).List(opts)
	if err != nil {
		glog.Errorf("Unable to list replica sets for %s: %s", domain.Name, err)
		return nil, err
	}

	for _, rc := range replicaSets.Items {
		if HasDomainNameLabel(rc.Labels, domain.Name) {
			return &rc, nil
		}
	}
	return nil, nil
}

// CreateReplicaSetForWebLogicDomain will create a new Kubernetes ReplicaSet based on a predefined template
func CreateReplicaSetForWebLogicDomain(clientset kubernetes.Interface, domain *types.WebLogicDomain, service *v1.Service) (controller *v1beta1.ReplicaSet, err error) {
	// Find ReplicaSet and if it does not exist create it
	existingReplicaSet, err := GetReplicaSetForWebLogicDomain(domain, clientset)
	if err != nil {
		glog.Errorf("Error finding replica set for domain: %v", err)
		return nil, err
	}

	if existingReplicaSet != nil {
		glog.V(2).Infof("Replica set with label %s already exists", getLabelSelectorForDomain(domain))
		return existingReplicaSet, nil
	}

	glog.V(4).Infof("Creating a new replica set for domain %s", domain.Name)
	rs := replicasets.NewForDomain(domain, service.Name)

	glog.V(4).Infof("Creating domain %+v", rs)
	return clientset.ExtensionsV1beta1().ReplicaSets(domain.Namespace).Create(rs)
}

// DeleteReplicaSetForWebLogicDomain will delete a replica set by name
func DeleteReplicaSetForWebLogicDomain(clientset kubernetes.Interface, domain *types.WebLogicDomain) error {
	replicaSet, err := GetReplicaSetForWebLogicDomain(domain, clientset)
	if err != nil || replicaSet == nil {
		glog.Errorf("Could not delete replica set: %s", err)
		return err
	}

	glog.V(4).Infof("Deleting replica set %s", replicaSet.Name)
	var policy = metav1.DeletePropagationBackground
	return clientset.ExtensionsV1beta1().
		ReplicaSets(domain.Namespace).
		Delete(replicaSet.Name, &metav1.DeleteOptions{PropagationPolicy: &policy})
}

func createWebLogicDomain(domain *types.WebLogicDomain, kubeClient kubernetes.Interface, restClient *rest.RESTClient) error {
	domain.EnsureDefaults()

	// Validate that a label is set on the domain
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

	_, err = CreateReplicaSetForWebLogicDomain(kubeClient, domain, domainService)
	if err != nil {
		return err
	}

	return nil
}

func updateWebLogicDomain(domain *types.WebLogicDomain, restClient *rest.RESTClient) error {
	result := restClient.Put().
		Resource(constants.WebLogicDomainResourceKindPlural).
		Namespace(domain.Namespace).
		Name(domain.Name).
		Body(domain).
		Do()
	return result.Error()
}

// When delete domain is called we will delete the replica set (which also deletes the associated service)
//TODO handling to call stopWeblogic.sh needs to be done here
func deleteWebLogicDomain(domain *types.WebLogicDomain, kubeClient kubernetes.Interface, restClient *rest.RESTClient) error {
	err := DeleteReplicaSetForWebLogicDomain(kubeClient, domain)
	if err != nil {
		return err
	}

	err = DeleteServiceForWebLogicDomain(kubeClient, domain)
	if err != nil {
		return err
	}

	return nil
}

// GetServiceForWebLogicDomain returns the associated service for a given domain
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

// CreateServiceForWebLogicDomain will create a new Kubernetes Service based on a predefined template
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

// DeleteServiceForWebLogicDomain deletes the Service associated with a Weblogic domain.
func DeleteServiceForWebLogicDomain(clientset kubernetes.Interface, domain *types.WebLogicDomain) error {
	service, err := GetServiceForWebLogicDomain(domain, clientset)
	if err != nil || service == nil {
		glog.Errorf("Could not delete service: %s", err)
		return err
	}
	glog.V(4).Infof("Deleting service %s", service.Name)
	return clientset.CoreV1().Services(domain.Namespace).Delete(service.Name, nil)
}

func GetDomainForReplicaSet(replicaset *v1beta1.ReplicaSet, restClient *rest.RESTClient) (domain *types.WebLogicDomain, err error) {
	if weblogicDomainName, ok := replicaset.Labels[constants.WebLogicDomainLabel]; ok {
		domain = &types.WebLogicDomain{}
		result := restClient.Get().
			Resource(constants.WebLogicDomainResourceKindPlural).
			Namespace(replicaset.Namespace).
			Name(weblogicDomainName).
			Do().
			Into(domain)
		return domain, result
	}
	return nil, fmt.Errorf("unable to get Label %s from replicaset. Not part of domain", constants.WebLogicDomainLabel)
}

func updateDomainWithReplicaSet(domain *types.WebLogicDomain, replicaSet *v1beta1.ReplicaSet, kubeClient kubernetes.Interface, restClient *rest.RESTClient) (err error) {
	err = PopulateServerDetailsForWebLogicDomain(domain, restClient)
	if err != nil {
		return err
	}
	return nil
}

func PopulateServerDetailsForWebLogicDomain(domain *types.WebLogicDomain, restClient *rest.RESTClient) error {
	serverListFile := "/u01/oracle/user_projects/domains/" + domain.Name + "/serverList.json"

	file, err := ioutil.ReadFile(serverListFile)
	if err != nil {
		glog.V(4).Infof(err.Error())
	}

	err = json.Unmarshal(file, &domain.Spec.ServersAvailable)
	if err != nil {
		glog.V(4).Infof(err.Error())
	}

	//err = json.Unmarshal(file, &domain.Spec.ServersRunning)
	//if err != nil {
	//	glog.V(4).Infof(err.Error())
	//}

	err = updateWebLogicDomain(domain, restClient)
	if err != nil {
		return err
	}

	return nil
}
