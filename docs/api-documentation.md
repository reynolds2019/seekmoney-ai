# SeekMoney AI 服务接口清单文档

## 一、概述

### 1.1 系统说明
SeekMoney AI 是一个基于社交媒体数据采集、语义聚类与大模型分析的自动化商业痛点挖掘系统。前端基于 Next.js 15 App Router 开发，提供国际化页面支持；后端通过 Next.js Route Handlers 提供异步任务式的 RESTful API 接口。

### 1.2 服务基础信息
- **部署访问地址**：`http://192.168.8.206:3000`
- **页面入口地址**：`http://192.168.8.206:3000/zh`（中文）/ `http://192.168.8.206:3000/en`（英文）
- **接口协议**：HTTP / HTTPS
- **数据传输格式**：JSON (`Content-Type: application/json`)
- **交互模式**：对于耗时较长的数据爬取、语义聚类和大模型分析，采用 **「异步任务创建 (202 Accepted) + 客户端轮询 (SWR/GET)」** 机制。

---

## 二、接口总览表

| 接口分类 | 请求方式 | 接口路径 | 功能说明 | 环境权限 |
| :--- | :--- | :--- | :--- | :--- |
| **痛点挖掘分析** | `POST` | `/api/analyze` | 创建痛点挖掘与聚类分析异步任务 | 全部环境 |
| **痛点挖掘分析** | `GET` | `/api/jobs/[jobId]` | 查询痛点挖掘任务状态、进度与分析报告 | 全部环境 |
| **AI 产品创意分析** | `POST` | `/api/analyze-ai-product` | 创建 AI 产品创意挖掘异步任务 | 全部环境 |
| **AI 产品创意分析** | `GET` | `/api/ai-product-jobs/[jobId]` | 查询 AI 产品分析任务状态、进度与产品方案 | 全部环境 |
| **系统监控** | `GET` | `/api/health` | 服务健康检查与任务负载统计 | 全部环境（开发模式包含任务明细） |
| **运维诊断** | `GET` | `/api/test-env` | 环境变量配置完整性与连通性自检 | 仅开发环境 (`development`) |
| **运维诊断** | `GET` | `/api/test-tikhub` | TikHub API 直连测试与网络连通性排查 | 仅开发环境 (`development`) |
| **运维诊断** | `GET` | `/api/debug-tikhub` | TikHub 爬虫与数据解析全链路分步诊断 | 仅开发环境 (`development`) |

---

## 三、核心业务接口详细定义

### 1. 创建痛点挖掘分析任务

- **接口地址**：`POST /api/analyze`
- **接口描述**：提交目标搜索关键词及抓取配置，后端校验入参并初始化异步任务队列，立即返回任务唯一标识 `jobId`。
- **请求头**：
  ```http
  Content-Type: application/json
  ```
- **请求体 (Request Body)**：

| 字段名 | 类型 | 必填 | 默认值 | 限制条件 / 取值范围 | 说明 |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `keywords` | `string[]` | 是 | 无 | 1 ~ 10 个元素，每个词最多 50 字符 | 目标搜索关键词列表 |
| `limit` | `number` | 否 | `200` | 1 ~ 500 | 抓取数据总量上限 |
| `dataSource` | `string` | 否 | `'tiktok'` | 枚举：`'tikhub'`, `'tiktok'`, `'bilibili'`, `'wechat'`, `'youtube'`, `'xiaohongshu'` | 数据来源平台 |
| `deepCrawl` | `boolean` | 否 | `false` | `true` / `false` | 是否开启深度抓取（包含评论） |
| `maxVideos` | `number` | 否 | `10` | 1 ~ 50 | 深度抓取时的最大视频数 |
| `locale` | `string` | 否 | `'zh'` | `'zh'` / `'en'` | 输出分析报告的语言 |
| `tiktokConfig`| `object` | 否 | 无 | 仅 TikTok/抖音 数据源生效 | 深度采集细项参数配置 |
| └ `enableComments` | `boolean` | 否 | `true` | `true` / `false` | 是否抓取视频评论 |
| └ `maxVideos` | `number` | 否 | `15` | 1 ~ 50 | 采集的最大视频数 |
| └ `maxCommentsPerVideo` | `number` | 否 | `20` | 1 ~ 100 | 单个视频抓取评论条数 |
| └ `enableSubComments` | `boolean` | 否 | `false` | `true` / `false` | 是否抓取二级评论（楼中楼） |

- **请求示例**：
  ```json
  {
    "keywords": ["AI变现", "自媒体运营"],
    "limit": 100,
    "dataSource": "tiktok",
    "deepCrawl": true,
    "locale": "zh",
    "tiktokConfig": {
      "enableComments": true,
      "maxVideos": 10,
      "maxCommentsPerVideo": 30,
      "enableSubComments": false
    }
  }
  ```

- **响应状态码**：
  - `202 Accepted`：任务创建成功，已提交后台异步执行
  - `400 Bad Request`：入参校验未通过（关键词为空/超长、不支持的数据源等）
  - `500 Internal Server Error`：服务器内部错误

- **响应示例 (202 Accepted)**：
  ```json
  {
    "jobId": "f784d1bc-6ea9-42b4-82a1-638c35bb5f51"
  }
  ```

---

### 2. 查询痛点挖掘任务状态与结果

- **接口地址**：`GET /api/jobs/[jobId]`
- **接口描述**：根据 `jobId` 查询痛点分析任务的执行阶段、百分比进度，以及执行完毕后的完整聚类与大模型深度商业分析报告。
- **路径参数 (Path Parameters)**：
  - `jobId` (`string`, 必需)：创建任务时返回的 UUID。

- **响应状态码**：
  - `200 OK`：成功查询到任务状态或结果
  - `400 Bad Request`：未传递 `jobId`
  - `404 Not Found`：任务不存在或服务器重启导致内存记录丢失
  - `500 Internal Server Error`：服务器异常

- **进行中响应示例 (200 OK - processing)**：
  ```json
  {
    "jobId": "f784d1bc-6ea9-42b4-82a1-638c35bb5f51",
    "status": "processing",
    "progress": "正在进行语义聚类...",
    "progressStage": "clustering",
    "progressPercent": 60,
    "keywords": ["AI变现", "自媒体运营"]
  }
  ```
  > `progressStage` 阶段枚举值：`'init'`（初始化） -> `'validating'`（参数校验） -> `'crawling'`（抓取数据） -> `'clustering'`（语义聚类） -> `'analyzing'`（大模型深度分析） -> `'completed'`（已完成） / `'failed'`（已失败）。

- **完成时响应示例 (200 OK - completed)**：
  ```json
  {
    "jobId": "f784d1bc-6ea9-42b4-82a1-638c35bb5f51",
    "status": "completed",
    "progress": "分析完成",
    "progressStage": "completed",
    "progressPercent": 100,
    "keywords": ["AI变现", "自媒体运营"],
    "results": [
      {
        "id": "cluster_0",
        "size": 38,
        "representative_texts": ["新手不知道怎么用AI写文案变现", "尝试了很多工具都赚不到钱"],
        "analysis": {
          "one_line_pain": "缺乏系统化将AI提示词转化为具体交付成果的能力",
          "paid_interest": "High",
          "rationale": "大量用户表达付费购买保姆级工作流模板的强烈意愿",
          "potential_product": "垂直行业爆款图文一键生成器及提示词知识库",
          "pain_depth": {
            "surface_pain": "不知道用什么AI工具高效写稿",
            "root_causes": ["信息差大", "缺少商业闭环指引", "缺乏实际交付物模板"],
            "user_scenarios": ["副业创业", "自媒体图文代运营"],
            "emotional_intensity": 4
          },
          "market_landscape": {
            "existing_solutions": [
              { "name": "通用大模型对话框", "limitation": "缺少针对垂直平台的优化" }
            ],
            "unmet_needs": ["开箱即用的工作流模板", "对标爆款拆解工具"],
            "opportunity": "打造垂直人群的微型 SaaS 工具"
          },
          "mvp_plan": {
            "core_features": ["小红书爆款文案仿写", "多轮自动化图文生成"],
            "validation_hypotheses": [
              { "hypothesis": "用户愿意为垂直工作流付费", "test_method": "社群预售" }
            ],
            "first_users": "自媒体交流社群群友",
            "timeline": "2周上线 MVP",
            "estimated_cost": "¥500 部署与测试成本"
          },
          "keyword_relevance": 92
        },
        "priority_score": {
          "demand_intensity": 85,
          "market_size": 75,
          "competition": 60,
          "overall": 80,
          "level": "High"
        }
      }
    ],
    "rawData": {
      "videos": [],
      "comments": [],
      "rawTexts": []
    },
    "clusteredData": []
  }
  ```

---

### 3. 创建 AI 产品创意分析任务

- **接口地址**：`POST /api/analyze-ai-product`
- **接口描述**：针对特定关键词，跳过传统多层分析，直接驱动大模型从用户反馈中提炼产品定位、核心功能矩阵、AI 能力架构与落地 Roadmap。
- **请求头**：
  ```http
  Content-Type: application/json
  ```
- **请求体 (Request Body)**：

| 字段名 | 类型 | 必填 | 默认值 | 限制条件 / 取值范围 | 说明 |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `keywords` | `string[]` | 是 | 无 | 1 ~ 10 个元素，每个词最多 50 字符 | 目标搜索关键词列表 |
| `limit` | `number` | 否 | `50` | 1 ~ 200 | 抓取分析的文本条数上限 |
| `dataSource` | `string` | 否 | `'tiktok'` | 枚举：`'tikhub'`, `'tiktok'`, `'bilibili'`, `'wechat'`, `'youtube'`, `'xiaohongshu'` | 数据来源平台 |
| `locale` | `string` | 否 | `'zh'` | `'zh'` / `'en'` | 输出分析报告的语言 |

- **请求示例**：
  ```json
  {
    "keywords": ["跨境电商独立站"],
    "limit": 50,
    "dataSource": "tiktok",
    "locale": "zh"
  }
  ```

- **响应状态码**：
  - `202 Accepted`：任务创建成功
  - `400 Bad Request`：参数格式或数量异常
  - `500 Internal Server Error`：服务器内部错误

- **响应示例 (202 Accepted)**：
  ```json
  {
    "jobId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  }
  ```

---

### 4. 查询 AI 产品创意分析任务状态与结果

- **接口地址**：`GET /api/ai-product-jobs/[jobId]`
- **接口描述**：根据 `jobId` 查询 AI 产品创意分析任务的执行进度及最终生成的完整产品定义文档。
- **路径参数 (Path Parameters)**：
  - `jobId` (`string`, 必需)：AI 产品分析任务的 UUID。

- **响应状态码**：
  - `200 OK`：成功获取任务状态或分析结果
  - `400 Bad Request`：缺少 `jobId`
  - `404 Not Found`：任务不存在
  - `500 Internal Server Error`：服务器内部错误

- **完成时响应示例 (200 OK)**：
  ```json
  {
    "jobId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "status": "completed",
    "progress": "分析完成",
    "progressStage": "completed",
    "progressPercent": 100,
    "keywords": ["跨境电商独立站"],
    "limit": 50,
    "dataSource": "tiktok",
    "locale": "zh",
    "startTime": 1740000000000,
    "results": [
      {
        "id": "prod_1",
        "size": 25,
        "representative_texts": ["独立站选品和投放太难了，素材消耗极快"],
        "analysis": {
          "product_name": "AdGenius AI - 跨平台电商素材生成器",
          "product_category": "SaaS / 电商营销自动化",
          "target_users": "中小跨境独立站卖家与投放投手",
          "core_features": [
            "一键抓取竞品爆款视频分镜头",
            "自动生成多语言脚本与混剪模板",
            "批量生成符合 Meta / TikTok 规范的广告图文"
          ],
          "ai_capabilities": [
            "多模态视频内容理解",
            "跨语言营销文案本地化适配",
            "基于历史转化数据的自动 A/B 预测"
          ],
          "value_proposition": "将单条爆款视频的复刻与测试周期从 3 天缩短至 10 分钟",
          "implementation_difficulty": "Medium",
          "market_potential": "High",
          "technical_stack": [
            "Next.js",
            "FastAPI",
            "FFmpeg",
            "GLM-4 / Claude",
            "Whisper"
          ],
          "development_roadmap": [
            "第 1-2 周：上线基础脚本生成器与分镜头模板",
            "第 3-4 周：打通视频自动化混剪与渲染管线",
            "第 5-6 周：集成 TikTok Ads 接口实现一键发布"
          ]
        }
      }
    ]
  }
  ```

---

## 四、系统监控与运维诊断接口

### 1. 服务健康检查

- **接口地址**：`GET /api/health`
- **接口描述**：检查服务运行时状态、正常运行时间（Uptime）；开发模式下额外返回任务队列统计。
- **环境限制**：全环境可用。
- **响应示例 (200 OK - 生产环境)**：
  ```json
  {
    "status": "healthy",
    "timestamp": "2026-10-12T12:35:00.000Z",
    "uptime": 3600.45
  }
  ```
- **响应示例 (200 OK - 开发环境)**：
  ```json
  {
    "status": "healthy",
    "timestamp": "2026-10-12T12:35:00.000Z",
    "uptime": 3600.45,
    "jobs": {
      "stats": {
        "total": 5,
        "processing": 1,
        "completed": 3,
        "failed": 1
      },
      "recentJobs": ["job-uuid-1", "job-uuid-2"]
    }
  }
  ```

---

### 2. 环境变量加载检测 (Dev Only)

- **接口地址**：`GET /api/test-env`
- **接口描述**：验证服务端关键环境变量（`TIKHUB_API_TOKEN`、`GLM_API_KEY` 等）是否已注入，并模拟调用基础测试。
- **环境限制**：**仅开发环境有效**。若当前为生产环境 (`NODE_ENV === 'production'`)，直接返回 `404 Not Found`。
- **响应示例 (200 OK)**：
  ```json
  {
    "env": {
      "TIKHUB_API_TOKEN": "已设置",
      "TIKHUB_API_KEY": "未设置",
      "GLM_API_KEY": "已设置",
      "NODE_ENV": "development"
    },
    "tests": [
      { "name": "直接 env Token", "success": true, "code": 200 },
      { "name": "createTikHubClient", "success": true, "code": 200 }
    ]
  }
  ```

---

### 3. TikHub API 连通性测试 (Dev Only)

- **接口地址**：`GET /api/test-tikhub`
- **接口描述**：使用当前配置的 Token，向 TikHub 接口（`https://api.tikhub.io`）发起搜索请求测试连通性与鉴权。
- **环境限制**：**仅开发环境有效**。生产环境下返回 `404 Not Found`。
- **响应示例 (200 OK)**：
  ```json
  {
    "success": true,
    "results": [
      { "test": "axios.post", "status": 200, "data": { "code": 200, "msg": "success" } },
      { "test": "axios.create", "status": 200, "data": { "code": 200, "msg": "success" } }
    ]
  }
  ```

---

### 4. TikHub 全链路分步调试 (Dev Only)

- **接口地址**：`GET /api/debug-tikhub`
- **接口描述**：执行包括数据源创建、服务可用性检查、基础搜索、带评论搜索、数据清洗在内的全套诊断流程。
- **环境限制**：**仅开发环境有效**。生产环境下返回 `404 Not Found`。
- **响应示例 (200 OK)**：
  ```json
  {
    "steps": [
      { "name": "创建数据源", "success": true },
      { "name": "检查可用性", "success": true, "available": true },
      { "name": "基础搜索", "success": true, "videoCount": 5, "textCount": 5 },
      { "name": "带评论搜索", "success": true, "videoCount": 2, "commentCount": 15 },
      { "name": "数据清洗", "success": true, "original": 20, "cleaned": 18 }
    ]
  }
  ```

---

## 五、错误代码与异常处理规范

各接口遵循标准 HTTP 状态码规范：

| HTTP 状态码 | 业务场景 | 统一错误返回结构 |
| :--- | :--- | :--- |
| `400 Bad Request` | 客户端入参验证失败（如缺失必填参数、参数越界） | `{"error": "具体参数错误提示"}` |
| `404 Not Found` | 任务不存在，或在生产环境下访问调试接口 | `{"error": "任务不存在 / Not found"}` |
| `500 Internal Server Error` | 服务端未捕获的内部异常、三方网络中断 | `{"error": "服务器内部错误"}` |

---

## 六、快速调用与测试示例 (cURL)

### 1. 提交痛点分析任务
```bash
curl -X POST "http://192.168.8.206:3000/api/analyze" \
  -H "Content-Type: application/json" \
  -d '{
    "keywords": ["AI写作", "自媒体变现"],
    "limit": 50,
    "dataSource": "tiktok",
    "deepCrawl": false,
    "locale": "zh"
  }'
```

### 2. 轮询任务状态及获取结果
```bash
# 替换为实际返回的 jobId
curl -X GET "http://192.168.8.206:3000/api/jobs/f784d1bc-6ea9-42b4-82a1-638c35bb5f51"
```

### 3. 服务健康检查
```bash
curl -X GET "http://192.168.8.206:3000/api/health"
```
