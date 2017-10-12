#!/bin/bash

echo ------------------------------------------------------------------------------------------
echo Kubernetes Stop Managed Server Begin
echo ------------------------------------------------------------------------------------------

if [ -d ${DOMAIN_HOME} ]; then

    LOCKFILE=${DOMAIN_HOME}/ms.lok
    if [ -e ${LOCKFILE} ] && kill -0 `cat ${LOCKFILE}`; then
        echo "Already Running"
        exit
    fi

    trap "rm -f ${LOCKFILE}; exit" INT TERM EXIT
    echo $$ > ${LOCKFILE}

    json=${DOMAIN_HOME}/serverList.json
    file=${DOMAIN_HOME}/serverList.json.tmp

    podkey="podName"
    podre="\"($podkey)\": \"([^\"]*)\""

    mskey="serverName"
    msre="\"($mskey)\": \"([^\"]*)\""

    mypodname=${MY_POD_NAME}
    msname=""
    foundms=false

    while IFS='' read -r line || [[ -n "$line" ]]; do

        if [ "$foundms" = false ]; then
            if [[ $line =~ $podre ]]; then
                podname="${BASH_REMATCH[2]}"

                if [[ $podname = ${mypodname} ]]; then
                    if [[ $line =~ $msre ]]; then
                        msname="${BASH_REMATCH[2]}"
                        if [[ $msname != "AdminServer" ]]; then
                            echo "Found Managed Server $msname Running in Pod $mypodname"
                            foundms=true
                            line=${line/\"podName\": \"$mypodname\"/\"podName\": \"\"}
                        else
                            msname=""
                        fi
                    fi
                fi
            fi
        fi

        echo $line >> $file

    done < "$json"

    mv $file $json

    rm -f ${LOCKFILE}

    if [[ ! -z "${msname// }" ]]; then
        echo "Stopping $msname..."

        ${DOMAIN_HOME}/bin/stopManagedWebLogic.sh ${msname} "t3://${DOMAIN_NAME}:7001"
    fi
fi

echo ------------------------------------------------------------------------------------------
echo Kubernetes Stop Managed Server End
echo ------------------------------------------------------------------------------------------
