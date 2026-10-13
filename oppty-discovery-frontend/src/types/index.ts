/**
 * 数据源类型
 */
export type DataSourceType = 'tiktok' | 'douyin';

/**
 * 任务主状态
 */
export type JobStatusType = 'processing' | 'completed' | 'failed';

/**
 * 任务流转细分阶段
 */
export type ProgressStageType = 
  | 'init' 
  | 'validating' 
  | 'crawling' 
  | 'clustering' 
  | 'analyzing' 
  | 'completed' 
  | 'failed';

/**
 * 痛点深度剖析
 */
export interface PainDepth {
  surface_pain: string;
  root_causes: string[];
  user_scenarios: string[];
  emotional_intensity: number;
}

/**
 * 竞品与市场格局
 */
export interface MarketLandscape {
  existing_solutions: Array<{
    name: string;
    limitation: string;
  }>;
  unmet_needs: string[];
  opportunity: string;
}

/**
 * MVP落地执行方案
 */
export interface MVPPlan {
  core_features: string[];
  validation_hypotheses: Array<{
    hypothesis: string;
    test_method: string;
  }>;
  first_users: string;
  timeline: string;
  estimated_cost: string;
}

/**
 * 单条商机聚类报告视图
 */
export interface OpptyReportVO {
  id: number;
  clusterCode: string;
  signalSize: number;
  oneLinePain: string;
  paidInterest: 'High' | 'Medium' | 'Low';
  rationale: string;
  potentialProduct: string;
  painDepth: PainDepth;
  marketLandscape: MarketLandscape;
  mvpPlan: MVPPlan;
  keywordRelevance: number;
  priorityScore: number;
  representativeTexts: string[];
  minioMarkdownUrl?: string;
  minioHtmlUrl?: string;
}

/**
 * 任务详情与报告综合视图
 */
export interface OpptyJobDetailVO {
  jobId: string;
  status: JobStatusType;
  progressStage: ProgressStageType;
  progressPercent: number;
  progressMessage: string;
  keywords: string[];
  dataSource: DataSourceType;
  deepCrawl: boolean;
  maxVideos: number;
  maxCommentsPerVideo: number;
  estimatedCost: number;
  createdAt: string;
  results: OpptyReportVO[];
  errorMessage?: string;
}

/**
 * 商机历史记录摘要项
 */
export interface OpptyHistoryItemVO {
  jobId: string;
  keywords: string[];
  dataSource: DataSourceType;
  signalCount: number;
  reportCount: number;
  primaryPain: string;
  status: JobStatusType;
  estimatedCost: number;
  createdAt: string;
}

/**
 * 用户档案
 */
export interface UserProfileVO {
  id: number;
  email: string;
  nickname: string;
  avatarUrl: string;
  isNewUser: boolean;
}

/**
 * 统一后端返回包装
 */
export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
  timestamp: number;
}
