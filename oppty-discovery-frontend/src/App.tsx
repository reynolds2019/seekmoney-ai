import React, { useState, useEffect, useRef } from 'react';
import { Navbar } from './components/layout/Navbar';
import { LoginModal } from './components/auth/LoginModal';
import { WorkspaceForm } from './components/workspace/WorkspaceForm';
import { ProgressCard } from './components/workspace/ProgressCard';
import { ReportView } from './components/report/ReportView';
import { FullscreenReportModal } from './components/report/FullscreenReportModal';
import { HistoryView } from './components/history/HistoryView';
import { opptyApi } from './api/client';
import { OpptyJobDetailVO, OpptyReportVO, DataSourceType } from './types';
import { Sparkles, ShieldCheck, Zap } from 'lucide-react';

export const App: React.FC = () => {
  const [currentTab, setCurrentTab] = useState<'workspace' | 'history'>('workspace');
  const [currentJobId, setCurrentJobId] = useState<string | null>(null);
  const [currentJob, setCurrentJob] = useState<OpptyJobDetailVO | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [fullscreenReport, setFullscreenReport] = useState<OpptyReportVO | null>(null);
  const pollTimerRef = useRef<number | null>(null);

  // 轮询当前任务状态
  useEffect(() => {
    if (!currentJobId) return;

    const fetchStatus = async () => {
      try {
        const detail = await opptyApi.getJobDetail(currentJobId);
        setCurrentJob(detail);

        if (detail.status === 'completed' || detail.status === 'failed') {
          setIsSubmitting(false);
          if (pollTimerRef.current) {
            clearInterval(pollTimerRef.current);
            pollTimerRef.current = null;
          }
        }
      } catch (err) {
        console.error('轮询任务状态失败:', err);
      }
    };

    fetchStatus();
    pollTimerRef.current = setInterval(fetchStatus, 2000);

    return () => {
      if (pollTimerRef.current) {
        clearInterval(pollTimerRef.current);
        pollTimerRef.current = null;
      }
    };
  }, [currentJobId]);

  // 处理提交智能分析
  const handleSubmitJob = async (params: {
    keywords: string[];
    dataSource: DataSourceType;
    deepCrawl: boolean;
    maxVideos: number;
    maxCommentsPerVideo: number;
  }) => {
    setIsSubmitting(true);
    setCurrentJob(null);
    try {
      const resp = await opptyApi.createJob(params);
      setCurrentJobId(resp.jobId);
    } catch (err: any) {
      alert(err.message || '创建任务失败，请检查服务状态');
      setIsSubmitting(false);
    }
  };

  // 从历史记录中选择查看某份商机报告
  const handleSelectHistoryJob = (jobId: string) => {
    setCurrentJobId(jobId);
    setCurrentTab('workspace');
  };

  return (
    <div className="min-h-screen flex flex-col relative">
      {/* 背景柔和氛围光晕 */}
      <div className="bg-ambient" />

      {/* 顶部导航 */}
      <Navbar currentTab={currentTab} onTabChange={setCurrentTab} />

      {/* 主体工作区容器 */}
      <main className="flex-1 max-w-7xl mx-auto w-full px-4 sm:px-6 lg:px-8 py-6 sm:py-10">
        {currentTab === 'workspace' ? (
          <div className="space-y-8">
            {/* Slogan 标语横幅 */}
            <div className="text-center max-w-3xl mx-auto space-y-2 mb-8">
              <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-blue-50 border border-blue-200/80 text-blue-700 text-xs font-semibold mb-1 shadow-sm">
                <Sparkles className="w-3.5 h-3.5 text-blue-600" />
                <span>智能商机雷达 · 深度评论挖掘</span>
              </div>
              <h1 className="text-2xl sm:text-4xl font-extrabold text-gray-900 tracking-tight">
                痛点一键透视，商机即刻成型
              </h1>
              <p className="text-xs sm:text-sm text-gray-500 max-w-2xl mx-auto leading-relaxed">
                一款集“社区评论挖掘、痛点线索提炼与商业方案构思”于一体的智能商机雷达。帮助创业者在几秒内从海量社交原声中锁定真正值得做的爆款方向。
              </p>
            </div>

            {/* 左右分栏：左侧配置表单，右侧状态感知与分析报告 */}
            <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
              <div className="lg:col-span-5">
                <WorkspaceForm isLoading={isSubmitting} onSubmit={handleSubmitJob} />
              </div>

              <div className="lg:col-span-7 space-y-6">
                {currentJob && <ProgressCard job={currentJob} />}

                {currentJob && currentJob.results && currentJob.results.length > 0 ? (
                  <ReportView
                    reports={currentJob.results}
                    onOpenFullscreen={(rep) => setFullscreenReport(rep)}
                  />
                ) : !currentJob ? (
                  <div className="card-glass p-8 sm:p-12 text-center text-gray-400 space-y-4">
                    <div className="w-16 h-16 rounded-2xl bg-gray-50 border border-gray-100 flex items-center justify-center mx-auto text-blue-500 shadow-inner">
                      <Zap className="w-8 h-8" />
                    </div>
                    <div>
                      <h3 className="text-base font-bold text-gray-800">等待启动商机雷达</h3>
                      <p className="text-xs text-gray-400 mt-1 max-w-md mx-auto">
                        在左侧选择数据源并设定关键词，点击“开始智能分析”即可在此实时查看多维商业报告与全屏研读。
                      </p>
                    </div>
                    <div className="flex justify-center gap-4 text-[11px] text-gray-500 pt-2">
                      <span className="flex items-center gap-1">
                        <ShieldCheck className="w-3.5 h-3.5 text-emerald-500" /> 真实社群原声
                      </span>
                      <span className="flex items-center gap-1">
                        <ShieldCheck className="w-3.5 h-3.5 text-blue-500" /> 语义密度聚类
                      </span>
                      <span className="flex items-center gap-1">
                        <ShieldCheck className="w-3.5 h-3.5 text-indigo-500" /> 全屏沉浸研读
                      </span>
                    </div>
                  </div>
                ) : null}
              </div>
            </div>
          </div>
        ) : (
          <HistoryView onSelectJob={handleSelectHistoryJob} />
        )}
      </main>

      {/* 全屏沉浸模式模态框 */}
      {fullscreenReport && (
        <FullscreenReportModal
          report={fullscreenReport}
          onClose={() => setFullscreenReport(null)}
        />
      )}

      {/* 登录即注册无感弹窗 */}
      <LoginModal />

      {/* 页脚 - 彻底清除第三方旧品牌 */}
      <footer className="mt-16 py-6 border-t border-gray-200 bg-white/50 text-center text-xs text-gray-400 space-y-1">
        <p className="font-semibold text-gray-600">
          商机发现助手 (Oppty Discovery Assistant) · 痛点一键透视，商机即刻成型
        </p>
        <p className="text-[11px]">
          From pain point to profit in seconds. © 2026 All rights reserved.
        </p>
      </footer>
    </div>
  );
};
