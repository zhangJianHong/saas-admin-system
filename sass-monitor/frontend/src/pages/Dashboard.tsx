import React, { useState, useEffect } from 'react';
import {
  Row,
  Col,
  Card,
  Statistic,
  Table,
  Progress,
  Tag,
  Space,
  Button,
  Typography,
  Alert,
  Spin,
  Divider,
  Badge,
  Tooltip,
} from 'antd';
import {
  DatabaseOutlined,
  UserOutlined,
  TeamOutlined,
  CloudServerOutlined,
  ReloadOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
  WarningOutlined,
  CloseCircleOutlined,
  DollarOutlined,
  RiseOutlined,
  FallOutlined,
  ClockCircleOutlined,
  MailOutlined,
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { dashboardApi } from '../services/dashboard';
import { Overview, Organization, DatabaseStatus } from '../types';

const { Title, Text } = Typography;

interface DashboardData {
  overview: Overview;
  databaseStatus: DatabaseStatus;
  organizations: Organization[];
}

const Dashboard: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [data, setData] = useState<DashboardData | null>(null);
  const [error, setError] = useState<string | null>(null);

  const fetchData = async (showLoading = false) => {
    try {
      if (showLoading) {
        setRefreshing(true);
      } else {
        setLoading(true);
      }
      setError(null);

      const [overview, databaseStatus, organizations] = await Promise.all([
        dashboardApi.getOverview(),
        dashboardApi.getDatabaseStatus(),
        dashboardApi.getOrganizations(),
      ]);

      setData({ overview, databaseStatus, organizations });
    } catch (err: any) {
      setError(err.message || '获取数据失败');
      console.error('Dashboard fetch error:', err);
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  };

  useEffect(() => {
    fetchData();
    // 每30秒自动刷新
    const interval = setInterval(() => fetchData(true), 30000);
    return () => clearInterval(interval);
  }, []);

  const handleRefresh = () => {
    fetchData(true);
  };

  const getDatabaseStatusIcon = (status: string) => {
    return status === 'healthy' ? (
      <CheckCircleOutlined style={{ color: '#52c41a' }} />
    ) : (
      <ExclamationCircleOutlined style={{ color: '#ff4d4f' }} />
    );
  };

  const getDatabaseStatusColor = (status: string) => {
    return status === 'healthy' ? 'success' : 'error';
  };

  // 计算订阅统计
  const calculateSubscriptionStats = (orgs: Organization[]) => {
    const stats = {
      active: 0,
      expiring_soon: 0,
      expired: 0,
      none: 0,
      total_active_subscriptions: 0,
    };

    // 确保 orgs 是数组
    if (!Array.isArray(orgs)) {
      return stats;
    }

    orgs.forEach(org => {
      const status = org.subscription_status || 'none';
      stats[status as keyof typeof stats] = (stats[status as keyof typeof stats] || 0) + 1;
      stats.total_active_subscriptions += org.active_subscription_count || 0;
    });

    return stats;
  };

  const organizationColumns = [
    {
      title: '组织名称',
      dataIndex: 'name',
      key: 'name',
      width: 200,
      render: (name: string, record: Organization) => (
        <Space>
          <TeamOutlined />
          <div>
            <div style={{ fontWeight: 500 }}>{name}</div>
            <Text type="secondary" style={{ fontSize: 12 }}>
              {record.user_count} 用户 · {record.active_subscription_count || 0} 订阅
            </Text>
          </div>
        </Space>
      ),
    },
    {
      title: '订阅状态',
      key: 'subscription_status',
      width: 150,
      render: (record: Organization) => {
        const status = record.subscription_status || 'none';
        const days = record.days_until_expiration;

        const statusConfig: Record<string, { color: string; icon: React.ReactNode; text: string }> = {
          active: {
            color: 'success',
            icon: <CheckCircleOutlined />,
            text: days !== undefined ? `${days}天后到期` : '正常'
          },
          expiring_soon: {
            color: 'warning',
            icon: <WarningOutlined />,
            text: days !== undefined ? `${days}天后到期` : '即将到期'
          },
          expired: {
            color: 'error',
            icon: <CloseCircleOutlined />,
            text: '已到期'
          },
          none: {
            color: 'default',
            icon: <ClockCircleOutlined />,
            text: '无订阅'
          }
        };

        const config = statusConfig[status] || statusConfig.none;

        return <Tag color={config.color} icon={config.icon}>{config.text}</Tag>;
      },
    },
    {
      title: '存储使用',
      dataIndex: 'storage_usage',
      key: 'storage_usage',
      width: 120,
      render: (usage: number) => {
        const usageMB = usage / (1024 * 1024);
        return `${usageMB.toFixed(1)} MB`;
      },
    },
    {
      title: '操作',
      key: 'action',
      width: 80,
      render: (record: Organization) => (
        <Button type="link" size="small" onClick={() => navigate(`/organizations/${record.id}`)}>
          查看详情
        </Button>
      ),
    },
  ];

  if (loading) {
    return (
      <div style={{ padding: 24, textAlign: 'center' }}>
        <Spin size="large" />
        <div style={{ marginTop: 16 }}>加载仪表板数据...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div style={{ padding: 24 }}>
        <Alert
          message="加载失败"
          description={error}
          type="error"
          showIcon
          action={
            <Button size="small" type="primary" onClick={() => fetchData()}>
              重试
            </Button>
          }
        />
      </div>
    );
  }

  if (!data) {
    return null;
  }

  const { overview, databaseStatus, organizations } = data;

  // 确保 organizations 是数组
  const orgsArray = Array.isArray(organizations) ? organizations : [];
  const subscriptionStats = calculateSubscriptionStats(orgsArray);

  // 获取即将到期的组织
  const expiringOrgs = orgsArray
    .filter(org => org.subscription_status === 'expiring_soon')
    .sort((a, b) => (a.days_until_expiration || 999) - (b.days_until_expiration || 999))
    .slice(0, 5);

  // 获取用户最多的组织
  const topOrganizations = [...orgsArray]
    .sort((a, b) => b.user_count - a.user_count)
    .slice(0, 5);

  return (
    <div style={{ padding: 24, background: 'var(--theme-bg-primary)', minHeight: '100vh' }}>
      {/* 页面标题和刷新按钮 */}
      <div style={{ marginBottom: 24, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div>
          <Title level={2} style={{ margin: 0, color: 'var(--theme-text-primary)' }}>
            仪表盘
          </Title>
          <Text type="secondary">系统概览与监控</Text>
        </div>
        <Button
          icon={<ReloadOutlined spin={refreshing} />}
          onClick={handleRefresh}
          loading={refreshing}
        >
          刷新
        </Button>
      </div>

      {/* 第一行：核心统计卡片 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={12} lg={6}>
          <Card hoverable>
            <Statistic
              title="总组织数"
              value={overview.total_organizations}
              prefix={<TeamOutlined />}
              valueStyle={{ color: '#1890ff' }}
              suffix={
                <Tooltip title="查看所有组织">
                  <Button
                    type="link"
                    size="small"
                    onClick={() => navigate('/organizations')}
                    style={{ fontSize: 12 }}
                  >
                    详情
                  </Button>
                </Tooltip>
              }
            />
            <Divider style={{ margin: '12px 0' }} />
            <Space split={<Divider type="vertical" />}>
              <Text type="secondary" style={{ fontSize: 12 }}>
                活跃: {orgsArray.filter(o => o.user_count > 0).length}
              </Text>
              <Text type="secondary" style={{ fontSize: 12 }}>
                待激活: {orgsArray.filter(o => o.user_count === 0).length}
              </Text>
            </Space>
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card hoverable>
            <Statistic
              title="总用户数"
              value={overview.total_users}
              prefix={<UserOutlined />}
              valueStyle={{ color: '#52c41a' }}
              suffix={
                <Tooltip title="查看所有用户">
                  <Button
                    type="link"
                    size="small"
                    onClick={() => navigate('/users')}
                    style={{ fontSize: 12 }}
                  >
                    详情
                  </Button>
                </Tooltip>
              }
            />
            <Divider style={{ margin: '12px 0' }} />
            <Text type="secondary" style={{ fontSize: 12 }}>
              平均每组织: {overview.total_organizations > 0
                ? (overview.total_users / overview.total_organizations).toFixed(1)
                : 0} 用户
            </Text>
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card hoverable>
            <Statistic
              title="总订阅数"
              value={overview.total_subscriptions}
              prefix={<CloudServerOutlined />}
              valueStyle={{ color: '#fa8c16' }}
              suffix={
                <Tooltip title="查看订阅套餐">
                  <Button
                    type="link"
                    size="small"
                    onClick={() => navigate('/subscription-plans')}
                    style={{ fontSize: 12 }}
                  >
                    详情
                  </Button>
                </Tooltip>
              }
            />
            <Divider style={{ margin: '12px 0' }} />
            <Text type="secondary" style={{ fontSize: 12 }}>
              活跃订阅: {subscriptionStats.total_active_subscriptions}
            </Text>
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card hoverable>
            <Statistic
              title="系统健康度"
              value={Object.values(overview.system_health).filter(status => status === 'healthy').length}
              suffix={`/ ${Object.keys(overview.system_health).length}`}
              prefix={<DatabaseOutlined />}
              valueStyle={{
                color: Object.values(overview.system_health).every(s => s === 'healthy')
                  ? '#52c41a'
                  : '#faad14'
              }}
            />
            <Divider style={{ margin: '12px 0' }} />
            <Text type="secondary" style={{ fontSize: 12 }}>
              {Object.values(overview.system_health).every(s => s === 'healthy')
                ? '所有系统运行正常'
                : '部分系统需要注意'}
            </Text>
          </Card>
        </Col>
      </Row>

      {/* 第二行：订阅到期警告和快速统计 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        {/* 订阅到期统计 */}
        <Col xs={24} lg={12}>
          <Card
            title={
              <Space>
                <WarningOutlined style={{ color: '#faad14' }} />
                <span>订阅到期状态</span>
              </Space>
            }
            extra={
              <Button type="link" onClick={() => navigate('/organizations')}>
                查看全部
              </Button>
            }
          >
            <Row gutter={16}>
              <Col span={12}>
                <Card size="small" style={{ textAlign: 'center', background: '#f0f5ff', border: '1px solid #d6e4ff' }}>
                  <Statistic
                    title="正常"
                    value={subscriptionStats.active}
                    valueStyle={{ color: '#52c41a', fontSize: 28 }}
                    prefix={<CheckCircleOutlined />}
                  />
                </Card>
              </Col>
              <Col span={12}>
                <Card size="small" style={{ textAlign: 'center', background: '#fffbe6', border: '1px solid #ffe58f' }}>
                  <Statistic
                    title="即将到期"
                    value={subscriptionStats.expiring_soon}
                    valueStyle={{ color: '#faad14', fontSize: 28 }}
                    prefix={<WarningOutlined />}
                  />
                </Card>
              </Col>
              <Col span={12} style={{ marginTop: 16 }}>
                <Card size="small" style={{ textAlign: 'center', background: '#fff1f0', border: '1px solid #ffccc7' }}>
                  <Statistic
                    title="已到期"
                    value={subscriptionStats.expired}
                    valueStyle={{ color: '#ff4d4f', fontSize: 28 }}
                    prefix={<CloseCircleOutlined />}
                  />
                </Card>
              </Col>
              <Col span={12} style={{ marginTop: 16 }}>
                <Card size="small" style={{ textAlign: 'center', background: '#fafafa', border: '1px solid #d9d9d9' }}>
                  <Statistic
                    title="无订阅"
                    value={subscriptionStats.none}
                    valueStyle={{ color: '#8c8c8c', fontSize: 28 }}
                    prefix={<ClockCircleOutlined />}
                  />
                </Card>
              </Col>
            </Row>
          </Card>
        </Col>

        {/* 数据库状态 */}
        <Col xs={24} lg={12}>
          <Card
            title={
              <Space>
                <DatabaseOutlined />
                <span>数据库状态</span>
              </Space>
            }
          >
            <Space direction="vertical" style={{ width: '100%' }} size="large">
              {/* PostgreSQL */}
              <div>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8 }}>
                  <Space>
                    <Badge status={databaseStatus.postgresql.status === 'healthy' ? 'success' : 'error'} />
                    <Text strong>PostgreSQL</Text>
                  </Space>
                  <Tag color={getDatabaseStatusColor(databaseStatus.postgresql.status)}>
                    {databaseStatus.postgresql.status}
                  </Tag>
                </div>
                <div style={{ marginBottom: 4 }}>
                  <Text type="secondary" style={{ fontSize: 12 }}>
                    连接数: {databaseStatus.postgresql.connections} / {databaseStatus.postgresql.max_connections}
                  </Text>
                </div>
                <Progress
                  percent={Math.round((databaseStatus.postgresql.connections / databaseStatus.postgresql.max_connections) * 100)}
                  size="small"
                  status={
                    (databaseStatus.postgresql.connections / databaseStatus.postgresql.max_connections) > 0.8
                      ? 'exception'
                      : 'normal'
                  }
                />
              </div>

              {/* ClickHouse */}
              <div>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8 }}>
                  <Space>
                    <Badge status="success" />
                    <Text strong>ClickHouse</Text>
                  </Space>
                  <Tag color="success">运行中</Tag>
                </div>
                <Text type="secondary" style={{ fontSize: 12 }}>
                  {Object.keys(databaseStatus.clickhouse).length} 个数据库正在运行
                </Text>
              </div>

              {/* Redis */}
              <div>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8 }}>
                  <Space>
                    <Badge status={databaseStatus.redis.status === 'healthy' ? 'success' : 'error'} />
                    <Text strong>Redis</Text>
                  </Space>
                  <Tag color={getDatabaseStatusColor(databaseStatus.redis.status)}>
                    {databaseStatus.redis.status}
                  </Tag>
                </div>
                <div style={{ marginBottom: 4 }}>
                  <Text type="secondary" style={{ fontSize: 12 }}>
                    内存使用: {(databaseStatus.redis.used_memory / 1024 / 1024).toFixed(1)} MB
                    {databaseStatus.redis.max_memory
                      ? ` / ${(databaseStatus.redis.max_memory / 1024 / 1024).toFixed(0)} MB`
                      : ''
                    }
                  </Text>
                </div>
                {databaseStatus.redis.max_memory && (
                  <Progress
                    percent={Math.round((databaseStatus.redis.used_memory / databaseStatus.redis.max_memory) * 100)}
                    size="small"
                    status={
                      (databaseStatus.redis.used_memory / databaseStatus.redis.max_memory) > 0.8
                        ? 'exception'
                        : 'normal'
                    }
                  />
                )}
              </div>
            </Space>
          </Card>
        </Col>
      </Row>

      {/* 第三行：即将到期的组织 */}
      {expiringOrgs.length > 0 && (
        <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
          <Col span={24}>
            <Alert
              message={
                <Space>
                  <WarningOutlined />
                  <span>订阅到期提醒</span>
                </Space>
              }
              description={
                <div>
                  <Text>以下组织的订阅即将到期，请及时提醒续费：</Text>
                  <Table
                    columns={[
                      {
                        title: '组织',
                        key: 'org',
                        render: (record: Organization) => (
                          <Space>
                            <TeamOutlined />
                            <span>{record.name}</span>
                          </Space>
                        ),
                      },
                      {
                        title: '剩余天数',
                        key: 'days',
                        render: (record: Organization) => (
                          <Tag color="warning">
                            <ClockCircleOutlined /> {record.days_until_expiration} 天
                          </Tag>
                        ),
                      },
                      {
                        title: '到期时间',
                        key: 'end_date',
                        render: (record: Organization) => (
                          record.subscription_end_date
                            ? new Date(record.subscription_end_date).toLocaleDateString('zh-CN')
                            : '-'
                        ),
                      },
                      {
                        title: '操作',
                        key: 'action',
                        render: (record: Organization) => (
                          <Space>
                            <Button
                              type="link"
                              size="small"
                              icon={<MailOutlined />}
                              onClick={() => navigate(`/organizations/${record.id}`)}
                            >
                              发送提醒
                            </Button>
                            <Button
                              type="link"
                              size="small"
                              onClick={() => navigate(`/organizations/${record.id}`)}
                            >
                              查看详情
                            </Button>
                          </Space>
                        ),
                      },
                    ]}
                    dataSource={expiringOrgs}
                    rowKey="id"
                    pagination={false}
                    size="small"
                    style={{ marginTop: 12 }}
                  />
                </div>
              }
              type="warning"
              showIcon
            />
          </Col>
        </Row>
      )}

      {/* 第四行：活跃组织Top 5 */}
      <Row gutter={[16, 16]}>
        <Col span={24}>
          <Card
            title={
              <Space>
                <RiseOutlined style={{ color: '#1890ff' }} />
                <span>活跃组织 Top 5</span>
              </Space>
            }
            extra={
              <Button type="link" onClick={() => navigate('/organizations')}>
                查看全部组织
              </Button>
            }
          >
            <Table
              columns={organizationColumns}
              dataSource={topOrganizations}
              rowKey="id"
              pagination={false}
              size="middle"
            />
          </Card>
        </Col>
      </Row>

      {/* 底部提示信息 */}
      <Row style={{ marginTop: 16 }}>
        <Col span={24}>
          <Text type="secondary" style={{ fontSize: 12 }}>
            <ClockCircleOutlined /> 数据每30秒自动刷新 · 最后更新: {new Date().toLocaleString('zh-CN')}
          </Text>
        </Col>
      </Row>
    </div>
  );
};

export default Dashboard;
