package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"oppty-discovery-backend/internal/config"
)

/**
 * desc: TikHub 社交媒体数据源 API 客户端封装
 * date: 2026-10-13 11:50:00
 * @author reynolds
 * @since v1.0.0
 */
type TikHubClient interface {
	SearchVideos(ctx context.Context, keyword string, limit int) ([]TikHubVideoItem, error)
	GetVideoComments(ctx context.Context, awemeID string, count int) ([]string, error)
}

type TikHubVideoItem struct {
	AwemeID     string `json:"awemeId"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Likes       string `json:"likes"`
	ShareURL    string `json:"shareUrl"`
	CommentCount int   `json:"commentCount"`
}

type tikHubClientImpl struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
}

func NewTikHubClient(cfg *config.TikHubConfig) TikHubClient {
	timeout := 30 * time.Second
	if cfg.TimeoutSeconds > 0 {
		timeout = time.Duration(cfg.TimeoutSeconds) * time.Second
	}
	return &tikHubClientImpl{
		baseURL:  strings.TrimRight(cfg.BaseURL, "/"),
		apiToken: cfg.APIToken,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// SearchVideos 调用 TikHub 搜索视频接口
func (c *tikHubClientImpl) SearchVideos(ctx context.Context, keyword string, limit int) ([]TikHubVideoItem, error) {
	endpoint := fmt.Sprintf("%s/api/v1/douyin/search/fetch_general_search_v1", c.baseURL)
	log.Printf("[TikHub API] 正在发起全网视频搜索，关键词: [%s], 限制条数: %d, 端点: %s", keyword, limit, endpoint)

	reqPayload := map[string]interface{}{
		"keyword":         keyword,
		"cursor":          0,
		"sort_type":       "0",
		"publish_time":    "0",
		"filter_duration": "0",
		"content_type":    "1",
		"search_id":       "",
		"backtrace":       "",
	}

	bodyBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("序列化 TikHub 请求体失败: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("构建 TikHub 请求失败: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 TikHub 搜索接口网络异常: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取 TikHub 响应体失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[TikHub API 警告] 响应状态码非 200: %d, 响应: %s", resp.StatusCode, string(respBytes))
		return nil, fmt.Errorf("TikHub API 响应错误 (HTTP %d): %s", resp.StatusCode, string(respBytes))
	}

	var rawRes map[string]interface{}
	if err := json.Unmarshal(respBytes, &rawRes); err != nil {
		return nil, fmt.Errorf("解析 TikHub JSON 响应失败: %w", err)
	}

	var items []interface{}
	if dataMap, ok := rawRes["data"].(map[string]interface{}); ok {
		if innerData, ok := dataMap["data"].([]interface{}); ok {
			items = innerData
		} else if innerMap, ok := dataMap["data"].(map[string]interface{}); ok {
			if deeperData, ok := innerMap["data"].([]interface{}); ok {
				items = deeperData
			}
		}
	} else if dataArr, ok := rawRes["data"].([]interface{}); ok {
		items = dataArr
	}

	var results []TikHubVideoItem
	for _, it := range items {
		itemMap, ok := it.(map[string]interface{})
		if !ok {
			continue
		}

		awemeInfo, ok := itemMap["aweme_info"].(map[string]interface{})
		if !ok {
			awemeInfo = itemMap
		}

		awemeID, _ := awemeInfo["aweme_id"].(string)
		desc, _ := awemeInfo["desc"].(string)
		shareURL, _ := awemeInfo["share_url"].(string)

		nickname := "社区创作者"
		if authorMap, ok := awemeInfo["author"].(map[string]interface{}); ok {
			if name, ok := authorMap["nickname"].(string); ok && name != "" {
				nickname = name
			}
		}

		diggCountStr := "100"
		commentCount := 20
		if stats, ok := awemeInfo["statistics"].(map[string]interface{}); ok {
			if digg, ok := stats["digg_count"].(float64); ok {
				diggCountStr = fmt.Sprintf("%.0f", digg)
			}
			if cCount, ok := stats["comment_count"].(float64); ok {
				commentCount = int(cCount)
			}
		}

		if awemeID != "" || desc != "" {
			results = append(results, TikHubVideoItem{
				AwemeID:      awemeID,
				Title:        desc,
				Author:       nickname,
				Likes:        diggCountStr,
				ShareURL:     shareURL,
				CommentCount: commentCount,
			})
		}

		if len(results) >= limit {
			break
		}
	}

	log.Printf("[TikHub API] 成功从全网抓取并解析到 %d 条视频数据", len(results))
	return results, nil
}

/**
 * desc: 调用 TikHub 抓取单个短视频下的真实用户原声评论
 * date: 2026-10-15 11:25:00
 * @author reynolds
 * @since v1.0.0
 * @param ctx 上下文
 * @param awemeID 视频ID
 * @param count 限制评论条数
 * @return []string 真实评论文本切片
 * @return error 异常信息
 */
func (c *tikHubClientImpl) GetVideoComments(ctx context.Context, awemeID string, count int) ([]string, error) {
	if awemeID == "" {
		return nil, nil
	}

	if count <= 0 {
		count = 20
	}

	params := url.Values{}
	params.Set("aweme_id", awemeID)
	params.Set("cursor", "0")
	params.Set("count", fmt.Sprintf("%d", count))

	endpoint := fmt.Sprintf("%s/api/v1/douyin/web/fetch_video_comments?%s", c.baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("创建 TikHub 评论请求失败: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 TikHub 评论接口网络异常: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取 TikHub 评论响应体失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[TikHub API 警告] 视频 [%s] 评论接口响应非200: %d, 响应: %s", awemeID, resp.StatusCode, string(respBytes))
		return nil, fmt.Errorf("TikHub 获取评论接口异常 (HTTP %d)", resp.StatusCode)
	}

	var resMap map[string]interface{}
	if err := json.Unmarshal(respBytes, &resMap); err != nil {
		return nil, fmt.Errorf("解析 TikHub 评论 JSON 响应失败: %w", err)
	}

	var comments []string
	var commentList []interface{}

	// 稳健支持 TikHub 各版本层级解析: data.comments / data.data.comments / comments
	if dataMap, ok := resMap["data"].(map[string]interface{}); ok {
		if cl, ok := dataMap["comments"].([]interface{}); ok {
			commentList = cl
		} else if innerMap, ok := dataMap["data"].(map[string]interface{}); ok {
			if cl, ok := innerMap["comments"].([]interface{}); ok {
				commentList = cl
			}
		}
	} else if cl, ok := resMap["comments"].([]interface{}); ok {
		commentList = cl
	}

	for _, cItem := range commentList {
		cMap, ok := cItem.(map[string]interface{})
		if !ok {
			continue
		}
		var text string
		if t, ok := cMap["text"].(string); ok {
			text = t
		} else if c, ok := cMap["content"].(string); ok {
			text = c
		}
		text = strings.TrimSpace(text)
		if text != "" {
			comments = append(comments, text)
		}
	}

	log.Printf("[TikHub API] 视频 [%s] 成功解析到 %d 条真实社群评论", awemeID, len(comments))
	return comments, nil
}

