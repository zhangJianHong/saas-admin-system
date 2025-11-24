import React, { useEffect, useState } from 'react';
import { ConfigProvider, App as AntdApp } from 'antd';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import ThemeManager from './utils/theme';
import Layout from './components/Layout';
import Dashboard from './pages/Dashboard';
import Login from './pages/Login';
import ProtectedRoute from './components/ProtectedRoute';
import Organizations from './pages/Organizations';
import OrganizationDetail from './pages/OrganizationDetail';
import SubscriptionPlansPage from './pages/SubscriptionPlansPage';
import Users from './pages/Users';
import UserDetail from './pages/UserDetail';
import Monitoring from './pages/Monitoring';
import Settings from './pages/Settings';
import './App.css';

const App: React.FC = () => {
  const [loading, setLoading] = useState(true);

  // 初始化应用
  useEffect(() => {
    const initializeApp = async () => {
      try {
        // 初始化主题
        ThemeManager.initialize();
      } catch (error) {
        console.error('App initialization error:', error);
      } finally {
        setLoading(false);
      }
    };

    initializeApp();
  }, []);

  // 获取Ant Design主题配置
  const antdTheme = ThemeManager.getAntdTheme();

  // 获取当前主题类名
  const themeClassName = ThemeManager.getThemeClassName();

  if (loading) {
    return (
      <div className={`app-loading ${themeClassName}`}>
        <div className="loading-spinner" />
        <p>正在初始化应用...</p>
      </div>
    );
  }

  return (
    <div className={`app ${themeClassName}`}>
      <ConfigProvider theme={antdTheme}>
        <AntdApp>
          <Router>
            <Routes>
              {/* 登录页面 */}
              <Route path="/login" element={<Login />} />

              {/* 受保护的路由 */}
              <Route
                path="/*"
                element={
                  <ProtectedRoute>
                    <Layout>
                      <Routes>
                        <Route path="/dashboard" element={<Dashboard />} />
                        <Route path="/organizations" element={<Organizations />} />
                        <Route path="/organizations/:id" element={<OrganizationDetail />} />
                        <Route path="/subscription-plans" element={<SubscriptionPlansPage />} />
                        <Route path="/users" element={<Users />} />
                        <Route path="/users/:id" element={<UserDetail />} />
                        <Route path="/monitoring/*" element={<Monitoring />} />
                        <Route path="/settings/*" element={<Settings />} />
                        <Route path="/" element={<Navigate to="/dashboard" replace />} />
                      </Routes>
                    </Layout>
                  </ProtectedRoute>
                }
              />
            </Routes>
          </Router>
        </AntdApp>
      </ConfigProvider>
    </div>
  );
};

export default App;