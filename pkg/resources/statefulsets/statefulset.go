package statefulsets

import (
	"fmt"
	"k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/sczachariah/weblogic-operator/pkg/constants"
	"github.com/sczachariah/weblogic-operator/pkg/types"
)

// WeblogicImageName is the base Docker image for the operator.
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

// Builds the Weblogic operator container for a server
func weblogicOperatorContainer(server *types.WeblogicServer) v1.Container {
	return v1.Container{
		//TODO : Use different container names ???
		Name: "weblogic",
		Image: fmt.Sprintf("%s:%s", WeblogicImageName, server.Spec.Version),
		Ports: []v1.ContainerPort{{ContainerPort: 7001}},
		Env: []v1.EnvVar{
			serverNameEnvVar(server),
			namespaceEnvVar(),
		},
	}
}

// NewForServer creates a new StatefulSet for the given WeblogicServer.
func NewForServer(server *types.WeblogicServer, serviceName string) *v1beta1.StatefulSet {
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
					Containers:   containers,
				},
			},
			ServiceName: serviceName,
		},
	}

	return ss
}
