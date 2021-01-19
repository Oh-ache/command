#! /bin/sh
#
# translate.sh
# Copyright (C) 2021 ache <1751987128@qq.com>
#
# Distributed under terms of the MIT license.
#

baiDuUrl='http://api.fanyi.baidu.com/api/trans/vip/translate'
baiDuAppId="你的appid"
baiDuSercert="你的密钥"

rand=$RANDOM
q=$1

md5Res=`md5 -s "$baiDuAppId$1$rand$baiDuSercert" | awk -F " " '{print $NF}'`

data=`curl -s -H 'Content-Type:application/x-www-form-urlencoded' -X POST -d "q=$q&from=auto&to=zh&appid=$baiDuAppId&salt=$rand&sign=$md5Res" "$baiDuUrl"`

# brew安装jq
response=`echo $data | jq . |  awk -F "[dst]" '/dst/{print$0}'`

result=(${response//": "/ })
echo ${result[1]}
