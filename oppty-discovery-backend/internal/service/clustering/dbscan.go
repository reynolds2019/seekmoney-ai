package clustering

import (
	"math"
	"sort"
)

/**
 * desc: 纯 Go 内存余弦距离计算与 DBSCAN 密度聚类算法引擎
 * date: 2026-10-15 15:35:00
 * @author reynolds
 * @since v1.0.0
 */

// CosineDistance 计算两个高维向量间的余弦距离 (1 - 余弦相似度)
// 距离范围在 [0, 2]，数值越小说明语义越接近
func CosineDistance(a, b []float64) float64 {
	if len(a) == 0 || len(b) == 0 || len(a) != len(b) {
		return 1.0
	}

	var dotProduct, normA, normB float64
	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 1.0
	}

	similarity := dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
	// 避免浮点数精度轻微越界
	if similarity > 1.0 {
		similarity = 1.0
	} else if similarity < -1.0 {
		similarity = -1.0
	}

	return 1.0 - similarity
}

type pointStatus int

const (
	statusUnvisited pointStatus = iota
	statusNoise
	statusClustered
)

// DBSCANClusterer 纯 Go 内存密度聚类执行器
type DBSCANClusterer struct {
	Eps     float64 // 邻域半径 (余弦距离阈值，通常在 0.25 ~ 0.40 之间)
	MinPts  int     // 形成高密度簇的最小样本数 (通常 2 ~ 3)
}

// NewDBSCANClusterer 实例化聚类执行器
func NewDBSCANClusterer(eps float64, minPts int) *DBSCANClusterer {
	if eps <= 0 {
		eps = 0.35 // 默认余弦距离半径
	}
	if minPts <= 0 {
		minPts = 2 // 最小 2 个点成簇
	}
	return &DBSCANClusterer{
		Eps:    eps,
		MinPts: minPts,
	}
}

// ClusterGroup 聚类完成后的簇结构
type ClusterGroup struct {
	ClusterID      int          // 簇编号
	Indices        []int        // 样本索引列表
	Texts          []string     // 评论文本列表
	Representative string       // 最具代表性的核心痛点评论原声
	Size           int          // 簇内样本总数
	AvgDistance    float64      // 簇内紧密度 (平均余弦距离)
}

// Run 执行 DBSCAN 密度聚类并返回所有簇的点索引数组及离群噪点索引
func (d *DBSCANClusterer) Run(dataset [][]float64) (clusters [][]int, noise []int) {
	n := len(dataset)
	if n == 0 {
		return [][]int{}, []int{}
	}

	statuses := make([]pointStatus, n)
	clusters = make([][]int, 0)
	noise = make([]int, 0)

	// 预先计算或按需查询邻域
	regionQuery := func(p int) []int {
		var neighbors []int
		for i := 0; i < n; i++ {
			if CosineDistance(dataset[p], dataset[i]) <= d.Eps {
				neighbors = append(neighbors, i)
			}
		}
		return neighbors
	}

	for i := 0; i < n; i++ {
		if statuses[i] != statusUnvisited {
			continue
		}

		neighbors := regionQuery(i)
		if len(neighbors) < d.MinPts {
			// 点密度不足，暂标为噪点
			statuses[i] = statusNoise
		} else {
			// 开辟新簇并进行密度扩散
			currentCluster := []int{i}
			statuses[i] = statusClustered

			// 扩散种子队列
			seeds := make([]int, 0, len(neighbors))
			for _, neighbor := range neighbors {
				if neighbor != i {
					seeds = append(seeds, neighbor)
				}
			}

			for len(seeds) > 0 {
				curr := seeds[0]
				seeds = seeds[1:]

				if statuses[curr] == statusNoise {
					// 之前被标为噪点的边界点，现在可被吸纳为簇成员
					statuses[curr] = statusClustered
					currentCluster = append(currentCluster, curr)
				}

				if statuses[curr] != statusUnvisited {
					continue
				}

				statuses[curr] = statusClustered
				currentCluster = append(currentCluster, curr)

				currNeighbors := regionQuery(curr)
				if len(currNeighbors) >= d.MinPts {
					// 密度足够，将新近邻加入种子探索队列
					for _, cn := range currNeighbors {
						if statuses[cn] == statusUnvisited || statuses[cn] == statusNoise {
							seeds = append(seeds, cn)
						}
					}
				}
			}

			clusters = append(clusters, currentCluster)
		}
	}

	for i := 0; i < n; i++ {
		if statuses[i] == statusNoise {
			noise = append(noise, i)
		}
	}

	return clusters, noise
}

// OrganizeResults 整理聚类簇，提取代表性原声并排序
func OrganizeResults(clusters [][]int, cleanedItems []CleanedItem, embeddings [][]float64) []ClusterGroup {
	var results []ClusterGroup

	for idx, memberIndices := range clusters {
		if len(memberIndices) == 0 {
			continue
		}

		// 1. 计算簇的中心均值向量
		dim := len(embeddings[memberIndices[0]])
		center := make([]float64, dim)
		for _, memberIdx := range memberIndices {
			vec := embeddings[memberIdx]
			for d := 0; d < dim; d++ {
				center[d] += vec[d]
			}
		}
		for d := 0; d < dim; d++ {
			center[d] /= float64(len(memberIndices))
		}

		// 2. 计算各点到中心的距离，挑选最具代表性的评论
		// 规则：综合“距离中心最近”与“包含高频痛点白名单关键词加权”
		var bestScore = -99999.0
		var bestRepresentative = cleanedItems[memberIndices[0]].CleanText
		var totalDistance float64
		texts := make([]string, 0, len(memberIndices))

		for _, memberIdx := range memberIndices {
			item := cleanedItems[memberIdx]
			texts = append(texts, item.CleanText)

			dist := CosineDistance(embeddings[memberIdx], center)
			totalDistance += dist

			// 代表性综合打分：基础质量分 - 距离惩罚（越近分越高）
			repScore := item.Score - (dist * 2.0)
			if repScore > bestScore {
				bestScore = repScore
				bestRepresentative = item.CleanText
			}
		}

		results = append(results, ClusterGroup{
			ClusterID:      idx,
			Indices:        memberIndices,
			Texts:          texts,
			Representative: bestRepresentative,
			Size:           len(texts),
			AvgDistance:    totalDistance / float64(len(memberIndices)),
		})
	}

	// 按包含样本数量从大到小排序，高频核心痛点置顶
	sort.Slice(results, func(i, j int) bool {
		return results[i].Size > results[j].Size
	})

	for i := range results {
		results[i].ClusterID = i
	}

	return results
}
