#!/bin/bash

# 当前 shell 的执行目录
SHELL_FOLDER=$(cd "$(dirname "$0")";pwd)
FOLDER="/home/ubuntu/blog/blotter"
LOG_FILE="${FOLDER}/action.log"

function func_log() {
    echo $@ &>>$LOG_FILE
}

# 判断当前是脚本自己递归调用还是其他程序调用
if [[ -z $OHYEE_ACTION_AFTER_PULL ]]; then 
    # 接收到请求，判断是否需要执行更新
    func_log $(env)
    func_log "Commit ${OHYEE_ACTION_REPO_GIT}@${OHYEE_ACTION_BRANCH} from ${OHYEE_ACTION_PUSHER} $OHYEE_ACTION_COMMIT $(date -u '+%Y-%m-%d %H:%M:%S')"
    func_log "Need update?"
    if [[ ${OHYEE_ACTION_BRANCH} == "refs/heads/master" && ${OHYEE_ACTION_TYPE} == "gitea" ]]; then
        func_log "Update files"

        # 强制更新
        cd $FOLDER
        git stash
        git fetch --all
        git reset --hard origin/master
        git pull
        git stash pop

        export OHYEE_ACTION_AFTER_PULL="true"
        bash $0
        unset OHYEE_ACTION_AFTER_PULL
    else
        func_log "No need"
    fi
    func_log "Finished at $(date -u '+%Y-%m-%d %H:%M:%S')"
else
    # 递归执行，直接执行更新后命令
    func_log "run action script"

    # 重新编译，并关闭原有程序，在 screen 启动当前程序
    cd $FOLDER

    SCREEN_NAME="back"
    bash $FOLDER/build.bash &>>$LOG_FILE
    screen -S ${SCREEN_NAME} -X quit &>>$LOG_FILE
    screen -dmS ${SCREEN_NAME} ./blotter &>>$LOG_FILE

    func_log "rebuild at $(date -u '+%Y-%m-%d %H:%M:%S')"
fi
