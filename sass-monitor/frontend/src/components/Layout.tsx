import React, { useState } from 'react';
import {
  Layout as AntLayout,
  Menu,
  Button,
  Typography,
  Space,
  Avatar,
  Dropdown,
  theme,
} from 'antd';
import { useNavigate, useLocation } from 'react-router-dom';
import {
  DashboardOutlined,
  DatabaseOutlined,
  TeamOutlined,
  DollarOutlined,
  SettingOutlined,
  BellOutlined,
  LogoutOutlined,
  UserOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
} from '@ant-design/icons';
import { useTheme } from '../hooks/useTheme';
import ThemeToggle from './ThemeToggle';

const { Header, Sider, Content } = AntLayout;
const { Text } = Typography;

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const { isDark } = useTheme();
  const [collapsed, setCollapsed] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const {
    token: { colorBgContainer },
  } = theme.useToken();

  const menuItems = [
    {
      key: '/',
      icon: <DashboardOutlined />,
      label: '仪表板',
    },
    {
      key: '/organizations',
      icon: <TeamOutlined />,
      label: '组织管理',
    },
    {
      key: '/users',
      icon: <UserOutlined />,
      label: '用户管理',
    },
    {
      key: '/subscription-plans',
      icon: <DollarOutlined />,
      label: '订阅计划',
    },
    {
      key: '/monitoring',
      icon: <DatabaseOutlined />,
      label: '监控数据',
      children: [
        {
          key: '/monitoring/metrics',
          label: '实时指标',
        },
        {
          key: '/monitoring/alerts',
          label: '告警管理',
        },
      ],
    },
    {
      key: '/settings',
      icon: <SettingOutlined />,
      label: '系统设置',
      children: [
        {
          key: '/settings/database',
          label: '数据库配置',
        },
        {
          key: '/settings/theme',
          label: '主题设置',
        },
      ],
    },
  ];

  const userMenuItems = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人信息',
    },
    {
      key: 'settings',
      icon: <SettingOutlined />,
      label: '用户设置',
    },
    {
      type: 'divider' as const,
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
    },
  ];

  // 获取当前选中的菜单项
  const getSelectedKeys = (pathname: string): string[] => {
    if (pathname === '/' || pathname === '/dashboard') {
      return ['/'];
    }

    // 监控模块的路径处理
    if (pathname.startsWith('/monitoring')) {
      const path = pathname.split('/')[2];
      return path ? [`/monitoring/${path}`] : ['/monitoring'];
    }

    // 设置模块的路径处理
    if (pathname.startsWith('/settings')) {
      const path = pathname.split('/')[2];
      return path ? [`/settings/${path}`] : ['/settings'];
    }

    return [pathname];
  };

  const handleMenuClick = ({ key }: { key: string }) => {
    // 路由处理
    if (key === 'logout') {
      // 处理退出登录
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      localStorage.removeItem('user_info');
      window.location.href = '/login';
    } else if (key === 'profile') {
      // 个人信息页面
      navigate('/profile');
    } else {
      // 处理菜单导航
      navigate(key);
    }
  };

  return (
    <AntLayout style={{ minHeight: '100vh' }}>
      {/* 侧边栏 */}
      <Sider
        trigger={null}
        collapsible
        collapsed={collapsed}
        theme={isDark ? 'dark' : 'light'}
        style={{
          background: isDark ? '#001529' : '#ffffff',
          borderRight: `1px solid ${isDark ? '#303030' : '#f0f0f0'}`,
        }}
      >
        {/* Logo区域 */}
        <div
          style={{
            height: 32,
            margin: 16,
            display: 'flex',
            alignItems: 'center',
            justifyContent: collapsed ? 'center' : 'flex-start',
          }}
        >
          <DatabaseOutlined style={{ fontSize: 24, color: '#1890ff' }} />
          {!collapsed && (
            <Text strong style={{ marginLeft: 8, color: isDark ? '#ffffff' : '#262626' }}>
              SaaS Monitor
            </Text>
          )}
        </div>

        {/* 菜单 */}
        <Menu
          theme={isDark ? 'dark' : 'light'}
          mode="inline"
          selectedKeys={getSelectedKeys(location.pathname)}
          items={menuItems}
          onClick={handleMenuClick}
        />
      </Sider>

      {/* 主内容区 */}
      <AntLayout>
        {/* 顶部导航栏 */}
        <Header
          style={{
            padding: '0 16px',
            background: colorBgContainer,
            borderBottom: `1px solid ${isDark ? '#303030' : '#f0f0f0'}`,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
          }}
        >
          {/* 左侧：折叠按钮 */}
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
            style={{
              fontSize: '16px',
              width: 64,
              height: 64,
            }}
          />

          {/* 右侧：工具栏 */}
          <Space size="middle">
            {/* 主题切换 */}
            <ThemeToggle />

            {/* 通知 */}
            <Button
              type="text"
              icon={<BellOutlined />}
              style={{
                fontSize: '16px',
                color: isDark ? '#ffffff' : '#262626',
              }}
            />

            {/* 用户头像和菜单 */}
            <Dropdown
              menu={{
                items: userMenuItems,
                onClick: handleMenuClick,
              }}
              placement="bottomRight"
            >
              <Space style={{ cursor: 'pointer' }}>
                <Avatar icon={<UserOutlined />} />
                <Text>管理员</Text>
              </Space>
            </Dropdown>
          </Space>
        </Header>

        {/* 内容区域 */}
        <Content
          style={{
            margin: 0,
            padding: 0,
            background: colorBgContainer,
            overflow: 'auto',
          }}
        >
          {children}
        </Content>
      </AntLayout>
    </AntLayout>
  );
};

export default Layout;