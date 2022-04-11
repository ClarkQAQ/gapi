#!/bin/bash

# 当前脚本所在目录
SCRIPT_DIR=$(cd $(dirname ${BASH_SOURCE[0]}); pwd)

# 进入脚本所在目录
cd $SCRIPT_DIR

# 删除测试文件
# rm -rf test
# mkdir -p test

# 清屏
clear

echo "workdir: `pwd`"

go run .

