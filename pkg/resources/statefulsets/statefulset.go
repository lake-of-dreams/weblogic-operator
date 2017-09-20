package statefulsets

import (
	"fmt"
	"k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"weblogic-operator/pkg/constants"
	"weblogic-operator/pkg/types"
)

func serverNameEnvVar(server *types.WeblogicServer) v1.EnvVar {
	return v1.EnvVar{Name: "WEBLOGIC_SERVER_NAME", Value: server.Name}
}

func namespaceEnvVar() v1.EnvVar {
	return v1.EnvVar{
		Name: "POD_NAMESPACE",
		ValueFrom: &v1.EnvVarSource{
			FieldRef: &v1.ObjectFieldSelector{
				FieldPath: "metadata.namespace",
			},
		},
	}
}

// Builds the WeblogicServer container
func weblogicOperatorContainer(server *types.WeblogicServer) v1.Container {
	return v1.Container{
		//TODO : Use different container names ???
		Name:            server.Name,
		Image:           fmt.Sprintf("%s:%s", constants.WeblogicImageName, server.Spec.Version),
		ImagePullPolicy: v1.PullAlways,
		Ports:           []v1.ContainerPort{{ContainerPort: 7001}},
		Env: []v1.EnvVar{
			serverNameEnvVar(server),
			namespaceEnvVar(),
		},
		//Lifecycle: &v1.Lifecycle{
		//	PreStop: &v1.Handler{
		//		Exec: &v1.ExecAction{
		//			Command: []string{"/u01/oracle/user_projects/domains/base_domain/bin/stopWebLogic.sh"},
		//		},
		//	},
		//},
	}
}

// NewForServer creates a new StatefulSet for the given WeblogicServer.
func NewForServer(server *types.WeblogicServer, serviceName string) *v1beta1.StatefulSet {
	var timeOut int64 = 120
	containers := []v1.Container{weblogicOperatorContainer(server)}

	ss := &v1beta1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: server.Namespace,
			Name:      server.Name,
			Labels: map[string]string{
				constants.WeblogicServerLabel: server.Name,
			},
		},
		Spec: v1beta1.StatefulSetSpec{
			Replicas: &server.Spec.Replicas,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						constants.WeblogicServerLabel: server.Name,
					},
				},
				Spec: v1.PodSpec{
					NodeSelector: server.Spec.NodeSelector,
					ImagePullSecrets: []v1.LocalObjectReference{
						{
							Name: "weblogic-docker-store",
						},
					},
					Containers:                    containers,
					TerminationGracePeriodSeconds: &timeOut,
				},
			},
			ServiceName: serviceName,
		},
	}

	return ss
}
