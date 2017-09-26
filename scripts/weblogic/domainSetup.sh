#!/bin/bash
echo ------------------------------------------------------------------------------------------
echo Kubernetes Domain Setup
echo ------------------------------------------------------------------------------------------
echo Start - Domain Setup
echo $ORACLE_HOME/oracle_common/common/bin/wlst.sh kubeCreateDomain.py $1 $2 $3 $4 $5 $6 $7
echo End - Domain Setup
echo ------------------------------------------------------------------------------------------
echo 
echo ------------------------------------------------------------------------------------------
echo Start - Pack Domain
echo $ORACLE_HOME/oracle_common/common/bin/pack.sh $1
echo End - Pack Domain
echo ------------------------------------------------------------------------------------------