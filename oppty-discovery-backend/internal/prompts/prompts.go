package prompts

import (
	"fmt"
	"strings"
)

/**
 * desc: 集中管理商机发现助手所有大模型提示词 (Prompts)
 * date: 2026-10-13 11:25:00
 * @author reynolds
 * @since v1.0.0
 */

// 系统级角色设定
const SystemPromptAnalyst = "你是一位顶级的创业导师、产品经理与商业分析专家，精通从社交媒体真实用户评论与吐槽中挖掘高商业价值的痛点线索，并评估其付费转化潜力与MVP落地路径。请严格按照纯JSON格式输出分析结果，绝不包含任何多余文本或markdown标记。"

// 构建商机聚类深度分析提示词
func BuildClusterAnalysisPrompt(keywords []string, texts []string, dataSize int) string {
	qualityLevel := "初步验证"
	if dataSize < 50 {
		qualityLevel = "小样本探索"
	} else if dataSize >= 200 {
		qualityLevel = "可靠样本"
	}

	joinedKeywords := strings.Join(keywords, "、")
	joinedTexts := strings.Join(texts, "\n\n")

	prompt := fmt.Sprintf(`你是顶级的创业导师和市场分析师。你的核心使命是帮助创业者从真实用户评论中提炼高价值商机，判断方向是否值得落地。

【分析背景】
• 用户检索关键词：「%s」
• 数据来源：社交媒体真实评论与社区讨论
• 数据规模：%d 条（%s）

【核心任务】
从“需求痛点验证”与“商业转化”角度深度透视该聚类评论，输出结构化商机洞察报告。

【分析维度】
1️⃣ 痛点深度挖掘（pain_depth）
- surface_pain: 表面痛点是什么？（一句话精炼）
- root_causes: 连问 3 个“为什么”，找到根本原因（数组，包含 3 个字符串）
- user_scenarios: 这是什么具体场景下的痛点？（数组，1~3 个场景）
- emotional_intensity: 用户情绪强度评分（整数 1~5，1=轻微吐槽，5=极度愤怒或强烈迫切）

2️⃣ 市场与竞品分析（market_landscape）
- existing_solutions: 现有替代方案或竞品（数组，格式：[{"name":"产品或方案名称","limitation":"具体局限性与槽点"}]）
- unmet_needs: 哪些核心需求尚未被市面产品满足？（数组，2~3 个字符串）
- opportunity: 差异化市场机会在哪里？（一句话明确切入点）

3️⃣ MVP 执行方案（mvp_plan）
- core_features: 核心功能清单（数组，3~5 个，一句话说明）
- validation_hypotheses: 需验证的核心商业假设（数组，格式：[{"hypothesis":"核心假设","test_method":"低成本验证方法"}]）
- first_users: 第一批种子用户从哪来？（一句话获客渠道）
- timeline: 预估开发上线周期（如：“2周开发 + 1周测试”）
- estimated_cost: 预估测试落地成本（如：“¥500 种子流量与服务器部署成本”）

4️⃣ 商业评估与相关度
- market_size_score: 潜在受众市场规模评分（整数 1~5 分，1=极小众，3=中等规模几十万，5=超大众级）
- keyword_relevance: 该聚类与搜索关键词「%s」的相关度评分（整数 0~100 分）
- paid_interest: 付费意愿强烈程度（枚举："High" / "Medium" / "Low"）
- one_line_pain: 一句话核心痛点提炼
- rationale: 判断付费意愿与痛点真实性的依据（引用原文依据）
- potential_product: 建议构思的产品概念（一句话说明产品形态与核心卖点）

【待分析真实评论文本】
%s

【严格输出格式要求】
必须严格返回标准 JSON 对象，禁止包裹 markdown 代码块标记，不要多余解释。格式如下：
{
  "one_line_pain": "核心痛点一句话总结",
  "paid_interest": "High",
  "rationale": "判断理由与原文依据",
  "potential_product": "产品概念与定位一句话",
  "pain_depth": {
    "surface_pain": "表面痛点描述",
    "root_causes": ["根本原因1", "根本原因2", "根本原因3"],
    "user_scenarios": ["场景1", "场景2"],
    "emotional_intensity": 4
  },
  "market_landscape": {
    "existing_solutions": [
      {"name": "现有产品A", "limitation": "局限性与不足"},
      {"name": "现有方案B", "limitation": "局限性与不足"}
    ],
    "unmet_needs": ["未满足需求1", "未满足需求2"],
    "opportunity": "切入机会描述"
  },
  "mvp_plan": {
    "core_features": ["核心功能1", "核心功能2", "核心功能3"],
    "validation_hypotheses": [
      {"hypothesis": "假设1", "test_method": "测试验证方法1"}
    ],
    "first_users": "种子用户获取渠道",
    "timeline": "2周上线 MVP",
    "estimated_cost": "¥500 部署与测试成本"
  },
  "market_size_score": 4,
  "keyword_relevance": 90
}`, joinedKeywords, dataSize, qualityLevel, joinedKeywords, joinedTexts)

	return prompt
}
