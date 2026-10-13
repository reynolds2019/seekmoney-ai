# Nginx (Alpine) 容器化部署与运维手册

本文档记录了服务器 `192.168.8.206` 上基于 Docker 部署的最新版 Nginx 服务关键配置、拓扑架构与日常运维操作指南。

---

## 1. 基本信息与部署架构

- **软件名称**：Nginx
- **版本镜像**：`nginx:alpine`
- **容器名称**：`common-nginx`
- **重启策略**：`--restart=unless-stopped`（开机自启）
- **网络模式与端口映射**：
  - HTTP 访问：`0.0.0.0:80 -> 80/tcp`
  - HTTPS 访问：`0.0.0.0:443 -> 443/tcp`
- **访问测试地址**：`http://192.168.8.206/`

---

## 2. 宿主机持久化目录结构

所有核心数据、配置与日志均持久化挂载至数据盘 `/data/nginx` 目录下：

```
/data/nginx/
├── conf/
│   └── nginx.conf            # 【主配置文件】(只读挂载至 /etc/nginx/nginx.conf)
├── conf.d/
│   └── default.conf          # 【虚拟主机与子配置目录】(只读挂载至 /etc/nginx/conf.d/)
├── html/
│   └── index.html            # 【静态网页根目录】(读写挂载至 /usr/share/nginx/html/)
├── logs/
│   ├── access.log            # 【访问日志】(读写挂载至 /var/log/nginx/)
│   └── error.log             # 【错误日志】(读写挂载至 /var/log/nginx/)
└── nginx_ops_manual.md       # 【本运维手册】
```

---

## 3. 部署启动命令参考

```bash
docker run -d \
  --name common-nginx \
  -p 80:80 \
  -p 443:443 \
  -v /data/nginx/conf/nginx.conf:/etc/nginx/nginx.conf:ro \
  -v /data/nginx/conf.d:/etc/nginx/conf.d:ro \
  -v /data/nginx/html:/usr/share/nginx/html:rw \
  -v /data/nginx/logs:/var/log/nginx:rw \
  --restart=unless-stopped \
  nginx:alpine
```

---

## 4. 日常运维与管理命令

### 4.1 容器生命周期管理
```bash
# 查看容器运行状态
docker ps -f name=common-nginx

# 启动容器
docker start common-nginx

# 停止容器
docker stop common-nginx

# 重启容器
docker restart common-nginx
```

### 4.2 配置热重载与语法检测
修改 `/data/nginx/conf.d/*.conf` 或 `/data/nginx/conf/nginx.conf` 后，建议在不重启容器的情况下进行重载：
```bash
# 1. 检测配置语法是否正确
docker exec common-nginx nginx -t

# 2. 语法无误后平滑重载配置
docker exec common-nginx nginx -s reload
```

### 4.3 日志查看与实时追踪
```bash
# 查看实时访问日志（宿主机直接查看）
tail -f /data/nginx/logs/access.log

# 查看实时错误日志
tail -f /data/nginx/logs/error.log

# 查看容器标准输出日志
docker logs -f --tail 100 common-nginx
```

### 4.4 静态资源部署与更新
将前端打包后的静态资源（如 `dist/` 目录内容）直接放置到宿主机 `/data/nginx/html/` 目录下即可即时生效，无需重启 Nginx 容器。

---

## 5. 常见问题排查与注意事项

1. **403 Forbidden 报错**：
   - 检查 `/data/nginx/html/` 目录下是否存在 `index.html`，或检查文件权限是否对容器内 `nginx` 用户可读。
2. **反向代理配置注意**：
   - 若在 `conf.d/` 中配置反向代理至同服务器的其他 Docker 容器，建议通过宿主机 IP（如 `http://192.168.8.206:端口`）或统一加入同一 Docker 网络通过容器名互联。
