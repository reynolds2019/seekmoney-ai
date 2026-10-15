import React, { useState, useEffect } from 'react';
import { History, Search, Trash2, ArrowRight, Loader2, Sparkles, Calendar, Globe2, AlertCircle } from 'lucide-react';
import { OpptyHistoryItemVO } from '../../types';
import { opptyApi } from '../../api/client';

interface HistoryViewProps {
  onSelectJob: (jobId: string) => void;
}

export const HistoryView: React.FC<HistoryViewProps> = ({ onSelectJob }) => {
  const [historyList, setHistoryList] = useState<OpptyHistoryItemVO[]>([]);
  const [total, setTotal] = useState(0);
  const [pageNo, setPageNo] = useState(1);
  const [keyword, setKeyword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState('');

  const fetchHistory = async () => {
    setIsLoading(true);
    setErrorMsg('');
    try {
      const resp = await opptyApi.getHistory(pageNo, 10, keyword);
      setHistoryList(resp.records || []);
      setTotal(resp.total || 0);
    } catch (err: any) {
      setErrorMsg(err.message || '获取商机历史列表失败');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchHistory();
  }, [pageNo]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setPageNo(1);
    fetchHistory();
  };

  const handleDelete = async (jobId: string, e: React.MouseEvent) => {
    e.stopPropagation();
    if (!window.confirm('确定要删除该条商机分析历史记录吗？')) {
      return;
    }
    try {
      await opptyApi.deleteHistory(jobId);
      setHistoryList(historyList.filter((item) => item.jobId !== jobId));
      setTotal((prev) => Math.max(0, prev - 1));
    } catch (err: any) {
      alert(err.message || '删除失败');
    }
  };

  return (
    <div className="space-y-6">
      {/* 头部标题与搜索过滤 */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h2 className="text-lg font-extrabold text-gray-900 flex items-center gap-2">
            <History className="w-5 h-5 text-blue-600" />
            商机发现历史中心 ({total})
          </h2>
          <p className="text-xs text-gray-500 mt-0.5">
            所有生成的商机雷达报告均已自动云端归档，点击即可随时复现与全屏研读
          </p>
        </div>

        <form onSubmit={handleSearch} className="flex gap-2 w-full sm:w-auto">
          <div className="relative flex-1 sm:w-64">
            <Search className="w-4 h-4 text-gray-400 absolute left-3 top-3 pointer-events-none" />
            <input
              type="text"
              placeholder="搜索历史关键词..."
              value={keyword}
              onChange={(e) => setKeyword(e.target.value)}
              className="input-field pl-9 text-xs sm:text-sm py-2"
            />
          </div>
          <button type="submit" className="btn btn-secondary text-xs px-3">
            搜索
          </button>
        </form>
      </div>

      {errorMsg && (
        <div className="p-3 bg-red-50 border border-red-200 text-red-700 text-xs rounded-xl flex items-center gap-2">
          <AlertCircle className="w-4 h-4 shrink-0" />
          <span>{errorMsg}</span>
        </div>
      )}

      {/* 列表区域 */}
      {isLoading ? (
        <div className="card-glass p-12 text-center text-gray-400 flex flex-col items-center justify-center gap-3">
          <Loader2 className="w-6 h-6 animate-spin text-blue-600" />
          <span className="text-xs">正在加载历史商机档案...</span>
        </div>
      ) : historyList.length === 0 ? (
        <div className="card-glass p-12 text-center text-gray-400 space-y-2">
          <Sparkles className="w-8 h-8 text-gray-300 mx-auto" />
          <p className="text-sm font-semibold text-gray-700">暂无商机历史记录</p>
          <p className="text-xs text-gray-400">去工作台发起一次智能分析，系统将自动沉淀报告</p>
        </div>
      ) : (
        <div className="space-y-3">
          {historyList.map((item) => (
            <div
              key={item.jobId}
              onClick={() => onSelectJob(item.jobId)}
              className="card-glass p-4 sm:p-5 flex flex-col sm:flex-row sm:items-center justify-between gap-4 hover:border-blue-400 cursor-pointer transition-all group"
            >
              <div className="space-y-1.5 flex-1 min-w-0">
                <div className="flex flex-wrap items-center gap-2">
                  <span className="text-xs font-bold text-gray-900 bg-blue-50 text-blue-700 px-2 py-0.5 rounded-md border border-blue-200">
                    {item.keywords?.join(' / ')}
                  </span>
                  <span className="text-[11px] text-gray-400 flex items-center gap-1">
                    <Globe2 className="w-3 h-3" />
                    {item.dataSource === 'tiktok' ? 'TikTok' : '抖音'}
                  </span>
                  <span className="text-[11px] text-gray-400 flex items-center gap-1">
                    <Calendar className="w-3 h-3" />
                    {item.createdAt}
                  </span>
                </div>

                <h3 className="text-sm font-bold text-gray-800 group-hover:text-blue-600 transition truncate">
                  {item.primaryPain}
                </h3>
              </div>

              <div className="flex items-center justify-between sm:justify-end gap-3 pt-2 sm:pt-0 border-t sm:border-t-0 border-gray-100">
                <div className="text-right">
                  <span className="text-xs font-mono font-bold text-gray-700 block">
                    {item.reportCount} 组商机
                  </span>
                  <span className="text-[10px] text-gray-400 font-mono">
                    成本: ${item.estimatedCost}
                  </span>
                </div>

                <button
                  type="button"
                  onClick={(e) => handleDelete(item.jobId, e)}
                  className="p-1.5 rounded-lg text-gray-400 hover:text-red-600 hover:bg-red-50 transition"
                  title="删除记录"
                >
                  <Trash2 className="w-4 h-4" />
                </button>

                <div className="text-gray-400 group-hover:text-blue-600 group-hover:translate-x-1 transition-all">
                  <ArrowRight className="w-4 h-4" />
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};
