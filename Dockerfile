# ================== 1. 构建阶段 (Builder) ==================
# 使用 Go 1.23 (对应你的 go.mod 要求)
FROM golang:1.23-alpine AS builder

# 优化 1: 替换 Alpine 软件源为阿里云 (加速基础软件安装)
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 优化 2: 安装必要工具
RUN apk add --no-cache git ca-certificates tzdata

# 优化 3: 设置 Go 编译环境 (使用国内七牛云代理加速模块下载)
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /build

# 优化 4: 利用 Docker 缓存机制，先只拷贝依赖描述文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 拷贝源码
COPY . .

# 优化 5: 整理依赖 (防止版本不一致报错)
RUN go mod tidy

# 编译 (输出为 govote)
RUN go build -o govote .

# ================== 2. 运行阶段 (Runner) ==================
FROM alpine:latest

# 优化 6: 再次替换源并安装时区/证书数据
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add --no-cache ca-certificates tzdata

# 优化 7: 设置时区为上海 (解决日志时间差 8 小时问题)
ENV TZ=Asia/Shanghai

WORKDIR /app

# 从构建阶段拷贝编译好的二进制文件
COPY --from=builder /build/govote .

# 关键: 拷贝配置文件目录 (否则启动会报错)
COPY config ./config

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["./govote"]
