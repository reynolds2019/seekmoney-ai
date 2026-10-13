package constant

/**
 * desc: 系统全局常量与字典定义
 * date: 2026-10-13 11:25:00
 * @author reynolds
 * @since v1.0.0
 */

// 数据源类型枚举
type DataSourceType string

const (
	DataSourceTikTok DataSourceType = "tiktok" // TikTok 国际版数据源
	DataSourceDouyin DataSourceType = "douyin" // 抖音国内版数据源
)

// 获取数据源中文说明
func (d DataSourceType) Desc() string {
	switch d {
	case DataSourceTikTok:
		return "TikTok 国际版"
	case DataSourceDouyin:
		return "抖音短视频"
	default:
		return "未知数据源"
	}
}

// 任务执行主状态枚举
type JobStatus string

const (
	JobStatusProcessing JobStatus = "processing" // 进行中
	JobStatusCompleted  JobStatus = "completed"  // 分析完成
	JobStatusFailed     JobStatus = "failed"     // 执行失败
)

// 任务细分流转阶段枚举
type ProgressStage string

const (
	StageInit       ProgressStage = "init"       // 初始化阶段
	StageValidating ProgressStage = "validating" // 参数校验
	StageCrawling   ProgressStage = "crawling"   // 社交媒体数据抓取中
	StageClustering ProgressStage = "clustering" // 评论语义密度聚类分析中
	StageAnalyzing  ProgressStage = "analyzing"  // 大模型深度商业洞察提炼中
	StageCompleted  ProgressStage = "completed"  // 全量分析已完成
	StageFailed     ProgressStage = "failed"     // 执行失败中断
)

// 付费意愿强烈程度枚举
type PaidInterest string

const (
	PaidInterestHigh   PaidInterest = "High"   // 高付费意愿
	PaidInterestMedium PaidInterest = "Medium" // 中等付费意愿
	PaidInterestLow    PaidInterest = "Low"    // 低付费意愿
)

// 系统公共字典常量
const (
	AppNameZh        = "商机发现助手"
	AppNameEn        = "Oppty Discovery Assistant"
	AppSloganZh      = "痛点一键透视，商机即刻成型。"
	AppSloganEn      = "From pain point to profit in seconds."
	AppDescriptionZh = "一款集“社区评论挖掘、痛点线索提炼与商业方案构思”于一体的智能商机雷达。"

	DefaultMinioBucket = "oppty-discovery" // MinIO 专用存储桶
	DefaultDatabase    = "oppty_discovery" // 数据库名称

	HeaderAuthorization = "Authorization"
	BearerPrefix        = "Bearer "
	ContextKeyUserID    = "current_user_id"
	ContextKeyUserEmail = "current_user_email"
)
