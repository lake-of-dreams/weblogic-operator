package services

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"weblogic-operator/pkg/constants"
	"weblogic-operator/pkg/types"
	"github.com/golang/glog"
	"fmt"
)

// NewServiceForServer will return a new NodePort Kubernetes service for a WeblogicManagedServer
func NewServiceForServer(server *types.WebLogicManagedServer) *v1.Service {
	var startPort int32 = 7001
	//var weblogicPorts []v1.ServicePort
	weblogicPorts := make([]v1.ServicePort, server.Spec.Domain.Spec.ManagedServerCount)

	for i := 1; i <= server.Spec.Domain.Spec.ManagedServerCount; i++ {
		var port = startPort + int32(i*2)
		glog.V(4).Info("Calculated port ", fmt.Sprint(port))
		weblogicPorts[i-1] = v1.ServicePort{
			Name: fmt.Sprint("managedserver-", i-1, "port"),
			Port: port,
		}
	}
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				constants.WebLogicManagedServerLabel: server.Name,
				server.Spec.DomainName:               "managedserver",
			},
			Name:      server.Name,
			Namespace: server.Namespace,
		},
		Spec: v1.ServiceSpec{
			Type:  v1.ServiceTypeNodePort,
			Ports: weblogicPorts,
			Selector: map[string]string{
				constants.WebLogicManagedServerLabel: server.Name,
				server.Spec.DomainName:               "managedserver",
			},
		},
	}
	return svc
}

func NewServiceForDomain(domain *types.WebLogicDomain) *v1.Service {
	weblogicPort := v1.ServicePort{Port: 7001}
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				constants.WebLogicDomainLabel: domain.Name,
			},
			Name:      domain.Name,
			Namespace: domain.Namespace,
		},
		Spec: v1.ServiceSpec{
			Type:  v1.ServiceTypeNodePort,
			Ports: []v1.ServicePort{weblogicPort},
			Selector: map[string]string{
				constants.WebLogicDomainLabel: domain.Name,
			},
		},
	}
	return svc
}

func NewHeadlessServiceForDomain(domain *types.WebLogicDomain) *v1.Service {
	weblogicPort := v1.ServicePort{
		Name:     domain.Name,
		Port:     7001,
		Protocol: v1.ProtocolTCP,
	}
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				constants.WebLogicDomainLabel: domain.Name,
			},
			Name:      domain.Name,
			Namespace: domain.Namespace,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				weblogicPort},
			Selector: map[string]string{
				constants.WebLogicDomainLabel: domain.Name,
			},
			ClusterIP: v1.ClusterIPNone,
		},
	}
	return svc
}
