import React from 'react';
import { Loader2, CheckCircle2, AlertTriangle, ArrowRight, Sparkles, Network, DownloadCloud } from 'lucide-react';
import { OpptyJobDetailVO } from '../../types';

interface ProgressCardProps {
  job: OpptyJobDetailVO;
}

export const ProgressCard: React.FC<ProgressCardProps> = ({ job }) => {
  const steps = [
    { key: 'crawling', name: '真实声量采集', desc: '采集短视频原声与真实评论', icon: DownloadCloud },
    { key: 'clustering', name: '用户吐槽聚合', desc: '提炼高频痛点与核心诉求', icon: Network },
    { key: 'analyzing', name: '商机方案成型', desc: '梳理市场空白与行动方案', icon: Sparkles },
  ];

  const formatProgressMessage = (msg: string) => {
    if (!msg) return '';
    return msg
      .replace(/Qwen\s*大模型深度商业洞察/g, '正在深度提炼商业方案与行动路线')
      .replace(/Qwen|大模型/g, '智能分析')
      .replace(/语义密度聚类/g, '用户反馈整理')
      .replace(/MVP/g, '轻量启动方案');
  };

  const getStepStatus = (stepKey: string) => {
    if (job.status === 'completed') return 'done';
    if (job.status === 'failed') return 'failed';

    const stageOrder = ['init', 'validating', 'crawling', 'clustering', 'analyzing', 'completed'];
    const currentIndex = stageOrder.indexOf(job.progressStage);
    const stepIndex = stageOrder.indexOf(stepKey);

    if (currentIndex > stepIndex) return 'done';
    if (currentIndex === stepIndex) return 'active';
    return 'pending';
  };

  return (
    <div className="card-glass p-6 sm:p-7 relative overflow-hidden">
      {/* 装饰发光条 */}
      <div className="absolute top-0 left-0 right-0 h-1 bg-gradient-to-r from-blue-500 via-indigo-500 to-amber-500 animate-pulse"></div>

      <div className="flex items-center justify-between mb-5">
        <div className="flex items-center gap-3">
          <div className="w-9 h-9 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center">
            {job.status === 'failed' ? (
              <AlertTriangle className="w-5 h-5 text-red-500" />
            ) : job.status === 'completed' ? (
              <CheckCircle2 className="w-5 h-5 text-emerald-600" />
            ) : (
              <Loader2 className="w-5 h-5 animate-spin text-blue-600" />
            )}
          </div>
          <div>
            <h3 className="text-sm sm:text-base font-bold text-gray-900">
              {job.status === 'completed'
                ? '商机全维分析已完成'
                : job.status === 'failed'
                ? '任务分析中断'
                : '正在从海量真实社群中寻找高价值商机...'}
            </h3>
            <p className="text-xs text-gray-500 font-mono mt-0.5">
              任务标识: {job.jobId.slice(0, 8)}... | 检索词: {job.keywords?.join('、')}
            </p>
          </div>
        </div>
        <div className="text-right">
          <span className="text-lg sm:text-2xl font-black font-mono text-blue-600">
            {job.progressPercent}%
          </span>
        </div>
      </div>

      {/* 进度条 */}
      <div className="w-full bg-gray-100 rounded-full h-2 mb-6 overflow-hidden">
        <div
          className={`h-full rounded-full transition-all duration-500 ease-out ${
            job.status === 'failed' ? 'bg-red-500' : 'bg-gradient-to-r from-blue-600 to-indigo-600'
          }`}
          style={{ width: `${Math.max(job.progressPercent, 5)}%` }}
        ></div>
      </div>

      {/* 步骤条卡片 */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
        {steps.map((step, idx) => {
          const status = getStepStatus(step.key);
          const Icon = step.icon;

          return (
            <div
              key={idx}
              className={`p-3.5 rounded-xl border transition-all flex items-start gap-3 ${
                status === 'done'
                  ? 'bg-emerald-50/50 border-emerald-200'
                  : status === 'active'
                  ? 'bg-blue-50/70 border-blue-300 ring-2 ring-blue-500/20'
                  : 'bg-gray-50/50 border-gray-100 opacity-60'
              }`}
            >
              <div
                className={`w-7 h-7 rounded-lg flex items-center justify-center shrink-0 mt-0.5 ${
                  status === 'done'
                    ? 'bg-emerald-600 text-white'
                    : status === 'active'
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-200 text-gray-500'
                }`}
              >
                {status === 'done' ? (
                  <CheckCircle2 className="w-4 h-4" />
                ) : status === 'active' ? (
                  <Loader2 className="w-4 h-4 animate-spin" />
                ) : (
                  <Icon className="w-4 h-4" />
                )}
              </div>
              <div className="min-w-0">
                <div className="text-xs font-bold text-gray-900">{step.name}</div>
                <div className="text-[11px] text-gray-500 truncate">{step.desc}</div>
              </div>
            </div>
          );
        })}
      </div>

      {/* 实时阶段消息描述 */}
      <div className="mt-4 pt-3 border-t border-gray-100 flex items-center justify-between text-xs text-gray-500">
        <div className="flex items-center gap-1.5 truncate">
          <ArrowRight className="w-3.5 h-3.5 text-blue-500 shrink-0" />
          <span className="font-medium text-gray-700 truncate">
            {formatProgressMessage(job.progressMessage)}
          </span>
        </div>
        <span className="text-gray-400 font-mono shrink-0">预估成本: ${job.estimatedCost}</span>
      </div>
    </div>
  );
};
