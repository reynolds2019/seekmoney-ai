package clustering

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"oppty-discovery-backend/internal/config"
)

/**
 * desc: 大模型文本向量化 (Embedding) 客户端封装，支持多批次分片调用
 * date: 2026-10-15 15:32:00
 * @author reynolds
 * @since v1.0.0
 */

// EmbeddingClient 向量化客户端接口
type EmbeddingClient interface {
	GetEmbeddings(ctx context.Context, texts []string) ([][]float64, error)
	GetModelName() string
}

type embeddingClientImpl struct {
	cfg        *config.LLMConfig
	httpClient *http.Client
}

// NewEmbeddingClient 实例化 Embedding 客户端
func NewEmbeddingClient(cfg *config.LLMConfig) EmbeddingClient {
	timeout := 60 * time.Second
	if cfg.TimeoutSeconds > 0 {
		timeout = time.Duration(cfg.TimeoutSeconds) * time.Second
	}
	return &embeddingClientImpl{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *embeddingClientImpl) GetModelName() string {
	if c.cfg.EmbeddingModel != "" {
		return c.cfg.EmbeddingModel
	}
	return "text-embedding-v4"
}

// embeddingReq OpenAI 兼容格式入参
type embeddingReq struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

// embeddingResp OpenAI 兼容格式出参
type embeddingResp struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Index     int       `json:"index"`
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// GetEmbeddings 分批获取文本切片的高维向量表示
func (c *embeddingClientImpl) GetEmbeddings(ctx context.Context, texts []string) ([][]float64, error) {
	if len(texts) == 0 {
		return [][]float64{}, nil
	}

	batchSize := c.cfg.EmbeddingBatchSize
	if batchSize <= 0 {
		batchSize = 10 // 通义千问限制单次最大 10 条
	}

	total := len(texts)
	allEmbeddings := make([][]float64, total)

	apiURL := c.cfg.EmbeddingURL
	if apiURL == "" {
		apiURL = fmt.Sprintf("%s/embeddings", strings.TrimRight(c.cfg.BaseURL, "/"))
	}
	modelName := c.GetModelName()

	log.Printf("[向量化引擎] 开始批量获取 %d 条文本的高维向量，分批大小: %d，调用模型: %s", total, batchSize, modelName)

	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}
		subTexts := texts[i:end]

		reqPayload := embeddingReq{
			Model: modelName,
			Input: subTexts,
		}
		payloadBytes, err := json.Marshal(reqPayload)
		if err != nil {
			return nil, fmt.Errorf("序列化向量化请求参数失败: %w", err)
		}

		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewBuffer(payloadBytes))
		if err != nil {
			return nil, fmt.Errorf("构建 HTTP 请求失败: %w", err)
		}
		httpReq.Header.Set("Content-Type", "application/json")
		if c.cfg.APIKey != "" {
			httpReq.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
		}

		resp, err := c.httpClient.Do(httpReq)
		if err != nil {
			return nil, fmt.Errorf("调用向量化网关网络超时或失败: %w", err)
		}

		respBody, readErr := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if readErr != nil {
			return nil, fmt.Errorf("读取向量化响应数据体失败: %w", readErr)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("向量化网关返回非 200 状态码 [%d]: %s", resp.StatusCode, string(respBody))
		}

		var parsedResp embeddingResp
		if err := json.Unmarshal(respBody, &parsedResp); err != nil {
			return nil, fmt.Errorf("反序列化向量化响应失败: %w，原始报文: %s", err, string(respBody))
		}

		if parsedResp.Error != nil && parsedResp.Error.Message != "" {
			return nil, fmt.Errorf("向量化模型返回错误: %s (code: %s)", parsedResp.Error.Message, parsedResp.Error.Code)
		}

		if len(parsedResp.Data) != len(subTexts) {
			return nil, fmt.Errorf("向量化返回结果数量不匹配，输入 %d 条，返回 %d 条", len(subTexts), len(parsedResp.Data))
		}

		for idx, item := range parsedResp.Data {
			targetIdx := i + idx
			allEmbeddings[targetIdx] = item.Embedding
		}
	}

	log.Printf("[向量化引擎] 全部 %d 条文本向量提取完毕，单向量维度: %d", total, len(allEmbeddings[0]))
	return allEmbeddings, nil
}
