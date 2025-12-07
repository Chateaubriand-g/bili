FROM golang:1.25-alpine AS builder

WORKDIR /app

#COPY go.mod go.sum ./
COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go get github.com/Chateaubriand-g/bili/common && \
    go mod tidy

#COPY . .

# 编译网关服务：
# 1. CGO_ENABLED=0：禁用 CGO，生成静态二进制（不依赖系统库，Alpine 可直接运行）
# 2. -ldflags "-w -s"：移除调试信息，减小二进制体积
# 3. -o gateway：指定输出文件名为 gateway（与启动命令一致）
RUN CGO_ENABLED=0 go build -ldflags "-w -s" -o auth_service main.go


FROM alpine:3.20

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# 从编译阶段复制编译好的二进制文件（仅复制必要文件，减小镜像体积）
COPY --from=builder /app/auth_service ./
COPY --from=builder /app/config /app/config

# 赋予二进制执行权限（关键：避免 permission denied 错误）
RUN chmod +x ./auth_service && \
    # 将工作目录权限赋予非 root 用户
    chown -R appuser:appgroup /app

# 切换到非 root 用户运行服务（生产环境强制要求，降低安全风险）
USER appuser

# 暴露网关服务的端口（根据你的代码监听端口修改，如 8081）
EXPOSE 8081

# 启动命令（运行编译好的二进制）
CMD ["./auth_service"]
