import React from 'react';
import {
  Radar,
  Sparkles,
  MessageSquareQuote,
  Flame,
  Rocket,
  ShieldCheck,
  ArrowRight,
} from 'lucide-react';

interface AboutViewProps {
  onStartDiscovery: () => void;
}

export const AboutView: React.FC<AboutViewProps> = ({ onStartDiscovery }) => {
  return (
    <div className="space-y-12 py-4 sm:py-8 max-w-5xl mx-auto">
      {/* 头部品牌横幅区 */}
      <section className="text-center space-y-4">
        <div className="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full bg-blue-50 border border-blue-200 text-blue-700 text-xs font-semibold shadow-sm">
          <Radar className="w-4 h-4 text-blue-600 animate-pulse" />
          <span>新一代智能商机雷达 · Oppty Discovery Assistant</span>
        </div>

        <h1 className="text-3xl sm:text-5xl font-black text-gray-900 tracking-tight leading-tight">
          痛点一键透视，<span className="text-transparent bg-clip-text bg-gradient-to-r from-blue-600 to-indigo-600">商机即刻成型</span>
        </h1>

        <p className="text-sm sm:text-base font-medium text-gray-400 font-sans tracking-wide">
          From pain point to profit in seconds.
        </p>

        <p className="text-sm sm:text-base text-gray-600 max-w-2xl mx-auto leading-relaxed pt-2">
          一款集“社区评论挖掘、痛点线索提炼与商业方案构思”于一体的智能商机雷达。帮助创业者、独立开发者和产品经理在几秒内从海量社交原声中锁定真正值得做的爆款方向。
        </p>

        <div className="pt-4 flex items-center justify-center gap-4">
          <button
            onClick={onStartDiscovery}
            className="btn btn-primary text-sm py-2.5 px-6 shadow-lg shadow-blue-500/25 flex items-center gap-2"
          >
            <Sparkles className="w-4 h-4" />
            <span>立即开启商机发现</span>
            <ArrowRight className="w-4 h-4" />
          </button>
        </div>
      </section>

      {/* 核心价值三大支柱 */}
      <section className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="card-glass p-6 sm:p-7 space-y-3 hover:border-blue-300 transition-all">
          <div className="w-10 h-10 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center shadow-sm">
            <MessageSquareQuote className="w-5 h-5" />
          </div>
          <h3 className="text-base font-bold text-gray-900">真实社群原声捕获</h3>
          <p className="text-xs text-gray-600 leading-relaxed">
            告别主观臆断和失真的问卷调研，直接下潜至主流视频平台的真实评论区，捕捉最原始、最强烈的一线用户槽点。
          </p>
        </div>

        <div className="card-glass p-6 sm:p-7 space-y-3 hover:border-indigo-300 transition-all">
          <div className="w-10 h-10 rounded-xl bg-indigo-50 text-indigo-600 flex items-center justify-center shadow-sm">
            <Flame className="w-5 h-5" />
          </div>
          <h3 className="text-base font-bold text-gray-900">痛点分级与意愿透视</h3>
          <p className="text-xs text-gray-600 leading-relaxed">
            自动过滤噪音，按照情绪烈度、出现频次和付费意愿对用户痛点进行多维解构，精准定位用户“愿意掏钱解决”的问题。
          </p>
        </div>

        <div className="card-glass p-6 sm:p-7 space-y-3 hover:border-amber-300 transition-all">
          <div className="w-10 h-10 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center shadow-sm">
            <Rocket className="w-5 h-5" />
          </div>
          <h3 className="text-base font-bold text-gray-900">快速验证与启动方案</h3>
          <p className="text-xs text-gray-600 leading-relaxed">
            不只是列出问题，更直接输出差异化市场切入点、必备核心功能与种子获客渠道，助您将一个抱怨转化为可落地的商业产品。
          </p>
        </div>
      </section>

      {/* 三步找到你的商业爆款 */}
      <section className="card-glass p-6 sm:p-10 space-y-8 bg-gradient-to-br from-white via-blue-50/20 to-indigo-50/30">
        <div className="text-center space-y-2">
          <h2 className="text-xl sm:text-2xl font-black text-gray-900">三步从用户吐槽到商业方案</h2>
          <p className="text-xs sm:text-sm text-gray-500">零门槛操作，几分钟内完成一次细分赛道的敏捷市场洞察</p>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-3 gap-6 relative">
          <div className="p-5 bg-white rounded-2xl border border-gray-100 shadow-sm space-y-3 relative">
            <div className="w-8 h-8 rounded-lg bg-blue-600 text-white font-black text-sm flex items-center justify-center">
              1
            </div>
            <h4 className="text-sm font-bold text-gray-900">输入赛道关键词</h4>
            <p className="text-xs text-gray-600 leading-relaxed">
              输入你感兴趣的行业或话题（如“短视频剪辑”、“跨境电商运营”），系统自动绑定目标搜索域。
            </p>
          </div>

          <div className="p-5 bg-white rounded-2xl border border-gray-100 shadow-sm space-y-3 relative">
            <div className="w-8 h-8 rounded-lg bg-indigo-600 text-white font-black text-sm flex items-center justify-center">
              2
            </div>
            <h4 className="text-sm font-bold text-gray-900">设定抓取规模与预算</h4>
            <p className="text-xs text-gray-600 leading-relaxed">
              按需拖动视频数与每视频评论数滑块，系统动态透明预估消耗成本，按需开启深度评论挖掘。
            </p>
          </div>

          <div className="p-5 bg-white rounded-2xl border border-gray-100 shadow-sm space-y-3 relative">
            <div className="w-8 h-8 rounded-lg bg-amber-600 text-white font-black text-sm flex items-center justify-center">
              3
            </div>
            <h4 className="text-sm font-bold text-gray-900">获取全屏商业报告</h4>
            <p className="text-xs text-gray-600 leading-relaxed">
              即刻获取结构化商机透视矩阵，查看真实声量原声佐证，支持全屏沉浸研读并永久留存于历史中心。
            </p>
          </div>
        </div>
      </section>

      {/* 底部信任与安全保障 */}
      <section className="p-6 bg-gray-50 rounded-2xl border border-gray-200/80 flex flex-col sm:flex-row items-center justify-between gap-4">
        <div className="flex items-center gap-3">
          <ShieldCheck className="w-6 h-6 text-emerald-600 shrink-0" />
          <div>
            <div className="text-xs sm:text-sm font-bold text-gray-900">隐私与数据安全保障</div>
            <p className="text-xs text-gray-500">
              数据仅用于当前任务商业提炼，个人历史报告严格隔离，支持随时查看与一键归档删除。
            </p>
          </div>
        </div>
        <button
          onClick={onStartDiscovery}
          className="btn btn-primary text-xs py-2 px-5 shrink-0"
        >
          开始体验
        </button>
      </section>
    </div>
  );
};
