package clustering_test

import (
	"testing"

	"oppty-discovery-backend/internal/service/clustering"
)

/**
 * desc: 评论文本数据清洗与痛点打分单元测试
 * date: 2026-10-15 15:45:00
 * @author reynolds
 * @since v1.0.0
 */

// 测试单条文本噪音识别功能
func TestDataCleaner_IsNoise(t *testing.T) {
	cleaner := clustering.NewDataCleaner(4)

	t.Run("测试极短文本与纯数字应判定为噪音", func(t *testing.T) {
		shortTexts := []string{"", " ", "好", "对", "666", "12345", "8888"}
		for _, text := range shortTexts {
			if !cleaner.IsNoise(text) {
				t.Errorf("文本 [%s] 期望被判定为噪音，但判定为非噪音", text)
			}
		}
	})

	t.Run("测试纯标点与重复语气词应判定为噪音", func(t *testing.T) {
		punctTexts := []string{"????", "！！！", "。。。", "哈哈哈哈哈哈", "啊啊啊啊"}
		for _, text := range punctTexts {
			if !cleaner.IsNoise(text) {
				t.Errorf("语气词标点 [%s] 期望被判定为噪音，但判定为非噪音", text)
			}
		}
	})

	t.Run("测试高频客套口头禅应判定为噪音", func(t *testing.T) {
		noisePhrases := []string{"谢谢", "感谢", "点赞", "收藏", "已关", "加微信", "牛逼", "太强了"}
		for _, phrase := range noisePhrases {
			if !cleaner.IsNoise(phrase) {
				t.Errorf("黑名单短语 [%s] 期望被判定为噪音，但判定为非噪音", phrase)
			}
		}
	})

	t.Run("测试真实高质量痛点评论不应判定为噪音", func(t *testing.T) {
		validTexts := []string{
			"这个软件怎么导出高清视频？一直闪退报错！",
			"电池掉电太快了，刚买一个月就发烫严重",
			"市面上课程太贵了，都是割韭菜，根本学不会",
		}
		for _, text := range validTexts {
			if cleaner.IsNoise(text) {
				t.Errorf("优质痛点评论 [%s] 误判为噪音", text)
			}
		}
	})
}

// 测试痛点关键词识别与质量评分
func TestDataCleaner_CalculateQualityScore(t *testing.T) {
	cleaner := clustering.NewDataCleaner(4)

	t.Run("测试包含痛点词的评论获得更高评分", func(t *testing.T) {
		normalText := "今天天气挺好的，出去散步看风景"
		painText := "为什么每次打开都卡顿闪退？到底怎么解决？"

		normalScore := cleaner.CalculateQualityScore(normalText)
		painScore := cleaner.CalculateQualityScore(painText)

		if painScore <= normalScore {
			t.Errorf("痛点文本评分 (%f) 应该严格大于普通文本评分 (%f)", painScore, normalScore)
		}
		if painScore < 3.0 {
			t.Errorf("痛点文本评分期望 >= 3.0，实际为 %f", painScore)
		}
	})
}

// 测试批量清洗功能
func TestDataCleaner_CleanBatch(t *testing.T) {
	cleaner := clustering.NewDataCleaner(4)

	t.Run("测试批量清洗成功剔除噪点并保留有效样本", func(t *testing.T) {
		rawInputs := []string{
			"666",
			"电池耗电太快了，充不进电，难用死了",
			"点赞",
			"请问小白该如何快速起步自媒体剪辑？有没有推荐的工具？",
			"哈哈哈哈",
		}

		cleaned := cleaner.CleanBatch(rawInputs)
		if len(cleaned) != 2 {
			t.Fatalf("批量清洗期望保留 2 条有效评论，实际保留 %d 条", len(cleaned))
		}

		if !cleaned[0].HasPainPoint && !cleaned[1].HasPainPoint {
			t.Errorf("保留的评论期望包含痛点特征标记")
		}
	})
}
