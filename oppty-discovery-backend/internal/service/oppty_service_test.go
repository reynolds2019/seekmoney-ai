package service_test

import (
	"strings"
	"testing"

	"oppty-discovery-backend/internal/config"
	"oppty-discovery-backend/internal/prompts"
	"oppty-discovery-backend/internal/service"
)

/**
 * desc: 商机核心业务逻辑与成本计算单元测试
 * date: 2026-10-13 11:25:00
 * @author reynolds
 * @since v1.0.0
 */

// 测试预估成本计算公式
func TestCalculateEstimatedCost(t *testing.T) {
	t.Run("测试基础成本_不开启评论抓取", func(t *testing.T) {
		svc := service.NewOpptyService(&config.AppConfig{}, nil, nil)
		// 10 个视频，不抓评论：10 / 20 * 0.01 = 0.005 -> 0.01 元
		cost := svc.CalculateEstimatedCost(10, 20, false)
		if cost < 0.005 || cost > 0.015 {
			t.Errorf("预估成本计算偏差，期望约 0.01，实际为: %f", cost)
		}
	})

	t.Run("测试深度抓取成本_视频10评论20", func(t *testing.T) {
		svc := service.NewOpptyService(&config.AppConfig{}, nil, nil)
		// 视频成本: 10/20*0.01 = 0.005; 评论成本: 10 * (20/20)*0.01 = 0.10; 总计: 0.105 -> 0.11 元
		cost := svc.CalculateEstimatedCost(10, 20, true)
		if cost != 0.11 {
			t.Errorf("深度成本计算偏差，期望 0.11，实际为: %f", cost)
		}
	})

	t.Run("测试大样本深度抓取成本_视频20评论50", func(t *testing.T) {
		svc := service.NewOpptyService(&config.AppConfig{}, nil, nil)
		// 视频成本: 20/20*0.01 = 0.01; 评论成本: 20 * (50/20)*0.01 = 0.50; 总计: 0.51 元
		cost := svc.CalculateEstimatedCost(20, 50, true)
		if cost != 0.51 {
			t.Errorf("大样本深度成本计算偏差，期望 0.51，实际为: %f", cost)
		}
	})
}

// 测试提示词生成完整性
func TestBuildClusterAnalysisPrompt(t *testing.T) {
	t.Run("测试分析提示词格式包含必须约束与字段", func(t *testing.T) {
		keywords := []string{"AI写作", "自媒体变现"}
		texts := []string{"做自媒体不知道怎么选题", "市面上的课程都是割韭菜"}
		prompt := prompts.BuildClusterAnalysisPrompt(keywords, texts, 20)

		if !strings.Contains(prompt, "AI写作、自媒体变现") {
			t.Errorf("提示词中未包含搜索关键词")
		}
		if !strings.Contains(prompt, "pain_depth") {
			t.Errorf("提示词中未包含 pain_depth 维度")
		}
		if !strings.Contains(prompt, "market_landscape") {
			t.Errorf("提示词中未包含 market_landscape 维度")
		}
		if !strings.Contains(prompt, "mvp_plan") {
			t.Errorf("提示词中未包含 mvp_plan 维度")
		}
	})
}
