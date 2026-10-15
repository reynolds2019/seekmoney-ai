# 服务接口清单 (API Specification)

> **文档版本**：v1.0.7  
> **协议规范**：HTTP/1.1 RESTful (JSON)  
> **状态**：已基线化 / 增量变更同步完成  
> **工程标识**：`oppty-discovery`

### 版本变更记录 (Revision History)

| 版本号 | 变更日期 | 变更类型（三级决策树） | 变更摘要 | 责任人 |
| :--- | :--- | :--- | :--- | :--- |
| **v1.0.0** | 2026-10-13 | 初始基线 | 初始 API 契约基线化，包含鉴权、分析、详情、历史及健康检查接口 | reynolds |
| **v1.0.1** | 2026-10-13 | 功能级变更 (A1) & 模块内扩展 (A2) | 1. 补充后台异步协程调用的外部 TikHub API 契约细节（搜索接口与评论接口协议定义）；<br>2. 补充调用超时、容错降级及状态机映射。 | reynolds |
| **v1.0.4** | 2026-10-13 | 契约级扩展 (A2) | 1. `POST /api/v1/analyze`：数据源枚举扩展为 `tiktok, douyin, xiaohongshu, bilibili, wechat, youtube`；<br>2. 新增免鉴权公开分享接口 `GET /api/v1/share/jobs/:jobId`；<br>3. 报告出参增加 `sourceVideos` 数组（含标题、作者、获赞、原文外链）。 | reynolds |
| **v1.0.7** | 2026-10-13 | 功能级变更 (A1) & 交互级优化 (A3) | 1. **登录 UI 交互优化**：修复登录弹窗邮箱/验证码左侧图标与输入文字重叠的样式缺陷；清空预填码并增加查收邮箱提醒；<br>2. **Supabase 真实邮件外发集成**：后端读取外部环境生产配置段接入 Supabase OTP，通过 `POST /auth/v1/otp` 投递真实邮件验证码，支持 429 频控防刷友好提示；<br>3. **三重认证登录闭环**：`POST /api/v1/auth/verify-code` 支持“测试万能码 -> 5分钟内存外发码校验核销 -> Supabase 在线验证”三重保障。 | reynolds |
| **v1.0.8** | 2026-10-15 | 功能级变更 (A1) | 1. **评论数据契约真实性强化**：明确 `representativeTexts` 必须来自真实外部 TikHub API 数据源，严禁注入 mock 模版评论；<br>2. **TikHub 抓取契约去降级**：统一通过 `GET /api/v1/douyin/web/fetch_video_comments` 抓取短视频真实社群原声评论，移除假数据兜底分支。 | reynolds |

---

## 1. 统一响应规范 (Global Result)

所有 API 出参均严格包裹为标准统一结构：

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {},
  "timestamp": 1755792000000
}
```

Go 统一封装结构体 (`pkg/response/response.go`)：
```go
package response

// Result 统一API响应结构体
// desc: 统一前后端JSON交互协议
// date: 2026-10-13 11:20:00
// @author reynolds
// @since v1.0.0
type Result struct {
    Code      int         `json:"code"`
    Message   string      `json:"message"`
    Data      interface{} `json:"data"`
    Timestamp int64       `json:"timestamp"`
}
```

---

## 2. 全局错误码定义表 (Error Codes)

| 错误码 | 错误常量名称 | 中文描述 | 应对建议 |
| :--- | :--- | :--- | :--- |
| `200` | `CodeSuccess` | 操作成功 | 正常渲染业务数据 |
| `202` | `CodeAccepted` | 任务已接收并提交异步处理 | 前端启动定时轮询状态感知 |
| `4001` | `CodeParamInvalid` | 请求参数校验不通过 | 检查输入字段边界与格式 |
| `4003` | `CodeUnauthorized` | 登录凭证无效或已过期 | 引导退出并重定向至登录页 |
| `4004` | `CodeJobNotFound` | 未找到指定商机分析任务 | 刷新列表或检查任务标识有效性 |
| `4029` | `CodeRateLimit` | 验证码或分析任务提交过于频繁 | 提示用户稍后重试 |
| `5001` | `CodeCrawlFailed` | 社交媒体数据源抓取异常 | 检查网络连通性或数据源状态 |
| `5002` | `CodeLLMAnalyzeFailed`| 大模型商业分析服务异常 | 检查模型配置或配额状态 |
| `5003` | `CodeStorageError` | MinIO 对象存储或数据库访问异常| 检查底层服务健康状况 |

---

## 3. 详细接口清单

### 3.1 发送邮箱验证码 (OTP)
- **接口路径**：`POST /api/v1/auth/send-code`
- **功能描述**：输入邮箱地址，调用 Supabase Auth OTP 投递 6 位数字验证码
- **权限级别**：公开无鉴权
- **Go 请求结构体**：
```go
type SendCodeReq struct {
    Email string `json:"email" binding:"required,email"`
}
```
- **响应出参**：
```json
{
  "code": 200,
  "message": "验证码已成功发送至您的邮箱，请注意查收",
  "data": {
    "cooldownSeconds": 60
  },
  "timestamp": 1755792000000
}
```

---

### 3.2 邮箱验证码校验与「登录即注册」
- **接口路径**：`POST /api/v1/auth/verify-code`
- **功能描述**：校验邮箱与 6 位验证码，自动建档（首次）并签发会话 Token
- **权限级别**：公开无鉴权
- **Go 请求结构体**：
```go
type VerifyCodeReq struct {
    Email string `json:"email" binding:"required,email"`
    Code  string `json:"code" binding:"required,len=6"`
}
```
- **Go 响应出参 (AuthVO)**：
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "entrepreneur@example.com",
      "nickname": "商机探索者",
      "avatarUrl": "",
      "isNewUser": true
    }
  },
  "timestamp": 1755792000000
}
```

---

### 3.3 退出登录
- **接口路径**：`POST /api/v1/auth/logout`
- **功能描述**：注销当前会话
- **权限级别**：需登录鉴权 (`Authorization: Bearer <Token>`)
- **响应出参**：
```json
{
  "code": 200,
  "message": "已成功退出登录",
  "data": null,
  "timestamp": 1755792000000
}
```

---

### 3.4 提交商机发现分析任务
- **接口路径**：`POST /api/v1/analyze`
- **功能描述**：提交数据源、关键词与抓取参数，创建异步分析任务
- **权限级别**：需登录鉴权 (`Authorization: Bearer <Token>`)
- **Go 请求结构体**：
```go
type AnalyzeReq struct {
    Keywords            []string `json:"keywords" binding:"required,min=1,max=5,dive,min=1,max=50"`
    DataSource          string   `json:"dataSource" binding:"required,oneof=tiktok douyin"`
    DeepCrawl           bool     `json:"deepCrawl"`
    MaxVideos           int      `json:"maxVideos" binding:"required,min=5,max=20"`
    MaxCommentsPerVideo int      `json:"maxCommentsPerVideo" binding:"required,min=10,max=100"`
}
```
- **响应出参 (202 Accepted)**：
```json
{
  "code": 202,
  "message": "商机挖掘任务已创建，正在后台进行全网数据深度透视",
  "data": {
    "jobId": "f784d1bc-6ea9-42b4-82a1-638c35bb5f51",
    "status": "processing",
    "estimatedCost": 0.11
  },
  "timestamp": 1755792000000
}
```

---

### 3.5 查询商机分析任务状态与报告详情
- **接口路径**：`GET /api/v1/jobs/:jobId`
- **功能描述**：根据 `jobId` 轮询任务进度，或在完成后获取完整商机报告
- **权限级别**：需登录鉴权
- **路径参数**：`jobId` (UUID 字符串)
- **进行中响应示例 (processing)**：
```json
{
  "code": 200,
  "message": "任务执行中",
  "data": {
    "jobId": "f784d1bc-6ea9-42b4-82a1-638c35bb5f51",
    "status": "processing",
    "progressStage": "clustering",
    "progressPercent": 65,
    "progressMessage": "已完成数据抓取，正在进行评论语义密度聚类分析...",
    "keywords": ["AI变现", "自媒体运营"],
    "results": []
  },
  "timestamp": 1755792000000
}
```
- **完成响应示例 (completed)**：
```json
{
  "code": 200,
  "message": "分析完成",
  "data": {
    "jobId": "f784d1bc-6ea9-42b4-82a1-638c35bb5f51",
    "status": "completed",
    "progressStage": "completed",
    "progressPercent": 100,
    "progressMessage": "商机分析全量完成",
    "keywords": ["AI变现", "自媒体运营"],
    "dataSource": "tiktok",
    "deepCrawl": true,
    "maxVideos": 10,
    "maxCommentsPerVideo": 30,
    "estimatedCost": 0.16,
    "createdAt": "2026-10-13T11:30:00Z",
    "results": [
      {
        "id": "1",
        "clusterCode": "cluster_0",
        "signalSize": 42,
        "oneLinePain": "新手创业者缺乏将AI提示词沉淀为自动化可交付交付物的标准化工作流",
        "paidInterest": "High",
        "rationale": "超60%的评论集中吐槽市面教程假大空，愿意为开箱即用的工作流模板付费",
        "potentialProduct": "垂类行业爆款文案自动化仿写与图文一键排版小工具",
        "painDepth": {
          "surface_pain": "不知道用什么AI工具高效写出高互动文案",
          "root_causes": ["信息噪音大", "缺少商业闭环指引", "缺乏实际交付物模板"],
          "user_scenarios": ["副业创业", "自媒体图文代运营"],
          "emotional_intensity": 4
        },
        "marketLandscape": {
          "existing_solutions": [
            { "name": "通用大模型对话框", "limitation": "缺少针对垂直平台的提示词微调与输出约束" }
          ],
          "unmet_needs": ["开箱即用的工作流模板", "对标爆款拆解工具"],
          "opportunity": "打造垂直人群的微型 SaaS 工具"
        },
        "mvpPlan": {
          "core_features": ["爆款文案特征提取", "自动化批量图文生成", "一键多平台排版"],
          "validation_hypotheses": [
            { "hypothesis": "用户愿意为垂直工作流付费", "test_method": "微信小群预售" }
          ],
          "first_users": "自媒体交流社群群友与副业小白",
          "timeline": "2周上线 MVP",
          "estimated_cost": "¥500 部署与测试成本"
        },
        "keywordRelevance": 92,
        "priorityScore": 86,
        "representativeTexts": [
          "新手不知道怎么用AI写文案变现，求保姆级教程",
          "尝试了很多工具都赚不到钱，有没有真正能跑通的案例？"
        ],
        "minioMarkdownUrl": "http://192.168.8.206:9000/oppty-discovery/reports/f784d1bc.md",
        "minioHtmlUrl": "http://192.168.8.206:9000/oppty-discovery/reports/f784d1bc.html"
      }
    ]
  },
  "timestamp": 1755792000000
}
```

---

### 3.6 个人商机历史列表分页查询
- **接口路径**：`GET /api/v1/history`
- **功能描述**：按时间倒序获取当前用户的商机历史任务列表
- **权限级别**：需登录鉴权
- **Query 参数**：
  - `pageNo` (int, 默认 1)
  - `pageSize` (int, 默认 10, 最大 50)
  - `keyword` (string, 可选模糊匹配)
- **响应出参**：
```json
{
  "code": 200,
  "message": "获取商机历史成功",
  "data": {
    "total": 5,
    "pageNo": 1,
    "pageSize": 10,
    "records": [
      {
        "jobId": "f784d1bc-6ea9-42b4-82a1-638c35bb5f51",
        "keywords": ["AI变现", "自媒体运营"],
        "dataSource": "tiktok",
        "signalCount": 42,
        "reportCount": 2,
        "primaryPain": "新手创业者缺乏将AI提示词沉淀为自动化可交付交付物的标准化工作流",
        "status": "completed",
        "estimatedCost": 0.16,
        "createdAt": "2026-10-13T11:30:00Z"
      }
    ]
  },
  "timestamp": 1755792000000
}
```

---

### 3.7 删除商机历史记录
- **接口路径**：`DELETE /api/v1/history/:jobId`
- **功能描述**：删除并归档用户指定的商机历史记录
- **权限级别**：需登录鉴权（只能删除属于当前用户的记录）
- **响应出参**：
```json
{
  "code": 200,
  "message": "历史记录删除成功",
  "data": null,
  "timestamp": 1755792000000
}
```

---

### 3.8 服务健康检查
- **接口路径**：`GET /api/health`
- **功能描述**：检查应用运行状态、PostgreSQL 与 MinIO 连通性
- **权限级别**：公开无鉴权
- **响应出参**：
```json
{
  "code": 200,
  "message": "商机发现助手服务健康",
  "data": {
    "service": "oppty-discovery-backend",
    "status": "healthy",
    "database": "connected",
    "minio": "connected",
    "timestamp": 1755792000000
  },
  "timestamp": 1755792000000
}
```

---

## 4. 外部集成接口契约 (TikHub Social API Specification)

由后端异步协程调用，非直接暴露给前端。

### 4.1 抖音综合搜索接口
- **端点**：`POST https://api.tikhub.io/api/v1/douyin/search/fetch_general_search_v1`
- **请求头**：
  - `Authorization: Bearer <TIKHUB_API_KEY>`
  - `Content-Type: application/json`
- **请求体 (JSON)**：
```json
{
  "keyword": "AI副业",
  "offset": 0,
  "count": 10,
  "sort_type": "0",
  "publish_time": "0"
}
```
- **核心响应提取**：
  - 提取 `data.data[].aweme_info.aweme_id`（视频 ID）；
  - 提取 `data.data[].aweme_info.desc`（视频标题与文案）；
  - 提取 `data.data[].aweme_info.statistics.digg_count`（点赞数）。

### 4.2 抖音视频评论采集接口
- **端点**：`GET https://api.tikhub.io/api/v1/douyin/web/fetch_video_comments?aweme_id=<AWEME_ID>&cursor=0&count=20`
- **请求头**：
  - `Authorization: Bearer <TIKHUB_API_KEY>`
- **核心响应提取**：
  - 提取 `data.comments[].cid`（评论 ID）；
  - 提取 `data.comments[].text`（评论正文）；
  - 提取 `data.comments[].digg_count`（点赞热度）；
  - 提取 `data.comments[].user.nickname`（评论者昵称）。

### 4.3 高可用降级策略 (Fault Tolerance & Fallback)
1. **调用超时**：单次 HTTP 请求超时设定为 15 秒；
2. **熔断与降级**：若三方 API 触发限频（429）、未配置有效 Token 或网络无法连通，系统捕获异常并自动注入经过语义验证的典型高价值原声样本作为保底，确保整个商机挖掘流水线（语义聚类与 LLM 洞察提炼）始终 100% 成功交付。

---

## 5. 公开分享与视频溯源接口

### 5.1 公开免登录获取商机透视报告
- **端点**：`GET /api/v1/share/jobs/:jobId`
- **认证要求**：**免登录 / 无需 Authorization 请求头**（公开只读查阅）
- **路径参数**：
  - `jobId`：任务 UUID
- **响应示例 (JSON)**：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "jobId": "c1a59c4a-2081-47f3-b9bb-bca736e2c3b9",
    "status": "completed",
    "progressStage": "completed",
    "progressPercent": 100,
    "progressMessage": "商机深度透视报告已生成！",
    "keywords": ["AI变现", "自媒体运营"],
    "dataSource": "douyin",
    "deepCrawl": true,
    "maxVideos": 10,
    "maxCommentsPerVideo": 20,
    "estimatedCost": 0.11,
    "createdAt": "2026-10-13 14:00:00",
    "results": [
      {
        "id": 1,
        "clusterCode": "cluster_0",
        "signalSize": 12,
        "oneLinePain": "缺乏开箱即用的垂直类自动化生产SOP与集成工具",
        "paidInterest": "High",
        "rationale": "高频出现由于工具割裂、人工耗时长导致的低效吐槽...",
        "potentialProduct": "全流程自媒体爆款批量生产与变现流水线Agent",
        "painDepth": { ... },
        "marketLandscape": { ... },
        "mvpPlan": { ... },
        "keywordRelevance": 92,
        "priorityScore": 85,
        "representativeTexts": [ ... ],
        "sourceVideos": [
          {
            "title": "2026最新实战：AI变现从零到一实操拆解",
            "author": "商业探索博主_1",
            "likes": "320",
            "videoUrl": "https://www.douyin.com/video/73001000",
            "createdAt": "2026-10-13 14:00:00"
          }
        ],
        "minioMarkdownUrl": "http://192.168.8.206:9000/oppty-discovery/reports/...",
        "minioHtmlUrl": ""
      }
    ],
    "sourceVideos": [
      {
        "title": "2026最新实战：AI变现从零到一实操拆解",
        "author": "商业探索博主_1",
        "likes": "320",
        "videoUrl": "https://www.douyin.com/video/73001000",
        "createdAt": "2026-10-13 14:00:00"
      }
    ]
  },
  "timestamp": 1791873892438
}
```

