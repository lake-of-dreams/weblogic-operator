package replicasets

import (
	"fmt"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"weblogic-operator/pkg/constants"
	"k8s.io/api/extensions/v1beta1"
	"weblogic-operator/pkg/types"
)

func oracleHomeEnvVar() v1.EnvVar {
	return v1.EnvVar{Name: "ORACLE_HOME", Value: "/u01/oracle"}
}

func domainNameEnvVar(domain *types.WebLogicDomain) v1.EnvVar {
	return v1.EnvVar{Name: "DOMAIN_NAME", Value: domain.Name}
}

func domainHomeEnvVar(domain *types.WebLogicDomain) v1.EnvVar {
	return v1.EnvVar{Name: "DOMAIN_HOME", Value: "/u01/oracle/user_projects/domains/" + domain.Name}
}

func managedServerCountEnvVar(domain *types.WebLogicDomain) v1.EnvVar {
	return v1.EnvVar{Name: "MANAGED_SERVER_COUNT", Value: domain.Spec.ManagedServerCount}
}

func domainNamespaceEnvVar() v1.EnvVar {
	return v1.EnvVar{
		Name: "POD_NAMESPACE",
		ValueFrom: &v1.EnvVarSource{
			FieldRef: &v1.ObjectFieldSelector{
				FieldPath: "metadata.namespace",
			},
		},
	}
}

// Builds the WebLogicDomain container
func weblogicDomainContainer(domain *types.WebLogicDomain) v1.Container {
	return v1.Container{
		Name:            domain.Name + "-adminserver",
		Image:           fmt.Sprintf("%s:%s", constants.WeblogicImageName, domain.Spec.Version),
		ImagePullPolicy: v1.PullAlways,
		Ports: []v1.ContainerPort{{
			ContainerPort: 7001},
		},
		VolumeMounts: []v1.VolumeMount{{
			Name:      domain.Name + "-storage",
			MountPath: "/u01/oracle/user_projects"},
		},
		Env: []v1.EnvVar{
			oracleHomeEnvVar(),
			domainNameEnvVar(domain),
			domainHomeEnvVar(domain),
			managedServerCountEnvVar(domain),
			domainNamespaceEnvVar(),
		},
		Command: []string{"/u01/oracle/user_projects/domainSetup.sh"},
		Lifecycle: &v1.Lifecycle{
			//PostStart: &v1.Handler{
			//	Exec: &v1.ExecAction{
			//		Command: []string{"echo Hello World!!"},
			//	},
			//},
			PreStop: &v1.Handler{
				Exec: &v1.ExecAction{
					Command: []string{"/u01/oracle/user_projects/domains/" + domain.Name + "/bin/stopWebLogic.sh"},
				},
			},
		},
	}
}

// NewForDomain creates a new ReplicationController for the given WebLogicDomain.
func NewForDomain(domain *types.WebLogicDomain, serviceName string) *v1beta1.ReplicaSet {
	containers := []v1.Container{weblogicDomainContainer(domain)}

	rs := &v1beta1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: domain.Namespace,
			Name:      domain.Name,
			Labels: map[string]string{
				constants.WebLogicDomainLabel: domain.Name,
			},
		},
		Spec: v1beta1.ReplicaSetSpec{
			Replicas:        &domain.Spec.Replicas,
			MinReadySeconds: 0,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					constants.WebLogicDomainLabel: domain.Name,
					domain.Name:                   "adminserver",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: domain.Name + "-adminserver",
					Labels: map[string]string{
						constants.WebLogicDomainLabel: domain.Name,
						domain.Name:                   "adminserver",
					},
				},
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{{
						Name: domain.Name + "-storage",
						VolumeSource: v1.VolumeSource{
							PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
								ClaimName: "weblogic-operator-claim",
							},
						},
					},
					},
					NodeSelector: domain.Spec.NodeSelector,
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
