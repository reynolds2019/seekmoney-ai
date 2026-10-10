# 问题修复与维护记录 (RECORD.md)

---

## 记录 01：智能分析启动后进度条闪退消失无提示问题排查与修复

- **记录日期**：2026-10-10
- **涉及服务**：`seekmoney-ai` Web 前端与后端服务
- **运行环境**：Docker 容器（宿主机：`192.168.8.206`，端口：`3000`）

---

### 一、问题现象

用户在 Web 前端页面（`http://192.168.8.206:3000/zh`）点击【开始智能分析】后：
1. 页面短暂展示进度条动画（约 1~2 秒）；
2. 随后进度条突然消失，页面直接重置回初始的“暂无分析结果”空白占位图；
3. 界面未给出任何报错提示或失败原因，任务没有正常执行。

---

### 二、排查过程与根因分析

#### 1. 容器日志排查
查看容器运行日志 `docker logs seekmoney-ai`，发现任务执行在数据抓取步骤中断：
```
[Xiaohongshu Service] 开始搜索关键词: 智能录音卡片, 限制: 10
[Xiaohongshu Service] 搜索第 1 页
[TikHub API] GET /api/v1/xiaohongshu/app_v2/search_notes
[TikHub API] Error 402: {
  detail: {
    code: 402,
    request_id: '93096ec0-28b3-4548-8c84-5d5ed54a6dff',
    message: 'Insufficient balance, this endpoint requires payment and does not accept free credit, to top up balance or get free credit, go: https://user.tikhub.io/users/add_credit',
    message_zh: '免费额度以及余额不足，此路由需要付费，并且不接受使用免费额度请求，前往充值或获取免费额度：https://user.tikhub.io/users/add_credit',
    router: '/api/v1/xiaohongshu/app_v2/search_notes',
    params: { keyword: '智能录音卡片', page: '1' }
  }
}
[Xiaohongshu Service] 搜索失败: Error: Xiaohongshu API 搜索失败: 402 Request failed with status code 402
```

#### 2. TikHub 账户状态确认
调用 TikHub 官方接口 `/api/v1/tikhub/user/get_user_info` 查询 Token 对应的账户资产状态：
```json
{
  "user_data": {
    "balance": 0.0,
    "free_credit": 0.05,
    "email_verified": true,
    "is_active": true
  }
}
```
**根本触发原因**：
- 账户当前 `balance`（付费余额）为 `0.0`，虽有 `0.05` 免费额度，但 TikHub 平台对小红书路由（`/api/v1/xiaohongshu/app_v2/search_notes`）限制为**仅限付费余额扣费，不接受使用免费试用额度**；
- 导致调用小红书数据抓取时直接触发 HTTP 402 拒绝。

#### 3. 前端界面交互缺陷（导致闪退感知）
分析前端 `src/app/[locale]/page.tsx` 代码逻辑：
- 前端使用 `useSWR` 轮询 `/api/jobs/[jobId]`；
- 当服务端返回 `status === "failed"` 时，回调函数直接执行 `setIsLoading(false)`；
- 此时右侧展示区的判断条件为：
  ```tsx
  {results.length === 0 && !isLoading && (
    /* 直接展示空状态占位图，未包含对 status === "failed" 的任何处理 */
  )}
  ```
- 由于结果列表为空且 `isLoading` 为 false，界面瞬间切换回空状态占位图，将后端返回的详细错误信息掩盖，导致用户直观感觉为“进度条刚出来就闪退消失”。

#### 4. 错误信息提取不全
- 原 `tikhub-client.ts` 错误处理中直接读取 `error.response?.data?.message`；
- TikHub 接口返回的中文诊断位于 `data.detail.message_zh`，导致关键诊断信息丢失，只向上层抛出模糊的 `Request failed with status code 402`。

---

### 三、修复方案与代码改动

#### 1. 前端错误诊断卡片与重试交互 (`src/app/[locale]/page.tsx`)
- 新增 `jobData?.status === "failed" && !isLoading` 条件分支；
- 渲染醒目的红底错误诊断卡片，完整展示服务端返回的 `error` 诊断信息，包括错误来源与充值链接；
- 提供【重置并重新分析】按钮，点击后重置任务状态，彻底消除“进度条闪退空白”现象。

#### 2. TikHub 错误深度解析 (`lib/services/tikhub-client.ts`)
- 封装 `formatTikHubError` 辅助函数，优先提取 `data.detail.message_zh`；
- 针对 HTTP 402 状态码，自动格式化为可读性强的中文提示：
  `小红书 API 搜索失败 (HTTP 402 余额不足)：免费额度以及余额不足，此路由需要付费，并且不接受使用免费额度请求，前往充值或获取免费额度：https://user.tikhub.io/users/add_credit`；
- 同步覆盖小红书、TikTok 等数据源的搜索与评论接口。

#### 3. 数据源表单提示优化 (`src/components/AnalysisForm.tsx`)
- 在数据源下拉选单中小红书被选中时，显示明确的提示标签：
  `⚠️ 提示：小红书为 TikHub 平台付费路由，不接受免费试用额度，需账户具备充值余额。`，避免用户在零余额时产生困惑。

#### 4. 大模型 Base URL 兼容适配 (`lib/services/glm-service.ts`, `lib/services/ai-product-service.ts`)
- 自动适配外部传入的 `GLM_BASE_URL`（如 `http://192.168.8.206:3002/v1`），智能补齐 `/chat/completions` 后缀，防止直接 POST 到 `/v1` 报 404。

---

### 四、验证与上线结果

1. **代码同步**：通过 WSL rsync 将改动文件同步至服务器 `/data/seekmoney-ai`；
2. **镜像重构**：复用已缓存的基础依赖层，执行 `docker compose build && docker compose up -d`；
3. **验证结果**：
   - 容器服务状态正常（`Up`，HTTP 200）；
   - 在选择未充值的小红书数据源时，界面不再闪退，而是清晰呈现 402 余额不足诊断面板与充值指引；
   - 切换至支持免费额度的数据源（如抖音、TikTok 等）或充值后，分析流程可完整顺畅运行。

---

## 记录 02：智能分析进度卡在 60%（评论聚类步骤）问题排查与修复

- **记录日期**：2026-10-10
- **涉及服务**：`seekmoney-ai` 后端服务（聚类分析模块与 Embedding 客户端）
- **运行环境**：Docker 容器（宿主机：`192.168.8.206`，端口：`3000`）

---

### 一、问题现象

用户在 Web 前端页面（`http://192.168.8.206:3000/zh`）点击【开始智能分析】后：
1. 进度条正常从 0% 执行至 60%；
2. 状态文案停留在 `正在对 75 条评论内容进行聚类...`；
3. 随后进度条长时间静止，不再向前推进，既不报错也不结束。

---

### 二、排查过程与根因分析

#### 1. 容器运行日志排查
检查容器后台日志发现如下关键错误：
```
400 <400> ***.***.InvalidParameter: Value error, batch size is invalid, it should not be larger than 10.: ***.contents
```
- 系统配置使用的 Embedding 模型为通义千问 `text-embedding-v4`；
- 阿里云 DashScope 对 Embedding 请求单批次文本数量有强制硬性限制：**单次不得超过 10 条**；
- 代码中 `OpenAIEmbeddingProvider` 写死了 `batchSize = 100`，对清洗后的 33 条评论一次性打包发送，直接被通义千问 API 拦截并返回 HTTP 400 错误。

#### 2. 异步步骤缺少 Reject 导致 Promise 永久挂起（死锁根因）
排查 `lib/services/clustering/ClusteringService.ts` 中的 `stepAsync` 实现：
```typescript
private stepAsync<T>(name: string, fn: () => Promise<T>): Promise<T> {
  return new Promise(async (resolve) => {
    const start = Date.now();
    const result = await fn(); // 这里抛出 400 异常
    const duration = Date.now() - start;
    this.steps.push({ name, duration });
    console.log(`[ClusteringService] ${name}: ${duration}ms`);
    resolve(result); // 永远无法到达，且没有 reject 参数！
  });
}
```
- `new Promise(async (resolve) => { ... })` **没有接收 `reject` 参数，且未包含 `try-catch`**；
- 当 `await fn()`（即 Embedding 请求）抛出 400 异常时，该 Promise 既未 resolve 也未 reject，**永久处于 Pending 状态**；
- 外层业务流程 `await this.stepAsync(...)` 被无限期阻塞挂起，任务状态始终无法流转到 `completed` 或 `failed`，导致前台一直显示 60%。

---

### 三、修复方案与代码改动

#### 1. 自适应 Batch Size 与向量维度支持 (`lib/services/clustering/EmbeddingProvider.ts`)
- **模型智能识别**：当模型名称包含 `v4`、`qwen` 或 `dashscope` 时，自动将 `batchSize` 调整为 `10`（适配通义千问限制）；其他模型如 OpenAI 保持 100；
- **环境变量支持**：支持通过 `OPENAI_EMBEDDING_BATCH_SIZE` 手动灵活指定批大小；
- **动态维度感知**：从返回的向量结果中动态识别向量维度，避免由于写死维度引起的兼容性问题。

#### 2. 重构异步步骤包装函数 (`lib/services/clustering/ClusteringService.ts`)
- 将 `stepAsync` 改为标准原生 `async/await` 与 `try-catch` 结构：
  ```typescript
  private async stepAsync<T>(name: string, fn: () => Promise<T>): Promise<T> {
    const start = Date.now();
    try {
      const result = await fn();
      const duration = Date.now() - start;
      this.steps.push({ name, duration });
      console.log(`[ClusteringService] ${name}: ${duration}ms`);
      return result;
    } catch (error) {
      const duration = Date.now() - start;
      console.error(`[ClusteringService] ${name} 步骤失败 (${duration}ms):`, error);
      throw error;
    }
  }
  ```
- 彻底杜绝 Promise 挂起死锁问题，确保任何阶段异常都能立即捕获并通知任务状态机。

---

### 四、验证与上线结果

1. **代码同步**：通过 WSL rsync 将修复代码同步至服务器 `/data/seekmoney-ai/`；
2. **容器构建**：重新编译 Docker 镜像并平滑重启容器；
3. **验证结果**：
   - 评论文本分批（每批 10 条）顺利完成向量化；
   - 聚类算法 DBSCAN 正常计算簇分类；
   - 进度条平稳度过 60% 并顺利进入后续痛点提炼与智能方案生成，流程完整走通。

