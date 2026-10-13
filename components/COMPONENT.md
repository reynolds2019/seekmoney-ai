# 服务器公共技术组件复用清单

本文档汇总了服务器 `192.168.8.206` 上可供后续业务应用直接复用的核心技术组件与连接配置信息。

---

## 1. Web 与反向代理：`common-nginx`

- **组件定位**：统一 Web 网关、静态页面托管与反向代理服务
- **容器名称**：`common-nginx`（镜像：`nginx:alpine`）
- **监听端口**：`80` (HTTP), `443` (HTTPS)
- **宿主机持久化路径**：
  - 子配置文件目录：`/data/nginx/conf.d/`（新增应用代理配置置于此处）
  - 主配置文件：`/data/nginx/conf/nginx.conf`
  - 静态页面根目录：`/data/nginx/html/`
  - 日志目录：`/data/nginx/logs/`
- **配置生效命令**：
  ```bash
  docker exec common-nginx nginx -t && docker exec common-nginx nginx -s reload
  ```
- **复用指引**：后续新开发应用只需在 `/data/nginx/conf.d/` 下新建 `xxx.conf` 添加 `location /app-name/ { proxy_pass ...; }`，并执行重载即可对外提供服务。

---

## 2. 关系型与向量数据库：`AnyDocs-postgres`

- **组件定位**：PostgreSQL 17 数据库（集成了 ParadeDB 向量/全文检索扩展与 pgvector 插件）
- **容器名称**：`AnyDocs-postgres`（镜像：`paradedb/paradedb:v0.22.2-pg17`）
- **对外服务端口**：**`0.0.0.0:5432 -> 5432`**（已暴露宿主机 5432 端口，供远程客户端直连）
- **连接凭证**：
  - **服务器 IP**：`192.168.8.206`
  - **端口 (Port)**：**`5432`**（标准 PostgreSQL 端口）
  - **超级用户**：`postgres`
  - **访问密码**：`postgres123!@#`
  - **默认数据库**：`WeKnora`
- **连接字符串**：
  ```
  postgresql://postgres:postgres123!@#@192.168.8.206:5432/WeKnora
  ```
- **复用指引**：
  - **容器间复用**：新应用容器在 `docker-compose.yml` 中加入 `networks: - AnyDocs-network`，即可通过主机名 `AnyDocs-postgres:5432` 连通，可创建独立 Database/Schema 隔离业务数据。
  - **远程 GUI 客户端直连 (DBeaver / Navicat / DataGrip / pgAdmin)**：直接填写 Host: `192.168.8.206`，Port: `5432`，User: `postgres`，Password: `postgres123!@#` 即可免 SSH 隧道直接连接。

---

## 3. 分布式缓存：`AnyDocs-redis`

- **组件定位**：Redis 7.0 分布式高速缓存、会话存储与消息队列
- **容器名称**：`AnyDocs-redis`（镜像：`redis:7.0-alpine`）
- **容器网络**：`AnyDocs-network`（内部 IP：`172.24.0.2`，端口：`6379`）
- **连接凭证**：
  - **主机地址**：`AnyDocs-redis`（同一 Docker 网络内）或 `172.24.0.2`
  - **端口**：`6379`
  - **认证密码**：`redis123!@#`
- **连接字符串**：
  ```
  redis://:redis123!@#@AnyDocs-redis:6379/0
  ```
- **复用指引**：新应用接入 `AnyDocs-network` 后，连接 `AnyDocs-redis:6379`，建议各应用使用独立 `db` 索引（如 1~15）或统一增加 Key 前缀进行命名空间隔离。

---

## 4. 对象存储 OSS：`minio-server`

- **组件定位**：S3 兼容的高性能对象存储服务
- **容器名称**：`minio-server`（镜像：`minio/minio`）
- **访问入口**：
  - **S3 API 端点**：`http://192.168.8.206:9000`（容器内网络：`http://minio-server:9000`）
  - **Web Console 控制台**：`http://192.168.8.206:9001`
- **管理凭证**：
  - **Access Key (用户名)**：`admin`
  - **Secret Key (密码)**：`Admin@206`
- **宿主机存储目录**：`/data/minio`
- **复用指引**：新应用可通过 Web 控制台（`:9001`）创建专属 Bucket（如 `app-files`）及专用 AccessKey，使用标准 S3 SDK（AWS S3 / MinIO SDK）进行文件与多媒体资源存储。
