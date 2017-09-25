package replicasets

import (
	"fmt"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"weblogic-operator/pkg/constants"
	"k8s.io/api/extensions/v1beta1"
	"weblogic-operator/pkg/types"
)

// WeblogicImageName is the base Docker image used by the operator.
const WeblogicImageName = "store/oracle/weblogic"

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
		Name:  "weblogic",
		Image: fmt.Sprintf("%s:%s", WeblogicImageName, server.Spec.Version),
		Ports: []v1.ContainerPort{{ContainerPort: 7001}},
		Env: []v1.EnvVar{
			serverNameEnvVar(server),
			namespaceEnvVar(),
		},
		Resources: server.Spec.Resources,
		Lifecycle: &v1.Lifecycle{
			PreStop: &v1.Handler{
				Exec: &v1.ExecAction{
					Command: []string{"/u01/oracle/user_projects/domains/base_domain/bin/stopWebLogic.sh"},
				},
			},
		},
	}
}

// NewForServer creates a new ReplicationController for the given WeblogicServer.
func NewForServer(server *types.WeblogicServer, serviceName string) *v1beta1.ReplicaSet {
	var timeOut int64 = 120
	containers := []v1.Container{weblogicOperatorContainer(server)}

	rs := &v1beta1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: server.Namespace,
			Name:      server.Name,
			Labels: map[string]string{
				constants.WeblogicServerLabel: server.Name,
			},
		},
		Spec: v1beta1.ReplicaSetSpec{
			Replicas: &server.Spec.Replicas,
			MinReadySeconds: 0,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					constants.WeblogicServerLabel: server.Name,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						constants.WeblogicServerLabel: server.Name,
					},
				},
				Spec: v1.PodSpec{
					NodeSelector:                  server.Spec.NodeSelector,
					Containers:                    containers,
					TerminationGracePeriodSeconds: &timeOut,
				},
			},
		},
	}

	return rs
}