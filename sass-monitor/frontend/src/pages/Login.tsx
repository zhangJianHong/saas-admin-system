import React, { useState } from 'react';
import { Form, Input, Button, Card, Typography, message } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { AuthService } from '../services/auth';
import ThemeManager from '../utils/theme';

const { Title, Text } = Typography;

interface LoginForm {
  username: string;
  password: string;
}

const Login: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const onFinish = async (values: LoginForm) => {
    setLoading(true);
    try {
      const response = await AuthService.login(values);

      // 保存认证信息
      localStorage.setItem('access_token', response.token);
      localStorage.setItem('refresh_token', response.refresh_token);
      localStorage.setItem('user_info', JSON.stringify(response.user));

      message.success('登录成功！');

      // 跳转到首页
      navigate('/');
    } catch (error: any) {
      console.error('Login error:', error);
      message.error(error.response?.data?.error || '登录失败，请检查用户名和密码');
    } finally {
      setLoading(false);
    }
  };

  const themeClassName = ThemeManager.getThemeClassName();

  return (
    <div className={`login-page ${themeClassName}`}>
      <div className="login-container">
        <Card className="login-card" bordered={false}>
          <div className="login-header">
            <Title level={2} className="login-title">
              SaaS 后台监控系统
            </Title>
            <Text type="secondary" className="login-subtitle">
              请登录您的管理员账号
            </Text>
          </div>

          <Form
            name="login"
            className="login-form"
            initialValues={{ remember: true }}
            onFinish={onFinish}
            size="large"
          >
            <Form.Item
              name="username"
              rules={[{ required: true, message: '请输入用户名！' }]}
            >
              <Input
                prefix={<UserOutlined />}
                placeholder="用户名"
                defaultValue="admin"
              />
            </Form.Item>

            <Form.Item
              name="password"
              rules={[{ required: true, message: '请输入密码！' }]}
            >
              <Input.Password
                prefix={<LockOutlined />}
                placeholder="密码"
                defaultValue="admin123"
              />
            </Form.Item>

            <Form.Item>
              <Button
                type="primary"
                htmlType="submit"
                className="login-button"
                loading={loading}
                block
              >
                登录
              </Button>
            </Form.Item>
          </Form>

          <div className="login-footer">
            <Text type="secondary">
              默认账号：admin / admin123
            </Text>
          </div>
        </Card>
      </div>

      <style>{`
        .login-page {
          min-height: 100vh;
          display: flex;
          align-items: center;
          justify-content: center;
          background: linear-gradient(135deg, var(--theme-bg-primary) 0%, var(--theme-bg-secondary) 100%);
          padding: 20px;
        }

        .login-container {
          width: 100%;
          max-width: 400px;
        }

        .login-card {
          box-shadow: 0 8px 32px var(--theme-card-shadow);
          border-radius: var(--theme-border-radius);
          backdrop-filter: blur(10px);
          background: var(--theme-bg-primary);
          border: 1px solid var(--theme-border-color);
        }

        .login-header {
          text-align: center;
          margin-bottom: 32px;
        }

        .login-title {
          color: var(--theme-text-primary) !important;
          margin-bottom: 8px !important;
        }

        .login-subtitle {
          color: var(--theme-text-secondary);
          font-size: 14px;
        }

        .login-form {
          margin-top: 24px;
        }

        .login-button {
          height: 44px;
          border-radius: var(--theme-border-radius);
          font-size: 16px;
          font-weight: 500;
        }

        .login-footer {
          text-align: center;
          margin-top: 24px;
          padding-top: 16px;
          border-top: 1px solid var(--theme-border-color-split);
        }

        /* Ant Design 样式覆盖 */
        .login-page .ant-input-affix-wrapper {
          background: var(--theme-bg-secondary) !important;
          border-color: var(--theme-border-color) !important;
          border-radius: var(--theme-border-radius) !important;
        }

        .login-page .ant-input-affix-wrapper:hover,
        .login-page .ant-input-affix-wrapper-focused {
          border-color: var(--theme-primary) !important;
          box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.1) !important;
        }

        .login-page .ant-input {
          background: transparent !important;
          color: var(--theme-text-primary) !important;
        }

        .login-page .ant-input::placeholder {
          color: var(--theme-text-disabled) !important;
        }

        .login-page .ant-btn-primary {
          background: var(--theme-primary) !important;
          border-color: var(--theme-primary) !important;
        }

        .login-page .ant-btn-primary:hover {
          background: var(--theme-primary-light) !important;
          border-color: var(--theme-primary-light) !important;
        }
      `}</style>
    </div>
  );
};

export default Login;