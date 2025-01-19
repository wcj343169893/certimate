# 使用官方的 Golang 镜像作为基础镜像
FROM golang:1.20-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制所有文件到工作目录
COPY . .

# 下载依赖并构建项目
RUN go mod download
RUN go build -o certimate main.go

# 使用一个更小的基础镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制二进制文件到最终镜像
COPY --from=builder /app/certimate .

# 暴露端口
EXPOSE 8090

# 运行二进制文件
CMD ["./certimate", "serve", "--http", "0.0.0.0:8090"]