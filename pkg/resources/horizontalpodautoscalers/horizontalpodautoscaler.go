package horizontalpodautoscalers

import (
	"weblogic-operator/pkg/constants"
	"weblogic-operator/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/autoscaling/v1"
)

// NewForHorizontalPodAutoscaling creates a new HPA for the given WeblogicServer - managedserver.
func NewForHorizontalPodAutoscaling(server *types.WeblogicServer, serviceName string) *v1.HorizontalPodAutoscaler {
	var minReplicas int32 = 1
	var maxReplicas int32 = 5
	var targetCPUUtilization int32 = 50
	hpaMinReplicas := &minReplicas
	hpaTargetCPUUtilization := &targetCPUUtilization

	hpa := &v1.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      constants.HorizontalPodAutoscalerName,
		},
		Spec: v1.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: v1.CrossVersionObjectReference{
				Kind: constants.HorizontalPodAutoscalerKind,
				Name: constants.HorizontalPodAutoscalerTargetLabel,
			},
			MinReplicas: hpaMinReplicas,
			MaxReplicas: maxReplicas,
			TargetCPUUtilizationPercentage: hpaTargetCPUUtilization,
		},
	}

	return hpa
}