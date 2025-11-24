import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Card,
  Descriptions,
  Table,
  Tabs,
  Tag,
  Button,
  message,
  Space,
  Statistic,
  Row,
  Col,
  Typography,
  Spin,
} from 'antd';
import {
  ArrowLeftOutlined,
  TeamOutlined,
  CloudServerOutlined,
  DollarOutlined,
  CalendarOutlined,
  UserOutlined,
  MailOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  WarningOutlined,
  ClockCircleOutlined,
} from '@ant-design/icons';
import OrganizationService from '../services/organizationService';
import { Organization, OrganizationSubscription } from '../types';
import type { ColumnType } from 'antd/es/table';

const { Title, Text } = Typography;
const { TabPane } = Tabs;

const OrganizationDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [organization, setOrganization] = useState<Organization | null>(null);
  const [subscriptions, setSubscriptions] = useState<OrganizationSubscription[]>([]);
  const [subscriptionsLoading, setSubscriptionsLoading] = useState(false);
  const [subscriptionsPage, setSubscriptionsPage] = useState(1);
  const [subscriptionsTotal, setSubscriptionsTotal] = useState(0);

  useEffect(() => {
    if (id) {
      fetchOrganizationDetail();
      fetchSubscriptions();
    }
  }, [id]);

  const fetchOrganizationDetail = async () => {
    try {
      setLoading(true);
      const data = await OrganizationService.getOrganizationById(id!);
      setOrganization(data);
    } catch (error: any) {
      message.error('获取组织详情失败');
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  const fetchSubscriptions = async () => {
    try {
      setSubscriptionsLoading(true);
      const response = await OrganizationService.getOrganizationSubscriptions(id!, {
        page: subscriptionsPage,
        page_size: 10,
      });
      setSubscriptions(response.data || []);
      setSubscriptionsTotal(response.total || 0);
    } catch (error: any) {
      message.error('获取订阅信息失败');
      console.error(error);
    } finally {
      setSubscriptionsLoading(false);
    }
  };

  const handleSendReminder = async () => {
    try {
      await OrganizationService.sendExpiryReminder(id!);
      message.success('到期提醒已发送（预留功能）');
    } catch (error: any) {
      message.error('发送提醒失败');
    }
  };

  const getSubscriptionStatusTag = () => {
    if (!organization) return null;

    const status = organization.subscription_status || 'none';
    const days = organization.days_until_expiration;

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
  };

  // 订阅列表表格列
  const subscriptionColumns: ColumnType<OrganizationSubscription>[] = [
    {
      title: '订阅用户',
      key: 'user',
      render: (record: OrganizationSubscription) => (
        <Space direction="vertical" size={0}>
          <Space>
            <UserOutlined />
            <Text strong>{record.username || '未知用户'}</Text>
          </Space>
          {record.user_email && (
            <Text type="secondary" style={{ fontSize: 12 }}>
              <MailOutlined /> {record.user_email}
            </Text>
          )}
        </Space>
      ),
    },
    {
      title: '订阅套餐',
      key: 'plan',
      render: (record: OrganizationSubscription) => (
        <Space direction="vertical" size={0}>
          <Text strong>{record.plan_name}</Text>
          <Text type="secondary" style={{ fontSize: 12 }}>
            <DollarOutlined /> ¥{record.plan_pricing?.toFixed(2) || '0.00'}/{record.billing_cycle === 'monthly' ? '月' : record.billing_cycle === 'yearly' ? '年' : '季'}
          </Text>
        </Space>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const statusConfig: Record<string, { color: string; text: string }> = {
          active: { color: 'success', text: '活跃' },
          trial: { color: 'processing', text: '试用中' },
          expired: { color: 'error', text: '已到期' },
          cancelled: { color: 'default', text: '已取消' },
        };
        const config = statusConfig[status] || { color: 'default', text: status };
        return <Tag color={config.color}>{config.text}</Tag>;
      },
    },
    {
      title: '开始时间',
      dataIndex: 'start_date',
      key: 'start_date',
      render: (date: string) => new Date(date).toLocaleDateString('zh-CN'),
    },
    {
      title: '到期时间',
      dataIndex: 'end_date',
      key: 'end_date',
      render: (date: string | null, record: OrganizationSubscription) => {
        if (!date) return <Text type="secondary">-</Text>;
        const daysLeft = record.days_until_expiry;
        return (
          <Space direction="vertical" size={0}>
            <Text>{new Date(date).toLocaleDateString('zh-CN')}</Text>
            {daysLeft !== undefined && (
              <Text type={daysLeft < 7 ? 'danger' : daysLeft < 30 ? 'warning' : 'secondary'} style={{ fontSize: 12 }}>
                {daysLeft > 0 ? `剩余${daysLeft}天` : daysLeft === 0 ? '今天到期' : `已过期${Math.abs(daysLeft)}天`}
              </Text>
            )}
          </Space>
        );
      },
    },
    {
      title: '付款方式',
      dataIndex: 'payment_method',
      key: 'payment_method',
      render: (method: string | null) => method || <Text type="secondary">-</Text>,
    },
    {
      title: '试用天数',
      dataIndex: 'trial_days_used',
      key: 'trial_days_used',
      render: (days: number | null) => days !== null && days !== undefined ? `${days}天` : <Text type="secondary">-</Text>,
    },
  ];

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '100px 0' }}>
        <Spin size="large" />
      </div>
    );
  }

  if (!organization) {
    return <div>组织不存在</div>;
  }

  return (
    <div style={{ padding: 24 }}>
      <Button
        icon={<ArrowLeftOutlined />}
        onClick={() => navigate('/organizations')}
        style={{ marginBottom: 16 }}
      >
        返回列表
      </Button>

      <Card>
        <Space direction="vertical" size="large" style={{ width: '100%' }}>
          <div>
            <Title level={2}>
              <TeamOutlined /> {organization.name}
            </Title>
            <Text type="secondary">组织ID: {organization.id}</Text>
          </div>

          {/* 统计卡片 */}
          <Row gutter={16}>
            <Col span={6}>
              <Card>
                <Statistic
                  title="用户数量"
                  value={organization.user_count}
                  prefix={<UserOutlined />}
                />
              </Card>
            </Col>
            <Col span={6}>
              <Card>
                <Statistic
                  title="订阅数量"
                  value={organization.subscription_count}
                  prefix={<CloudServerOutlined />}
                />
              </Card>
            </Col>
            <Col span={6}>
              <Card>
                <Statistic
                  title="活跃订阅"
                  value={organization.active_subscription_count || 0}
                  prefix={<CheckCircleOutlined />}
                  valueStyle={{ color: '#3f8600' }}
                />
              </Card>
            </Col>
            <Col span={6}>
              <Card>
                <Statistic
                  title="工作空间"
                  value={organization.workspace_count || 0}
                  prefix={<TeamOutlined />}
                />
              </Card>
            </Col>
          </Row>

          {/* 基本信息 */}
          <Descriptions title="基本信息" bordered column={2}>
            <Descriptions.Item label="组织名称">{organization.name}</Descriptions.Item>
            <Descriptions.Item label="拥有者ID">{organization.owner_id}</Descriptions.Item>
            <Descriptions.Item label="描述" span={2}>
              {organization.description || <Text type="secondary">无</Text>}
            </Descriptions.Item>
            <Descriptions.Item label="订阅状态">
              {getSubscriptionStatusTag()}
            </Descriptions.Item>
            <Descriptions.Item label="订阅到期时间">
              {organization.subscription_end_date ? (
                <Space>
                  <CalendarOutlined />
                  <Text>{new Date(organization.subscription_end_date).toLocaleDateString('zh-CN')}</Text>
                </Space>
              ) : (
                <Text type="secondary">-</Text>
              )}
            </Descriptions.Item>
            <Descriptions.Item label="创建时间">
              {new Date(organization.created_at).toLocaleString('zh-CN')}
            </Descriptions.Item>
            <Descriptions.Item label="更新时间">
              {organization.updated_at ? new Date(organization.updated_at).toLocaleString('zh-CN') : <Text type="secondary">-</Text>}
            </Descriptions.Item>
          </Descriptions>

          {/* 操作按钮 */}
          {organization.subscription_status === 'expiring_soon' && (
            <Card>
              <Space>
                <WarningOutlined style={{ color: '#faad14' }} />
                <Text>订阅即将到期</Text>
                <Button type="primary" icon={<MailOutlined />} onClick={handleSendReminder}>
                  发送到期提醒
                </Button>
              </Space>
            </Card>
          )}

          {/* 订阅信息标签页 */}
          <Tabs defaultActiveKey="subscriptions">
            <TabPane tab={`订阅列表 (${subscriptionsTotal})`} key="subscriptions">
              <Table
                loading={subscriptionsLoading}
                dataSource={subscriptions}
                columns={subscriptionColumns}
                rowKey="id"
                pagination={{
                  current: subscriptionsPage,
                  pageSize: 10,
                  total: subscriptionsTotal,
                  onChange: (page) => {
                    setSubscriptionsPage(page);
                    fetchSubscriptions();
                  },
                  showSizeChanger: false,
                  showTotal: (total) => `共 ${total} 条记录`,
                }}
              />
            </TabPane>
          </Tabs>
        </Space>
      </Card>
    </div>
  );
};

export default OrganizationDetail;
