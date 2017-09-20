package services

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"weblogic-operator/pkg/constants"
	"weblogic-operator/pkg/types"
)

// NewForServer will return a new NodePort Kubernetes service for a Weblogic Server
func NewForServer(server *types.WeblogicServer) *v1.Service {
	weblogicPort := v1.ServicePort{Port: 7001}
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Labels:    map[string]string{constants.WeblogicServerLabel: server.Name},
			Name:      server.Name,
			Namespace: server.Namespace,
		},
		Spec: v1.ServiceSpec{
			Type:  v1.ServiceTypeNodePort,
			Ports: []v1.ServicePort{weblogicPort},
			Selector: map[string]string{
				constants.WeblogicServerLabel: server.Name,
			},
		},
	}
	return svc
}

func NewServiceForDomain(domain *types.WeblogicDomain) *v1.Service {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Labels:    map[string]string{constants.WeblogicDomainLabel: domain.Name},
			Name:      domain.Name,
			Namespace: domain.Namespace,
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{
				constants.WeblogicDomainLabel: domain.Name,
			},
		},
	}
	return svc
}
