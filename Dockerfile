# server builder

FROM golang:1.24 AS server_builder

ENV APP_HOME=/code/bbs-go/server
WORKDIR "$APP_HOME"

COPY ./server ./
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
RUN CGO_ENABLED=0 go build -v -o tkdev main.go && chmod +x tkdev


# 注意：前端文件已在Windows环境下构建完成
# 确保在构建前已复制以下内容到Linux环境：
# - site/dist/ 目录 (前台静态文件)
# - admin/dist/ 目录 (后台静态文件)

# run
FROM node:20-alpine

ENV APP_HOME=/app/bbs-go
WORKDIR "$APP_HOME"

COPY --from=server_builder /code/bbs-go/server/tkdev ./server/tkdev
COPY --from=server_builder /code/bbs-go/server/migrations ./server/migrations
COPY --from=server_builder /code/bbs-go/server/locales ./server/locales
COPY --from=server_builder /code/bbs-go/server/bbs-go.yaml ./server/
# 直接使用已构建好的前端静态文件（确保这些目录存在）
COPY ./site/dist ./site/dist
COPY ./admin/dist ./server/admin

COPY ./start.sh ${APP_HOME}/start.sh
RUN chmod +x ${APP_HOME}/start.sh

EXPOSE 8082 3000

CMD ["./start.sh"]
