import React, { createContext, useContext, useState, useEffect } from 'react';
import { UserProfileVO } from '../types';

interface AuthContextType {
  user: UserProfileVO | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoginModalOpen: boolean;
  openLoginModal: () => void;
  closeLoginModal: () => void;
  setAuthSession: (token: string, user: UserProfileVO) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<UserProfileVO | null>(() => {
    const cached = localStorage.getItem('oppty_user');
    return cached ? JSON.parse(cached) : null;
  });
  const [token, setToken] = useState<string | null>(() => {
    return localStorage.getItem('oppty_token');
  });
  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false);

  useEffect(() => {
    // 若初次进入没有用户，可默认配置一个快速体验会话或保持为空
    if (!token) {
      // 保持未登录状态
    }
  }, [token]);

  const setAuthSession = (newToken: string, newUser: UserProfileVO) => {
    localStorage.setItem('oppty_token', newToken);
    localStorage.setItem('oppty_user', JSON.stringify(newUser));
    setToken(newToken);
    setUser(newUser);
    setIsLoginModalOpen(false);
  };

  const logout = () => {
    localStorage.removeItem('oppty_token');
    localStorage.removeItem('oppty_user');
    setToken(null);
    setUser(null);
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        isAuthenticated: !!token,
        isLoginModalOpen,
        openLoginModal: () => setIsLoginModalOpen(true),
        closeLoginModal: () => setIsLoginModalOpen(false),
        setAuthSession,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth 必须在 AuthProvider 内部使用');
  }
  return context;
};
