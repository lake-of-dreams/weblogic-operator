#!/bin/bash
echo ------------------------------------------------------------------------------------------
echo Kubernetes Domain Setup
echo ------------------------------------------------------------------------------------------
echo Start - Domain Setup
$ORACLE_HOME/oracle_common/common/bin/wlst.sh -skipWLSModuleScanning /u01/oracle/user_projects/kubeCreateDomain.py >> /u01/oracle/user_projects/domainSetup.log 2>&1
echo End - Domain Setup
echo ------------------------------------------------------------------------------------------
echo
mkdir -p ${DOMAIN_HOME}/servers/AdminServer/security/
echo "username=weblogic" > ${DOMAIN_HOME}/servers/AdminServer/security/boot.properties
echo "password=welcome1" >> ${DOMAIN_HOME}/servers/AdminServer/security/boot.properties
${DOMAIN_HOME}/bin/setDomainEnv.sh
echo ------------------------------------------------------------------------------------------
# Start Admin Server and tail the logs
${DOMAIN_HOME}/bin/startWebLogic.sh
touch ${DOMAIN_HOME}/servers/AdminServer/logs/AdminServer.log
tail -f ${DOMAIN_HOME}/servers/AdminServer/logs/AdminServer.log &