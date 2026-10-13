# Supabase 开源 BaaS 平台运维手册

> **软件名称**：Supabase（开源 Firebase 替代品 & 统一 BaaS 平台）  
> **部署方式**：Docker / Docker Compose (`--restart=unless-stopped`)  
> **目标服务器**：`192.168.8.206`（SSH 端口：22）  
> **服务器部署目录**：`/data/supabase/`  
> **本地工作区目录**：`d:\ai-ops\component-ops-206\Supabase\`  
> **最后更新日期**：2026-09-03  

---

## 1. 核心服务访问信息

| 服务模块 | 访问端点 / 端口 | 协议 / 鉴权 | 说明 |
| :--- | :--- | :--- | :--- |
| **Supabase Studio 控制台** | **`http://192.168.8.206:8000/`** | Web (HTTP Basic Auth) | Supabase 可视化管理控制台<br>账号：`supabase`<br>密码：`f1189f0b2d60e4e1158d28ea565c351c` |
| **PostgreSQL 数据库直连** | **`192.168.8.206:15433`** | PostgreSQL 17 | Supabase 原生数据库直连端口<br>账号：`postgres`<br>密码：`6128de580367207db74b707749b8fb04` |
| **Supavisor 事务连接池** | **`192.168.8.206:6543`** | Transaction Pooler | 高并发短连接池（适合 Serverless / Edge Functions） |
| **Supavisor 会话连接池** | **`192.168.8.206:15432`** | Session Pooler | 长连接会话池（已避开 5432 端口让给 `AnyDocs-postgres` 直连） |
| **Kong API Gateway** | **`http://192.168.8.206:8000/`** | REST / Realtime | REST API、Auth、Storage、Realtime 统一接入入口 |

---

## 2. 数据库与技术组件复用

1. **PostgreSQL 17 定制实例** (`supabase-db`)：
   - 包含 `pgvector`、`pg_graphql`、`pg_cron`、`pgjwt`、`vault` 等专有扩展；
   - 原生直连端口为 **`15433`**，包含独立业务库（如 `letta`）。
2. **MinIO 对象存储复用** (`supabase-storage`)：
   - 自动对接本机 `minio-server`（`192.168.8.206:9000`）；
   - 使用隔离文件桶 `supabase-storage` 进行统一存储。

---

## 3. 常用运维指令

```bash
# 进入部署目录
cd /data/supabase

# 查看 11 个容器运行状态
docker compose ps

# 重启 Supabase 全套服务
docker compose restart

# 查看 API Gateway 网关日志
docker compose logs -f --tail=100 kong

# 验证 Supavisor 会话连接池 (15432 端口)
PGPASSWORD=6128de580367207db74b707749b8fb04 psql -h 127.0.0.1 -p 15432 -U postgres.supabase-db -d postgres
```
