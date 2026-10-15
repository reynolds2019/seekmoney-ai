import React, { useEffect, useState } from 'react';
import {
  X,
  Sparkles,
  Flame,
  Target,
  ShieldAlert,
  Compass,
  Rocket,
  CheckCircle,
  HelpCircle,
  Clock,
  Coins,
  MessageSquare,
  Share2,
  Video,
  ThumbsUp,
  ExternalLink,
  Check,
  AlertCircle,
} from 'lucide-react';
import { OpptyReportVO } from '../../types';
import { copyToClipboard } from '../../utils/clipboard';

interface FullscreenReportModalProps {
  report: OpptyReportVO | null;
  jobId?: string;
  onClose: () => void;
}

export const FullscreenReportModal: React.FC<FullscreenReportModalProps> = ({ report, jobId, onClose }) => {
  const [shareCopied, setShareCopied] = useState(false);

  // 支持按 ESC 键退出全屏
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        onClose();
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [onClose]);

  if (!report) return null;

  const handleShare = async () => {
    const currentOrigin = window.location.origin;
    const currentPath = window.location.pathname;
    const shareUrl = `${currentOrigin}${currentPath}?share=${jobId || ''}`;
    const ok = await copyToClipboard(shareUrl);
    if (ok) {
      setShareCopied(true);
      setTimeout(() => setShareCopied(false), 3000);
    }
  };

  const { painDepth, marketLandscape, mvpPlan } = report;

  return (
    <div className="modal-fullscreen bg-gray-50 flex flex-col">
      {/* 顶部工具栏 */}
      <header className="sticky top-0 z-20 bg-white/90 backdrop-blur-md border-b border-gray-200 px-4 sm:px-8 py-3.5 flex items-center justify-between shadow-sm">
        <div className="flex items-center gap-3">
          <div className="w-9 h-9 rounded-xl bg-blue-600 text-white flex items-center justify-center font-bold text-sm shadow-md shadow-blue-500/20">
            <Sparkles className="w-5 h-5" />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <span className="font-extrabold text-sm sm:text-base text-gray-900">
                商机深度全屏透视报告
              </span>
              <span className="text-xs font-mono px-2 py-0.5 bg-gray-100 text-gray-600 rounded">
                {report.clusterCode}
              </span>
            </div>
            <p className="text-[11px] text-gray-500">痛点一键透视，商机即刻成型 · 支持 ESC 键随时退出</p>
          </div>
        </div>

        <div className="flex items-center gap-2.5">
          <button
            type="button"
            onClick={handleShare}
            className={`btn text-xs py-2 px-3 flex items-center gap-1.5 transition-all ${
              shareCopied
                ? 'bg-emerald-600 hover:bg-emerald-700 text-white'
                : 'btn-secondary text-gray-700 hover:text-blue-600'
            }`}
            title="复制在线公开查阅链接"
          >
            {shareCopied ? (
              <>
                <Check className="w-3.5 h-3.5 text-white" />
                <span>已复制！</span>
              </>
            ) : (
              <>
                <Share2 className="w-3.5 h-3.5 text-blue-600" />
                <span>分享报告</span>
              </>
            )}
          </button>

          <button
            onClick={onClose}
            className="btn btn-secondary text-xs py-2 px-3.5 flex items-center gap-1.5"
          >
            <X className="w-4 h-4" />
            <span>退出全屏 (ESC)</span>
          </button>
        </div>
      </header>

      {/* 主体研读内容 */}
      <main className="flex-1 max-w-5xl mx-auto w-full px-4 sm:px-8 py-8 space-y-8">
        {/* 1. 核心综述横幅 */}
        <div className="card-glass p-6 sm:p-8 bg-gradient-to-br from-white via-blue-50/20 to-indigo-50/30 border-blue-200/60 shadow-lg">
          <div className="flex flex-wrap items-center gap-3 mb-4">
            <span
              className={`badge ${
                report.paidInterest === 'High'
                  ? 'badge-high'
                  : report.paidInterest === 'Medium'
                  ? 'badge-medium'
                  : 'badge-low'
              } text-xs px-3 py-1`}
            >
              <Flame className="w-3.5 h-3.5" />
              {report.paidInterest === 'High'
                ? '强付费意愿验证'
                : report.paidInterest === 'Medium'
                ? '中等付费意愿'
                : '低付费意愿'}
            </span>
            <span className="text-xs font-medium text-gray-500">
              关联社区吐槽声量：<b className="text-gray-900 font-mono">{report.signalSize}</b> 条
            </span>
            <span className="text-xs font-medium text-gray-500">
              商业优先级评分：<b className="text-blue-600 font-mono">{report.priorityScore}</b> / 100
            </span>
          </div>

          <h1 className="text-xl sm:text-2xl font-black text-gray-900 leading-snug mb-4">
            {report.oneLinePain}
          </h1>

          <div className="p-4 bg-white/90 rounded-2xl border border-blue-100 shadow-sm flex items-start gap-3">
            <Rocket className="w-5 h-5 text-blue-600 shrink-0 mt-0.5" />
            <div>
              <div className="text-xs font-bold text-blue-800 uppercase tracking-wider mb-0.5">
                推荐商业产品构思
              </div>
              <p className="text-sm font-bold text-gray-900">{report.potentialProduct}</p>
            </div>
          </div>

          <p className="mt-4 text-xs sm:text-sm text-gray-600 leading-relaxed">
            <b className="text-gray-900 font-semibold">商业判断依据：</b>
            {report.rationale}
          </p>
        </div>

        {/* 2. 痛点四维深度解构 (Pain Depth) */}
        <section className="card-glass p-6 sm:p-8 space-y-4">
          <div className="flex items-center gap-2 pb-3 border-b border-gray-100">
            <ShieldAlert className="w-5 h-5 text-red-500" />
            <h2 className="text-base font-bold text-gray-900">痛点四维深度解构</h2>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="p-4 bg-gray-50 rounded-xl border border-gray-100">
              <span className="text-xs font-bold text-gray-500 block mb-1">表象浅层痛点</span>
              <p className="text-sm font-semibold text-gray-900">
                {painDepth?.surface_pain || report.oneLinePain}
              </p>
            </div>

            <div className="p-4 bg-gray-50 rounded-xl border border-gray-100">
              <span className="text-xs font-bold text-gray-500 block mb-1">情绪烈度评分 (1-5)</span>
              <div className="flex items-center gap-2">
                <div className="flex gap-1 text-amber-500">
                  {Array.from({ length: 5 }).map((_, i) => (
                    <Flame
                      key={i}
                      className={`w-4 h-4 ${
                        i < (painDepth?.emotional_intensity || 4)
                          ? 'fill-amber-500 text-amber-500'
                          : 'text-gray-200'
                      }`}
                    />
                  ))}
                </div>
                <span className="text-xs font-mono font-bold text-gray-700">
                  {painDepth?.emotional_intensity || 4} 级烈度
                </span>
              </div>
            </div>
          </div>

          {/* 3个根本原因与典型场景 */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pt-2">
            <div>
              <span className="text-xs font-bold text-gray-600 block mb-2 flex items-center gap-1.5">
                <HelpCircle className="w-3.5 h-3.5 text-blue-500" />
                连问 3 次为什么（根本原因）：
              </span>
              <div className="space-y-2">
                {painDepth?.root_causes?.map((cause, i) => (
                  <div key={i} className="text-xs p-2.5 bg-white rounded-lg border border-gray-200 flex items-start gap-2">
                    <span className="w-4 h-4 rounded-full bg-blue-100 text-blue-700 text-[10px] font-bold flex items-center justify-center shrink-0">
                      {i + 1}
                    </span>
                    <span className="text-gray-800">{cause}</span>
                  </div>
                ))}
              </div>
            </div>

            <div>
              <span className="text-xs font-bold text-gray-600 block mb-2 flex items-center gap-1.5">
                <Compass className="w-3.5 h-3.5 text-indigo-500" />
                典型爆发场景：
              </span>
              <div className="space-y-2">
                {painDepth?.user_scenarios?.map((scen, i) => (
                  <div key={i} className="text-xs p-2.5 bg-white rounded-lg border border-gray-200 flex items-start gap-2">
                    <span className="w-1.5 h-1.5 rounded-full bg-indigo-500 mt-1.5 shrink-0" />
                    <span className="text-gray-800">{scen}</span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </section>

        {/* 3. 市场与竞品格局分析 */}
        <section className="card-glass p-6 sm:p-8 space-y-4">
          <div className="flex items-center gap-2 pb-3 border-b border-gray-100">
            <Compass className="w-5 h-5 text-indigo-500" />
            <h2 className="text-base font-bold text-gray-900">市场格局与竞品局限</h2>
          </div>

          {/* 竞品局限 */}
          <div>
            <span className="text-xs font-bold text-gray-500 uppercase tracking-wider block mb-2">
              现有替代方案与用户槽点
            </span>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
              {marketLandscape?.existing_solutions?.map((sol, i) => (
                <div key={i} className="p-3.5 bg-gray-50 rounded-xl border border-gray-100">
                  <div className="text-xs font-bold text-gray-900 mb-1">{sol.name}</div>
                  <div className="text-xs text-red-600 flex items-start gap-1.5">
                    <span className="font-semibold shrink-0">局限：</span>
                    <span>{sol.limitation}</span>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* 未满足需求与机会切入 */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pt-2">
            <div>
              <span className="text-xs font-bold text-gray-500 uppercase tracking-wider block mb-2">
                尚未被满足的需求
              </span>
              <ul className="space-y-1.5">
                {marketLandscape?.unmet_needs?.map((need, i) => (
                  <li key={i} className="text-xs text-gray-800 flex items-start gap-2">
                    <CheckCircle className="w-3.5 h-3.5 text-emerald-500 shrink-0 mt-0.5" />
                    <span>{need}</span>
                  </li>
                ))}
              </ul>
            </div>

            <div className="p-4 bg-emerald-50/70 rounded-xl border border-emerald-200">
              <span className="text-xs font-bold text-emerald-800 uppercase tracking-wider block mb-1">
                差异化市场切入机会
              </span>
              <p className="text-xs font-semibold text-emerald-950 leading-relaxed">
                {marketLandscape?.opportunity}
              </p>
            </div>
          </div>
        </section>

        {/* 4. 快速验证与启动方案 */}
        <section className="card-glass p-6 sm:p-8 space-y-4">
          <div className="flex items-center gap-2 pb-3 border-b border-gray-100">
            <Target className="w-5 h-5 text-blue-600" />
            <h2 className="text-base font-bold text-gray-900">快速验证与启动方案</h2>
          </div>

          {/* 核心功能 */}
          <div>
            <span className="text-xs font-bold text-gray-500 uppercase tracking-wider block mb-2">
              首期必备功能清单
            </span>
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
              {mvpPlan?.core_features?.map((feat, i) => (
                <div key={i} className="p-3 bg-blue-50/50 rounded-xl border border-blue-100 text-xs font-medium text-gray-900">
                  <span className="font-bold text-blue-600 mr-1.5">0{i + 1}.</span>
                  {feat}
                </div>
              ))}
            </div>
          </div>

          {/* 落地周期与成本 */}
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 pt-2">
            <div className="p-3.5 bg-gray-50 rounded-xl border border-gray-100 flex items-center gap-3">
              <div className="w-8 h-8 rounded-lg bg-white shadow-sm flex items-center justify-center text-blue-600">
                <Compass className="w-4 h-4" />
              </div>
              <div>
                <div className="text-[11px] text-gray-500">首批种子用户渠道</div>
                <div className="text-xs font-bold text-gray-900">{mvpPlan?.first_users}</div>
              </div>
            </div>

            <div className="p-3.5 bg-gray-50 rounded-xl border border-gray-100 flex items-center gap-3">
              <div className="w-8 h-8 rounded-lg bg-white shadow-sm flex items-center justify-center text-amber-600">
                <Clock className="w-4 h-4" />
              </div>
              <div>
                <div className="text-[11px] text-gray-500">预估开发上线周期</div>
                <div className="text-xs font-bold text-gray-900">{mvpPlan?.timeline}</div>
              </div>
            </div>

            <div className="p-3.5 bg-gray-50 rounded-xl border border-gray-100 flex items-center gap-3">
              <div className="w-8 h-8 rounded-lg bg-white shadow-sm flex items-center justify-center text-emerald-600">
                <Coins className="w-4 h-4" />
              </div>
              <div>
                <div className="text-[11px] text-gray-500">预估测试验证成本</div>
                <div className="text-xs font-bold text-gray-900">{mvpPlan?.estimated_cost}</div>
              </div>
            </div>
          </div>
        </section>

        {/* 5. 真实原声代表文本 */}
        <section className="card-glass p-6 sm:p-8 space-y-3">
          <div className="flex items-center gap-2 pb-2 border-b border-gray-100">
            <MessageSquare className="w-5 h-5 text-gray-500" />
            <h2 className="text-base font-bold text-gray-900">
              真实社群原声评论文本样本 ({report.representativeTexts?.length || 0})
            </h2>
          </div>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 pt-2">
            {report.representativeTexts?.map((text, i) => (
              <div key={i} className="p-3.5 bg-gray-50 rounded-xl border border-gray-100 text-xs text-gray-700 leading-relaxed flex items-start justify-between gap-2">
                <span>{text}</span>
              </div>
            ))}
          </div>
        </section>

        {/* 6. 抓取视频与原文链接溯源 */}
        <section className="card-glass p-6 sm:p-8 space-y-4">
          <div className="flex items-center justify-between pb-3 border-b border-gray-100">
            <div className="flex items-center gap-2">
              <Video className="w-5 h-5 text-indigo-600" />
              <h2 className="text-base font-bold text-gray-900">
                原始抓取视频与内容原文链接溯源 ({report.sourceVideos?.length || 0})
              </h2>
            </div>
            <span className="text-xs text-gray-400">客观证据链直达真实页面</span>
          </div>

          {!report.sourceVideos || report.sourceVideos.length === 0 ? (
            <div className="p-6 bg-gray-50/80 rounded-xl border border-gray-200/80 text-center space-y-2">
              <div className="w-10 h-10 rounded-full bg-amber-50 text-amber-500 flex items-center justify-center mx-auto">
                <AlertCircle className="w-5 h-5" />
              </div>
              <p className="text-sm font-bold text-gray-800">无法追溯至原文链接</p>
              <p className="text-xs text-gray-500 max-w-md mx-auto">
                当前数据源平台暂未返回可公开直接访问的内容原链或受平台防盗链保护，系统已停止虚假链接降级处理，确保分析依据客观真实。
              </p>
            </div>
          ) : (
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3.5 pt-2">
              {report.sourceVideos.map((video, vIdx) => (
                <div
                  key={vIdx}
                  className="p-4 bg-gray-50/80 border border-gray-200/80 rounded-xl flex flex-col justify-between gap-2 hover:bg-indigo-50/20 hover:border-indigo-200 transition"
                >
                  <div>
                    <h3 className="font-bold text-xs text-gray-900 leading-snug line-clamp-2 mb-1.5">
                      {video.title}
                    </h3>
                    <div className="flex flex-wrap items-center gap-3 text-[11px] text-gray-500">
                      <span>发布者: <b className="text-gray-700">{video.author}</b></span>
                      <span className="flex items-center gap-1">
                        <ThumbsUp className="w-3 h-3 text-amber-500" />
                        {video.likes}
                      </span>
                      {video.createdAt && <span>时间: {video.createdAt}</span>}
                    </div>
                  </div>

                  <div className="pt-2 border-t border-gray-100/80 flex items-center justify-end">
                    {video.videoUrl ? (
                      <a
                        href={video.videoUrl}
                        target="_blank"
                        rel="noreferrer"
                        className="text-xs font-semibold text-blue-600 hover:text-blue-700 flex items-center gap-1 px-2.5 py-1 rounded-lg hover:bg-blue-50 transition"
                      >
                        <span>访问原文</span>
                        <ExternalLink className="w-3.5 h-3.5" />
                      </a>
                    ) : (
                      <span className="text-xs text-gray-400 bg-gray-100 px-2.5 py-1 rounded-lg">
                        无法追溯至原文链接
                      </span>
                    )}
                  </div>
                </div>
              ))}
            </div>
          )}
        </section>
      </main>
    </div>
  );
};
