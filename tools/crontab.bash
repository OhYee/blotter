#!/bin/bash
# 运行该文件将相关定时任务追加至 crontab 中

SHELL_FOLDER=$(cd "$(dirname "$0")";pwd)

crontab -l > ./crontab.temp
echo "
# blotter
30 * * * * /bin/bash $SHELL_FOLDER/baidupush/push.crontab.bash
0  3 * * * /bin/bash $SHELL_FOLDER/spider/spider.crontab.bash
10 5 * * * /bin/bash $SHELL_FOLDER/backup/backup.crontab.bash
" >> ./crontab.temp
crontab ./crontab.temp
rm crontab.temp

crontab -l