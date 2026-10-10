# ==========================================
# 阶段 1: 依赖安装 (基于服务器现有的 node:22-alpine 镜像)
# ==========================================
FROM node:22-alpine AS deps
WORKDIR /app

# 安装可能需要的原生构建基础工具
RUN apk add --no-cache libc6-compat

COPY package.json package-lock.json ./
RUN npm ci

# ==========================================
# 阶段 2: 项目构建
# ==========================================
FROM node:22-alpine AS builder
WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY . .

ENV NEXT_TELEMETRY_DISABLED=1
RUN npm run build

# ==========================================
# 阶段 3: 生产运行时镜像
# ==========================================
FROM node:22-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

# 创建数据存储与日志目录
RUN mkdir -p /app/analysis-results /app/logs

# 从构建产物中复制静态资源与独立运行时
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

# 暴露运行端口
EXPOSE 3000

CMD ["node", "server.js"]
