FROM golang:1.20 as build

# 容器环境变量添加，会覆盖默认的变量值
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

# 设置工作区
WORKDIR /go/release

# 把全部文件添加到/go/release目录
ADD .. .

# 编译：把cmd/main.go编译成可执行的二进制文件，命名为app
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o app main.go

# 运行：使用alpine作为基础镜像
FROM alpine as prod

RUN apk add --no-cache ca-certificates

# 在build阶段复制可执行的go二进制文件app
COPY --from=build /go/release/app /

# 启动服务
CMD ["/app"]