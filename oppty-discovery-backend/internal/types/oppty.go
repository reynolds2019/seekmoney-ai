package types

/**
 * desc: 商机挖掘、报告呈现与历史记录请求与响应 DTO/VO
 * date: 2026-10-13 11:25:00
 * @author reynolds
 * @since v1.0.0
 */

// 提交商机分析请求入参
type AnalyzeReq struct {
	Keywords            []string `json:"keywords" binding:"required,min=1,max=5,dive,min=1,max=50"`
	DataSource          string   `json:"dataSource" binding:"required,oneof=tiktok douyin"`
	DeepCrawl           bool     `json:"deepCrawl"`
	MaxVideos           int      `json:"maxVideos" binding:"required,min=5,max=20"`
	MaxCommentsPerVideo int      `json:"maxCommentsPerVideo" binding:"required,min=10,max=100"`
}

// 提交商机分析任务响应出参
type AnalyzeRespVO struct {
	JobID         string  `json:"jobId"`
	Status        string  `json:"status"`
	EstimatedCost float64 `json:"estimatedCost"`
}

// 痛点深度四维剖析
type PainDepthVO struct {
	SurfacePain        string   `json:"surface_pain"`
	RootCauses         []string `json:"root_causes"`
	UserScenarios      []string `json:"user_scenarios"`
	EmotionalIntensity int      `json:"emotional_intensity"`
}

// 现有解决方案项
type ExistingSolutionVO struct {
	Name       string `json:"name"`
	Limitation string `json:"limitation"`
}

// 市场与竞品格局
type MarketLandscapeVO struct {
	ExistingSolutions []ExistingSolutionVO `json:"existing_solutions"`
	UnmetNeeds        []string             `json:"unmet_needs"`
	Opportunity       string               `json:"opportunity"`
}

// 假设验证项
type ValidationHypothesisVO struct {
	Hypothesis string `json:"hypothesis"`
	TestMethod string `json:"test_method"`
}

// MVP 商业执行规划
type MVPPlanVO struct {
	CoreFeatures         []string                 `json:"core_features"`
	ValidationHypotheses []ValidationHypothesisVO `json:"validation_hypotheses"`
	FirstUsers           string                   `json:"first_users"`
	Timeline             string                   `json:"timeline"`
	EstimatedCost        string                   `json:"estimated_cost"`
}

// 单条商机聚类报告视图
type OpptyReportVO struct {
	ID                  uint64            `json:"id"`
	ClusterCode         string            `json:"clusterCode"`
	SignalSize          int               `json:"signalSize"`
	OneLinePain         string            `json:"oneLinePain"`
	PaidInterest        string            `json:"paidInterest"`
	Rationale           string            `json:"rationale"`
	PotentialProduct    string            `json:"potentialProduct"`
	PainDepth           PainDepthVO       `json:"painDepth"`
	MarketLandscape     MarketLandscapeVO `json:"marketLandscape"`
	MVPPlan             MVPPlanVO         `json:"mvpPlan"`
	KeywordRelevance    int               `json:"keywordRelevance"`
	PriorityScore       int               `json:"priorityScore"`
	RepresentativeTexts []string          `json:"representativeTexts"`
	MinioMarkdownURL    string            `json:"minioMarkdownUrl"`
	MinioHTMLURL        string            `json:"minioHtmlUrl"`
}

// 任务详情与报告综合视图
type OpptyJobDetailVO struct {
	JobID               string          `json:"jobId"`
	Status              string          `json:"status"`
	ProgressStage       string          `json:"progressStage"`
	ProgressPercent     int             `json:"progressPercent"`
	ProgressMessage     string          `json:"progressMessage"`
	Keywords            []string        `json:"keywords"`
	DataSource          string          `json:"dataSource"`
	DeepCrawl           bool            `json:"deepCrawl"`
	MaxVideos           int             `json:"maxVideos"`
	MaxCommentsPerVideo int             `json:"maxCommentsPerVideo"`
	EstimatedCost       float64         `json:"estimatedCost"`
	CreatedAt           string          `json:"createdAt"`
	Results             []OpptyReportVO `json:"results"`
	ErrorMessage        string          `json:"errorMessage,omitempty"`
}

// 商机历史列表项
type OpptyHistoryItemVO struct {
	JobID         string   `json:"jobId"`
	Keywords      []string `json:"keywords"`
	DataSource    string   `json:"dataSource"`
	SignalCount   int      `json:"signalCount"`
	ReportCount   int      `json:"reportCount"`
	PrimaryPain   string   `json:"primaryPain"`
	Status        string   `json:"status"`
	EstimatedCost float64  `json:"estimatedCost"`
	CreatedAt     string   `json:"createdAt"`
}

// 商机历史分页响应
type OpptyHistoryPageRespVO struct {
	Total    int64                `json:"total"`
	PageNo   int                  `json:"pageNo"`
	PageSize int                  `json:"pageSize"`
	Records  []OpptyHistoryItemVO `json:"records"`
}

// 健康检查出参
type HealthRespVO struct {
	Service   string `json:"service"`
	Status    string `json:"status"`
	Database  string `json:"database"`
	Minio     string `json:"minio"`
	Timestamp int64  `json:"timestamp"`
}
