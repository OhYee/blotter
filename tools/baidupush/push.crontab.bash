#!/bin/bash
SHELL_FOLDER=$(cd "$(dirname "$0")";pwd)

curl https://www.oyohyee.com/sitemap.txt > $SHELL_FOLDER/urls.txt 1>>$SHELL_FOLDER/log.log 2>>$SHELL_FOLDER/error.log
curl -H 'Content-Type:text/plain' --data-binary @$SHELL_FOLDER/urls.txt "http://data.zz.baidu.com/urls?site=www.oyohyee.com&token=Ah36uCNcnwq2q7LX" 1>>$SHELL_FOLDER/log.log 2>>$SHELL_FOLDER/error.log
