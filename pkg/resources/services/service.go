package services

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"weblogic-operator/pkg/constants"
	"weblogic-operator/pkg/types"
)

// NewServiceForServer will return a new NodePort Kubernetes service for a WeblogicManagedServer
func NewServiceForServer(server *types.WebLogicManagedServer) *v1.Service {
	weblogicPort := v1.ServicePort{Port: 7001}
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{constants.WebLogicManagedServerLabel: server.Name,
				constants.WebLogicDomainLabel: server.Spec.Domain.Name},
			Name:      server.Name,
			Namespace: server.Namespace,
		},
		Spec: v1.ServiceSpec{
			Type:  v1.ServiceTypeNodePort,
			Ports: []v1.ServicePort{weblogicPort},
			Selector: map[string]string{
				constants.WebLogicManagedServerLabel: server.Name,
				constants.WebLogicDomainLabel:        server.Spec.Domain.Name,
			},
		},
	}
	return svc
}

func NewServiceForDomain(domain *types.WebLogicDomain) *v1.Service {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Labels:    map[string]string{constants.WebLogicDomainLabel: domain.Name},
			Name:      domain.Name,
			Namespace: domain.Namespace,
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{
				constants.WebLogicDomainLabel: domain.Name,
			},
		},
	}
	return svc
}
