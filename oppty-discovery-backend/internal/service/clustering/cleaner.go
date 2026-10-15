package clustering

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

/**
 * desc: 社交网络评论文本数据清洗器与痛点质量评分模块
 * date: 2026-10-15 15:30:00
 * @author reynolds
 * @since v1.0.0
 */

// CleanedItem 清洗后的文本项与特征评分
type CleanedItem struct {
	OriginalText string  // 原始文本
	CleanText    string  // 清洗去除首尾空格后的文本
	Score        float64 // 质量评分
	HasPainPoint bool    // 是否命中痛点高频关键词
}

// DataCleaner 数据清洗器
type DataCleaner struct {
	minLength     int
	noiseRegexes  []*regexp.Regexp
	noisePhrases  map[string]struct{}
	whitelistKeywords []string
}

// NewDataCleaner 实例化数据清洗器
func NewDataCleaner(minLength int) *DataCleaner {
	if minLength <= 0 {
		minLength = 4
	}

	// 常见噪音正则模式（纯数字、纯符号、单字符多重重复、无意义标点）
	noisePatterns := []string{
		`^[0-9]+$`,                                               // 纯数字
		`^[!@#$%^&*()_+=\-\[\]{};':"\\|,.<>\/?~` + "`" + `]+$`,  // 纯英文标点符号
		`^[。！？、；：""（）《》【】……——·]+$`,                      // 纯中文标点符号
		`^(.)\1{3,}$`,                                            // 单字符连续重复4次及以上
		`^[!?！？\s]+$`,                                           // 连续感叹号问号
		`^[。，、\s]+$`,                                           // 连续逗号句号
		`^(哈|呵|嘿|嘻|哈哈|呵呵|嘿嘿|嘻嘻)+$`,                      // 连续笑声词
		`^(啊|呀|哦|唉|哎|恩|嗯)+$`,                                // 连续感叹语气词
	}

	compiledRegexes := make([]*regexp.Regexp, 0, len(noisePatterns))
	for _, pattern := range noisePatterns {
		if re, err := regexp.Compile(pattern); err == nil {
			compiledRegexes = append(compiledRegexes, re)
		}
	}

	// 常见无意义极短高频短语库（黑名单）
	noisePhrasesList := []string{
		"好", "行", "嗯", "哦", "呵", "哈", "嘿", "嗨", "哟", "哎",
		"666", "999", "888", "111", "000",
		"牛逼", "牛啊", "太强了", "绝了", "太棒了",
		"哈哈哈哈", "哈哈哈", "呵呵呵", "嘻嘻嘻",
		"对", "是", "是的", "对的", "没错", "同意", "赞同", "认可", "不错", "可以", "还行", "确实", "实在", "真的",
		"啥", "啥？", "什么？", "啊？", "真的吗", "是吗", "真的假的", "为什么", "咋回事", "怎么了",
		"谢谢", "感谢", "感谢感谢", "欢迎", "欢迎欢迎", "关注", "回关", "互关", "已关",
		"加微信", "v我", "看主页", "看简介", "详情见主页", "点赞", "收藏", "转发", "关注我", "关注博主",
		"不知道", "不清楚", "不了解", "没看", "没看懂", "看不懂", "路人", "路过", "飘过", "默默路过",
	}

	noisePhrasesMap := make(map[string]struct{}, len(noisePhrasesList))
	for _, phrase := range noisePhrasesList {
		noisePhrasesMap[strings.ToLower(strings.TrimSpace(phrase))] = struct{}{}
	}

	// 高价值商业痛点特征关键词库（白名单加权）
	whitelistKeywords := []string{
		// 疑问与求助
		"怎么", "如何", "怎样", "怎么弄", "怎么搞", "怎么用", "求教", "请教",
		// 阻碍与否定
		"不会", "不懂", "不明白", "不知道", "不理解", "用不来",
		// 缺陷与抱怨
		"问题", "错误", "失败", "不行", "不能", "无法", "难用", "不好用", "用不了", "不好", "差",
		// 性能与稳定性
		"卡", "卡顿", "慢", "闪退", "崩溃", "报错", "发烫", "发热", "死机", "掉帧",
		// 成本与收费
		"多少钱", "价格", "贵", "太贵", "便宜", "收费", "割韭菜", "性价比", "收费太高",
		// 替代与选购
		"推荐", "建议", "哪个好", "哪个牌子", "有替代吗", "平替",
		// 电池续航痛点
		"耗电", "掉电快", "不耐用", "充不进电", "续航差",
	}

	return &DataCleaner{
		minLength:         minLength,
		noiseRegexes:      compiledRegexes,
		noisePhrases:      noisePhrasesMap,
		whitelistKeywords: whitelistKeywords,
	}
}

// IsNoise 判断某单条文本是否属于无价值噪音水帖
func (c *DataCleaner) IsNoise(text string) bool {
	trimmed := strings.TrimSpace(text)
	charCount := utf8.RuneCountInString(trimmed)

	// 1. 字符长度小于门槛
	if charCount < c.minLength {
		return true
	}

	// 2. 匹配噪音黑名单完全匹配
	if _, ok := c.noisePhrases[strings.ToLower(trimmed)]; ok {
		return true
	}

	// 3. 匹配噪音正则表达式
	for _, re := range c.noiseRegexes {
		if re.MatchString(trimmed) {
			return true
		}
	}

	return false
}

// HasPainPointKeyword 检查是否包含高价值商业痛点词
func (c *DataCleaner) HasPainPointKeyword(text string) bool {
	for _, kw := range c.whitelistKeywords {
		if strings.Contains(text, kw) {
			return true
		}
	}
	return false
}

// CalculateQualityScore 计算该文本的商业价值质量打分
func (c *DataCleaner) CalculateQualityScore(text string) float64 {
	score := 1.0
	charLen := utf8.RuneCountInString(text)

	// 命中痛点词汇加权 (+2.0)
	if c.HasPainPointKeyword(text) {
		score += 2.0
	}

	// 文本长度适中（15~150个字最详实）加权 (+1.0)
	if charLen >= 15 && charLen <= 150 {
		score += 1.0
	}

	return score
}

// CleanBatch 批量执行文本清洗，过滤噪音并计算质量分数
func (c *DataCleaner) CleanBatch(texts []string) []CleanedItem {
	var results []CleanedItem
	for _, t := range texts {
		trimmed := strings.TrimSpace(t)
		if trimmed == "" {
			continue
		}
		if c.IsNoise(trimmed) {
			continue
		}

		hasPain := c.HasPainPointKeyword(trimmed)
		score := c.CalculateQualityScore(trimmed)

		results = append(results, CleanedItem{
			OriginalText: t,
			CleanText:    trimmed,
			Score:        score,
			HasPainPoint: hasPain,
		})
	}
	return results
}
