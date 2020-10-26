#!/bin/bash
SHELL_FOLDER=$(cd "$(dirname "$0")";pwd)

mongodump -o $SHELL_FOLDER/backup/`date '+%Y_%m_%d'` 1>>$SHELL_FOLDER/log.log 2>>$SHELL_FOLDER/error.log
