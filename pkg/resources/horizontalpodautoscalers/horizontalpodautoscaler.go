package horizontalpodautoscalers

import (
	"k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"weblogic-operator/pkg/constants"
	"weblogic-operator/pkg/types"
)

// NewForHorizontalPodAutoscaling creates a new HPA for the given WebLogicManagedServer - managedserver.
func NewForHorizontalPodAutoscaling(server *types.WebLogicManagedServer, serviceName string) *v1.HorizontalPodAutoscaler {
	var minReplicas int32 = 0
	var maxReplicas int32 = int32(server.Spec.Domain.Spec.ManagedServerCount)
	var targetCPUUtilization int32 = 50
	hpaMinReplicas := &minReplicas
	hpaTargetCPUUtilization := &targetCPUUtilization

	hpa := &v1.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name: constants.HorizontalPodAutoscalerName,
		},
		Spec: v1.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: v1.CrossVersionObjectReference{
				Kind: constants.HorizontalPodAutoscalerKind,
				Name: server.Name,
			},
			MinReplicas:                    hpaMinReplicas,
			MaxReplicas:                    maxReplicas,
			TargetCPUUtilizationPercentage: hpaTargetCPUUtilization,
		},
	}

	return hpa
}
