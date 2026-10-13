import React, { useState, useEffect } from 'react';
import { X, Mail, KeyRound, Loader2, Sparkles, CheckCircle2 } from 'lucide-react';
import { useAuth } from '../../context/AuthContext';
import { authApi } from '../../api/client';

export const LoginModal: React.FC = () => {
  const { isLoginModalOpen, closeLoginModal, setAuthSession } = useAuth();
  const [email, setEmail] = useState('');
  const [code, setCode] = useState('');
  const [cooldown, setCooldown] = useState(0);
  const [isSendingCode, setIsSendingCode] = useState(false);
  const [isLoggingIn, setIsLoggingIn] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');
  const [successNotice, setSuccessNotice] = useState('');

  useEffect(() => {
    let timer: number;
    if (cooldown > 0) {
      timer = setInterval(() => setCooldown((prev) => prev - 1), 1000);
    }
    return () => clearInterval(timer);
  }, [cooldown]);

  if (!isLoginModalOpen) return null;

  const handleSendCode = async () => {
    if (!email || !email.includes('@')) {
      setErrorMessage('请输入合法的邮箱地址');
      return;
    }
    setErrorMessage('');
    setIsSendingCode(true);
    try {
      const resp = await authApi.sendCode(email);
      setCooldown(resp.cooldownSeconds || 60);
      setSuccessNotice('验证码已发送至您的邮箱，支持测试码 123456');
    } catch (err: any) {
      setErrorMessage(err.message || '发送验证码失败');
    } finally {
      setIsSendingCode(false);
    }
  };

  const handleVerify = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email) {
      setErrorMessage('请输入邮箱');
      return;
    }
    if (!code || code.length !== 6) {
      setErrorMessage('请输入6位数字验证码');
      return;
    }
    setErrorMessage('');
    setIsLoggingIn(true);
    try {
      const resp = await authApi.verifyCode(email, code);
      setAuthSession(resp.token, resp.user);
      closeLoginModal();
    } catch (err: any) {
      setErrorMessage(err.message || '登录失败');
    } finally {
      setIsLoggingIn(false);
    }
  };

  return (
    <div className="modal-overlay">
      <div className="card-glass w-full max-w-md p-6 sm:p-8 relative shadow-2xl animate-in">
        {/* 关闭按钮 */}
        <button
          onClick={closeLoginModal}
          className="absolute top-5 right-5 text-gray-400 hover:text-gray-600 p-1 rounded-lg"
        >
          <X className="w-5 h-5" />
        </button>

        {/* 头部标题与简介 */}
        <div className="text-center mb-6">
          <div className="w-12 h-12 bg-blue-100 text-blue-600 rounded-2xl flex items-center justify-center mx-auto mb-3 shadow-inner">
            <Sparkles className="w-6 h-6" />
          </div>
          <h2 className="text-xl font-bold text-gray-900">登录即注册 · 无感开启</h2>
          <p className="text-xs text-gray-500 mt-1">
            首次登录将自动为您创建专属商机发现档案，无需复杂注册
          </p>
        </div>

        {/* 错误与成功通知 */}
        {errorMessage && (
          <div className="mb-4 p-3 bg-red-50 border border-red-200 text-red-700 text-xs rounded-xl flex items-center gap-2">
            <span>{errorMessage}</span>
          </div>
        )}
        {successNotice && (
          <div className="mb-4 p-3 bg-emerald-50 border border-emerald-200 text-emerald-700 text-xs rounded-xl flex items-center gap-2">
            <CheckCircle2 className="w-4 h-4 shrink-0" />
            <span>{successNotice}</span>
          </div>
        )}

        <form onSubmit={handleVerify} className="space-y-4">
          {/* 邮箱输入与获取验证码 */}
          <div>
            <label className="block text-xs font-semibold text-gray-700 mb-1.5">电子邮箱</label>
            <div className="relative flex items-center">
              <Mail className="w-4 h-4 text-gray-400 absolute left-3 pointer-events-none" />
              <input
                type="email"
                placeholder="请输入您的常用邮箱 (如 name@example.com)"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="input-field pl-9 pr-28 text-sm"
                required
              />
              <button
                type="button"
                onClick={handleSendCode}
                disabled={isSendingCode || cooldown > 0}
                className="absolute right-1.5 px-3 py-1.5 text-xs font-semibold text-blue-600 bg-blue-50 hover:bg-blue-100 rounded-lg disabled:opacity-50 transition"
              >
                {isSendingCode ? (
                  <Loader2 className="w-3.5 h-3.5 animate-spin" />
                ) : cooldown > 0 ? (
                  `${cooldown}s`
                ) : (
                  '获取验证码'
                )}
              </button>
            </div>
          </div>

          {/* 验证码输入 */}
          <div>
            <label className="block text-xs font-semibold text-gray-700 mb-1.5">6位验证码</label>
            <div className="relative flex items-center">
              <KeyRound className="w-4 h-4 text-gray-400 absolute left-3 pointer-events-none" />
              <input
                type="text"
                maxLength={6}
                placeholder="请输入6位验证码 (自测可用 123456)"
                value={code}
                onChange={(e) => setCode(e.target.value.trim())}
                className="input-field pl-9 tracking-widest font-mono text-sm"
                required
              />
            </div>
          </div>

          {/* 提交按钮 */}
          <button
            type="submit"
            disabled={isLoggingIn}
            className="w-full btn btn-primary py-2.5 mt-2 flex items-center justify-center gap-2"
          >
            {isLoggingIn && <Loader2 className="w-4 h-4 animate-spin" />}
            <span>验证并无感登录</span>
          </button>
        </form>

        <div className="mt-5 text-center">
          <p className="text-[11px] text-gray-400">
            登录即代表同意用户服务协议与隐私条款 · 数据均已安全加密存储
          </p>
        </div>
      </div>
    </div>
  );
};
