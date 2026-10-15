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
	FetchPlatformData(ctx context.Context, platform string, keyword string, maxItems int, maxCommentsPerItem int) ([]TikHubVideoItem, []string, error)
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

/**
 * desc: 通用全平台真实社交媒体数据与原声评论采集调度分发器
 * date: 2026-10-15 14:05:00
 * @author reynolds
 * @since v1.0.0
 * @param ctx 上下文
 * @param platform 数据源平台标识 (douyin, tiktok, xiaohongshu, bilibili, wechat, youtube)
 * @param keyword 搜索关键词
 * @param maxItems 最大抓取条目数
 * @param maxCommentsPerItem 每条内容抓取评论上限
 * @return []TikHubVideoItem 真实内容实体切片
 * @return []string 真实原声评论文本切片
 * @return error 异常信息
 */
func (c *tikHubClientImpl) FetchPlatformData(ctx context.Context, platform string, keyword string, maxItems int, maxCommentsPerItem int) ([]TikHubVideoItem, []string, error) {
	if maxItems <= 0 {
		maxItems = 5
	}
	if maxCommentsPerItem <= 0 {
		maxCommentsPerItem = 10
	}

	var items []TikHubVideoItem
	var comments []string
	var err error

	log.Printf("[全域采集引擎] 正在针对平台 [%s] 发起关键词 [%s] 的真实数据采集 (限制内容: %d, 单项评论上限: %d)", platform, keyword, maxItems, maxCommentsPerItem)

	switch platform {
	case "xiaohongshu":
		searchLimit := maxItems * 3
		if searchLimit < 15 {
			searchLimit = 15
		}
		candidateItems, sErr := c.searchXiaohongshu(ctx, keyword, searchLimit)
		err = sErr
		if err == nil && len(candidateItems) > 0 {
			var selectedItems []TikHubVideoItem
			for _, it := range candidateItems {
				if it.AwemeID != "" {
					cmts, cErr := c.getXiaohongshuComments(ctx, it.AwemeID, maxCommentsPerItem)
					if cErr == nil && len(cmts) > 0 {
						comments = append(comments, cmts...)
					} else if cErr != nil {
						log.Printf("[TikHub 小红书告警] 笔记 [%s] 评论获取失败: %v", it.AwemeID, cErr)
					}
				}
				if len(selectedItems) < maxItems {
					selectedItems = append(selectedItems, it)
				}
				// 若已收集到足够的真实评论且内容卡片数已满足用户要求，可提前返回
				if len(selectedItems) >= maxItems && len(comments) >= maxCommentsPerItem*2 {
					break
				}
			}
			items = selectedItems
		}
	case "bilibili":
		items, err = c.searchBilibili(ctx, keyword, maxItems)
		if err == nil && len(items) > 0 {
			for _, it := range items {
				if it.AwemeID != "" {
					cmts, cErr := c.getBilibiliComments(ctx, it.AwemeID, maxCommentsPerItem)
					if cErr == nil && len(cmts) > 0 {
						comments = append(comments, cmts...)
					}
				}
			}
		}
	case "wechat":
		items, err = c.searchWeChat(ctx, keyword, maxItems)
		if err == nil && len(items) > 0 {
			for _, it := range items {
				if it.AwemeID != "" {
					cmts, cErr := c.getWeChatComments(ctx, it.AwemeID, maxCommentsPerItem)
					if cErr == nil && len(cmts) > 0 {
						comments = append(comments, cmts...)
					}
				}
			}
		}
	case "youtube":
		items, err = c.searchYouTube(ctx, keyword, maxItems)
		if err == nil && len(items) > 0 {
			for _, it := range items {
				if it.AwemeID != "" {
					cmts, cErr := c.getYouTubeComments(ctx, it.AwemeID, maxCommentsPerItem)
					if cErr == nil && len(cmts) > 0 {
						comments = append(comments, cmts...)
					}
				}
			}
		}
	default: // 抖音 / TikTok 及默认短视频
		items, err = c.SearchVideos(ctx, keyword, maxItems)
		if err == nil && len(items) > 0 {
			for _, it := range items {
				if it.AwemeID != "" {
					cmts, cErr := c.GetVideoComments(ctx, it.AwemeID, maxCommentsPerItem)
					if cErr == nil && len(cmts) > 0 {
						comments = append(comments, cmts...)
					}
				}
			}
		}
	}

	if err != nil {
		log.Printf("[全域采集引擎告警] 平台 [%s] 真实采集遇到异常: %v", platform, err)
		return nil, nil, err
	}

	log.Printf("[全域采集引擎就绪] 平台 [%s] 成功获取 %d 条真实内容与 %d 条社群原声评论", platform, len(items), len(comments))
	return items, comments, nil
}

// searchXiaohongshu 调用 TikHub 搜索小红书笔记
func (c *tikHubClientImpl) searchXiaohongshu(ctx context.Context, keyword string, limit int) ([]TikHubVideoItem, error) {
	endpoint := fmt.Sprintf("%s/api/v1/xiaohongshu/app_v2/search_notes?keyword=%s&page=1", c.baseURL, url.QueryEscape(keyword))
	log.Printf("[TikHub 小红书] 发起笔记搜索: [%s], 端点: %s", keyword, endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("创建小红书搜索请求失败: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求小红书搜索接口网络异常: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取小红书响应体失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[TikHub 小红书 API 告警] 搜索接口响应非200: %d, 响应: %s", resp.StatusCode, string(respBytes))
		return nil, fmt.Errorf("小红书搜索接口响应状态码非 200: %d, 详情: %s", resp.StatusCode, string(respBytes))
	}

	var resMap map[string]interface{}
	if err := json.Unmarshal(respBytes, &resMap); err != nil {
		return nil, fmt.Errorf("解析小红书 JSON 响应失败: %w", err)
	}

	var rawItems []interface{}
	if dataMap, ok := resMap["data"].(map[string]interface{}); ok {
		if innerData, ok := dataMap["data"].(map[string]interface{}); ok {
			if items, ok := innerData["items"].([]interface{}); ok {
				rawItems = items
			}
		} else if items, ok := dataMap["items"].([]interface{}); ok {
			rawItems = items
		}
	} else if items, ok := resMap["items"].([]interface{}); ok {
		rawItems = items
	}

	var results []TikHubVideoItem
	for _, it := range rawItems {
		itemMap, ok := it.(map[string]interface{})
		if !ok {
			continue
		}

		// 适配 note_card 或 note 嵌套结构
		targetMap := itemMap
		if nc, ok := itemMap["note_card"].(map[string]interface{}); ok {
			targetMap = nc
		} else if n, ok := itemMap["note"].(map[string]interface{}); ok {
			targetMap = n
		}

		noteID, _ := targetMap["id"].(string)
		if noteID == "" {
			noteID, _ = targetMap["note_id"].(string)
		}
		if noteID == "" {
			noteID, _ = itemMap["id"].(string)
		}

		title, _ := targetMap["display_title"].(string)
		if title == "" {
			title, _ = targetMap["title"].(string)
		}
		if title == "" {
			title, _ = targetMap["desc"].(string)
		}

		nickname := "小红书博主"
		if userMap, ok := targetMap["user"].(map[string]interface{}); ok {
			if n, ok := userMap["nickname"].(string); ok && n != "" {
				nickname = n
			}
		} else if n, ok := targetMap["nickname"].(string); ok && n != "" {
			nickname = n
		}

		likes := "100"
		commentCount := 10
		if interact, ok := targetMap["interact_info"].(map[string]interface{}); ok {
			if lStr, ok := interact["liked_count"].(string); ok && lStr != "" {
				likes = lStr
			} else if lNum, ok := interact["liked_count"].(float64); ok {
				likes = fmt.Sprintf("%.0f", lNum)
			}
			if cNum, ok := interact["comment_count"].(float64); ok {
				commentCount = int(cNum)
			}
		} else if lNum, ok := targetMap["liked_count"].(float64); ok {
			likes = fmt.Sprintf("%.0f", lNum)
		} else if lStr, ok := targetMap["liked_count"].(string); ok && lStr != "" {
			likes = lStr
		}

		shareURL := fmt.Sprintf("https://www.xiaohongshu.com/explore/%s", noteID)

		if noteID != "" || title != "" {
			results = append(results, TikHubVideoItem{
				AwemeID:      noteID,
				Title:        title,
				Author:       nickname,
				Likes:        likes,
				ShareURL:     shareURL,
				CommentCount: commentCount,
			})
		}

		if len(results) >= limit {
			break
		}
	}

	log.Printf("[TikHub 小红书] 成功解析到 %d 篇真实笔记", len(results))
	return results, nil
}

// getXiaohongshuComments 调用 TikHub 获取小红书笔记评论
func (c *tikHubClientImpl) getXiaohongshuComments(ctx context.Context, noteID string, count int) ([]string, error) {
	if noteID == "" {
		return nil, nil
	}
	endpoint := fmt.Sprintf("%s/api/v1/xiaohongshu/app_v2/get_note_comments?note_id=%s", c.baseURL, url.QueryEscape(noteID))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("[TikHub 小红书 API 告警] 笔记 [%s] 评论接口响应非200: %d, 响应: %s", noteID, resp.StatusCode, string(respBytes))
		return nil, fmt.Errorf("读取小红书评论异常 (HTTP %d): %s", resp.StatusCode, string(respBytes))
	}

	var resMap map[string]interface{}
	if err := json.Unmarshal(respBytes, &resMap); err != nil {
		return nil, err
	}

	var comments []string
	var commentList []interface{}

	// TikHub 小红书真实接口结构为: resMap["data"]["data"]["comments"]，同时向下兼容 resMap["data"]["comments"]
	if dataMap, ok := resMap["data"].(map[string]interface{}); ok {
		if innerData, ok := dataMap["data"].(map[string]interface{}); ok {
			if cl, ok := innerData["comments"].([]interface{}); ok {
				commentList = cl
			}
		} else if cl, ok := dataMap["comments"].([]interface{}); ok {
			commentList = cl
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
		if t, ok := cMap["content"].(string); ok {
			text = t
		} else if t, ok := cMap["text"].(string); ok {
			text = t
		}
		text = strings.TrimSpace(text)
		if text != "" {
			comments = append(comments, text)
		}

		// 递归提取小红书楼中楼/子评论 (sub_comments)
		if subList, ok := cMap["sub_comments"].([]interface{}); ok {
			for _, subItem := range subList {
				if subMap, ok := subItem.(map[string]interface{}); ok {
					var subText string
					if st, ok := subMap["content"].(string); ok {
						subText = st
					} else if st, ok := subMap["text"].(string); ok {
						subText = st
					}
					subText = strings.TrimSpace(subText)
					if subText != "" {
						comments = append(comments, subText)
					}
				}
			}
		}

		if len(comments) >= count {
			break
		}
	}

	log.Printf("[TikHub 小红书] 笔记 [%s] 成功解析到 %d 条真实评论", noteID, len(comments))
	return comments, nil
}

// searchBilibili 调用 TikHub 搜索 Bilibili 视频
func (c *tikHubClientImpl) searchBilibili(ctx context.Context, keyword string, limit int) ([]TikHubVideoItem, error) {
	endpoint := fmt.Sprintf("%s/api/v1/bilibili/web/fetch_general_search?keyword=%s", c.baseURL, url.QueryEscape(keyword))
	log.Printf("[TikHub B站] 发起视频搜索: [%s], 端点: %s", keyword, endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("B站搜索接口异常 (HTTP %d)", resp.StatusCode)
	}

	var resMap map[string]interface{}
	if err := json.Unmarshal(respBytes, &resMap); err != nil {
		return nil, err
	}

	var rawResults []interface{}
	if dataMap, ok := resMap["data"].(map[string]interface{}); ok {
		if innerData, ok := dataMap["data"].(map[string]interface{}); ok {
			if rArr, ok := innerData["result"].([]interface{}); ok {
				rawResults = rArr
			}
		} else if rArr, ok := dataMap["result"].([]interface{}); ok {
			rawResults = rArr
		}
	} else if rArr, ok := resMap["result"].([]interface{}); ok {
		rawResults = rArr
	}

	// 扁平化展开可能是按分类聚合的分组数据 (result_type: video)
	var videoItems []map[string]interface{}
	for _, it := range rawResults {
		itemMap, ok := it.(map[string]interface{})
		if !ok {
			continue
		}
		// 若为分类聚合项: {"result_type": "video", "data": [...]}
		if rType, ok := itemMap["result_type"].(string); ok && rType == "video" {
			if subData, ok := itemMap["data"].([]interface{}); ok {
				for _, subIt := range subData {
					if subMap, ok := subIt.(map[string]interface{}); ok {
						videoItems = append(videoItems, subMap)
					}
				}
				continue
			}
		}
		// 若为直接视频项或未细分项
		videoItems = append(videoItems, itemMap)
	}

	var results []TikHubVideoItem
	for _, itemMap := range videoItems {
		bvid, _ := itemMap["bvid"].(string)
		aidStr := ""
		if aidNum, ok := itemMap["aid"].(float64); ok {
			aidStr = fmt.Sprintf("%.0f", aidNum)
		} else if idNum, ok := itemMap["id"].(float64); ok {
			aidStr = fmt.Sprintf("%.0f", idNum)
		} else if idStr, ok := itemMap["id"].(string); ok {
			aidStr = idStr
		}

		targetID := bvid
		if targetID == "" {
			targetID = aidStr
		}

		title, _ := itemMap["title"].(string)
		// 清理可能包含的 <em> 与 HTML 标签
		title = strings.ReplaceAll(title, "<em class=\"keyword\">", "")
		title = strings.ReplaceAll(title, "</em>", "")

		author, _ := itemMap["author"].(string)
		if author == "" {
			author = "B站UP主"
		}

		likes := "100"
		if likeNum, ok := itemMap["like"].(float64); ok {
			likes = fmt.Sprintf("%.0f", likeNum)
		} else if favNum, ok := itemMap["favorites"].(float64); ok {
			likes = fmt.Sprintf("%.0f", favNum)
		}

		shareURL, _ := itemMap["arcurl"].(string)
		if shareURL == "" && bvid != "" {
			shareURL = fmt.Sprintf("https://www.bilibili.com/video/%s", bvid)
		}

		if targetID != "" || title != "" {
			results = append(results, TikHubVideoItem{
				AwemeID:      targetID,
				Title:        title,
				Author:       author,
				Likes:        likes,
				ShareURL:     shareURL,
				CommentCount: 20,
			})
		}

		if len(results) >= limit {
			break
		}
	}

	log.Printf("[TikHub B站] 成功解析到 %d 个真实视频", len(results))
	return results, nil
}

// getBilibiliComments 调用 TikHub 获取 B站视频评论
func (c *tikHubClientImpl) getBilibiliComments(ctx context.Context, aid string, count int) ([]string, error) {
	if aid == "" {
		return nil, nil
	}
	endpoint := fmt.Sprintf("%s/api/v1/bilibili/app/fetch_video_comments?av_id=%s&mode=3", c.baseURL, url.QueryEscape(aid))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("读取B站评论异常 (HTTP %d)", resp.StatusCode)
	}

	var resMap map[string]interface{}
	if err := json.Unmarshal(respBytes, &resMap); err != nil {
		return nil, err
	}

	var comments []string
	var replies []interface{}
	if dataMap, ok := resMap["data"].(map[string]interface{}); ok {
		if innerData, ok := dataMap["data"].(map[string]interface{}); ok {
			if r, ok := innerData["replies"].([]interface{}); ok {
				replies = r
			}
		} else if r, ok := dataMap["replies"].([]interface{}); ok {
			replies = r
		}
	} else if r, ok := resMap["replies"].([]interface{}); ok {
		replies = r
	}

	for _, rItem := range replies {
		rMap, ok := rItem.(map[string]interface{})
		if !ok {
			continue
		}
		if cObj, ok := rMap["content"].(map[string]interface{}); ok {
			if msg, ok := cObj["message"].(string); ok && strings.TrimSpace(msg) != "" {
				comments = append(comments, strings.TrimSpace(msg))
			}
		}
		if len(comments) >= count {
			break
		}
	}

	log.Printf("[TikHub B站] 视频 [%s] 成功解析到 %d 条真实评论", aid, len(comments))
	return comments, nil
}

// searchWeChat 调用 TikHub 搜索微信视频号
func (c *tikHubClientImpl) searchWeChat(ctx context.Context, keyword string, limit int) ([]TikHubVideoItem, error) {
	endpoint := fmt.Sprintf("%s/api/v1/wechat_search/v2/fetch_search_videos", c.baseURL)
	log.Printf("[TikHub 视频号] 发起搜索: [%s], 端点: %s", keyword, endpoint)

	payloadBytes, _ := json.Marshal(map[string]string{"keyword": keyword})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("微信视频号搜索异常 (HTTP %d)", resp.StatusCode)
	}

	var resMap map[string]interface{}
	if err := json.Unmarshal(respBytes, &resMap); err != nil {
		return nil, err
	}

	var rawVideos []interface{}
	if dataMap, ok := resMap["data"].(map[string]interface{}); ok {
		if resultsObj, ok := dataMap["results"].(map[string]interface{}); ok {
			if vArr, ok := resultsObj["data"].([]interface{}); ok {
				for _, box := range vArr {
					if bMap, ok := box.(map[string]interface{}); ok {
						if subBoxes, ok := bMap["subBoxes"].([]interface{}); ok {
							for _, sb := range subBoxes {
								if sbMap, ok := sb.(map[string]interface{}); ok {
									if items, ok := sbMap["items"].([]interface{}); ok {
										rawVideos = append(rawVideos, items...)
									}
								}
							}
						} else {
							rawVideos = append(rawVideos, box)
						}
					}
				}
			}
		} else if items, ok := dataMap["items"].([]interface{}); ok {
			rawVideos = items
		} else if videos, ok := dataMap["videos"].([]interface{}); ok {
			rawVideos = videos
		} else if innerData, ok := dataMap["data"].([]interface{}); ok {
			rawVideos = innerData
		}
	} else if items, ok := resMap["items"].([]interface{}); ok {
		rawVideos = items
	}

	var results []TikHubVideoItem
	for _, it := range rawVideos {
		vMap, ok := it.(map[string]interface{})
		if !ok {
			continue
		}

		vid, _ := vMap["object_id"].(string)
		if vid == "" {
			vid, _ = vMap["docID"].(string)
		}
		if vid == "" {
			vid, _ = vMap["id"].(string)
		}
		if vid == "" {
			vid, _ = vMap["hashDocID"].(string)
		}

		desc, _ := vMap["description"].(string)
		if desc == "" {
			desc, _ = vMap["title"].(string)
		}
		if desc == "" {
			desc, _ = vMap["desc"].(string)
		}

		nickname, _ := vMap["nickname"].(string)
		if nickname == "" {
			nickname, _ = vMap["nickName"].(string)
		}
		if nickname == "" {
			nickname, _ = vMap["author"].(string)
		}
		if nickname == "" {
			nickname = "视频号创作者"
		}

		likes := "100"
		if lNum, ok := vMap["like_count"].(float64); ok {
			likes = fmt.Sprintf("%.0f", lNum)
		} else if lNum, ok := vMap["likeNum"].(float64); ok {
			likes = fmt.Sprintf("%.0f", lNum)
		} else if lStr, ok := vMap["likeNum"].(string); ok && lStr != "" {
			likes = lStr
		} else if fNum, ok := vMap["fav_count"].(float64); ok {
			likes = fmt.Sprintf("%.0f", fNum)
		}

		if vid != "" || desc != "" {
			results = append(results, TikHubVideoItem{
				AwemeID:      vid,
				Title:        desc,
				Author:       nickname,
				Likes:        likes,
				ShareURL:     "https://channels.weixin.qq.com",
				CommentCount: 15,
			})
		}
		if len(results) >= limit {
			break
		}
	}

	log.Printf("[TikHub 视频号] 成功解析到 %d 条真实视频", len(results))
	return results, nil
}

// getWeChatComments 调用 TikHub 获取微信视频号评论
func (c *tikHubClientImpl) getWeChatComments(ctx context.Context, videoID string, count int) ([]string, error) {
	if videoID == "" {
		return nil, nil
	}
	endpoint := fmt.Sprintf("%s/api/v1/wechat_channels/v2/fetch_video_comments?id=%s", c.baseURL, url.QueryEscape(videoID))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("读取视频号评论异常 (HTTP %d)", resp.StatusCode)
	}

	var resMap map[string]interface{}
	if err := json.Unmarshal(respBytes, &resMap); err != nil {
		return nil, err
	}

	var comments []string
	var commentInfo []interface{}
	if dataMap, ok := resMap["data"].(map[string]interface{}); ok {
		if ci, ok := dataMap["comment_info"].([]interface{}); ok {
			commentInfo = ci
		} else if cl, ok := dataMap["comments"].([]interface{}); ok {
			commentInfo = cl
		}
	} else if cl, ok := resMap["comments"].([]interface{}); ok {
		commentInfo = cl
	}

	for _, cItem := range commentInfo {
		cMap, ok := cItem.(map[string]interface{})
		if !ok {
			continue
		}
		var text string
		if t, ok := cMap["content"].(string); ok {
			text = t
		} else if t, ok := cMap["text"].(string); ok {
			text = t
		}
		text = strings.TrimSpace(text)
		if text != "" {
			comments = append(comments, text)
		}
		if len(comments) >= count {
			break
		}
	}

	log.Printf("[TikHub 视频号] 视频 [%s] 成功解析到 %d 条真实评论", videoID, len(comments))
	return comments, nil
}

// searchYouTube 调用 TikHub 搜索 YouTube 视频
func (c *tikHubClientImpl) searchYouTube(ctx context.Context, keyword string, limit int) ([]TikHubVideoItem, error) {
	escapedKW := url.QueryEscape(keyword)
	endpoint := fmt.Sprintf("%s/api/v1/youtube/web/search_video?search_query=%s&keyword=%s&need_format=true", c.baseURL, escapedKW, escapedKW)
	log.Printf("[TikHub YouTube] 发起视频搜索: [%s], 端点: %s", keyword, endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("YouTube搜索接口异常 (HTTP %d)", resp.StatusCode)
	}

	var resMap map[string]interface{}
	if err := json.Unmarshal(respBytes, &resMap); err != nil {
		return nil, err
	}

	var videoArr []interface{}
	if dataMap, ok := resMap["data"].(map[string]interface{}); ok {
		if v, ok := dataMap["videos"].([]interface{}); ok {
			videoArr = v
		} else if cList, ok := dataMap["contents"].([]interface{}); ok {
			videoArr = cList
		}
	} else if v, ok := resMap["videos"].([]interface{}); ok {
		videoArr = v
	} else if cList, ok := resMap["contents"].([]interface{}); ok {
		videoArr = cList
	}

	var results []TikHubVideoItem
	for _, it := range videoArr {
		vMap, ok := it.(map[string]interface{})
		if !ok {
			continue
		}

		vid, _ := vMap["video_id"].(string)
		if vid == "" {
			vid, _ = vMap["id"].(string)
		}

		title, _ := vMap["title"].(string)
		channel, _ := vMap["channel_name"].(string)
		if channel == "" {
			channel, _ = vMap["author"].(string)
		}
		if channel == "" {
			channel, _ = vMap["channel_title"].(string)
		}

		// 适配原生渲染器结构 (videoRenderer)
		if renderer, ok := vMap["videoRenderer"].(map[string]interface{}); ok {
			if rVid, ok := renderer["videoId"].(string); ok && rVid != "" {
				vid = rVid
			}
			if titleObj, ok := renderer["title"].(map[string]interface{}); ok {
				if runs, ok := titleObj["runs"].([]interface{}); ok && len(runs) > 0 {
					if rMap, ok := runs[0].(map[string]interface{}); ok {
						title, _ = rMap["text"].(string)
					}
				}
			}
			if ownerObj, ok := renderer["ownerText"].(map[string]interface{}); ok {
				if runs, ok := ownerObj["runs"].([]interface{}); ok && len(runs) > 0 {
					if rMap, ok := runs[0].(map[string]interface{}); ok {
						channel, _ = rMap["text"].(string)
					}
				}
			}
		}

		if channel == "" {
			channel = "YouTube Creator"
		}

		vURL, _ := vMap["video_url"].(string)
		if vURL == "" && vid != "" {
			vURL = fmt.Sprintf("https://www.youtube.com/watch?v=%s", vid)
		}

		if vid != "" || title != "" {
			results = append(results, TikHubVideoItem{
				AwemeID:      vid,
				Title:        title,
				Author:       channel,
				Likes:        "1000",
				ShareURL:     vURL,
				CommentCount: 50,
			})
		}
		if len(results) >= limit {
			break
		}
	}

	log.Printf("[TikHub YouTube] 成功解析到 %d 个真实视频", len(results))
	return results, nil
}

// getYouTubeComments 调用 TikHub 获取 YouTube 视频评论
func (c *tikHubClientImpl) getYouTubeComments(ctx context.Context, videoID string, count int) ([]string, error) {
	if videoID == "" {
		return nil, nil
	}
	endpoint := fmt.Sprintf("%s/api/v1/youtube/web/get_video_comments?video_id=%s&need_format=true", c.baseURL, url.QueryEscape(videoID))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("读取YouTube评论异常 (HTTP %d)", resp.StatusCode)
	}

	var resMap map[string]interface{}
	if err := json.Unmarshal(respBytes, &resMap); err != nil {
		return nil, err
	}

	var comments []string
	var commentList []interface{}
	if dataMap, ok := resMap["data"].(map[string]interface{}); ok {
		if fmtData, ok := dataMap["formatted_data"].(map[string]interface{}); ok {
			if cl, ok := fmtData["comments"].([]interface{}); ok {
				commentList = cl
			}
		} else if cl, ok := dataMap["comments"].([]interface{}); ok {
			commentList = cl
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
		} else if t, ok := cMap["content"].(string); ok {
			text = t
		}
		text = strings.TrimSpace(text)
		if text != "" {
			comments = append(comments, text)
		}
		if len(comments) >= count {
			break
		}
	}

	log.Printf("[TikHub YouTube] 视频 [%s] 成功解析到 %d 条真实评论", videoID, len(comments))
	return comments, nil
}


