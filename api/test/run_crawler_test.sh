#!/bin/bash

# 显示执行的命令
set -x

# 检查Chrome/Chromium浏览器是否已安装
check_chrome() {
  echo "检查Chrome/Chromium浏览器是否安装..."
  which chromium-browser || which google-chrome || which chromium
  if [ $? -ne 0 ]; then
    echo "警告: 未找到Chrome或Chromium浏览器，测试可能会失败"
    echo "尝试安装必要的依赖..."
    sudo apt-get update && sudo apt-get install -y libnss3 libgbm1 libasound2 libatk1.0-0 libatk-bridge2.0-0 libcups2 libxkbcommon0 libxdamage1 libxfixes3 libxrandr2 libxcomposite1 libxcursor1 libxi6 libxtst6 libnss3 libnspr4 libpango1.0-0
    echo "已安装依赖，但仍需安装Chrome/Chromium浏览器以完全支持测试"
  else
    echo "已找到Chrome/Chromium浏览器"
  fi
}

# 切换到项目根目录
cd ..

# 检查环境
check_chrome

# 运行爬虫测试，使用超时限制
echo "开始运行爬虫测试..."
timeout 180s go test -v ./test/crawler_test.go -run "TestNewService|TestSearchWeb"
TEST_RESULT=$?

if [ $TEST_RESULT -eq 124 ]; then
  echo "测试超时终止"
  exit 1
elif [ $TEST_RESULT -ne 0 ]; then
  echo "测试失败，退出码: $TEST_RESULT"
  exit $TEST_RESULT
else
  echo "测试成功完成"
fi

echo "测试完成" 