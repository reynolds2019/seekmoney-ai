# ==============================================================================
# 商机发现助手 (Oppty Discovery Assistant) - 后端服务 Dockerfile
# ==============================================================================

FROM alpine:latest

# 安装基础运行依赖（CA证书 + 时区数据）
RUN apk add --no-cache ca-certificates tzdata

# 设置时区为上海
ENV TZ=Asia/Shanghai

WORKDIR /app

# 拷贝 Linux 静态可执行文件
COPY server_linux /app/server
RUN chmod +x /app/server

# 暴露服务运行端口
EXPOSE 8080

# 启动命令，配置文件与日志通过目录卷挂载持久化
ENTRYPOINT ["./server"]
