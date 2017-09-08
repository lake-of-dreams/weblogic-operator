package services

import (
	"weblogic-operator/pkg/constants"
	"weblogic-operator/pkg/types"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewForServer will return a new headless Kubernetes service for a Weblogic server
func NewForServer(server *types.WeblogicServer) *v1.Service {
	mysqlPort := v1.ServicePort{Port: 3306}
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Labels:    map[string]string{constants.WeblogicServerLabel: server.Name},
			Name:      server.Name,
			Namespace: server.Namespace,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{mysqlPort},
			Selector: map[string]string{
				constants.WeblogicServerLabel: server.Name,
			},
		},
	}

	return svc
}
