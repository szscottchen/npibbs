#!/bin/sh

export APP_HOME=/app/bbs-go

# 启动 bbs-go-server
cd ${APP_HOME}/server
${APP_HOME}/server/tkdev &

# 启动 bbs-go-site
cd ${APP_HOME}/site
# 使用serve静态文件服务器提供前台服务
npx serve -s dist -l 3000 &

# 保持容器运行
wait
