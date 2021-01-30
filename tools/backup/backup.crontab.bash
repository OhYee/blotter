#!/bin/bash
SHELL_FOLDER=$(cd "$(dirname "$0")";pwd)
NAME=`date '+%Y_%m_%d'`
FILENAME="${SHELL_FOLDER}/backup/${NAME}"
ONEDRIVE=$(cat tools/tools.conf | grep -E '^OneDrive\s+.+$' | tr -s " " | cut -d " " -f 2)
SERVERCHAN=$(cat tools/tools.conf | grep -E '^ServerChan\s+.+$' | tr -s " " | cut -d " " -f 2)

function notify() {
    if [[ -n ${SERVERCHAN} ]]; then
        curl -X POST "http://sc.ftqq.com/${SERVERCHAN}.send" \
            -G \
            --data-urlencode "text=${1}" \
            --data-urlencode "desp=${2}"
    fi
}

mongodump -o $FILENAME 1>>$SHELL_FOLDER/log.log 2>>$SHELL_FOLDER/error.log
if [[ $? -ne 0 ]]; then
    notify "Backup error" ""
else
    if [[ -n ${ONEDRIVE} ]]; then
        zip -r "${FILENAME}.zip" "${FILENAME}/" 1>>$SHELL_FOLDER/log.log 2>>$SHELL_FOLDER/error.log \
            && onedrivecmd put "${FILENAME}.zip" "od:${ONEDRIVE}" 1>>$SHELL_FOLDER/log.log 2>>$SHELL_FOLDER/error.log
    fi
    if [[ $? -ne 0 ]]; then
        notify "Backup upload error" ""
    fi
fi



