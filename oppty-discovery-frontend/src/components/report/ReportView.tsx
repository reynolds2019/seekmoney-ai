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
} from 'lucide-react';
import { OpptyReportVO } from '../../types';

interface ReportViewProps {
  reports: OpptyReportVO[];
  onOpenFullscreen: (report: OpptyReportVO) => void;
}

export const ReportView: React.FC<ReportViewProps> = ({ reports, onOpenFullscreen }) => {
  const [selectedCommentsReport, setSelectedCommentsReport] = useState<OpptyReportVO | null>(null);
  const [copiedText, setCopiedText] = useState(false);

  if (!reports || reports.length === 0) {
    return null;
  }

  const handleCopy = (text: string) => {
    navigator.clipboard.writeText(text);
    setCopiedText(true);
    setTimeout(() => setCopiedText(false), 2000);
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-lg font-extrabold text-gray-900 flex items-center gap-2">
            <Sparkles className="w-5 h-5 text-amber-500" />
            商机透视报告 ({reports.length} 项核心商机)
          </h2>
          <p className="text-xs text-gray-500 mt-0.5">
            基于真实评论密度与商业价值模型提炼的商机洞察
          </p>
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
                      MVP 核心功能路线
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
              <div className="pt-3 border-t border-gray-100 flex items-center justify-between text-xs">
                <button
                  type="button"
                  onClick={() => setSelectedCommentsReport(report)}
                  className="flex items-center gap-1.5 text-gray-600 hover:text-blue-600 font-medium py-1"
                >
                  <MessageCircle className="w-3.5 h-3.5 text-blue-500" />
                  <span>查看原声评论 ({report.representativeTexts?.length || 0})</span>
                </button>

                <div className="flex items-center gap-2">
                  {report.minioMarkdownUrl && (
                    <a
                      href={report.minioMarkdownUrl}
                      target="_blank"
                      rel="noreferrer"
                      className="text-gray-400 hover:text-gray-700 p-1"
                      title="打开 MinIO Markdown 报告"
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
    </div>
  );
};
