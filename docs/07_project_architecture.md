# 商机发现助手 (Oppty Discovery Assistant) 工程架构与目录结构手册

> **文档版本**：v1.0.3  
> **更新日期**：2026-10-13  
> **归档位置**：`docs/07_project_architecture.md`  
> **工程标识**：`oppty-discovery`

### 版本变更记录 (Revision History)

| 版本号 | 变更日期 | 变更类型（三级决策树） | 变更摘要 | 责任人 |
| :--- | :--- | :--- | :--- | :--- |
| **v1.0.0** | 2026-10-13 | 初始基线 | 初始工程架构手册创建，梳理前后端工程结构与外部组件集成机制 | reynolds |
| **v1.0.1** | 2026-10-13 | 功能级变更 (A1) & 模块内扩展 (A2) | 1. 扩充后端架构：新增 `internal/service/tikhub_client.go` 及单测，正式接入 TikHub RESTful 采集客户端与降级管道；<br>2. 扩充前端架构：补充本地 PostCSS 插件闭环隔离说明（规避父级 Monorepo 配置污染）。 | reynolds |
| **v1.0.2** | 2026-10-13 | 功能级变更 (A1) | 1. 扩充前端架构：新增 `src/components/about/AboutView.tsx` 独立产品介绍组件；<br>2. 首页空间重构：移除中央 Slogan，工作台表单长说明精简为 Tooltip 气泡；<br>3. 前后端加载与报告文案去 AI 化、去技术化。 | reynolds |
| **v1.0.3** | 2026-10-13 | 功能级变更 (A1) | 1. 前端导航架构优化：顶栏中央回归为纯双业务 Tab（商机发现与商机历史）；<br>2. 产品介绍解耦为独立外部链接并采用新标签页全屏展示架构。 | reynolds |

---

## 一、 应用概览与设计定位

### 1.1 品牌与产品矩阵
- **应用中文名**：商机发现助手
- **应用英文名**：Oppty Discovery Assistant
- **核心Slogan**：
  - 中文：痛点一键透视，商机即刻成型。
  - 英文：*From pain point to profit in seconds.*
- **产品定义**：一款集“社区评论挖掘、痛点线索提炼与商业方案构思”于一体的智能商机雷达。
- **品牌合规**：已彻底清洗原工程中所有第三方品牌标识（如 SeekMoney、SeekMoney-ai 等），统一替换为商机发现助手。

### 1.2 整体技术架构图
本项目采用**前后端完全分离架构**，前端自适应 PC 与移动端双端视口，后端提供异步高并发 RESTful 服务：

```mermaid
graph TD
    Client[React 18 + TS 前端 (PC / 移动端自适应)] -->|HTTP / JSON (RESTful)| API_GW[Go Gin 后端 API 网关 :8080]

    subgraph 后端核心服务层 (oppty-discovery-backend)
        API_GW --> AuthHandler[Auth 认证控制层]
        API_GW --> OpptyHandler[商机挖掘控制层]
        API_GW --> HistoryHandler[历史中心控制层]
        API_GW --> HealthHandler[健康监控控制层]

        AuthHandler --> AuthService[AuthService (登录即注册/JWT)]
        OpptyHandler --> OpptyService[OpptyService (挖掘流水线/成本预估)]
        OpptyService --> PromptMgr[Prompts 集中提示词模块]
        OpptyService --> MinioService[MinioService (S3报告归档)]
    end

    subgraph 底层技术组件 (192.168.8.206)
        AuthService -->|OTP 邮件验证| Supabase[Supabase BaaS :8000]
        OpptyService -->|GORM 持久化| PG[(PostgreSQL 17 :5432 / oppty_discovery)]
        MinioService -->|Markdown/HTML| MinIO[(MinIO OSS :9000 / oppty-discovery)]
        OpptyService -->|全网视频与评论| TikHub[TikHub 爬虫 API]
        OpptyService -->|深度商业洞察| QwenLLM[通义千问 LLM :3002 / qwen3.6-flash]
    end
```

---

## 二、 后端工程结构：`oppty-discovery-backend`

### 2.1 整体目录树
后端严格遵循 **Standard Go Project Layout**，分层清晰，杜绝“上帝对象”和跨层泄露：

```text
oppty-discovery-backend/
├── cmd/
│   └── server/
│       └── main.go                     # 服务启动入口（路由注入、中间件配置、CORS、服务监听）
├── configs/
│   └── config.yaml                     # 全局集中配置文件（包含所有外部组件IP、端口与凭据，带详尽中文注释）
├── internal/
│   ├── config/
│   │   └── config.go                   # 配置加载与结构体映射（支持动态加载与环境覆盖）
│   ├── constant/
│   │   └── constants.go                # 全局常量字典与枚举定义（数据源枚举、任务状态、系统品牌常量）
│   ├── prompts/
│   │   └── prompts.go                  # 商业分析与产品挖掘大模型提示词（统一抽取、版本化管理）
│   ├── model/                          # GORM 数据库持久化实体模型
│   │   ├── user.go                     # 用户表 (t_user) 实体
│   │   ├── oppty_job.go                # 商机任务表 (t_oppty_job) 实体
│   │   ├── oppty_report.go             # 商机报告明细表 (t_oppty_report) 实体
│   │   └── job_raw_data.go             # 抓取原声数据快照表 (t_job_raw_data) 实体
│   ├── types/                          # 前后端交互 Request DTO / Response VO（与底层 Model 严格隔离）
│   │   ├── auth.go                     # 认证接口请求入参与用户信息 VO
│   │   └── oppty.go                    # 挖掘请求入参、痛点四维解构、MVP方案、历史列表 VO
│   ├── repository/                     # 数据访问层（DAO / Repository）
│   │   ├── database.go                 # PostgreSQL 17 连接池初始化与自动迁移建表
│   │   ├── user_repository.go          # 用户持久化操作接口与实现
│   │   └── oppty_repository.go         # 任务/报告/历史记录持久化操作接口与实现
│   ├── service/                        # 核心业务逻辑层（领域服务）
│   │   ├── auth_service.go             # Supabase OTP 接入、无感注册自动建档、JWT 签发与解析
│   │   ├── minio_service.go            # MinIO S3 客户端初始化、自动建桶与报告导出上传
│   │   ├── oppty_service.go            # 成本预估计算器、异步协程挖掘管线、评论聚类、大模型调用
│   │   ├── oppty_service_test.go       # 核心业务单元测试（全中文用例，测试成本计算与提示词生成）
│   │   ├── tikhub_client.go            # TikHub RESTful 采集客户端（抖音综合搜索与视频评论真实调用）
│   │   └── tikhub_client_test.go       # TikHub 客户端单元测试（验证真实网络连通与 API Key 鉴权有效性）
│   ├── handler/                        # HTTP 控制层（入参校验、调用 Service、统一出参包装）
│   │   ├── auth_handler.go             # 发送验证码、校验并无感登录、退出登录
│   │   ├── oppty_handler.go            # 提交分析、查询任务状态与报告详情、实时成本测算
│   │   ├── history_handler.go          # 个人商机历史分页检索、删除历史
│   │   └── health_handler.go           # 服务健康状态与数据库/MinIO 连通性自检
│   └── middleware/
│       └── auth_middleware.go          # JWT 会话鉴权与当前用户会话上下文注入中间件
├── pkg/                                # 公共跨模块工具包
│   ├── response/
│   │   └── response.go                 # 统一 API 响应格式包装 (Result: code, message, data, timestamp)
│   └── errors/
│       └── errors.go                   # 全局统一业务错误码定义与中文错误映射
├── go.mod                              # Go 依赖描述文件 (Go 1.23+)
└── go.sum                              # 依赖校验锁文件
```

### 2.2 核心代码设计与规约
1. **统一配置中心 (`configs/config.yaml`)**：
   - 包含 `server`、`database`、`supabase`、`minio`、`llm`、`tikhub` 6 大配置节。
   - 所有 IP、端口、账号密码统一提取，禁止在业务代码中散落硬编码，便于未来一键迁移机房。
2. **商业提示词独立维护 (`internal/prompts/prompts.go`)**：
   - 所有 LLM 提示词集中收敛，规范了包含“表面痛点、3个为什么根本原因、爆发场景、情绪烈度、现有竞品缺陷、未满足需求、差异化切入点、MVP核心功能、验证假设、种子渠道、周期成本”的结构化约束。
3. **成本预估算法**：
   - 公式：$\text{成本} = \frac{\text{视频数}}{20} \times 0.01 + (\text{开启评论} ? \text{视频数} \times \frac{\text{评论数}}{20} \times 0.01 : 0)$。
4. **全中文规约**：所有代码注释、日志打印、错误消息及单元测试全部使用规范中文。

---

## 三、 前端工程结构：`oppty-discovery-frontend`

### 3.1 整体目录树
前端基于 **React 18 + TypeScript + Vite**，深度融入 **Hallmark** 与 **UI Skills** 现代高质感设计规范：

```text
oppty-discovery-frontend/
├── public/                             # 静态公共资源
├── src/
│   ├── api/
│   │   └── client.ts                   # 统一 HTTP 网络服务客户端（带 Authorization 拦截与类型推导）
│   ├── assets/                         # 图标与多媒体资源
│   ├── components/                     # 高内聚通用业务展示与交互组件
│   │   ├── about/
│   │   │   └── AboutView.tsx           # 独立产品介绍视图（系统呈现Slogan、三大核心能力、三步使用指引与安全背书）
│   │   ├── auth/
│   │   │   └── LoginModal.tsx          # 登录即注册无感弹窗（60s倒计时、自测码支持、首次建档提示）
│   │   ├── layout/
│   │   │   └── Navbar.tsx              # 顶部导航栏（品牌Logo、Slogan、自适应抽屉、登录态与退出）
│   │   ├── workspace/
│   │   │   ├── WorkspaceForm.tsx       # 商机挖掘配置表单（数据源切换、关键词标签、滑块联动、成本计算器）
│   │   │   └── ProgressCard.tsx        # 异步挖掘进度感知卡片（动态状态机指示、呼吸光晕、阶段详情）
│   │   ├── report/
│   │   │   ├── ReportView.tsx          # 商机深度报告矩阵卡片、付费意愿高亮、原声评论弹窗
│   │   │   └── FullscreenReportModal.tsx # 沉浸式全屏研读大屏模式（四维解构、竞品局限、MVP路径、ESC退出）
│   │   └── history/
│   │       └── HistoryView.tsx         # 商机历史中心（分页加载、关键词过滤、一键复现报告、历史软删除）
│   ├── context/
│   │   └── AuthContext.tsx             # 全局认证状态上下文（提供登录、注销、会话持久化与弹窗控制器）
│   ├── types/
│   │   └── index.ts                    # TypeScript 类型定义（与后端 DTO/VO 100% 镜像映射）
│   ├── App.tsx                         # 顶层主应用路由与 Tab 容器（整合发现工作台与历史中心）
│   ├── index.css                       # 现代设计系统 CSS（氛围光晕、玻璃拟态、微动效、自适应断点）
│   └── main.tsx                        # React 应用入口引导挂载
├── index.html                          # HTML 单页入口（品牌 Title、高质感 Web 字体、Meta 标签）
├── package.json                        # 前端依赖配置与构建脚本
├── tsconfig.json                       # TypeScript 编译器配置
├── tsconfig.node.json                  # Vite 节点构建 TS 配置
├── vite.config.ts                      # Vite 打包配置（端口 5173、反向代理 `/api` -> `localhost:8080`）
└── postcss.config.js                   # 本地 PostCSS 配置隔离
```

### 3.2 前端核心交互流程
1. **无感登录注册**：用户输入邮箱即可获取验证码，验证成功即完成登录，本地安全存储 Token；开发阶段内置测试通道（验证码 `123456`）。
2. **商机发现工作台**：
   - 切换 TikTok / 抖音短视频；
   - 关键词多标签添加（回车快速生成胶囊标签）；
   - 开启[爬取评论]开关后展开视频数（5~20）与每视频评论数（10~100）滑块；
   - 实时计算并动态渲染预估成本（微动效高亮）；
   - 点击“开始智能分析”进入后台异步流水线。
3. **状态机进度感知**：平滑展示“数据采集 $\rightarrow$ 语义聚类 $\rightarrow$ 商业透视”三阶段与百分比。
4. **全屏沉浸研读模式**：
   - 点击报告卡片“全屏研读”或右上角全屏图标进入；
   - 支持键盘 `ESC` 随时退出；
   - 排版涵盖：四维痛点剖析、现有竞品槽点、未满足需求、差异化切入点、MVP 路线与真实原声佐证。
5. **历史中心与快速回溯**：云端分页拉取历史报告，点击任意历史项即刻复现完整商机透视视图。

### 3.3 前端工程化编译与样式隔离保障（Tailwind CSS 深度集成）
在多工程复合工作区环境下，为防止 Vite 在查找 PostCSS 配置时向上递归继承根目录其它框架的 `postcss.config.mjs` 造成 Tailwind CSS 指令（`@tailwind base; @tailwind components; @tailwind utilities;`）未被解析、样式丢失的问题，工程采取了以下强制隔离机制：
1. **本地显式配置隔离**：在 `oppty-discovery-frontend/vite.config.ts` 的 `css.postcss` 配置项中，显式声明加载 `tailwindcss` 与 `autoprefixer` 插件；
2. **截断向上递归寻址**：在前端工程根目录建立本地 `postcss.config.js`，严格锁定 PostCSS 仅在当前子工程作用域内闭环运行；
3. **构建验证指标**：确保 `npm run build` 生成的 CSS 文件体积达 30KB+（包含完整的现代 UI 工具类、Hallmark 质感类、网格流式响应式类）。

---

## 四、 第三方技术组件拓扑与连接凭据清单

| 技术组件 | 对应服务 / 容器 | 访问端点 (宿主机 IP) | 凭证 / 鉴权 | 专属隔离命名 / 说明 |
| :--- | :--- | :--- | :--- | :--- |
| **关系型数据库** | `AnyDocs-postgres` | `192.168.8.206:5432` | `postgres` / `postgres123!@#` | 专属独立库：`oppty_discovery` |
| **BaaS 认证平台** | `supabase-kong` | `http://192.168.8.206:8000` | ServiceKey: `f1189f0...`<br>JWTSecret: `6128de...` | 提供无感邮箱 OTP 验证码发送与校验 |
| **S3 对象存储** | `minio-server` | `192.168.8.206:9000` (API)<br>`192.168.8.206:9001` (控制台) | `admin` / `Admin@206` | 专属独立文件桶：`oppty-discovery`（存放 Markdown 报告） |
| **大模型推理引擎**| `Qwen / GLM 服务`| `http://192.168.8.206:3002/v1` | `sk-b4R1jMvuDD3ySozNK28...` | 模型名称：`qwen3.6-flash` |
| **社交数据采集** | TikHub API 网关 | `https://api.tikhub.io` | API Token (已配置真实凭证) | 用于全网采集 TikTok 与抖音短视频及评论 |

---

## 五、 本地运行与构建指令速查

### 5.1 后端服务运行
```bash
# 1. 进入后端目录
cd oppty-discovery-backend

# 2. 运行自动化单元测试
go test -v ./internal/service

# 3. 编译验证
go build ./cmd/server

# 4. 启动服务 (默认监听 8080 端口)
go run ./cmd/server
```

### 5.2 前端工程运行
```bash
# 1. 进入前端目录
cd oppty-discovery-frontend

# 2. 启动开发热更新服务器 (监听 5173 端口)
npm run dev

# 3. 生产打包构建验证
npm run build
```

---

## 六、 维护与二次开发指南

1. **新增数据源支持（如小红书、Bilibili）**：
   - 在 `internal/constant/constants.go` 的 `DataSourceType` 枚举中追加类型；
   - 在 `internal/service/oppty_service.go` 的 `fetchSocialData` 中增加对应数据源的抓取适配器；
   - 在前端 `src/types/index.ts` 与 `src/components/workspace/WorkspaceForm.tsx` 中增加对应平台图标与单选按钮。
2. **微调商机洞察提示词**：
   - 直接修改 [`internal/prompts/prompts.go`](file:///d:/ai-github/seekmoney-ai/oppty-discovery-backend/internal/prompts/prompts.go) 中的 `BuildClusterAnalysisPrompt` 函数即可，业务代码与提示词已完全解耦。
3. **调整成本计算比例**：
   - 后端修改 `internal/service/oppty_service.go` 的 `CalculateEstimatedCost`；
   - 前端同步更新 `src/components/workspace/WorkspaceForm.tsx` 的计算逻辑，单测在 `oppty_service_test.go` 中同步校验。
