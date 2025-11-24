import React, { useState, useEffect } from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import {
  Typography,
  Card,
  Row,
  Col,
  Form,
  Input,
  Select,
  Switch,
  Button,
  Space,
  message,
  Divider,
} from 'antd';
import {
  DatabaseOutlined,
  SettingOutlined,
  EyeOutlined,
  ReloadOutlined,
  SaveOutlined,
  UserOutlined,
  LockOutlined,
} from '@ant-design/icons';
import ThemeToggle from '../components/ThemeToggle';
import MonitoringService from '../services/monitoring';
import { AuthService } from '../services/auth';

const { Title, Text } = Typography;
const { Option } = Select;

const Settings: React.FC = () => {
  const [databaseForm] = Form.useForm();
  const [systemForm] = Form.useForm();
  const [passwordForm] = Form.useForm();
  const [loading, setLoading] = useState(false);

  // 初始化表单数据
  useEffect(() => {
    // TODO: 从API获取设置数据
    databaseForm.setFieldsValue({
      postgresql_host: 'localhost',
      postgresql_port: 5432,
      postgresql_database: 'sass_monitor',
      postgresql_max_connections: 100,

      clickhouse_host: 'localhost',
      clickhouse_port: 9000,
      clickhouse_database: 'traces',

      redis_host: 'localhost',
      redis_port: 6379,
      redis_database: 0,
      redis_max_memory: '1GB',
    });

    systemForm.setFieldsValue({
      collect_interval: 5,
      retention_days: 30,
      alert_enabled: true,
      cpu_threshold: 80,
      memory_threshold: 85,
      disk_threshold: 90,
      connection_threshold: 100,
    });
  }, [databaseForm, systemForm]);

  // 保存数据库配置
  const handleDatabaseSave = async (values: any) => {
    try {
      setLoading(true);
      // 将数据库配置转换为后端期望的格式
      const databaseConfigs = {
        // PostgreSQL配置
        'postgresql_host': values.postgresql_host,
        'postgresql_port': values.postgresql_port?.toString(),
        'postgresql_database': values.postgresql_database,
        'postgresql_max_connections': values.postgresql_max_connections?.toString(),

        // ClickHouse配置
        'clickhouse_host': values.clickhouse_host,
        'clickhouse_port': values.clickhouse_port?.toString(),
        'clickhouse_database': values.clickhouse_database,

        // Redis配置
        'redis_host': values.redis_host,
        'redis_port': values.redis_port?.toString(),
        'redis_database': values.redis_database?.toString(),
        'redis_max_memory': values.redis_max_memory,
      };

      await MonitoringService.updateSystemConfigs(databaseConfigs);
      message.success('数据库配置保存成功');
    } catch (error: any) {
      console.error('Database config save error:', error);
      message.error(error.response?.data?.message || '保存失败');
    } finally {
      setLoading(false);
    }
  };

  // 保存系统配置
  const handleSystemSave = async (values: any) => {
    try {
      setLoading(true);
      // 将系统配置转换为后端期望的格式
      const systemConfigs = {
        'collect_interval': values.collect_interval?.toString(),
        'retention_days': values.retention_days?.toString(),
        'alert_enabled': values.alert_enabled ? 'true' : 'false',
        'cpu_threshold': values.cpu_threshold?.toString(),
        'memory_threshold': values.memory_threshold?.toString(),
        'disk_threshold': values.disk_threshold?.toString(),
        'connection_threshold': values.connection_threshold?.toString(),
      };

      await MonitoringService.updateSystemConfigs(systemConfigs);
      message.success('系统配置保存成功');
    } catch (error: any) {
      console.error('System config save error:', error);
      message.error(error.response?.data?.message || '保存失败');
    } finally {
      setLoading(false);
    }
  };

  // 修改密码
  const handlePasswordChange = async (values: any) => {
    try {
      setLoading(true);

      await AuthService.changePassword(values.old_password, values.new_password);

      message.success('密码修改成功,请重新登录');
      passwordForm.resetFields();

      // 3秒后跳转到登录页
      setTimeout(() => {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('user_info');
        window.location.href = '/login';
      }, 3000);
    } catch (error: any) {
      console.error('Password change error:', error);
      message.error(error.response?.data?.error || '密码修改失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ padding: 24 }}>
      {/* 页面标题 */}
      <div style={{ marginBottom: 24 }}>
        <Title level={2} style={{ margin: 0 }}>系统设置</Title>
        <Text type="secondary">配置数据库连接和系统参数</Text>
      </div>

      {/* 导航卡片 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={12}>
          <Card
            hoverable
            bodyStyle={{ padding: 24 }}
            onClick={() => window.location.href = '/settings/database'}
          >
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
              <div>
                <div style={{ display: 'flex', alignItems: 'center', marginBottom: 8 }}>
                  <DatabaseOutlined style={{ fontSize: 24, color: '#1890ff', marginRight: 12 }} />
                  <Title level={4} style={{ margin: 0 }}>数据库配置</Title>
                </div>
                <Text type="secondary">PostgreSQL、ClickHouse、Redis 连接设置</Text>
              </div>
              <EyeOutlined style={{ fontSize: 16, color: '#1890ff' }} />
            </div>
          </Card>
        </Col>
        <Col xs={24} sm={12}>
          <Card
            hoverable
            bodyStyle={{ padding: 24 }}
            onClick={() => window.location.href = '/settings/theme'}
          >
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
              <div>
                <div style={{ display: 'flex', alignItems: 'center', marginBottom: 8 }}>
                  <SettingOutlined style={{ fontSize: 24, color: '#52c41a', marginRight: 12 }} />
                  <Title level={4} style={{ margin: 0 }}>主题设置</Title>
                </div>
                <Text type="secondary">界面主题和外观设置</Text>
              </div>
              <EyeOutlined style={{ fontSize: 16, color: '#52c41a' }} />
            </div>
          </Card>
        </Col>
        <Col xs={24} sm={12}>
          <Card
            hoverable
            bodyStyle={{ padding: 24 }}
            onClick={() => window.location.href = '/settings/account'}
          >
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
              <div>
                <div style={{ display: 'flex', alignItems: 'center', marginBottom: 8 }}>
                  <UserOutlined style={{ fontSize: 24, color: '#fa8c16', marginRight: 12 }} />
                  <Title level={4} style={{ margin: 0 }}>账户设置</Title>
                </div>
                <Text type="secondary">修改密码和个人信息</Text>
              </div>
              <EyeOutlined style={{ fontSize: 16, color: '#fa8c16' }} />
            </div>
          </Card>
        </Col>
      </Row>

      {/* 子路由 */}
      <Routes>
        <Route path="/" element={<Navigate to="/settings/database" replace />} />
        <Route
          path="/database"
          element={
            <Card title="数据库配置">
              <Form
                form={databaseForm}
                layout="vertical"
                onFinish={handleDatabaseSave}
              >
                <Title level={4}>PostgreSQL 配置</Title>
                <Row gutter={16}>
                  <Col span={12}>
                    <Form.Item
                      label="主机地址"
                      name="postgresql_host"
                      rules={[{ required: true, message: '请输入主机地址' }]}
                    >
                      <Input />
                    </Form.Item>
                  </Col>
                  <Col span={12}>
                    <Form.Item
                      label="端口"
                      name="postgresql_port"
                      rules={[{ required: true, message: '请输入端口' }]}
                    >
                      <Input type="number" />
                    </Form.Item>
                  </Col>
                </Row>
                <Row gutter={16}>
                  <Col span={12}>
                    <Form.Item
                      label="数据库名"
                      name="postgresql_database"
                      rules={[{ required: true, message: '请输入数据库名' }]}
                    >
                      <Input />
                    </Form.Item>
                  </Col>
                  <Col span={12}>
                    <Form.Item
                      label="最大连接数"
                      name="postgresql_max_connections"
                      rules={[{ required: true, message: '请输入最大连接数' }]}
                    >
                      <Input type="number" />
                    </Form.Item>
                  </Col>
                </Row>

                <Divider />

                <Title level={4}>ClickHouse 配置</Title>
                <Row gutter={16}>
                  <Col span={12}>
                    <Form.Item
                      label="主机地址"
                      name="clickhouse_host"
                      rules={[{ required: true, message: '请输入主机地址' }]}
                    >
                      <Input />
                    </Form.Item>
                  </Col>
                  <Col span={12}>
                    <Form.Item
                      label="端口"
                      name="clickhouse_port"
                      rules={[{ required: true, message: '请输入端口' }]}
                    >
                      <Input type="number" />
                    </Form.Item>
                  </Col>
                </Row>
                <Row gutter={16}>
                  <Col span={12}>
                    <Form.Item
                      label="数据库名"
                      name="clickhouse_database"
                      rules={[{ required: true, message: '请输入数据库名' }]}
                    >
                      <Input />
                    </Form.Item>
                  </Col>
                </Row>

                <Divider />

                <Title level={4}>Redis 配置</Title>
                <Row gutter={16}>
                  <Col span={12}>
                    <Form.Item
                      label="主机地址"
                      name="redis_host"
                      rules={[{ required: true, message: '请输入主机地址' }]}
                    >
                      <Input />
                    </Form.Item>
                  </Col>
                  <Col span={12}>
                    <Form.Item
                      label="端口"
                      name="redis_port"
                      rules={[{ required: true, message: '请输入端口' }]}
                    >
                      <Input type="number" />
                    </Form.Item>
                  </Col>
                </Row>
                <Row gutter={16}>
                  <Col span={12}>
                    <Form.Item
                      label="数据库"
                      name="redis_database"
                      rules={[{ required: true, message: '请输入数据库' }]}
                    >
                      <Input type="number" />
                    </Form.Item>
                  </Col>
                  <Col span={12}>
                    <Form.Item
                      label="最大内存"
                      name="redis_max_memory"
                      rules={[{ required: true, message: '请输入最大内存' }]}
                    >
                      <Select>
                        <Option value="256MB">256MB</Option>
                        <Option value="512MB">512MB</Option>
                        <Option value="1GB">1GB</Option>
                        <Option value="2GB">2GB</Option>
                        <Option value="4GB">4GB</Option>
                      </Select>
                    </Form.Item>
                  </Col>
                </Row>

                <Form.Item style={{ marginTop: 32 }}>
                  <Space>
                    <Button type="primary" htmlType="submit" icon={<SaveOutlined />} loading={loading}>
                      保存配置
                    </Button>
                    <Button icon={<ReloadOutlined />} onClick={() => databaseForm.resetFields()}>
                      重置
                    </Button>
                  </Space>
                </Form.Item>
              </Form>
            </Card>
          }
        />
        <Route
          path="/theme"
          element={
            <Card title="主题设置">
              <div style={{ marginBottom: 32 }}>
                <Title level={4}>界面主题</Title>
                <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
                  <Text>主题切换：</Text>
                  <ThemeToggle />
                </div>
              </div>

              <div>
                <Title level={4}>系统参数配置</Title>
                <Form
                  form={systemForm}
                  layout="vertical"
                  onFinish={handleSystemSave}
                >
                  <Row gutter={16}>
                    <Col span={12}>
                      <Form.Item
                        label="数据采集间隔（分钟）"
                        name="collect_interval"
                        rules={[{ required: true, message: '请输入采集间隔' }]}
                      >
                        <Input type="number" />
                      </Form.Item>
                    </Col>
                    <Col span={12}>
                      <Form.Item
                        label="数据保留天数"
                        name="retention_days"
                        rules={[{ required: true, message: '请输入保留天数' }]}
                      >
                        <Input type="number" />
                      </Form.Item>
                    </Col>
                  </Row>

                  <Row gutter={16}>
                    <Col span={12}>
                      <Form.Item
                        label="CPU告警阈值（%）"
                        name="cpu_threshold"
                        rules={[{ required: true, message: '请输入CPU阈值' }]}
                      >
                        <Input type="number" />
                      </Form.Item>
                    </Col>
                    <Col span={12}>
                      <Form.Item
                        label="内存告警阈值（%）"
                        name="memory_threshold"
                        rules={[{ required: true, message: '请输入内存阈值' }]}
                      >
                        <Input type="number" />
                      </Form.Item>
                    </Col>
                  </Row>

                  <Row gutter={16}>
                    <Col span={12}>
                      <Form.Item
                        label="磁盘告警阈值（%）"
                        name="disk_threshold"
                        rules={[{ required: true, message: '请输入磁盘阈值' }]}
                      >
                        <Input type="number" />
                      </Form.Item>
                    </Col>
                    <Col span={12}>
                      <Form.Item
                        label="连接数告警阈值"
                        name="connection_threshold"
                        rules={[{ required: true, message: '请输入连接数阈值' }]}
                      >
                        <Input type="number" />
                      </Form.Item>
                    </Col>
                  </Row>

                  <Form.Item>
                    <Form.Item name="alert_enabled" valuePropName="checked" noStyle>
                      <Switch checkedChildren="启用" unCheckedChildren="禁用" />
                    </Form.Item>
                    <Text style={{ marginLeft: 8 }}>启用告警功能</Text>
                  </Form.Item>

                  <Form.Item style={{ marginTop: 32 }}>
                    <Space>
                      <Button type="primary" htmlType="submit" icon={<SaveOutlined />} loading={loading}>
                        保存配置
                      </Button>
                      <Button icon={<ReloadOutlined />} onClick={() => systemForm.resetFields()}>
                        重置
                      </Button>
                    </Space>
                  </Form.Item>
                </Form>
              </div>
            </Card>
          }
        />
        <Route
          path="/account"
          element={
            <Card title={<><UserOutlined /> 账户设置</>}>
              <Row gutter={24}>
                <Col span={12}>
                  <Title level={4}><LockOutlined /> 修改密码</Title>
                  <Text type="secondary" style={{ marginBottom: 16, display: 'block' }}>
                    为了您的账户安全,建议定期更换密码
                  </Text>
                  <Form
                    form={passwordForm}
                    layout="vertical"
                    onFinish={handlePasswordChange}
                    style={{ maxWidth: 400 }}
                  >
                    <Form.Item
                      label="当前密码"
                      name="old_password"
                      rules={[
                        { required: true, message: '请输入当前密码' },
                      ]}
                    >
                      <Input.Password
                        prefix={<LockOutlined />}
                        placeholder="请输入当前密码"
                      />
                    </Form.Item>

                    <Form.Item
                      label="新密码"
                      name="new_password"
                      rules={[
                        { required: true, message: '请输入新密码' },
                        { min: 6, message: '密码长度至少为6位' },
                      ]}
                    >
                      <Input.Password
                        prefix={<LockOutlined />}
                        placeholder="请输入新密码(至少6位)"
                      />
                    </Form.Item>

                    <Form.Item
                      label="确认新密码"
                      name="confirm_password"
                      dependencies={['new_password']}
                      rules={[
                        { required: true, message: '请确认新密码' },
                        ({ getFieldValue }) => ({
                          validator(_, value) {
                            if (!value || getFieldValue('new_password') === value) {
                              return Promise.resolve();
                            }
                            return Promise.reject(new Error('两次输入的密码不一致'));
                          },
                        }),
                      ]}
                    >
                      <Input.Password
                        prefix={<LockOutlined />}
                        placeholder="请再次输入新密码"
                      />
                    </Form.Item>

                    <Form.Item>
                      <Space>
                        <Button
                          type="primary"
                          htmlType="submit"
                          icon={<SaveOutlined />}
                          loading={loading}
                        >
                          修改密码
                        </Button>
                        <Button
                          icon={<ReloadOutlined />}
                          onClick={() => passwordForm.resetFields()}
                        >
                          重置
                        </Button>
                      </Space>
                    </Form.Item>
                  </Form>
                </Col>
              </Row>
            </Card>
          }
        />
      </Routes>
    </div>
  );
};

export default Settings;