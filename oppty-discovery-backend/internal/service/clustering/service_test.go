package clustering_test

import (
	"context"
	"testing"

	"oppty-discovery-backend/internal/config"
	"oppty-discovery-backend/internal/service/clustering"
)

/**
 * desc: 语义聚类领域服务整体流程与降级测试
 * date: 2026-10-15 15:50:00
 * @author reynolds
 * @since v1.0.0
 */

func TestClusteringService_ClusterTexts(t *testing.T) {
	ctx := context.Background()

	t.Run("测试空评论集合平滑返回默认簇", func(t *testing.T) {
		svc := clustering.NewClusteringService(&config.LLMConfig{})
		res, err := svc.ClusterTexts(ctx, []string{}, []string{"测试"})
		if err != nil {
			t.Fatalf("空评论输入期望不返回错误，但返回: %v", err)
		}
		if len(res.Clusters) != 1 {
			t.Fatalf("期望返回 1 个保底簇，实际为: %d", len(res.Clusters))
		}
	})

	t.Run("测试不可达Embedding地址时自动平滑降级", func(t *testing.T) {
		// 配置一个不可达的 Embedding 地址，验证是否平滑降级而不抛出 panic
		badCfg := &config.LLMConfig{
			BaseURL:        "http://127.0.0.1:59999/v1",
			EmbeddingURL:   "http://127.0.0.1:59999/v1/embeddings",
			TimeoutSeconds: 1,
		}
		svc := clustering.NewClusteringService(badCfg)

		sampleTexts := []string{
			"电池续航太差了，经常闪退卡顿报错",
			"软件怎么导出超高清视频？一直失败",
			"市面上课程价格太贵了，都是割韭菜",
		}

		res, err := svc.ClusterTexts(ctx, sampleTexts, []string{"视频", "电池"})
		if err != nil {
			t.Fatalf("降级分支期望不向上抛出错误，但返回: %v", err)
		}

		if !res.IsDegraded {
			t.Errorf("期望标记为降级状态 (IsDegraded=true)")
		}

		if len(res.Clusters) == 0 {
			t.Errorf("降级后期望产出至少 1 个聚类簇")
		}

		t.Logf("降级成功，降级原因: %s，产出簇数: %d", res.DegradeReason, len(res.Clusters))
	})
}
