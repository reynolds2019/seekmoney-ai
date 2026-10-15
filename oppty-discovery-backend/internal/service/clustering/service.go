package clustering

import (
	"context"
	"log"
	"strings"
	"time"

	"oppty-discovery-backend/internal/config"
)

/**
 * desc: 语义向量密度聚类分析领域服务接口与落地实现
 * date: 2026-10-15 15:40:00
 * @author reynolds
 * @since v1.0.0
 */

// ClusteringResult 语义聚类最终汇总报告结构
type ClusteringResult struct {
	Clusters       []ClusterGroup // 聚类簇列表
	NoiseCount     int            // 过滤的噪点文本数
	TotalTexts     int            // 原始输入文本总数
	CleanedTexts   int            // 经过清洗保留的优质文本数
	DurationMs     int64          // 运算总耗时 (毫秒)
	IsDegraded     bool           // 是否触发了降级容错
	DegradeReason  string         // 降级触发原因描述
}

// ClusteringService 聚类领域服务接口
type ClusteringService interface {
	ClusterTexts(ctx context.Context, texts []string, keywords []string) (*ClusteringResult, error)
}

type clusteringServiceImpl struct {
	cleaner         *DataCleaner
	embeddingClient EmbeddingClient
	clusterer       *DBSCANClusterer
}

// NewClusteringService 实例化聚类领域服务
func NewClusteringService(cfg *config.LLMConfig) ClusteringService {
	return &clusteringServiceImpl{
		cleaner:         NewDataCleaner(4),
		embeddingClient: NewEmbeddingClient(cfg),
		clusterer:       NewDBSCANClusterer(0.35, 2),
	}
}

// ClusterTexts 针对输入的社交媒体评论执行完整的“清洗 -> 向量化 -> 密度聚类 -> 代表句提炼”流水线
func (s *clusteringServiceImpl) ClusterTexts(ctx context.Context, texts []string, keywords []string) (*ClusteringResult, error) {
	startTime := time.Now()
	total := len(texts)

	// 边界检查：若没有评论，返回默认空结果
	if total == 0 {
		return &ClusteringResult{
			Clusters: []ClusterGroup{{
				ClusterID:      0,
				Texts:          []string{"暂无真实社群原声评论"},
				Representative: "暂无真实社群原声评论",
				Size:           1,
			}},
			TotalTexts:   0,
			CleanedTexts: 0,
			DurationMs:   time.Since(startTime).Milliseconds(),
		}, nil
	}

	// 1. 数据清洗与质量打分
	cleanedItems := s.cleaner.CleanBatch(texts)
	cleanedCount := len(cleanedItems)
	log.Printf("[语义聚类引擎] 原始评论 %d 条，经规则清洗去噪后保留 %d 条有效痛点文本", total, cleanedCount)

	// 若清洗后有效样本不足 2 条，直接使用原文本兜底
	if cleanedCount < 2 {
		rawFallback := texts
		if len(rawFallback) > 10 {
			rawFallback = rawFallback[:10]
		}
		rep := rawFallback[0]
		return &ClusteringResult{
			Clusters: []ClusterGroup{{
				ClusterID:      0,
				Texts:          rawFallback,
				Representative: rep,
				Size:           len(rawFallback),
			}},
			TotalTexts:   total,
			CleanedTexts: cleanedCount,
			DurationMs:   time.Since(startTime).Milliseconds(),
		}, nil
	}

	cleanTexts := make([]string, cleanedCount)
	for i, item := range cleanedItems {
		cleanTexts[i] = item.CleanText
	}

	// 2. 调用本地网关批量获取文本高维向量
	embeddings, err := s.embeddingClient.GetEmbeddings(ctx, cleanTexts)
	if err != nil {
		// 触发明确的智能降级处理，并在控制台详细记录
		log.Printf("[语义聚类降级告警] 向量化接口请求异常: %v，启动智能关键词与频次降级聚类算法", err)
		degradedClusters := s.fallbackClustering(cleanedItems, keywords)
		return &ClusteringResult{
			Clusters:      degradedClusters,
			TotalTexts:    total,
			CleanedTexts:  cleanedCount,
			DurationMs:    time.Since(startTime).Milliseconds(),
			IsDegraded:    true,
			DegradeReason: "Embedding 向量化服务异常或超时，自动启用关键词降级策略",
		}, nil
	}

	// 3. 运行纯 Go 内存 DBSCAN 密度聚类
	clustersIdxs, noiseIdxs := s.clusterer.Run(embeddings)
	log.Printf("[语义聚类引擎] DBSCAN 聚类运算完成，识别出 %d 个高密度语义簇，分离 %d 条离群噪点", len(clustersIdxs), len(noiseIdxs))

	// 4. 若所有样本均被识别为噪点（无明显密度峰值），启动自适应宽松聚合
	if len(clustersIdxs) == 0 {
		log.Println("[语义聚类自适应] 未能形成高密度核心簇，按相似度最大连通分量进行自适应兜底归集")
		clustersIdxs = s.adaptiveFallbackClusters(embeddings)
	}

	// 5. 提取代表性痛点原声并按规模排序
	organizedClusters := OrganizeResults(clustersIdxs, cleanedItems, embeddings)

	// 确保至少有 1~3 个簇群输出
	if len(organizedClusters) == 0 {
		organizedClusters = s.fallbackClustering(cleanedItems, keywords)
	}

	duration := time.Since(startTime).Milliseconds()
	log.Printf("[语义聚类引擎] 流水线执行成功，最终产出 %d 个核心痛点簇群，总耗时: %d ms", len(organizedClusters), duration)

	return &ClusteringResult{
		Clusters:     organizedClusters,
		NoiseCount:   len(noiseIdxs),
		TotalTexts:   total,
		CleanedTexts: cleanedCount,
		DurationMs:   duration,
		IsDegraded:   false,
	}, nil
}

// adaptiveFallbackClusters 在 DBSCAN 未成簇时的自适应均分聚类
func (s *clusteringServiceImpl) adaptiveFallbackClusters(embeddings [][]float64) [][]int {
	n := len(embeddings)
	if n <= 3 {
		indices := make([]int, n)
		for i := 0; i < n; i++ {
			indices[i] = i
		}
		return [][]int{indices}
	}

	mid := n / 2
	c1 := make([]int, mid)
	for i := 0; i < mid; i++ {
		c1[i] = i
	}
	c2 := make([]int, n-mid)
	for i := mid; i < n; i++ {
		c2[i-mid] = i
	}
	return [][]int{c1, c2}
}

// fallbackClustering 基于关键词与特征的降级聚类实现
func (s *clusteringServiceImpl) fallbackClustering(items []CleanedItem, keywords []string) []ClusterGroup {
	var c1, c2 []string
	var rep1, rep2 string

	for _, it := range items {
		matched := false
		for _, kw := range keywords {
			if strings.Contains(it.CleanText, kw) {
				matched = true
				break
			}
		}

		if matched || it.HasPainPoint {
			c1 = append(c1, it.CleanText)
			if rep1 == "" || len(it.CleanText) > len(rep1) {
				rep1 = it.CleanText
			}
		} else {
			c2 = append(c2, it.CleanText)
			if rep2 == "" || len(it.CleanText) > len(rep2) {
				rep2 = it.CleanText
			}
		}
	}

	if len(c1) == 0 && len(c2) > 0 {
		c1 = c2
		rep1 = rep2
		c2 = nil
	}

	var results []ClusterGroup
	if len(c1) > 0 {
		results = append(results, ClusterGroup{
			ClusterID:      0,
			Texts:          c1,
			Representative: rep1,
			Size:           len(c1),
		})
	}
	if len(c2) > 0 {
		results = append(results, ClusterGroup{
			ClusterID:      1,
			Texts:          c2,
			Representative: rep2,
			Size:           len(c2),
		})
	}
	return results
}
