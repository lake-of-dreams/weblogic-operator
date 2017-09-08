package services

import (
	"weblogic-operator/pkg/constants"
	"weblogic-operator/pkg/types"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewForServer will return a new headless Kubernetes service for a Weblogic Server
func NewForServer(server *types.WeblogicServer) *v1.Service {
	weblogicPort := v1.ServicePort{Port: 7001}
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Labels:    map[string]string{constants.WeblogicServerLabel: server.Name},
			Name:      server.Name,
			Namespace: server.Namespace,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{weblogicPort},
			Selector: map[string]string{
				constants.WeblogicServerLabel: server.Name,
			},
		},
	}
	return svc
}
