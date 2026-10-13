import React, { useState } from 'react';
import { Radar, History, User, LogOut, Menu, X, Sparkles, BookOpen, ExternalLink } from 'lucide-react';
import { useAuth } from '../../context/AuthContext';

interface NavbarProps {
  currentTab: 'workspace' | 'history';
  onTabChange: (tab: 'workspace' | 'history') => void;
}

export const Navbar: React.FC<NavbarProps> = ({ currentTab, onTabChange }) => {
  const { user, isAuthenticated, openLoginModal, logout } = useAuth();
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  return (
    <header className="sticky top-0 z-50 bg-white/80 backdrop-blur-md border-b border-gray-200">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        {/* Brand Logo & Name */}
        <div className="flex items-center gap-3 cursor-pointer" onClick={() => onTabChange('workspace')}>
          <div className="w-10 h-10 rounded-xl bg-gradient-to-tr from-blue-600 to-indigo-500 flex items-center justify-center text-white shadow-md shadow-blue-500/20">
            <Radar className="w-5 h-5 animate-pulse" />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <span className="font-extrabold text-lg text-gray-900 tracking-tight">商机发现助手</span>
              <span className="hidden md:inline-block text-xs font-semibold px-2 py-0.5 bg-blue-50 text-blue-700 rounded-full border border-blue-200">
                Oppty Discovery
              </span>
            </div>
            <p className="text-[11px] text-gray-500 font-medium hidden sm:block">
              痛点一键透视，商机即刻成型
            </p>
          </div>
        </div>

        {/* Desktop Navigation Tabs - 双核心业务Tab */}
        <nav className="hidden md:flex items-center gap-1 bg-gray-100 p-1 rounded-xl">
          <button
            onClick={() => onTabChange('workspace')}
            className={`flex items-center gap-2 px-4 py-1.5 rounded-lg text-sm font-semibold transition-all ${
              currentTab === 'workspace'
                ? 'bg-white text-blue-600 shadow-sm'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            <Sparkles className="w-4 h-4" />
            商机发现
          </button>
          <button
            onClick={() => onTabChange('history')}
            className={`flex items-center gap-2 px-4 py-1.5 rounded-lg text-sm font-semibold transition-all ${
              currentTab === 'history'
                ? 'bg-white text-blue-600 shadow-sm'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            <History className="w-4 h-4" />
            商机历史
          </button>
        </nav>

        {/* User Account & External About Link */}
        <div className="hidden md:flex items-center gap-4">
          {/* 产品介绍新标签页超链接 */}
          <a
            href="?view=about"
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center gap-1.5 text-xs font-semibold text-gray-600 hover:text-blue-600 px-3 py-1.5 rounded-lg hover:bg-gray-100 transition-colors"
            title="在新标签页中打开产品介绍"
          >
            <BookOpen className="w-3.5 h-3.5 text-blue-500" />
            <span>产品介绍</span>
            <ExternalLink className="w-3 h-3 text-gray-400" />
          </a>

          {isAuthenticated && user ? (
            <div className="flex items-center gap-3 bg-gray-50 border border-gray-200 px-3 py-1.5 rounded-full">
              <div className="w-7 h-7 rounded-full bg-blue-100 text-blue-700 font-bold flex items-center justify-center text-xs">
                {user.nickname ? user.nickname.slice(0, 1).toUpperCase() : 'U'}
              </div>
              <span className="text-xs font-medium text-gray-700 max-w-[120px] truncate">
                {user.email}
              </span>
              <button
                onClick={logout}
                title="退出登录"
                className="text-gray-400 hover:text-red-500 transition-colors p-1"
              >
                <LogOut className="w-3.5 h-3.5" />
              </button>
            </div>
          ) : (
            <button onClick={openLoginModal} className="btn btn-primary text-xs py-2 px-4">
              <User className="w-3.5 h-3.5" />
              登录 / 注册
            </button>
          )}
        </div>

        {/* Mobile Actions & Hamburger Button */}
        <div className="flex md:hidden items-center gap-2">
          {!isAuthenticated && (
            <button
              onClick={openLoginModal}
              className="btn btn-primary text-xs py-1.5 px-2.5"
            >
              <User className="w-3.5 h-3.5" />
              <span>登录</span>
            </button>
          )}
          {isAuthenticated ? (
            <div className="w-8 h-8 rounded-full bg-blue-100 text-blue-700 font-bold flex items-center justify-center text-xs">
              {user?.nickname ? user.nickname.slice(0, 1).toUpperCase() : 'U'}
            </div>
          ) : null}
          <button
            onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            className="p-2 rounded-lg text-gray-600 hover:bg-gray-100"
          >
            {mobileMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
          </button>
        </div>
      </div>

      {/* Mobile Drawer */}
      {mobileMenuOpen && (
        <div className="md:hidden border-t border-gray-200 bg-white px-4 pt-3 pb-5 space-y-3">
          <button
            onClick={() => {
              onTabChange('workspace');
              setMobileMenuOpen(false);
            }}
            className={`w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-semibold ${
              currentTab === 'workspace' ? 'bg-blue-50 text-blue-600' : 'text-gray-700'
            }`}
          >
            <Sparkles className="w-4 h-4" />
            商机发现
          </button>
          <button
            onClick={() => {
              onTabChange('history');
              setMobileMenuOpen(false);
            }}
            className={`w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-semibold ${
              currentTab === 'history' ? 'bg-blue-50 text-blue-600' : 'text-gray-700'
            }`}
          >
            <History className="w-4 h-4" />
            商机历史
          </button>
          <a
            href="?view=about"
            target="_blank"
            rel="noopener noreferrer"
            onClick={() => setMobileMenuOpen(false)}
            className="w-full flex items-center justify-between px-3 py-2 rounded-lg text-sm font-semibold text-gray-700 hover:bg-gray-50"
          >
            <span className="flex items-center gap-3">
              <BookOpen className="w-4 h-4 text-blue-600" />
              产品介绍
            </span>
            <ExternalLink className="w-3.5 h-3.5 text-gray-400" />
          </a>
          <div className="pt-2 border-t border-gray-100">
            {isAuthenticated && user ? (
              <div className="flex items-center justify-between py-2">
                <span className="text-xs text-gray-600">{user.email}</span>
                <button
                  onClick={() => {
                    logout();
                    setMobileMenuOpen(false);
                  }}
                  className="flex items-center gap-1 text-xs text-red-600 font-medium"
                >
                  <LogOut className="w-3.5 h-3.5" />
                  退出
                </button>
              </div>
            ) : (
              <button
                onClick={() => {
                  openLoginModal();
                  setMobileMenuOpen(false);
                }}
                className="w-full btn btn-primary text-xs"
              >
                <User className="w-3.5 h-3.5" />
                登录 / 注册
              </button>
            )}
          </div>
        </div>
      )}
    </header>
  );
};
