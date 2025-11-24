#######################################
# 2.1 指定基础镜像（FROM）
#######################################
# 使用官方 Go 运行时作为构建阶段的基础镜像

FROM golang:1.23-alpine AS builder

# 可选：加一些构建需要的工具
RUN apk add --no-cache git ca-certificates
ENV GOPROXY=https://goproxy.cn,direct

# 2.2 设置工作目录（WORKDIR）
WORKDIR /app



# 2.3 复制项目文件到容器（COPY）
# 先复制 go.mod 和 go.sum，下载 Go 依赖
COPY go.mod go.sum ./
RUN go mod download

# 再复制剩下的项目代码
COPY . .

# 2.4 编译 Go 二进制
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# -o 指定输出二进制文件名
RUN go build -o /app/bin/app  /app/cmd/web/main.go


#######################################
# 下面是第二阶段：构建运行镜像（更小、更干净）
#######################################
FROM alpine:3.19

# 安装运行时需要的证书、时区
RUN apk add --no-cache ca-certificates tzdata


# 2.2 再设置一次工作目录（运行阶段）
WORKDIR /app

# 2.3 从 builder 阶段复制构建好的二进制
COPY --from=builder /app/bin/app /app/app

COPY Config ./Config

# 2.6 设置环境变量（ENV）[可选]
ENV APP_ENV=production \
    APP_PORT=8080 \
    PATHEXT="sercrt"

# 2.5 暴露端口（EXPOSE）
EXPOSE 8080
EXPOSE 3008

RUN ls -R /app

# 2.7 定义默认启动命令（CMD）
CMD ["/app/app"]