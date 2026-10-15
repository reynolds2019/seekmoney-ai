package clustering_test

import (
	"math"
	"testing"

	"oppty-discovery-backend/internal/service/clustering"
)

/**
 * desc: 余弦距离与纯 Go 内存 DBSCAN 密度聚类算法单元测试
 * date: 2026-10-15 15:47:00
 * @author reynolds
 * @since v1.0.0
 */

// 测试余弦距离数学计算精度
func TestCosineDistance(t *testing.T) {
	t.Run("测试完全相同向量余弦距离应为0", func(t *testing.T) {
		v1 := []float64{1.0, 2.0, 3.0}
		v2 := []float64{1.0, 2.0, 3.0}
		dist := clustering.CosineDistance(v1, v2)
		if math.Abs(dist) > 1e-6 {
			t.Errorf("相同向量距离期望为 0，实际为: %f", dist)
		}
	})

	t.Run("测试完全正交向量余弦距离应为1", func(t *testing.T) {
		v1 := []float64{1.0, 0.0}
		v2 := []float64{0.0, 1.0}
		dist := clustering.CosineDistance(v1, v2)
		if math.Abs(dist-1.0) > 1e-6 {
			t.Errorf("正交向量距离期望为 1.0，实际为: %f", dist)
		}
	})

	t.Run("测试完全相反向量余弦距离应为2", func(t *testing.T) {
		v1 := []float64{1.0, 1.0}
		v2 := []float64{-1.0, -1.0}
		dist := clustering.CosineDistance(v1, v2)
		if math.Abs(dist-2.0) > 1e-6 {
			t.Errorf("相反向量距离期望为 2.0，实际为: %f", dist)
		}
	})
}

// 测试纯内存 DBSCAN 聚类算法执行
func TestDBSCANClusterer_Run(t *testing.T) {
	t.Run("测试高维聚类点群能够准确形成独立簇群并滤除孤立噪点", func(t *testing.T) {
		// 构造两簇明显聚集的样本：
		// 簇A：围绕 (1, 1, 0) 微小扰动
		// 簇B：围绕 (0, 0, 1) 微小扰动
		// 噪点C：孤立点 (-1, 0, 0)
		dataset := [][]float64{
			{1.0, 1.0, 0.0},
			{1.0, 0.95, 0.05},
			{0.95, 1.0, 0.02}, // 簇A成员 (3个)

			{0.0, 0.02, 1.0},
			{0.01, 0.0, 0.98},
			{0.0, 0.01, 1.02}, // 簇B成员 (3个)

			{-1.0, 0.0, 0.0},  // 噪点
		}

		clusterer := clustering.NewDBSCANClusterer(0.2, 2)
		clusters, noise := clusterer.Run(dataset)

		if len(clusters) != 2 {
			t.Fatalf("期望识别出 2 个聚类簇，实际识别出: %d", len(clusters))
		}

		if len(noise) != 1 || noise[0] != 6 {
			t.Errorf("期望正确分离 1 个噪点(索引6)，实际噪点索引: %v", noise)
		}
	})
}

// 测试簇整理与核心原声代表句提取
func TestOrganizeResults(t *testing.T) {
	t.Run("测试簇内代表性评论提取与规模排序", func(t *testing.T) {
		clusters := [][]int{
			{0, 1},
			{2, 3, 4},
		}

		cleanedItems := []clustering.CleanedItem{
			{CleanText: "短评A", Score: 1.0},
			{CleanText: "短评B", Score: 1.5},
			{CleanText: "痛点评论C：续航太短，充不进电，难用", Score: 4.0, HasPainPoint: true},
			{CleanText: "痛点评论D：手机发烫死机", Score: 3.5, HasPainPoint: true},
			{CleanText: "痛点评论E：发热严重", Score: 3.0, HasPainPoint: true},
		}

		embeddings := [][]float64{
			{1.0, 0.0},
			{0.98, 0.05},
			{0.0, 1.0},
			{0.02, 0.98},
			{0.01, 0.99},
		}

		results := clustering.OrganizeResults(clusters, cleanedItems, embeddings)

		if len(results) != 2 {
			t.Fatalf("整理结果期望 2 个簇，实际为: %d", len(results))
		}

		// 排序验证：包含 3 个成员的簇应该排在第 1 位
		if results[0].Size != 3 {
			t.Errorf("规模最大的簇期望排在首位(Size=3)，实际 Size 为: %d", results[0].Size)
		}

		// 代表句验证：簇内得分最高且在中心的评论应作为代表句
		if results[0].Representative != "痛点评论C：续航太短，充不进电，难用" {
			t.Errorf("代表句提取偏差，实际为: %s", results[0].Representative)
		}
	})
}
