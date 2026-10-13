# 功能与接口映射关系文档 (Feature-API Mapping)

> **文档版本**：v1.0.3  
> **关联文档**：`docs/02_features.md`, `docs/03_data_model.md`, `docs/04_api_spec.md`  
> **状态**：已基线化 / 增量变更同步完成  
> **工程标识**：`oppty-discovery`

### 版本变更记录 (Revision History)

| 版本号 | 变更日期 | 变更类型（三级决策树） | 变更摘要 | 责任人 |
| :--- | :--- | :--- | :--- | :--- |
| **v1.0.0** | 2026-10-13 | 初始基线 | 初始四向可追溯性映射矩阵建立，覆盖全部 24 个功能项 | reynolds |
| **v1.0.1** | 2026-10-13 | 功能级变更 (A1) & 模块内扩展 (A2) | 1. `FEAT-OPPTY-07`：细化 TikHubClient 外部 RESTful 调用端点与高可用容错机制；<br>2. `FEAT-SYS-04`：记录前端 Vite 本地 PostCSS 隔离配置与 TailwindCSS 样式编译闭环映射。 | reynolds |
| **v1.0.2** | 2026-10-13 | 功能级变更 (A1) | 1. `FEAT-ABOUT-01`：新增独立产品介绍页面组件映射与导航联动；<br>2. 优化 `FEAT-OPPTY-04` / `FEAT-OPPTY-06`：去技术化与 Tooltip 问号组件映射。 | reynolds |
| **v1.0.3** | 2026-10-13 | 功能级变更 (A1) | 1. `FEAT-ABOUT-01`：更新为顶部右侧外部超链接 (`target="_blank" href="?view=about"`)，在新标签页全屏独立渲染，主导航恢复纯双业务 Tab。 | reynolds |

---

## 1. 四向可追溯性映射矩阵 (Traceability Matrix)

| 功能 ID | 功能名称 | 前端组件 / 页面 | 后端 API 路径与 Method | 涉及数据表 / 实体 | 状态与覆盖校验 |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **FEAT-ABOUT-01** | 独立产品介绍页面 | `src/components/about/AboutView.tsx`<br>`src/components/layout/Navbar.tsx` (外链) | 前端 URL 参数 `?view=about` 路由切换 | 静态全屏展示与导流 | 已完全覆盖 (100%) |
| **FEAT-AUTH-01** | 邮箱验证码发送 (OTP) | `src/pages/login/LoginPage.tsx`<br>`src/components/auth/SendCodeButton.tsx` | `POST /api/v1/auth/send-code` | Supabase Auth API / 邮件代理 | 已完全覆盖 (100%) |
| **FEAT-AUTH-02** | 验证码校验与无感注册 | `src/pages/login/LoginPage.tsx`<br>`src/components/auth/LoginForm.tsx` | `POST /api/v1/auth/verify-code` | `t_user` (`model.User`) | 已完全覆盖 (100%) |
| **FEAT-AUTH-03** | 会话保持与注销 | `src/components/layout/Navbar.tsx`<br>`src/context/AuthContext.tsx` | `POST /api/v1/auth/logout` | 本地存储 + Supabase Session | 已完全覆盖 (100%) |
| **FEAT-AUTH-04** | 用户个人信息查询 | `src/components/layout/Navbar.tsx` | 附随于 Verify Code 或 Auth Token 校验 | `t_user` (`model.User`) | 已完全覆盖 (100%) |
| **FEAT-OPPTY-01** | 数据源切换 | `src/components/workspace/DataSourceSelect.tsx` | 随分析请求提交 `POST /api/v1/analyze` | `t_oppty_job` | 已完全覆盖 (100%) |
| **FEAT-OPPTY-02** | 关键词标签化输入与校验 | `src/components/workspace/KeywordInput.tsx` | 随分析请求提交 `POST /api/v1/analyze` | `t_oppty_job` | 已完全覆盖 (100%) |
| **FEAT-OPPTY-03** | 深度抓取配置（视频与评论滑块）| `src/components/workspace/CrawlConfigSlider.tsx` | 随分析请求提交 `POST /api/v1/analyze` | `t_oppty_job` | 已完全覆盖 (100%) |
| **FEAT-OPPTY-04** | 实时预估成本动态计算展示 | `src/components/workspace/CostEstimator.tsx` | 前端实时公式计算 + 后端入参校验 | 内存即时计算，持久化于 `t_oppty_job.estimated_cost` | 已完全覆盖 (100%) |
| **FEAT-OPPTY-05** | 异步商机挖掘任务创建 | `src/pages/workspace/WorkspacePage.tsx`<br>`src/components/workspace/SubmitButton.tsx` | `POST /api/v1/analyze` | `t_oppty_job` (`model.OpptyJob`) | 已完全覆盖 (100%) |
| **FEAT-OPPTY-06** | 异步任务进度轮询与状态感知 | `src/components/workspace/JobProgressCard.tsx` | `GET /api/v1/jobs/:jobId` | `t_oppty_job` | 已完全覆盖 (100%) |
| **FEAT-OPPTY-07** | 社交媒体数据源抓取适配 (TikHub) | 后台异步协程驱动 | `service.TikHubClient`<br>• `POST /api/v1/douyin/search/fetch_general_search_v1`<br>• `GET /api/v1/douyin/web/fetch_video_comments` | `t_job_raw_data` | 已完全覆盖 (100%) |
| **FEAT-OPPTY-08** | 评论语义密度聚类分析 | 后台异步协程驱动 | 内部 Service 聚类算法引擎 | 内存聚类，持久化于 `t_oppty_report` | 已完全覆盖 (100%) |
| **FEAT-OPPTY-09** | 大模型深度商业洞察提炼 | 后台异步协程驱动 | 内部 Service 驱动 Qwen/GLM LLM | `t_oppty_report` (`model.OpptyReport`) | 已完全覆盖 (100%) |
| **FEAT-REPO-01** | 痛点聚类矩阵与价值卡片 | `src/components/report/ReportMatrix.tsx` | `GET /api/v1/jobs/:jobId` | `t_oppty_report` | 已完全覆盖 (100%) |
| **FEAT-REPO-02** | 痛点深度四维剖析 | `src/components/report/PainDepthView.tsx` | `GET /api/v1/jobs/:jobId` | `t_oppty_report.pain_depth` | 已完全覆盖 (100%) |
| **FEAT-REPO-03** | 竞品与市场格局对比 | `src/components/report/MarketLandscapeView.tsx` | `GET /api/v1/jobs/:jobId` | `t_oppty_report.market_landscape` | 已完全覆盖 (100%) |
| **FEAT-REPO-04** | MVP 商业落地执行规划卡片 | `src/components/report/MVPPlanCard.tsx` | `GET /api/v1/jobs/:jobId` | `t_oppty_report.mvp_plan` | 已完全覆盖 (100%) |
| **FEAT-REPO-05** | 真实代表性原声评论弹窗 | `src/components/report/RawCommentsModal.tsx` | `GET /api/v1/jobs/:jobId` | `t_oppty_report.representative_texts` | 已完全覆盖 (100%) |
| **FEAT-REPO-06** | 商机报告沉浸式全屏查看模式 | `src/components/report/FullscreenReportModal.tsx` | 本地模态全屏状态切换 | `t_oppty_report` 前端全屏呈现 | 已完全覆盖 (100%) |
| **FEAT-REPO-07** | 报告 Markdown/HTML 归档至 MinIO | `src/components/report/ExportReportButton.tsx` | 后端自动导出并存储 MinIO | `minio-server:9000/oppty-discovery` | 已完全覆盖 (100%) |
| **FEAT-HIST-01** | 个人商机历史分页列表查询 | `src/pages/history/HistoryPage.tsx`<br>`src/components/history/HistoryTable.tsx` | `GET /api/v1/history` | `t_oppty_job`<br>`t_oppty_report` | 已完全覆盖 (100%) |
| **FEAT-HIST-02** | 历史商机报告详情即时载入 | `src/pages/history/HistoryPage.tsx` | `GET /api/v1/jobs/:jobId` | `t_oppty_report` | 已完全覆盖 (100%) |
| **FEAT-HIST-03** | 历史记录多维检索筛选 | `src/components/history/HistorySearchFilter.tsx`| `GET /api/v1/history?keyword=...` | `t_oppty_job` | 已完全覆盖 (100%) |
| **FEAT-HIST-04** | 废弃商机历史删除与归档 | `src/components/history/DeleteJobButton.tsx` | `DELETE /api/v1/history/:jobId` | `t_oppty_job` (软删除)<br>`t_oppty_report` | 已完全覆盖 (100%) |
| **FEAT-SYS-01** | 集中化配置管理（含中文注释） | - | 后端 `internal/config/config.go` + `config.yaml` | 配置文件管理 | 已完全覆盖 (100%) |
| **FEAT-SYS-02** | 商业提示词全量集中抽取 | - | 后端 `internal/service/prompts/prompts.go` | 集中管理模块 | 已完全覆盖 (100%) |
| **FEAT-SYS-03** | 第三方品牌全面清理 | 全局页面、页眉、页脚、Title 与组件 | 后端错误信息与响应文案 | - | 已完全覆盖 (100%) |
| **FEAT-SYS-04** | 双端自适应响应式布局与 Hallmark UI | Tailwind CSS v3 + `src/index.css` 全局设计系统 + 移动端折叠抽屉与直接登录 | 前端工程化编译管道 | - | 已完全覆盖 (100%) |
| **FEAT-SYS-05** | 服务健康检查与依赖监控 | - | `GET /api/health` | 数据库/MinIO 连通测试 | 已完全覆盖 (100%) |

---

## 2. 覆盖率与一致性自检清单

- [x] **功能点全覆盖**：所有在 `02_features.md` 中定义的 24 项 P0/P1 功能均已 100% 映射对应 API 与前端组件。
- [x] **无孤立 API**：所有在 `04_api_spec.md` 中定义的接口均具备明确的业务功能归属和对应前端交互组件。
- [x] **数据模型闭环**：所有接口涉及的读写字段均已在 `03_data_model.md` PostgreSQL 表结构及 GORM Model 中清晰定义。
- [x] **架构约束达标**：后端彻底分层（Handler / Service / Repository / Model），前端组件与 API 服务解耦，所有错误码与响应结构规范统一。
