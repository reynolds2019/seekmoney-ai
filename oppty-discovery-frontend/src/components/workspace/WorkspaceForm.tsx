import React, { useState, useId } from 'react';
import { Sparkles, Sliders, MessageSquare, Video, Info, Tag, Plus, X, Globe2, AlertCircle, Heart, Tv, Youtube, MessageCircle, Loader2 } from 'lucide-react';
import { DataSourceType } from '../../types';

interface WorkspaceFormProps {
  isLoading: boolean;
  onSubmit: (params: {
    keywords: string[];
    dataSource: DataSourceType;
    deepCrawl: boolean;
    maxVideos: number;
    maxCommentsPerVideo: number;
  }) => void;
}

export const WorkspaceForm: React.FC<WorkspaceFormProps> = ({ isLoading, onSubmit }) => {
  const [dataSource, setDataSource] = useState<DataSourceType>('tiktok');
  const [keywords, setKeywords] = useState<string[]>(['AI变现', '自媒体运营']);
  const [inputKeyword, setInputKeyword] = useState('');
  const [deepCrawl, setDeepCrawl] = useState(true);
  const [maxVideos, setMaxVideos] = useState(5); // 默认设为最小值 5 (范围 5 - 20)
  const [maxCommentsPerVideo, setMaxCommentsPerVideo] = useState(10); // 默认设为最小值 10 (范围 10 - 100)
  const [validationError, setValidationError] = useState('');

  const keywordInputId = useId();
  const maxVideosId = useId();
  const maxCommentsId = useId();

  // 实时根据公式计算预估成本
  // Cost = (V / 20 * 0.01) + (deepCrawl ? V * (C / 20) * 0.01 : 0)
  const videoCost = (maxVideos / 20.0) * 0.01;
  const commentCost = deepCrawl ? maxVideos * (maxCommentsPerVideo / 20.0) * 0.01 : 0;
  const totalCost = (Math.round((videoCost + commentCost) * 100) / 100).toFixed(2);

  const handleAddKeyword = () => {
    const trimmed = inputKeyword.trim();
    if (!trimmed) return;
    if (keywords.includes(trimmed)) {
      setValidationError('该关键词已添加');
      return;
    }
    if (keywords.length >= 5) {
      setValidationError('单次分析最多添加 5 个关键词');
      return;
    }
    setKeywords([...keywords, trimmed]);
    setInputKeyword('');
    setValidationError('');
  };

  const handleRemoveKeyword = (index: number) => {
    setKeywords(keywords.filter((_, i) => i !== index));
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleAddKeyword();
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (keywords.length === 0) {
      setValidationError('请至少输入 1 个目标分析关键词');
      return;
    }
    setValidationError('');
    onSubmit({
      keywords,
      dataSource,
      deepCrawl,
      maxVideos,
      maxCommentsPerVideo,
    });
  };

  return (
    <form onSubmit={handleSubmit} className="card-glass p-5 sm:p-7 space-y-6">
      {/* 头部标题与简洁Tooltip帮助 */}
      <div className="flex items-center justify-between pb-3.5 border-b border-gray-100">
        <div className="flex items-center gap-2">
          <Sliders className="w-5 h-5 text-blue-600" />
          <h2 className="text-base sm:text-lg font-bold text-gray-900">
            商机透视雷达配置
          </h2>
          {/* 帮助悬浮解释气泡 Tooltip */}
          <div className="relative group inline-flex items-center">
            <button
              type="button"
              className="w-5 h-5 rounded-full bg-gray-100 hover:bg-blue-100 text-gray-400 hover:text-blue-600 flex items-center justify-center text-xs font-bold transition-colors cursor-help"
              aria-label="查看配置说明"
            >
              ?
            </button>
            <div className="absolute left-0 bottom-full mb-2 hidden group-hover:flex flex-col w-64 p-2.5 bg-gray-900/95 text-white text-xs rounded-xl shadow-xl backdrop-blur-sm z-30 pointer-events-none leading-relaxed transition-all">
              <span>设定挖掘数据源与抓取规模，系统将实时计算预估成本并启动异步洞察。</span>
              <span className="text-[10px] text-gray-400 mt-1">支持按需调整抓取视频条数与评论深度。</span>
              <div className="absolute left-2.5 top-full w-2 h-2 bg-gray-900/95 rotate-45 -mt-1" />
            </div>
          </div>
        </div>
      </div>

      {validationError && (
        <div className="p-3 bg-red-50 border border-red-200 text-red-700 text-xs rounded-xl flex items-center gap-2">
          <AlertCircle className="w-4 h-4 shrink-0" />
          <span>{validationError}</span>
        </div>
      )}

      {/* 1. 选择数据源：2列3行紧凑网格 */}
      <div>
        <label className="block text-xs font-bold text-gray-700 uppercase tracking-wider mb-2">
          1. 目标数据源平台
        </label>
        <div className="grid grid-cols-2 gap-2.5">
          {/* TikTok */}
          <button
            type="button"
            onClick={() => setDataSource('tiktok')}
            className={`flex items-center gap-2.5 p-3 rounded-xl border text-left transition-all ${
              dataSource === 'tiktok'
                ? 'border-blue-600 bg-blue-50/70 ring-2 ring-blue-600/20'
                : 'border-gray-200 hover:border-gray-300 bg-white'
            }`}
          >
            <div className="w-8 h-8 rounded-lg bg-black text-white flex items-center justify-center font-bold text-sm shrink-0">
              <Globe2 className="w-4 h-4" />
            </div>
            <div className="min-w-0">
              <div className="text-xs font-bold text-gray-900 truncate">TikTok 国际版</div>
              <div className="text-[10px] text-gray-500 truncate">海外全球社群痛点</div>
            </div>
          </button>

          {/* 抖音 */}
          <button
            type="button"
            onClick={() => setDataSource('douyin')}
            className={`flex items-center gap-2.5 p-3 rounded-xl border text-left transition-all ${
              dataSource === 'douyin'
                ? 'border-blue-600 bg-blue-50/70 ring-2 ring-blue-600/20'
                : 'border-gray-200 hover:border-gray-300 bg-white'
            }`}
          >
            <div className="w-8 h-8 rounded-lg bg-gradient-to-tr from-pink-500 to-purple-600 text-white flex items-center justify-center font-bold text-sm shrink-0">
              <Video className="w-4 h-4" />
            </div>
            <div className="min-w-0">
              <div className="text-xs font-bold text-gray-900 truncate">抖音短视频</div>
              <div className="text-[10px] text-gray-500 truncate">国内高频带货热议</div>
            </div>
          </button>

          {/* 小红书 */}
          <button
            type="button"
            onClick={() => setDataSource('xiaohongshu')}
            className={`flex items-center gap-2.5 p-3 rounded-xl border text-left transition-all ${
              dataSource === 'xiaohongshu'
                ? 'border-blue-600 bg-blue-50/70 ring-2 ring-blue-600/20'
                : 'border-gray-200 hover:border-gray-300 bg-white'
            }`}
          >
            <div className="w-8 h-8 rounded-lg bg-gradient-to-tr from-red-500 to-rose-600 text-white flex items-center justify-center font-bold text-sm shrink-0">
              <Heart className="w-4 h-4" />
            </div>
            <div className="min-w-0">
              <div className="text-xs font-bold text-gray-900 truncate">小红书</div>
              <div className="text-[10px] text-gray-500 truncate">真实种草与买家排雷</div>
            </div>
          </button>

          {/* Bilibili */}
          <button
            type="button"
            onClick={() => setDataSource('bilibili')}
            className={`flex items-center gap-2.5 p-3 rounded-xl border text-left transition-all ${
              dataSource === 'bilibili'
                ? 'border-blue-600 bg-blue-50/70 ring-2 ring-blue-600/20'
                : 'border-gray-200 hover:border-gray-300 bg-white'
            }`}
          >
            <div className="w-8 h-8 rounded-lg bg-gradient-to-tr from-sky-400 to-blue-500 text-white flex items-center justify-center font-bold text-sm shrink-0">
              <Tv className="w-4 h-4" />
            </div>
            <div className="min-w-0">
              <div className="text-xs font-bold text-gray-900 truncate">Bilibili</div>
              <div className="text-[10px] text-gray-500 truncate">深度评测与极客吐槽</div>
            </div>
          </button>

          {/* 微信视频号 */}
          <button
            type="button"
            onClick={() => setDataSource('wechat')}
            className={`flex items-center gap-2.5 p-3 rounded-xl border text-left transition-all ${
              dataSource === 'wechat'
                ? 'border-blue-600 bg-blue-50/70 ring-2 ring-blue-600/20'
                : 'border-gray-200 hover:border-gray-300 bg-white'
            }`}
          >
            <div className="w-8 h-8 rounded-lg bg-gradient-to-tr from-emerald-500 to-green-600 text-white flex items-center justify-center font-bold text-sm shrink-0">
              <MessageCircle className="w-4 h-4" />
            </div>
            <div className="min-w-0">
              <div className="text-xs font-bold text-gray-900 truncate">微信视频号</div>
              <div className="text-[10px] text-gray-500 truncate">熟人商业圈层与私域</div>
            </div>
          </button>

          {/* YouTube */}
          <button
            type="button"
            onClick={() => setDataSource('youtube')}
            className={`flex items-center gap-2.5 p-3 rounded-xl border text-left transition-all ${
              dataSource === 'youtube'
                ? 'border-blue-600 bg-blue-50/70 ring-2 ring-blue-600/20'
                : 'border-gray-200 hover:border-gray-300 bg-white'
            }`}
          >
            <div className="w-8 h-8 rounded-lg bg-gradient-to-tr from-red-600 to-red-700 text-white flex items-center justify-center font-bold text-sm shrink-0">
              <Youtube className="w-4 h-4" />
            </div>
            <div className="min-w-0">
              <div className="text-xs font-bold text-gray-900 truncate">YouTube</div>
              <div className="text-[10px] text-gray-500 truncate">全球技术与独立出海</div>
            </div>
          </button>
        </div>
      </div>

      {/* 2. 关键词设置 */}
      <div>
        <div className="flex justify-between items-center mb-2">
          <label htmlFor={keywordInputId} className="text-xs font-bold text-gray-700 uppercase tracking-wider">
            2. 挖掘关键词 (1~5个)
          </label>
          <span className="text-[11px] text-gray-400">回车键即可快速添加</span>
        </div>

        {/* 标签展示区 */}
        <div className="flex flex-wrap gap-2 mb-2 min-h-[32px]">
          {keywords.map((kw, idx) => (
            <span
              key={idx}
              className="inline-flex items-center gap-1.5 px-3 py-1 bg-blue-50 text-blue-700 text-xs font-medium rounded-full border border-blue-200 shadow-sm"
            >
              <Tag className="w-3 h-3 text-blue-500" />
              {kw}
              <button
                type="button"
                onClick={() => handleRemoveKeyword(idx)}
                className="hover:text-blue-900 p-0.5 rounded-full"
                aria-label={`移除关键词 ${kw}`}
              >
                <X className="w-3 h-3" />
              </button>
            </span>
          ))}
        </div>

        <div className="flex gap-2">
          <input
            id={keywordInputId}
            type="text"
            placeholder="输入细分领域关键词，如：AI绘画、小红书副业、跨境独立站"
            value={inputKeyword}
            onChange={(e) => setInputKeyword(e.target.value)}
            onKeyDown={handleKeyDown}
            disabled={isLoading}
            className="input-field text-xs sm:text-sm"
          />
          <button
            type="button"
            onClick={handleAddKeyword}
            disabled={isLoading || !inputKeyword.trim()}
            className="btn btn-secondary text-xs px-3 shrink-0"
          >
            <Plus className="w-3.5 h-3.5" />
            添加
          </button>
        </div>
      </div>

      {/* 3. 深度爬取配置开关 */}
      <div className="p-4 bg-gray-50 rounded-2xl border border-gray-200/80 space-y-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <MessageSquare className="w-4 h-4 text-blue-600" />
            <span className="text-xs font-bold text-gray-900">开启[爬取评论]深度分析</span>
          </div>
          <label className="relative inline-flex items-center cursor-pointer">
            <input
              type="checkbox"
              checked={deepCrawl}
              onChange={(e) => setDeepCrawl(e.target.checked)}
              disabled={isLoading}
              className="sr-only peer"
            />
            <div className="w-11 h-6 bg-gray-300 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
          </label>
        </div>

        {/* 视频数量滑块 (5 - 20) */}
        <div>
          <div className="flex justify-between items-center mb-1.5">
            <label htmlFor={maxVideosId} className="text-xs font-medium text-gray-600">
              视频抓取数量：
            </label>
            <span className="text-xs font-bold font-mono text-gray-900 bg-white px-2 py-0.5 rounded-md border border-gray-200 shadow-sm">
              {maxVideos} 条
            </span>
          </div>
          <input
            id={maxVideosId}
            type="range"
            min={5}
            max={20}
            step={1}
            value={maxVideos}
            onChange={(e) => setMaxVideos(Number(e.target.value))}
            disabled={isLoading}
            className="w-full accent-blue-600 cursor-pointer"
          />
          <div className="flex justify-between text-[10px] text-gray-400 mt-0.5">
            <span>5条 (快速验证)</span>
            <span>20条 (深度透视)</span>
          </div>
        </div>

        {/* 每视频评论数滑块 (10 - 100) - 仅在开启爬取评论时生效 */}
        {deepCrawl && (
          <div className="pt-2 border-t border-gray-200/60">
            <div className="flex justify-between items-center mb-1.5">
              <label htmlFor={maxCommentsId} className="text-xs font-medium text-gray-600">
                每视频评论数：
              </label>
              <span className="text-xs font-bold font-mono text-gray-900 bg-white px-2 py-0.5 rounded-md border border-gray-200 shadow-sm">
                {maxCommentsPerVideo} 条 / 视频
              </span>
            </div>
            <input
              id={maxCommentsId}
              type="range"
              min={10}
              max={100}
              step={5}
              value={maxCommentsPerVideo}
              onChange={(e) => setMaxCommentsPerVideo(Number(e.target.value))}
              disabled={isLoading}
              className="w-full accent-blue-600 cursor-pointer"
            />
            <div className="flex justify-between text-[10px] text-gray-400 mt-0.5">
              <span>10条</span>
              <span>100条</span>
            </div>
          </div>
        )}

        {/* 成本预估动态计算器 */}
        <div className="pt-3 border-t border-gray-200 flex items-center justify-between">
          <div className="flex items-center gap-1.5 text-xs text-gray-600">
            <Info className="w-3.5 h-3.5 text-blue-500" />
            <span>系统预估成本：</span>
          </div>
          <div className="text-right">
            <span className="text-base font-extrabold font-mono text-blue-600">
              ~${totalCost}
            </span>
            <span className="text-[10px] text-gray-400 ml-1">/ 次分析</span>
          </div>
        </div>
      </div>

      {/* 4. 提交分析按钮 */}
      <button
        type="submit"
        disabled={isLoading || keywords.length === 0}
        className="w-full btn btn-primary py-3.5 text-sm sm:text-base shadow-lg shadow-blue-500/25 flex items-center justify-center gap-2 group disabled:opacity-60 disabled:cursor-not-allowed transition-all"
      >
        {isLoading ? (
          <>
            <Loader2 className="w-5 h-5 animate-spin text-white" />
            <span>商机雷达挖掘中...</span>
          </>
        ) : (
          <>
            <Sparkles className="w-4 h-4 group-hover:rotate-12 transition-transform" />
            <span>立即寻找商机</span>
          </>
        )}
      </button>
    </form>
  );
};
