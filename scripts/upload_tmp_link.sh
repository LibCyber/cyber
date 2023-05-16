#!/bin/bash

# 请求用户输入一个token值
echo -n "请输入TmpLink密钥："
read -r token

# 获取当前目录下的build文件夹下的所有文件名
files=(./build/*)

# 遍历文件夹下的所有文件
for file in "${files[@]}"
do
  # 检查是否为文件
  if [ -f "$file" ]
  then
    # 执行curl命令
    curl -k -F "file=@$file" -F "token=$token" -F "model=99" -F "mrid=646369abc972a" -X POST "https://connect.tmp.link/api_v2/cli_uploader"
  fi
done
