# 商机发现助手 (Oppty Discovery Assistant) 容器化部署与运维手册

> **应用标识**：`oppty-discovery`  
> **服务版本**：`v1.0.7`  
> **归档位置**：本地 `docs/oppty_discovery_ops_guide.md` / 服务器 `/data/oppty-discovery-backend/oppty_discovery_ops_guide.md`  
> **更新日期**：2026-10-13  
> **维护责任人**：reynolds  

---

## 1. 系统架构与网络拓扑

本项目采用**前后端微服务分离 + 统一网关托管**架构，前端复用服务器统一反向代理网关 `common-nginx`，后端采用极简自包含 Docker 容器，并深度复用宿主机公共技术组件：

```mermaid
graph TD
    User([终端用户浏览器]) -->|HTTP /opp/| Nginx[common-nginx :80]

    subgraph 前端静态资源
        Nginx -->|/opp/| FrontFiles["/usr/share/nginx/html/opp/ (Vite React SPA)"]
    end

    subgraph 后端微服务
        Nginx -->|/opp/api/ 转发| Backend["oppty-discovery-backend :8084 -> :8080"]
    end

    subgraph 公共基础设施 (192.168.8.206)
        Backend -->|持久化| PG[(PostgreSQL 17 :5432 / oppty_discovery)]
        Backend -->|报告归档| MinIO[(MinIO OSS :9000 / oppty-discovery)]
        Backend -->|OTP验证码外发| Supabase[Supabase BaaS :8000]
        Backend -->|LLM商业解构| QwenLLM[通义千问推理引擎 :3002]
        Backend -->|外部社媒爬取| TikHub[TikHub API 网关 :443]
    end
```

---

## 2. 服务访问入口与端口矩阵

| 服务角色 | 访问协议与地址 | 宿主机端口 | 容器内部端口 | 备注说明 |
| :--- | :--- | :--- | :--- | :--- |
| **Web 前端应用** | `http://192.168.8.206/opp` (自动301至 `/opp/`) | `80` (HTTP) | `80` | 托管于 `common-nginx` |
| **后端 API 网关** | `http://192.168.8.206/opp/api/` | `80` (代理) | `8080` | 由 Nginx 统一反向代理 |
| **后端直连端口** | `http://192.168.8.206:8084/api/` | `8084` | `8080` | 容器直连端口（内部自检使用） |
| **健康自检探活** | `http://192.168.8.206/opp/api/health` | `80` | `8080` | 返回 JSON 探针健康状态 |

---

## 3. 服务器宿主机持久化目录结构

为保证数据与配置的**升级无损、独立备份**，所有配置、密钥及日志均持久化于 `/data` 目录：

```text
/data/
├── nginx/                                      # 统一 Web 网关目录
│   ├── conf.d/
│   │   └── default.conf                        # Nginx 虚拟主机配置（包含 /opp/ 与 /opp/api/ 路由规则）
│   └── html/
│       └── opp/                                # 商机发现助手前端打包构建产物 (dist)
│           ├── index.html
│           └── assets/
└── oppty-discovery-backend/                    # 后端服务安装与持久化根目录
    ├── server_linux                            # Linux 64位静态无依赖二进制执行文件
    ├── Dockerfile                              # 生产镜像构建文件 (基于 alpine:latest)
    ├── .env                                    # 【核心私密文件】私密环境变量 (包含真实 TikHub Token)
    ├── configs/
    │   └── config.yaml                         # 核心业务配置文件 (数据库、MinIO、Supabase、LLM配置)
    ├── logs/                                   # 应用运行日志目录 (持久化读写卷)
    └── oppty_discovery_ops_guide.md            # 本运维手册
```

---

## 4. 核心配置与环境变量

### 4.1 私密环境变量文件 (`/data/oppty-discovery-backend/.env`)
已脱敏隔离，切勿提交至公开代码仓库：
```env
# TikHub API 访问凭证 Token
TIKHUB_API_TOKEN=your_tikhub_api_token_here
```

### 4.2 业务主配置文件 (`/data/oppty-discovery-backend/configs/config.yaml`)
关键组件连接配置段已复用服务器既有基础设施：
- **数据库**：`host: 192.168.8.206`, `port: 5432`, `dbname: oppty_discovery`
- **对象存储**：`endpoint: 192.168.8.206:9000`, `bucketName: oppty-discovery`
- **邮件服务**：`supabase_url: http://192.168.8.206:8000`, `code_expire_minutes: 5`
- **大模型引擎**：`baseUrl: http://192.168.8.206:3002/v1`, `modelName: qwen3.6-flash`

---

## 5. 日常运维管理命令

### 5.1 后端容器生命周期管理
```bash
# 查看容器运行状态
docker ps -f name=oppty-discovery-backend

# 查看后端实时日志
docker logs -f --tail 100 oppty-discovery-backend

# 重启后端服务
docker restart oppty-discovery-backend

# 停止并删除后端容器
docker stop oppty-discovery-backend && docker rm oppty-discovery-backend

# 基于宿主机挂载卷全新启动后端容器
docker run -d \
  --name oppty-discovery-backend \
  -p 8084:8080 \
  -v /data/oppty-discovery-backend/configs:/app/configs:ro \
  -v /data/oppty-discovery-backend/logs:/app/logs:rw \
  -v /data/oppty-discovery-backend/.env:/app/.env:ro \
  --restart=unless-stopped \
  oppty-discovery-backend:latest
```

### 5.2 前端静态资源发布与更新
更新前端界面时，无需重启容器，直接将打包编译产物解压覆盖至宿主机目录：
```bash
# 覆盖部署前端代码
tar -xzf /tmp/opp_dist.tar.gz -C /data/nginx/html/opp/

# 赋予正确的只读与执行权限
chmod -R 755 /data/nginx/html/opp/
```

### 5.3 Nginx 配置热重载
```bash
# 检测 Nginx 配置语法
docker exec common-nginx nginx -t

# 平滑热重载配置
docker exec common-nginx nginx -s reload
```

---

## 6. 连通性自检与监控探针

运维人员可在服务器或内网终端执行以下命令进行自检：

```bash
# 1. 验证前端 HTML 页面可达性 (预期返回 200 OK)
curl -I http://192.168.8.206/opp/

# 2. 验证后端服务探活接口 (预期返回健康 JSON)
curl http://192.168.8.206/opp/api/health
```

**健康响应示例**：
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "service": "商机发现助手",
    "status": "healthy",
    "database": "connected",
    "minio": "connected",
    "timestamp": 1791889523610
  },
  "timestamp": 1791889523610
}
```

---

## 7. 应急与排错指南

1. **页面提示网络异常或 502 Bad Gateway**：
   - 检查后端容器状态：`docker ps -f name=oppty-discovery-backend`；
   - 检查容器日志：`docker logs --tail 50 oppty-discovery-backend`，确认是否有数据库或外部凭证连接报错。
2. **验证码无法接收**：
   - 检查 Supabase 容器集群状态：`docker ps | grep supabase`；
   - 本地紧急测试支持万能码：`777777` 可直接放行登录。
3. **TikHub 抓取失败**：
   - 检查 `/data/oppty-discovery-backend/.env` 中的 `TIKHUB_API_TOKEN` 是否有效或额度耗尽。
