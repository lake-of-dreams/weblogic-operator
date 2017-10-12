#!/bin/bash
echo ------------------------------------------------------------------------------------------
echo Kubernetes Domain Setup
echo ------------------------------------------------------------------------------------------

echo Start - Domain Setup
if [ ! -d ${DOMAIN_HOME} ]; then
    $ORACLE_HOME/oracle_common/common/bin/wlst.sh -skipWLSModuleScanning /u01/oracle/user_projects/kubeCreateDomain.py \
                                                        $MY_POD_NAME $ORACLE_HOME $DOMAIN_NAME $DOMAIN_HOME $MANAGED_SERVER_COUNT "7001" "weblogic" "welcome1" \
                                                        >> /u01/oracle/user_projects/domainSetup"_${DOMAIN_NAME}".log 2>&1
fi
echo End - Domain Setup

echo ------------------------------------------------------------------------------------------

echo Start - Admin Start
if [ -d ${DOMAIN_HOME} ]; then
    mkdir -p ${DOMAIN_HOME}/servers/AdminServer/security/
    echo "username=weblogic" > ${DOMAIN_HOME}/servers/AdminServer/security/boot.properties
    echo "password=welcome1" >> ${DOMAIN_HOME}/servers/AdminServer/security/boot.properties
    ${DOMAIN_HOME}/bin/setDomainEnv.sh

    ${DOMAIN_HOME}/bin/startWebLogic.sh
    touch ${DOMAIN_HOME}/servers/AdminServer/logs/AdminServer.log
    tail -f ${DOMAIN_HOME}/servers/AdminServer/logs/AdminServer.log &
fi
echo Stop - Admin Start

echo ------------------------------------------------------------------------------------------

