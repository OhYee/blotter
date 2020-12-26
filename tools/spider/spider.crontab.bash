#!/bin/bash

# ROOT user need next line to set environment variables
# export PYTHONPATH="/usr/lib/python38.zip:/usr/lib/python3.8:/usr/lib/python3.8/lib-dynload:/home/ubuntu/.local/lib/python3.8/site-packages:/usr/local/lib/python3.8/dist-packages:/usr/lib/python3/dist-packages"

SHELL_FOLDER=$(cd "$(dirname "$0")";pwd)

python3 $SHELL_FOLDER/main.py 1>>$SHELL_FOLDER/log.log 2>>$SHELL_FOLDER/error.log