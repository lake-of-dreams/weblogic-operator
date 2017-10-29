#!/bin/bash

export ORACLE_HOME="/u01/oracle"

/u01/oracle/user_projects/groovy/bin/groovy -cp "$ORACLE_HOME/oracle_common/modules/features/oracle.fmwplatform.fmwprov_lib.jar:$ORACLE_HOME/oracle_common/modules/fmwplatform/common/util.jar:$ORACLE_HOME/mwhome/oracle_common/modules/features/cieCfg_wls_lib.jar" -Djava.util.logging.config.file=jul-logger.properties $*
