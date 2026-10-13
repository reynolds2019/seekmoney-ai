import React, { useState, useEffect, useRef } from 'react';
import { Navbar } from './components/layout/Navbar';
import { LoginModal } from './components/auth/LoginModal';
import { WorkspaceForm } from './components/workspace/WorkspaceForm';
import { ProgressCard } from './components/workspace/ProgressCard';
import { ReportView } from './components/report/ReportView';
import { FullscreenReportModal } from './components/report/FullscreenReportModal';
import { HistoryView } from './components/history/HistoryView';
import { AboutView } from './components/about/AboutView';
import { opptyApi } from './api/client';
import { OpptyJobDetailVO, OpptyReportVO, DataSourceType } from './types';
import { ShieldCheck, Zap, Radar, ArrowLeft } from 'lucide-react';

export const App: React.FC = () => {
  const isAboutPage = typeof window !== 'undefined' && new URLSearchParams(window.location.search).get('view') === 'about';
  const shareJobId = typeof window !== 'undefined' ? new URLSearchParams(window.location.search).get('share') : null;

  const [currentTab, setCurrentTab] = useState<'workspace' | 'history'>('workspace');
  const [currentJobId, setCurrentJobId] = useState<string | null>(null);
  const [currentJob, setCurrentJob] = useState<OpptyJobDetailVO | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [fullscreenReport, setFullscreenReport] = useState<OpptyReportVO | null>(null);
  const [shareLoading, setShareLoading] = useState(false);
  const [shareError, setShareError] = useState<string | null>(null);
  const pollTimerRef = useRef<number | null>(null);

  // 若处于公开分享直达模式，免登录拉取报告
  useEffect(() => {
    if (!shareJobId) return;

    const fetchShareReport = async () => {
      setShareLoading(true);
      setShareError(null);
      try {
        const detail = await opptyApi.getShareJobDetail(shareJobId);
        setCurrentJob(detail);
      } catch (err: any) {
        setShareError(err.message || '获取分享商机报告失败，该链接可能已失效');
      } finally {
        setShareLoading(false);
      }
    };

    fetchShareReport();
  }, [shareJobId]);

  // 轮询当前任务状态
  useEffect(() => {
    if (!currentJobId || shareJobId) return;

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
  }, [currentJobId, shareJobId]);

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

  // 针对 ?view=about 参数，作为独立新标签页全屏渲染展示
  if (isAboutPage) {
    return (
      <div className="min-h-screen flex flex-col relative bg-white">
        <div className="bg-ambient" />
        <header className="sticky top-0 z-50 bg-white/80 backdrop-blur-md border-b border-gray-200">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
            <div
              className="flex items-center gap-3 cursor-pointer"
              onClick={() => (window.location.href = window.location.pathname)}
            >
              <div className="w-10 h-10 rounded-xl bg-gradient-to-tr from-blue-600 to-indigo-500 flex items-center justify-center text-white shadow-md shadow-blue-500/20">
                <Radar className="w-5 h-5 animate-pulse" />
              </div>
              <div className="flex items-center gap-2">
                <span className="font-extrabold text-lg text-gray-900 tracking-tight">商机发现助手</span>
                <span className="text-xs font-semibold px-2 py-0.5 bg-blue-50 text-blue-700 rounded-full border border-blue-200">
                  产品介绍
                </span>
              </div>
            </div>
            <a
              href={window.location.pathname}
              className="btn btn-primary text-xs py-2 px-4 flex items-center gap-1.5"
            >
              <ArrowLeft className="w-3.5 h-3.5" />
              <span>进入工作台</span>
            </a>
          </div>
        </header>

        <main className="flex-1 max-w-7xl mx-auto w-full px-4 sm:px-6 lg:px-8 py-6">
          <AboutView onStartDiscovery={() => (window.location.href = window.location.pathname)} />
        </main>

        <footer className="py-6 border-t border-gray-200 bg-white/50 text-center text-xs text-gray-400 space-y-1">
          <p className="font-semibold text-gray-600">
            商机发现助手 (Oppty Discovery Assistant) · 痛点一键透视，商机即刻成型
          </p>
          <p className="text-[11px]">
            From pain point to profit in seconds. © 2026 All rights reserved.
          </p>
        </footer>
      </div>
    );
  }

  // 针对 ?share=:jobId 参数，作为免登录公开分享直达视图全屏研读展示
  if (shareJobId) {
    return (
      <div className="min-h-screen flex flex-col relative bg-slate-50">
        <div className="bg-ambient" />

        {/* 分享专属顶栏 */}
        <header className="sticky top-0 z-40 bg-white/90 backdrop-blur-md border-b border-gray-200">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
            <div
              className="flex items-center gap-3 cursor-pointer"
              onClick={() => (window.location.href = window.location.pathname)}
            >
              <div className="w-10 h-10 rounded-xl bg-gradient-to-tr from-blue-600 to-indigo-500 flex items-center justify-center text-white shadow-md shadow-blue-500/20">
                <Radar className="w-5 h-5 animate-pulse" />
              </div>
              <div className="flex items-center gap-2">
                <span className="font-extrabold text-base text-gray-900 tracking-tight">商机发现助手</span>
                <span className="text-xs font-semibold px-2 py-0.5 bg-emerald-50 text-emerald-700 rounded-full border border-emerald-200 flex items-center gap-1">
                  <ShieldCheck className="w-3 h-3 text-emerald-600" />
                  公开分享只读模式
                </span>
              </div>
            </div>

            <a
              href={window.location.pathname}
              className="btn btn-primary text-xs py-2 px-4 flex items-center gap-1.5 shadow-sm"
            >
              <span>体验自己的商机挖掘</span>
            </a>
          </div>
        </header>

        {/* 分享内容主体展示 */}
        <main className="flex-1 max-w-5xl mx-auto w-full px-4 sm:px-6 lg:px-8 py-8">
          {shareLoading ? (
            <div className="card-glass p-12 text-center text-gray-400 space-y-3">
              <div className="w-10 h-10 rounded-full border-2 border-blue-600 border-t-transparent animate-spin mx-auto" />
              <p className="text-sm font-medium text-gray-600">正在免登录加载商机透视报告，请稍候...</p>
            </div>
          ) : shareError ? (
            <div className="card-glass p-12 text-center text-red-600 space-y-3 border-red-200">
              <p className="text-sm font-bold">{shareError}</p>
              <a href={window.location.pathname} className="btn btn-secondary text-xs inline-block">
                返回工作台
              </a>
            </div>
          ) : currentJob && currentJob.results && currentJob.results.length > 0 ? (
            <div className="space-y-6">
              <ReportView
                reports={currentJob.results}
                jobId={currentJob.jobId}
                sourceVideos={currentJob.sourceVideos}
                onOpenFullscreen={(rep) => setFullscreenReport(rep)}
              />
            </div>
          ) : (
            <div className="card-glass p-12 text-center text-gray-400">
              <p className="text-sm">该商机分析仍在进行中或尚未生成有效透视结果</p>
            </div>
          )}
        </main>

        {fullscreenReport && (
          <FullscreenReportModal
            report={fullscreenReport}
            jobId={currentJob?.jobId}
            onClose={() => setFullscreenReport(null)}
          />
        )}

        <footer className="py-6 border-t border-gray-200 bg-white/50 text-center text-xs text-gray-400 space-y-1">
          <p className="font-semibold text-gray-600">
            商机发现助手 (Oppty Discovery Assistant) · 痛点一键透视，商机即刻成型
          </p>
          <p className="text-[11px]">
            From pain point to profit in seconds. © 2026 All rights reserved.
          </p>
        </footer>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex flex-col relative">
      {/* 背景柔和氛围光晕 */}
      <div className="bg-ambient" />

      {/* 顶部导航 */}
      <Navbar currentTab={currentTab} onTabChange={setCurrentTab} />

      {/* 主体工作区容器 */}
      <main className="flex-1 max-w-7xl mx-auto w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-8">
        {currentTab === 'workspace' ? (
          <div className="space-y-6">
            {/* 左右分栏：左侧配置表单收窄为 4 列，右侧报告与研读区拓展为 8 列 */}
            <div className="grid grid-cols-1 lg:grid-cols-12 gap-7 items-start">
              <div className="lg:col-span-4">
                <WorkspaceForm isLoading={isSubmitting} onSubmit={handleSubmitJob} />
              </div>

              <div className="lg:col-span-8 space-y-6">
                {/* 进度感知卡片仅在生成过程中或异常中断时显示，分析完成呈现报告后自动隐藏 */}
                {currentJob && currentJob.status !== 'completed' && <ProgressCard job={currentJob} />}

                {currentJob && currentJob.results && currentJob.results.length > 0 ? (
                  <ReportView
                    reports={currentJob.results}
                    jobId={currentJob.jobId}
                    sourceVideos={currentJob.sourceVideos}
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
                        在左侧选择数据源并设定关键词，点击“开始寻找商机”即可在此实时查看多维商业报告与全屏研读。
                      </p>
                    </div>
                    <div className="flex justify-center gap-4 text-[11px] text-gray-500 pt-2">
                      <span className="flex items-center gap-1">
                        <ShieldCheck className="w-3.5 h-3.5 text-emerald-500" /> 真实社群原声
                      </span>
                      <span className="flex items-center gap-1">
                        <ShieldCheck className="w-3.5 h-3.5 text-blue-500" /> 智能商业提炼
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
          jobId={currentJob?.jobId}
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
