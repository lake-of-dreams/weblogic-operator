package constants

// WeblogicServerLabel is applied to all components of a Weblogic Server
const (
	WeblogicServerLabel              = "WeblogicServer.v1.weblogic.oracle.com"
	WeblogicServerResourceKind       = "WeblogicServer"
	WeblogicServerResourceKindPlural = "weblogicservers"
	WeblogicServerGroupName          = "weblogic.oracle.com"
	WeblogicServerSchemeVersion      = "v1"
	WeblogicImageName                = "docker.io/store/oracle/weblogic"

	WebLogicDomainLabel              = "WebLogicDomain.v1.weblogic.oracle.com"
	WebLogicDomainResourceKind       = "WebLogicDomain"
	WebLogicDomainResourceKindPlural = "weblogicdomains"
	WebLogicGroupName		         = "weblogic.oracle.com"
	WeblogicDomainSchemeVersion      = "v1"
)
