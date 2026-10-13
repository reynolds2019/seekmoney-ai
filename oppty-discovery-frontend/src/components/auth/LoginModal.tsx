import React, { useState, useEffect } from 'react';
import { X, Mail, KeyRound, Loader2, Radar, CheckCircle2 } from 'lucide-react';
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
      setSuccessNotice('验证码已发送，请前往您的邮箱查收并填入');
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
      <div className="card-glass w-full max-w-sm p-6 relative shadow-2xl animate-in">
        {/* 关闭按钮 */}
        <button
          onClick={closeLoginModal}
          className="absolute top-4 right-4 text-gray-400 hover:text-gray-600 p-1 rounded-lg"
        >
          <X className="w-5 h-5" />
        </button>

        {/* 头部标题与简介 */}
        <div className="text-center mb-5">
          <div className="w-10 h-10 rounded-xl bg-gradient-to-tr from-blue-600 to-indigo-500 flex items-center justify-center text-white mx-auto mb-2.5 shadow-md shadow-blue-500/20">
            <Radar className="w-5 h-5 animate-pulse" />
          </div>
          <h2 className="text-lg font-bold text-gray-900">登录 / 注册</h2>
          <p className="text-xs text-gray-500 mt-0.5">
            输入邮箱以开启商机挖掘
          </p>
        </div>

        {/* 错误与成功通知 */}
        {errorMessage && (
          <div className="mb-3.5 p-2.5 bg-red-50 border border-red-200 text-red-700 text-xs rounded-xl flex items-center gap-2">
            <span>{errorMessage}</span>
          </div>
        )}
        {successNotice && (
          <div className="mb-3.5 p-2.5 bg-emerald-50 border border-emerald-200 text-emerald-700 text-xs rounded-xl flex items-center gap-2">
            <CheckCircle2 className="w-4 h-4 shrink-0" />
            <span>{successNotice}</span>
          </div>
        )}

        <form onSubmit={handleVerify} className="space-y-3.5">
          {/* 邮箱输入与获取验证码 */}
          <div>
            <label className="block text-xs font-semibold text-gray-700 mb-1">电子邮箱</label>
            <div className="relative flex items-center">
              <Mail className="w-4 h-4 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2 pointer-events-none" />
              <input
                type="email"
                placeholder="请输入您的邮箱"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="input-field !pl-10 !pr-28 text-sm"
                style={{ paddingLeft: '2.5rem', paddingRight: '7rem' }}
                required
              />
              <button
                type="button"
                onClick={handleSendCode}
                disabled={isSendingCode || cooldown > 0}
                className="absolute right-1 px-2.5 py-1 text-xs font-semibold text-blue-600 bg-blue-50 hover:bg-blue-100 rounded-lg disabled:opacity-50 transition"
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
            <label className="block text-xs font-semibold text-gray-700 mb-1">6位验证码</label>
            <div className="relative flex items-center">
              <KeyRound className="w-4 h-4 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2 pointer-events-none" />
              <input
                type="text"
                maxLength={6}
                placeholder="请输入收到的6位验证码"
                value={code}
                onChange={(e) => setCode(e.target.value.trim())}
                className="input-field !pl-10 tracking-widest font-mono text-sm"
                style={{ paddingLeft: '2.5rem' }}
                required
              />
            </div>
            <p className="text-[11px] text-gray-400 mt-1">请先点击“获取验证码”，前往您的邮箱查收后填入</p>
          </div>

          {/* 提交按钮 */}
          <button
            type="submit"
            disabled={isLoggingIn}
            className="w-full btn btn-primary py-2.5 mt-2 flex items-center justify-center gap-2"
          >
            {isLoggingIn && <Loader2 className="w-4 h-4 animate-spin" />}
            <span>登录 / 注册</span>
          </button>
        </form>
      </div>
    </div>
  );
};
