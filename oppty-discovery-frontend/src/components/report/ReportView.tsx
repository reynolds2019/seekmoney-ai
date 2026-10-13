import React, { useState } from 'react';
import {
  Maximize2,
  Sparkles,
  Target,
  Rocket,
  MessageCircle,
  Copy,
  Flame,
  X,
  ExternalLink,
  Share2,
  Video,
  ThumbsUp,
  Check,
} from 'lucide-react';
import { OpptyReportVO, VideoItemVO } from '../../types';

interface ReportViewProps {
  reports: OpptyReportVO[];
  jobId?: string;
  sourceVideos?: VideoItemVO[];
  onOpenFullscreen: (report: OpptyReportVO) => void;
}

export const ReportView: React.FC<ReportViewProps> = ({
  reports,
  jobId,
  sourceVideos = [],
  onOpenFullscreen,
}) => {
  const [selectedCommentsReport, setSelectedCommentsReport] = useState<OpptyReportVO | null>(null);
  const [selectedVideosReport, setSelectedVideosReport] = useState<OpptyReportVO | null>(null);
  const [copiedText, setCopiedText] = useState(false);
  const [shareCopied, setShareCopied] = useState(false);

  if (!reports || reports.length === 0) {
    return null;
  }

  const handleCopy = (text: string) => {
    navigator.clipboard.writeText(text);
    setCopiedText(true);
    setTimeout(() => setCopiedText(false), 2000);
  };

  const handleShareReport = () => {
    const currentOrigin = window.location.origin;
    const currentPath = window.location.pathname;
    const targetJobId = jobId || (reports[0] as any)?.jobId || '';
    const shareUrl = `${currentOrigin}${currentPath}?share=${targetJobId}`;
    navigator.clipboard.writeText(shareUrl);
    setShareCopied(true);
    setTimeout(() => setShareCopied(false), 3000);
  };

  return (
    <div className="space-y-6">
      {/* 头部标题与分享操作栏 */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 bg-white p-4 sm:p-5 rounded-2xl border border-gray-100 shadow-sm">
        <div>
          <h2 className="text-base sm:text-lg font-black text-gray-900 flex items-center gap-2">
            <Sparkles className="w-5 h-5 text-amber-500" />
            商机透视报告 ({reports.length} 项核心商机)
          </h2>
          <p className="text-xs text-gray-500 mt-1">
            基于真实评论密度、商业痛点与原文证据链提炼的高价值商机
          </p>
        </div>

        {/* 分享报告操作按钮 */}
        <div className="flex items-center gap-2">
          <button
            type="button"
            onClick={handleShareReport}
            className={`btn text-xs py-2 px-3.5 flex items-center gap-1.5 transition-all ${
              shareCopied
                ? 'bg-emerald-600 hover:bg-emerald-700 text-white'
                : 'btn-secondary text-gray-700 hover:text-blue-600'
            }`}
            title="生成免登录在线公开查阅链接并复制"
          >
            {shareCopied ? (
              <>
                <Check className="w-3.5 h-3.5 text-white" />
                <span>分享链接已复制！</span>
              </>
            ) : (
              <>
                <Share2 className="w-3.5 h-3.5 text-blue-600" />
                <span>分享商机报告</span>
              </>
            )}
          </button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {reports.map((report, idx) => {
          const interestBadgeClass =
            report.paidInterest === 'High'
              ? 'badge-high'
              : report.paidInterest === 'Medium'
              ? 'badge-medium'
              : 'badge-low';

          // 综合视频项：若单项有挂载则用单项，否则用整体 sourceVideos
          const reportVideos = report.sourceVideos && report.sourceVideos.length > 0 ? report.sourceVideos : sourceVideos;

          return (
            <div
              key={idx}
              className="card-glass p-5 sm:p-6 flex flex-col justify-between hover:border-blue-300 transition-all group relative"
            >
              {/* 顶部标签与全屏按钮 */}
              <div>
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-2">
                    <span className="text-xs font-mono font-bold text-gray-400">
                      #{idx + 1}
                    </span>
                    <span className={`badge ${interestBadgeClass}`}>
                      <Flame className="w-3 h-3" />
                      {report.paidInterest === 'High'
                        ? '强付费意愿'
                        : report.paidInterest === 'Medium'
                        ? '中等意愿'
                        : '弱意愿'}
                    </span>
                    <span className="text-[11px] text-gray-400">
                      关联信号: {report.signalSize} 条
                    </span>
                  </div>

                  <button
                    onClick={() => onOpenFullscreen(report)}
                    className="p-1.5 rounded-lg text-gray-400 hover:text-blue-600 hover:bg-blue-50 transition"
                    title="全屏查看报告"
                  >
                    <Maximize2 className="w-4 h-4" />
                  </button>
                </div>

                {/* 一句话核心痛点 */}
                <h3 className="text-base font-bold text-gray-900 mb-2 leading-snug group-hover:text-blue-600 transition-colors">
                  {report.oneLinePain}
                </h3>

                {/* 构思产品概念 */}
                <div className="p-3 bg-gradient-to-r from-blue-50/80 to-indigo-50/80 rounded-xl border border-blue-100/80 mb-4">
                  <div className="text-[11px] font-bold text-blue-700 flex items-center gap-1.5 mb-1">
                    <Rocket className="w-3.5 h-3.5" />
                    商业产品构思定位
                  </div>
                  <p className="text-xs font-semibold text-gray-900">
                    {report.potentialProduct}
                  </p>
                </div>

                {/* 判定依据与原文说明 */}
                <p className="text-xs text-gray-600 mb-4 line-clamp-2 leading-relaxed">
                  <span className="font-semibold text-gray-800">透视依据：</span>
                  {report.rationale}
                </p>

                {/* 核心功能亮点 */}
                {report.mvpPlan?.core_features && (
                  <div className="mb-4">
                    <div className="text-[11px] font-bold text-gray-500 uppercase tracking-wider mb-1.5 flex items-center gap-1">
                      <Target className="w-3 h-3 text-indigo-500" />
                      首期必备功能
                    </div>
                    <ul className="space-y-1">
                      {report.mvpPlan.core_features.slice(0, 3).map((feat: string, fIdx: number) => (
                        <li key={fIdx} className="text-xs text-gray-700 flex items-start gap-1.5">
                          <span className="w-1.5 h-1.5 rounded-full bg-blue-500 mt-1.5 shrink-0" />
                          <span>{feat}</span>
                        </li>
                      ))}
                    </ul>
                  </div>
                )}
              </div>

              {/* 底部操作工具栏 */}
              <div className="pt-3 border-t border-gray-100 flex flex-wrap items-center justify-between gap-2 text-xs">
                <div className="flex items-center gap-3">
                  <button
                    type="button"
                    onClick={() => setSelectedCommentsReport(report)}
                    className="flex items-center gap-1.5 text-gray-600 hover:text-blue-600 font-medium py-1 transition-colors"
                  >
                    <MessageCircle className="w-3.5 h-3.5 text-blue-500" />
                    <span>原声评论 ({report.representativeTexts?.length || 0})</span>
                  </button>

                  <button
                    type="button"
                    onClick={() => setSelectedVideosReport(report)}
                    className="flex items-center gap-1.5 text-gray-600 hover:text-indigo-600 font-medium py-1 transition-colors"
                  >
                    <Video className="w-3.5 h-3.5 text-indigo-500" />
                    <span>来源视频/原文 ({reportVideos.length})</span>
                  </button>
                </div>

                <div className="flex items-center gap-2">
                  {report.minioMarkdownUrl && (
                    <a
                      href={report.minioMarkdownUrl}
                      target="_blank"
                      rel="noreferrer"
                      className="text-gray-400 hover:text-gray-700 p-1"
                      title="打开 MinIO Markdown 原生报告"
                    >
                      <ExternalLink className="w-3.5 h-3.5" />
                    </a>
                  )}
                  <button
                    onClick={() => onOpenFullscreen(report)}
                    className="btn btn-secondary text-xs py-1 px-2.5"
                  >
                    全屏研读
                  </button>
                </div>
              </div>
            </div>
          );
        })}
      </div>

      {/* 原声评论查看模态框 */}
      {selectedCommentsReport && (
        <div className="modal-overlay">
          <div className="card-glass w-full max-w-xl max-h-[80vh] flex flex-col p-6 relative shadow-2xl">
            <div className="flex items-center justify-between pb-3 border-b border-gray-200">
              <div className="flex items-center gap-2">
                <MessageCircle className="w-4 h-4 text-blue-600" />
                <h4 className="text-sm font-bold text-gray-900">
                  真实用户原声评论与吐槽样本
                </h4>
              </div>
              <button
                onClick={() => setSelectedCommentsReport(null)}
                className="text-gray-400 hover:text-gray-600 p-1 rounded-lg"
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            <div className="flex-1 overflow-y-auto py-4 space-y-2.5">
              {selectedCommentsReport.representativeTexts?.map((text, cIdx) => (
                <div
                  key={cIdx}
                  className="p-3 bg-gray-50 border border-gray-100 rounded-xl text-xs text-gray-700 flex items-start justify-between gap-2 group hover:bg-blue-50/40 transition"
                >
                  <p className="leading-relaxed flex-1">{text}</p>
                  <button
                    onClick={() => handleCopy(text)}
                    className="text-gray-400 hover:text-blue-600 p-1 shrink-0 opacity-0 group-hover:opacity-100 transition"
                    title="复制评论"
                  >
                    <Copy className="w-3.5 h-3.5" />
                  </button>
                </div>
              ))}
            </div>

            <div className="pt-3 border-t border-gray-200 flex justify-between items-center text-xs text-gray-500">
              <span>{copiedText ? '已复制到剪贴板！' : '共展示代表性真实吐槽'}</span>
              <button
                onClick={() => setSelectedCommentsReport(null)}
                className="btn btn-secondary text-xs py-1 px-3"
              >
                关闭
              </button>
            </div>
          </div>
        </div>
      )}

      {/* 抓取视频与原文链接查看模态框 */}
      {selectedVideosReport && (
        <div className="modal-overlay">
          <div className="card-glass w-full max-w-2xl max-h-[85vh] flex flex-col p-6 relative shadow-2xl">
            <div className="flex items-center justify-between pb-3.5 border-b border-gray-200">
              <div className="flex items-center gap-2">
                <Video className="w-4 h-4 text-indigo-600" />
                <h4 className="text-sm font-bold text-gray-900">
                  抓取到的原始视频与内容原文链接清单
                </h4>
              </div>
              <button
                onClick={() => setSelectedVideosReport(null)}
                className="text-gray-400 hover:text-gray-600 p-1 rounded-lg"
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            <div className="flex-1 overflow-y-auto py-4 space-y-3">
              {(
                (selectedVideosReport.sourceVideos && selectedVideosReport.sourceVideos.length > 0
                  ? selectedVideosReport.sourceVideos
                  : sourceVideos) || []
              ).map((video, vIdx) => (
                <div
                  key={vIdx}
                  className="p-3.5 bg-gray-50 border border-gray-200/80 rounded-xl text-xs flex flex-col sm:flex-row sm:items-center justify-between gap-3 hover:bg-indigo-50/30 hover:border-indigo-200 transition"
                >
                  <div className="space-y-1 min-w-0 flex-1">
                    <div className="font-bold text-gray-900 leading-snug line-clamp-2">
                      {video.title}
                    </div>
                    <div className="flex flex-wrap items-center gap-3 text-[11px] text-gray-500">
                      <span>发布者: <b className="text-gray-700">{video.author}</b></span>
                      <span className="flex items-center gap-1">
                        <ThumbsUp className="w-3 h-3 text-amber-500" />
                        {video.likes}
                      </span>
                      {video.createdAt && <span>时间: {video.createdAt}</span>}
                    </div>
                  </div>

                  <a
                    href={video.videoUrl}
                    target="_blank"
                    rel="noreferrer"
                    className="btn btn-primary text-xs py-1.5 px-3 flex items-center justify-center gap-1.5 shrink-0 self-start sm:self-center"
                  >
                    <span>跳转原文</span>
                    <ExternalLink className="w-3.5 h-3.5" />
                  </a>
                </div>
              ))}
            </div>

            <div className="pt-3 border-t border-gray-200 flex justify-between items-center text-xs text-gray-500">
              <span>提供真实采集证据链，点击按钮可直接访问目标平台</span>
              <button
                onClick={() => setSelectedVideosReport(null)}
                className="btn btn-secondary text-xs py-1 px-3"
              >
                关闭
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
