#!/bin/bash
SHELL_FOLDER=$(cd "$(dirname "$0")";pwd)

curl https://www.oyohyee.com/sitemap.txt 1>$SHELL_FOLDER/urls.txt 2>>$SHELL_FOLDER/error.log

echo `date '+%Y-%m-%d %H:%M:%S'` >> log.log
curl -H 'Content-Type:text/plain' --data-binary @$SHELL_FOLDER/urls.txt "http://data.zz.baidu.com/urls?site=www.oyohyee.com&token=Ah36uCNcnwq2q7LX" 1>>$SHELL_FOLDER/log.log 2>>$SHELL_FOLDER/error.log
echo "" >> log.log
