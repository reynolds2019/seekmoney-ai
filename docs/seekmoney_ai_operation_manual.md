# SeekMoney AI 运维手册 (Operation Manual)

> 本手册记录 SeekMoney AI 在生产服务器上的关键配置、架构信息、日常运维指令及灾备维护流程。

---

## 一、部署基本信息

| 项目 | 配置信息 | 说明 |
| :--- | :--- | :--- |
| **软件名称** | seekmoney-ai | 社交媒体痛点挖掘与商业机会分析系统 |
| **服务器 IP** | `192.168.8.206` | 内网服务器 |
| **服务器安装目录** | `/data/seekmoney-ai` | 部署于 512GB 数据盘 |
| **容器名称** | `seekmoney-ai` | Docker 容器化运行 |
| **服务端口** | `3000` | 容器端口 `3000` 映射宿主机 `3000` |
| **访问地址** | `http://192.168.8.206:3000` | 默认跳转至 `http://192.168.8.206:3000/zh` |
| **基础依赖镜像** | `node:22-alpine` | 使用服务器已有镜像多阶段构建 |
| **自启策略** | `--restart=unless-stopped` | 宿主机重启或服务异常退出后自动拉起 |

---

## 二、持久化挂载与权限设计

为了保证容器升级无损及数据安全，所有核心数据、配置与日志均挂载至宿主机卷：

| 宿主机路径 | 容器内挂载路径 | 权限要求 | 挂载说明 |
| :--- | :--- | :--- | :--- |
| `/data/seekmoney-ai/.env` | `/app/.env` | `600` (仅属主读写) | 核心密钥与模型连接配置 |
| `/data/seekmoney-ai/analysis-results/` | `/app/analysis-results/` | `775` (属主 `1000:1000`) | 痛点分析任务结果持久化存储 |
| `/data/seekmoney-ai/logs/` | `/app/logs/` | `775` (属主 `1000:1000`) | 系统与业务运行日志 |

> **权限校准命令**（如遇权限问题可直接执行）：
> ```bash
> chown -R 1000:1000 /data/seekmoney-ai/analysis-results /data/seekmoney-ai/logs
> chmod -R 775 /data/seekmoney-ai/analysis-results /data/seekmoney-ai/logs
> chmod 600 /data/seekmoney-ai/.env
> ```

---

## 三、核心环境配置清单 (`.env`)

配置文件位于宿主机 `/data/seekmoney-ai/.env`：

```env
# ==========================================
# SeekMoney AI 生产环境变量配置
# ==========================================

# 1. TikHub API 令牌 (数据抓取服务)
TIKHUB_API_TOKEN=your_tikhub_api_token_here
TIKHUB_TIMEOUT=60000

# 2. Chat 分析模型接入通义千问 (Qwen)
GLM_API_KEY=your_dashscope_or_qwen_api_key_here
GLM_BASE_URL=http://192.168.8.206:3002/v1
GLM_MODEL_NAME=qwen3.6-flash

# 3. Embedding 聚类接入通义千问 (OpenAI 兼容协议)
EMBEDDING_PROVIDER=openai
OPENAI_API_KEY=your_dashscope_or_qwen_api_key_here
OPENAI_BASE_URL=http://192.168.8.206:3002/v1
OPENAI_EMBEDDING_MODEL=text-embedding-v4
# OPENAI_EMBEDDING_BATCH_SIZE=10 # 可选，通义千问模型系统已默认自适应为10，其他模型默认100

# 4. 系统运行环境配置
NODE_ENV=production
HEADLESS=true
LOG_LEVEL=info
PYTHONIOENCODING=utf-8
SEAGULL_STORAGE_DIR=/app/analysis-results
```

---

## 四、常用日常运维命令

所有命令均在服务器 `/data/seekmoney-ai` 目录下执行：

### 1. 服务启停与状态查看
```bash
cd /data/seekmoney-ai

# 查看容器运行状态
docker compose ps

# 查看实时日志 (最后 100 行，持续输出)
docker compose logs -f --tail=100

# 停止服务
docker compose down

# 启动服务 (后台运行)
docker compose up -d

# 重启服务 (重新加载代码/配置后)
docker compose restart
```

### 2. 配置变更后重新生效
修改 `/data/seekmoney-ai/.env` 文件后，必须使用 `up -d` 重新注入环境变量：
```bash
cd /data/seekmoney-ai
docker compose up -d
```

### 3. 代码更新与重新构建
当有新版本代码更新后：
```bash
cd /data/seekmoney-ai

# 1. 重新构建 Docker 镜像
docker compose build

# 2. 重启容器
docker compose up -d
```

### 4. 容器内调试与交互
```bash
# 进入容器 Shell
docker exec -it seekmoney-ai sh

# 查看容器当前加载的环境变量
docker exec seekmoney-ai env
```

---

## 五、健康检查与接口测试

### 1. 本地/内部健康检查
```bash
curl -i http://127.0.0.1:3000/api/health
```
正常响应示例：
```json
HTTP/1.1 200 OK
Content-Type: application/json

{"status":"healthy","timestamp":"2026-10-10T09:13:38.674Z","uptime":41.38}
```

### 2. 首页页面可达性测试
```bash
curl -I http://127.0.0.1:3000/zh
```
正常响应 HTTP 状态码为 `200 OK`。

---

## 六、数据备份与灾难恢复

### 1. 数据备份
建议将配置与分析数据定期备份至归档目录（如 `/data/back/`）：
```bash
# 创建备份目录
mkdir -p /data/back/seekmoney-ai

# 备份配置文件与分析数据 (带时间戳)
tar -czvf /data/back/seekmoney-ai/seekmoney_backup_$(date +%Y%m%d_%H%M%S).tar.gz \
    -C /data/seekmoney-ai .env analysis-results/
```

### 2. 灾难恢复
```bash
# 解压备份文件至目标安装目录
tar -xzvf /data/back/seekmoney-ai/seekmoney_backup_YYYYMMDD_HHMMSS.tar.gz -C /data/seekmoney-ai/

# 恢复权限并拉起服务
chown -R 1000:1000 /data/seekmoney-ai/analysis-results
docker compose up -d
```

---

## 七、常见问题与排查指引

| 问题现象 | 可能原因 | 解决办法 |
| :--- | :--- | :--- |
| **容器无法启动或 CrashLoop** | 端口 3000 发生冲突 | 编辑 `docker-compose.yml` 将宿主机端口修改为 `"3003:3000"`，重新执行 `docker compose up -d` |
| **分析结果无法保存** | 挂载目录无写入权限 | 执行 `chown -R 1000:1000 /data/seekmoney-ai/analysis-results && chmod -R 775 /data/seekmoney-ai/analysis-results` |
| **大模型调用报 401/404** | API Key 失效或 URL 错误 | 检查 `.env` 中的 `GLM_API_KEY` 与 `GLM_BASE_URL`，确认 `http://192.168.8.206:3002/v1` 的 new-api 服务状态 |
| **TikHub 抓取超时** | 网络延迟或 Token 无效 | 检查 `.env` 中的 `TIKHUB_API_TOKEN` 是否最新，适当增大 `TIKHUB_TIMEOUT=120000` |
