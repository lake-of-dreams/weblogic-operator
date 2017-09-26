package replicasets

import (
	"fmt"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"weblogic-operator/pkg/constants"
	"k8s.io/api/extensions/v1beta1"
	"weblogic-operator/pkg/types"
)

func serverNameEnvVar(server *types.WebLogicManagedServer) v1.EnvVar {
	return v1.EnvVar{Name: "SERVER_NAME", Value: server.Name}
}

func serverNamespaceEnvVar() v1.EnvVar {
	return v1.EnvVar{
		Name: "POD_NAMESPACE",
		ValueFrom: &v1.EnvVarSource{
			FieldRef: &v1.ObjectFieldSelector{
				FieldPath: "metadata.namespace",
			},
		},
	}
}

// Builds the WebLogicManagedServer container
func WebLogicManagedServerContainer(server *types.WebLogicManagedServer) v1.Container {
	return v1.Container{
		Name:            server.Spec.Domain.Name + "-managedserver",
		Image:           fmt.Sprintf("%s:%s", constants.WeblogicImageName, server.Spec.Version),
		ImagePullPolicy: v1.PullAlways,
		Ports: []v1.ContainerPort{{
			ContainerPort: 7001},
		},
		VolumeMounts: []v1.VolumeMount{{
			Name:      server.Spec.Domain.Name + "_storage",
			MountPath: "/u01/oracle/user_projects"},
		},
		Env: []v1.EnvVar{
			oracleHomeEnvVar(),
			serverNameEnvVar(server),
			serverNamespaceEnvVar(),
		},
		Command: []string{"/u01/oracle/weblogic-operator/startServer.sh",
			"/u01/oracle/user_projects/domains/" + server.Spec.Domain.Name,
			server.Name,
			"weblogic",
			"welcome1",
		},
		Lifecycle: &v1.Lifecycle{
			PostStart: &v1.Handler{
				Exec: &v1.ExecAction{
					Command: []string{"echo Hello World!!"},
				},
			},
			PreStop: &v1.Handler{
				Exec: &v1.ExecAction{
					Command: []string{"/u01/oracle/user_projects/domains/" + server.Spec.Domain.Name + "/bin/stopManagedWebLogic.sh",
						server.Name,
					},
				},
			},
		},
	}
}

// NewForServer creates a new ReplicationController for the given WebLogicManagedServer.
func NewForServer(server *types.WebLogicManagedServer, serviceName string) *v1beta1.ReplicaSet {
	containers := []v1.Container{WebLogicManagedServerContainer(server)}

	rs := &v1beta1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: server.Namespace,
			Name:      server.Name,
			Labels: map[string]string{
				constants.WebLogicManagedServerLabel: server.Name,
				constants.WebLogicDomainLabel:        server.Spec.Domain.Name,
				server.Spec.Domain.Name:              "managedserver",
			},
		},
		Spec: v1beta1.ReplicaSetSpec{
			Replicas:        &server.Spec.Replicas,
			MinReadySeconds: 0,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					constants.WebLogicManagedServerLabel: server.Name,
					constants.WebLogicDomainLabel:        server.Spec.Domain.Name,
					server.Spec.Domain.Name:              "managedserver",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						constants.WebLogicManagedServerLabel: server.Name,
						constants.WebLogicDomainLabel:        server.Spec.Domain.Name,
						server.Spec.Domain.Name:              "managedserver",
					},
				},
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{{
						Name: server.Spec.Domain.Name + "_storage",
						VolumeSource: v1.VolumeSource{
							PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
								ClaimName: "weblogic-operator-claim",
							},
						},
					},
					},
					NodeSelector: server.Spec.NodeSelector,
					ImagePullSecrets: []v1.LocalObjectReference{{
						Name: "weblogic-docker-store",
					},
					},
					Containers: containers,
				},
			},
		},
	}

	return rs
}
