package replicasets

import (
	"fmt"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"weblogic-operator/pkg/constants"
	"weblogic-operator/pkg/types"
)

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

func podNameEnvVar() v1.EnvVar {
	return v1.EnvVar{
		Name: "MY_POD_NAME",
		ValueFrom: &v1.EnvVarSource{
			FieldRef: &v1.ObjectFieldSelector{
				FieldPath: "metadata.name",
			},
		},
	}
}

// Builds the WebLogicManagedServer container
func WebLogicManagedServerContainer(server *types.WebLogicManagedServer) v1.Container {
	return v1.Container{
		Name:            server.Spec.DomainName + "-managedserver",
		Image:           fmt.Sprintf("%s:%s", constants.WeblogicImageName, server.Spec.Domain.Spec.Version),
		ImagePullPolicy: v1.PullIfNotPresent,
		//Ports: []v1.ContainerPort{{
		//	ContainerPort: 7001},
		//},
		VolumeMounts: []v1.VolumeMount{{
			Name:      server.Spec.DomainName + "-storage",
			MountPath: "/u01/oracle/user_projects"},
		},
		Env: []v1.EnvVar{
			oracleHomeEnvVar(),
			podNameEnvVar(),
			domainNameEnvVar(&server.Spec.Domain),
			domainHomeEnvVar(&server.Spec.Domain),
			serverNamespaceEnvVar(),
		},
		Command: []string{"/u01/oracle/user_projects/startServer.sh"},
		Lifecycle: &v1.Lifecycle{
			PreStop: &v1.Handler{
				Exec: &v1.ExecAction{
					Command: []string{"/u01/oracle/user_projects/stopServer.sh"},
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
				constants.WebLogicDomainLabel:        server.Spec.DomainName,
				server.Spec.DomainName:               "managedserver",
			},
		},
		Spec: v1beta1.ReplicaSetSpec{
			Replicas:        &server.Spec.ServersToRun,
			MinReadySeconds: 0,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					constants.WebLogicManagedServerLabel: server.Name,
					constants.WebLogicDomainLabel:        server.Spec.DomainName,
					server.Spec.DomainName:               "managedserver",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: server.Spec.DomainName + "-managedserver",
					Labels: map[string]string{
						constants.WebLogicManagedServerLabel: server.Name,
						constants.WebLogicDomainLabel:        server.Spec.DomainName,
						server.Spec.DomainName:               "managedserver",
					},
				},
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{
						{
							Name: server.Spec.DomainName + "-storage",
							VolumeSource: v1.VolumeSource{
								PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
									ClaimName: "weblogic-operator-claim",
								},
							},
						},
					},
					//TODO: refer to same selector of this.replicaset spec
					NodeSelector: server.Spec.NodeSelector,
					ImagePullSecrets: []v1.LocalObjectReference{
						{
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
