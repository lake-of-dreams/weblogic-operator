package constants

// WebLogicManagedServerLabel is applied to all components of a Weblogic Server
const (
	WebLogicGroupName = "weblogic.oracle.com"

	WebLogicManagedServerLabel              = "WebLogicManagedServer.v1.weblogic.oracle.com"
	WebLogicManagedServerResourceKind       = "WebLogicManagedServer"
	WebLogicManagedServerResourceKindPlural = "weblogicmanagedservers"
	WebLogicManagedServerSchemeVersion      = "v1"

	WebLogicDomainLabel              = "WebLogicDomain.v1.weblogic.oracle.com"
	WebLogicDomainResourceKind       = "WebLogicDomain"
	WebLogicDomainResourceKindPlural = "weblogicdomains"
	WebLogicDomainSchemeVersion      = "v1"

	//Constants for Horizontal Pod Autoscaling
	HorizontalPodAutoscalerKind        = "ReplicaSet"
	HorizontalPodAutoscalerKindPlural  = "replicasets"
	HorizontalPodAutoscalerName        = "managedserver-scaler"
	HorizontalPodAutoscalerTargetLabel = "managedserver"

	WeblogicImageName = "docker.io/store/oracle/weblogic"
)
