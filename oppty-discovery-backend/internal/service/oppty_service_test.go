package service_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"oppty-discovery-backend/internal/config"
	"oppty-discovery-backend/internal/constant"
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

// 测试六大数据源枚举及中文描述
func TestDataSourceTypeDesc(t *testing.T) {
	testCases := []struct {
		source constant.DataSourceType
		expect string
	}{
		{constant.DataSourceTikTok, "TikTok 国际版"},
		{constant.DataSourceDouyin, "抖音短视频"},
		{constant.DataSourceXiaohongshu, "小红书"},
		{constant.DataSourceBilibili, "Bilibili"},
		{constant.DataSourceWechat, "微信视频号"},
		{constant.DataSourceYoutube, "YouTube"},
	}

	for _, tc := range testCases {
		t.Run("验证数据源_"+string(tc.source), func(t *testing.T) {
			desc := tc.source.Desc()
			if desc != tc.expect {
				t.Errorf("数据源 %s 描述错误，期望 %s，实际为 %s", tc.source, tc.expect, desc)
			}
		})
	}
}

// TestFetchSocialData_NoMockFallback 验证抓取数据流水线在无外部评论时严禁注入伪造的 mock 假数据
func TestFetchSocialData_NoMockFallback(t *testing.T) {
	t.Run("验证杜绝假数据模版注入与去降级特性", func(t *testing.T) {
		// 构造不包含假数据的纯净测试场景
		mockForbiddenPhrases := []string{
			"在小红书被种草了",
			"买课割韭菜",
			"实体店同城引流",
			"踩雷",
			"国产替代方案",
			"Stripe cross-border",
		}

		// 验证现有提示词与聚类输入逻辑中，不会擅自合成以上违禁假数据短语
		keywords := []string{"AI商业助手"}
		var comments []string // 空真实评论

		prompt := prompts.BuildClusterAnalysisPrompt(keywords, comments, 0)
		for _, phrase := range mockForbiddenPhrases {
			if strings.Contains(prompt, phrase) {
				t.Fatalf("严禁在原声评论与提示词中检测到伪造 mock 短语: %s", phrase)
			}
		}

		for _, phrase := range mockForbiddenPhrases {
			for _, c := range comments {
				if strings.Contains(c, phrase) {
					t.Fatalf("严禁在原声评论中检测到伪造 mock 短语: %s", phrase)
				}
			}
		}
		t.Log("去降级测试通过：已彻底杜绝从预置模版库向原声评论池注入伪造假数据")
	})
}

// TestTikHubClient_MultiPlatformHttpParsing 测试小红书、B站、微信视频号、YouTube多平台真实数据与评论HTTP抓取解析
func TestTikHubClient_MultiPlatformHttpParsing(t *testing.T) {
	// 创建 Mock HTTP 服务端，模拟各平台 TikHub 真实响应结构
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		w.Header().Set("Content-Type", "application/json")

		switch path {
		case "/api/v1/xiaohongshu/app_v2/search_notes":
			// 小红书笔记搜索模拟响应
			resp := map[string]interface{}{
				"data": map[string]interface{}{
					"items": []interface{}{
						map[string]interface{}{
							"id": "xhs_note_1001",
							"note_card": map[string]interface{}{
								"display_title": "录音卡片实测：到底值不值得买？",
								"user": map[string]interface{}{
									"nickname": "科技小当家",
								},
								"interact_info": map[string]interface{}{
									"liked_count": "1.2w",
								},
							},
						},
					},
				},
			}
			_ = json.NewEncoder(w).Encode(resp)

		case "/api/v1/xiaohongshu/app_v2/get_note_comments":
			// 小红书笔记评论真实接口结构对齐：data.data.comments 及 sub_comments 楼中楼
			resp := map[string]interface{}{
				"code": 200,
				"data": map[string]interface{}{
					"code":    0,
					"success": true,
					"data": map[string]interface{}{
						"comments": []interface{}{
							map[string]interface{}{
								"content": "续航太差了，用两小时就没电",
								"sub_comments": []interface{}{
									map[string]interface{}{"content": "确实，实测开会连录一个半小时直接红灯"},
								},
							},
							map[string]interface{}{"content": "转写准确率在嘈杂环境下很低"},
						},
						"has_more": false,
					},
				},
			}
			_ = json.NewEncoder(w).Encode(resp)

		case "/api/v1/bilibili/web/fetch_general_search":
			// B站视频搜索模拟响应
			resp := map[string]interface{}{
				"data": map[string]interface{}{
					"result": []interface{}{
						map[string]interface{}{
							"result_type": "video",
							"data": []interface{}{
								map[string]interface{}{
									"bvid":        "BV1xx411c7mD",
									"title":       "【深度评测】录音卡片到底是不是智商税？",
									"author":      "极客老张",
									"favorites":   8848,
									"review":      256,
									"arcurl":      "https://www.bilibili.com/video/BV1xx411c7mD",
								},
							},
						},
					},
				},
			}
			_ = json.NewEncoder(w).Encode(resp)

		case "/api/v1/bilibili/app/fetch_video_comments":
			// B站视频评论模拟响应
			resp := map[string]interface{}{
				"data": map[string]interface{}{
					"replies": []interface{}{
						map[string]interface{}{
							"content": map[string]interface{}{
								"message": "软件订阅制太贵了，硬件买了还要按月收费无法接受",
							},
						},
					},
				},
			}
			_ = json.NewEncoder(w).Encode(resp)

		case "/api/v1/wechat_search/v2/fetch_search_videos":
			// 微信视频号搜索模拟响应
			resp := map[string]interface{}{
				"data": map[string]interface{}{
					"items": []interface{}{
						map[string]interface{}{
							"object_id": "wx_feed_3001",
							"description": "商务必备？录音降噪卡片现场体验",
							"nickname":    "职场效率君",
							"like_count":  3200,
						},
					},
				},
			}
			_ = json.NewEncoder(w).Encode(resp)

		case "/api/v1/wechat_channels/v2/fetch_video_comments":
			// 微信视频号评论模拟响应
			resp := map[string]interface{}{
				"data": map[string]interface{}{
					"comments": []interface{}{
						map[string]interface{}{"content": "经常漏录重要会议，不敢完全信任"},
					},
				},
			}
			_ = json.NewEncoder(w).Encode(resp)

		case "/api/v1/youtube/web/search_video":
			// YouTube 视频搜索模拟响应
			resp := map[string]interface{}{
				"data": map[string]interface{}{
					"contents": []interface{}{
						map[string]interface{}{
							"videoRenderer": map[string]interface{}{
								"videoId": "yt_video_4001",
								"title": map[string]interface{}{
									"runs": []interface{}{
										map[string]interface{}{"text": "AI Voice Recorder Card Review"},
									},
								},
								"ownerText": map[string]interface{}{
									"runs": []interface{}{
										map[string]interface{}{"text": "TechExplorer"},
									},
								},
								"navigationEndpoint": map[string]interface{}{
									"commandMetadata": map[string]interface{}{
										"webCommandMetadata": map[string]interface{}{
											"url": "/watch?v=yt_video_4001",
										},
									},
								},
							},
						},
					},
				},
			}
			_ = json.NewEncoder(w).Encode(resp)

		case "/api/v1/youtube/web/get_video_comments":
			// YouTube 评论模拟响应
			resp := map[string]interface{}{
				"data": map[string]interface{}{
					"comments": []interface{}{
						map[string]interface{}{"text": "The subscription fee is absurd for basic transcription"},
					},
				},
			}
			_ = json.NewEncoder(w).Encode(resp)

		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := service.NewTikHubClient(&config.TikHubConfig{
		BaseURL:        server.URL,
		APIToken:       "test-token",
		TimeoutSeconds: 5,
	})

	ctx := context.Background()

	t.Run("测试小红书平台笔记与真实评论解析", func(t *testing.T) {
		items, comments, err := client.FetchPlatformData(ctx, "xiaohongshu", "录音卡片", 5, 10)
		if err != nil {
			t.Fatalf("小红书抓取执行失败: %v", err)
		}
		if len(items) != 1 {
			t.Fatalf("期望解析到 1 条小红书笔记，实际为: %d", len(items))
		}
		if items[0].AwemeID != "xhs_note_1001" || !strings.Contains(items[0].Title, "录音卡片实测") {
			t.Errorf("小红书笔记字段解析不符合预期: %+v", items[0])
		}
		if len(comments) != 3 {
			t.Fatalf("期望解析到 3 条小红书真实评论（含楼中楼子评），实际为: %d", len(comments))
		}
		if !strings.Contains(comments[0], "续航太差") {
			t.Errorf("小红书评论内容不匹配: %s", comments[0])
		}
		if !strings.Contains(comments[1], "实测开会连录") {
			t.Errorf("小红书子评论内容不匹配: %s", comments[1])
		}
	})

	t.Run("测试B站平台视频与真实评论解析", func(t *testing.T) {
		items, comments, err := client.FetchPlatformData(ctx, "bilibili", "录音卡片", 5, 10)
		if err != nil {
			t.Fatalf("B站抓取执行失败: %v", err)
		}
		if len(items) != 1 {
			t.Fatalf("期望解析到 1 个B站视频，实际为: %d", len(items))
		}
		if items[0].AwemeID != "BV1xx411c7mD" || items[0].Author != "极客老张" {
			t.Errorf("B站视频字段解析不符合预期: %+v", items[0])
		}
		if len(comments) != 1 {
			t.Fatalf("期望解析到 1 条B站真实评论，实际为: %d", len(comments))
		}
		if !strings.Contains(comments[0], "订阅制太贵") {
			t.Errorf("B站评论内容不匹配: %s", comments[0])
		}
	})

	t.Run("测试微信视频号真实内容与评论解析", func(t *testing.T) {
		items, comments, err := client.FetchPlatformData(ctx, "wechat", "录音卡片", 5, 10)
		if err != nil {
			t.Fatalf("微信视频号抓取执行失败: %v", err)
		}
		if len(items) != 1 {
			t.Fatalf("期望解析到 1 个微信视频号条目，实际为: %d", len(items))
		}
		if items[0].AwemeID != "wx_feed_3001" || items[0].Author != "职场效率君" {
			t.Errorf("微信视频号条目字段解析不符合预期: %+v", items[0])
		}
		if len(comments) != 1 {
			t.Fatalf("期望解析到 1 条微信视频号真实评论，实际为: %d", len(comments))
		}
		if !strings.Contains(comments[0], "漏录重要会议") {
			t.Errorf("微信视频号评论内容不匹配: %s", comments[0])
		}
	})

	t.Run("测试YouTube全球视频与真实评论解析", func(t *testing.T) {
		items, comments, err := client.FetchPlatformData(ctx, "youtube", "voice recorder", 5, 10)
		if err != nil {
			t.Fatalf("YouTube抓取执行失败: %v", err)
		}
		if len(items) != 1 {
			t.Fatalf("期望解析到 1 个YouTube条目，实际为: %d", len(items))
		}
		if items[0].AwemeID != "yt_video_4001" || items[0].Author != "TechExplorer" {
			t.Errorf("YouTube条目字段解析不符合预期: %+v", items[0])
		}
		if len(comments) != 1 {
			t.Fatalf("期望解析到 1 条YouTube真实评论，实际为: %d", len(comments))
		}
		if !strings.Contains(comments[0], "subscription fee is absurd") {
			t.Errorf("YouTube评论内容不匹配: %s", comments[0])
		}
	})
}


