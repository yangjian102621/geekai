FROM golang:1.19.0 as go-deps
WORKDIR /src
COPY go.mod .
COPY go.sum .
# 预下载依赖
RUN go mod download

# 开始第二个构建阶段，命名为 go-builder
FROM golang:1.19.0 as go-builder
WORKDIR /src
# 被 go-deps 阶段预下载的依赖会存放在此
COPY --from=go-deps /go/pkg/mod /go/pkg/mod
COPY . .
# 构建应用
RUN go build -v -o app.out main.go && ls -al

# 以 debian:stable-slim 作为基础镜像
FROM debian:stable-slim
WORKDIR /var/www/app
# 从 go-builder 阶段拷贝构建产物
COPY --from=go-builder /src/app.out /var/www/app/app

# 设置容器启动时执行的命令
CMD ["/var/www/app/app"]
