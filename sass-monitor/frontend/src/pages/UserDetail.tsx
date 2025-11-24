import React, { useState, useEffect } from 'react';
import {
  Card,
  Descriptions,
  Table,
  Tag,
  Spin,
  message,
  Button,
  Row,
  Col,
  Statistic,
  Avatar,
  Space,
  Tabs,
} from 'antd';
import {
  ArrowLeftOutlined,
  UserOutlined,
  TeamOutlined,
  CloudServerOutlined,
  DollarOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
} from '@ant-design/icons';
import { useParams, useNavigate } from 'react-router-dom';
import UserService from '../services/userService';
import {
  LightAdminUser,
  UserOrganization,
  UserWorkspace,
  UserSubscription,
} from '../types';
import type { ColumnType } from 'antd/es/table';

const UserDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [user, setUser] = useState<LightAdminUser | null>(null);
  const [organizations, setOrganizations] = useState<UserOrganization[]>([]);
  const [workspaces, setWorkspaces] = useState<UserWorkspace[]>([]);
  const [subscriptions, setSubscriptions] = useState<UserSubscription[]>([]);
  const [organizationsTotal, setOrganizationsTotal] = useState(0);
  const [workspacesTotal, setWorkspacesTotal] = useState(0);
  const [subscriptionsTotal, setSubscriptionsTotal] = useState(0);

  // 获取用户详情
  const fetchUserDetail = async () => {
    if (!id) return;

    try {
      setLoading(true);
      const userData = await UserService.getUserById(id);
      setUser(userData);
    } catch (error: any) {
      message.error('获取用户详情失败');
      console.error('Fetch user detail error:', error);
    } finally {
      setLoading(false);
    }
  };

  // 获取用户组织
  const fetchOrganizations = async () => {
    if (!id) return;

    try {
      const response = await UserService.getUserOrganizations(id, {
        page: 1,
        page_size: 100,
      });
      setOrganizations(response.data || []);
      setOrganizationsTotal(response.total || 0);
    } catch (error: any) {
      console.error('Fetch organizations error:', error);
    }
  };

  // 获取用户工作空间
  const fetchWorkspaces = async () => {
    if (!id) return;

    try {
      const response = await UserService.getUserWorkspaces(id, {
        page: 1,
        page_size: 100,
      });
      setWorkspaces(response.data || []);
      setWorkspacesTotal(response.total || 0);
    } catch (error: any) {
      console.error('Fetch workspaces error:', error);
    }
  };

  // 获取用户订阅
  const fetchSubscriptions = async () => {
    if (!id) return;

    try {
      const response = await UserService.getUserSubscriptions(id, {
        page: 1,
        page_size: 100,
      });
      setSubscriptions(response.data || []);
      setSubscriptionsTotal(response.total || 0);
    } catch (error: any) {
      console.error('Fetch subscriptions error:', error);
    }
  };

  useEffect(() => {
    fetchUserDetail();
    fetchOrganizations();
    fetchWorkspaces();
    fetchSubscriptions();
  }, [id]);

  // 组织表格列
  const organizationColumns: ColumnType<UserOrganization>[] = [
    {
      title: '组织名称',
      dataIndex: 'organization_name',
      key: 'organization_name',
      render: (text: string) => <strong>{text}</strong>,
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      render: (text: string | null) => text || '-',
    },
    {
      title: '加入时间',
      dataIndex: 'joined_at',
      key: 'joined_at',
      render: (text: string) => new Date(text).toLocaleString('zh-CN'),
    },
  ];

  // 工作空间表格列
  const workspaceColumns: ColumnType<UserWorkspace>[] = [
    {
      title: '工作空间',
      dataIndex: 'workspace_name',
      key: 'workspace_name',
      render: (text: string) => <strong>{text}</strong>,
    },
    {
      title: '所属组织',
      dataIndex: 'organization_name',
      key: 'organization_name',
    },
    {
      title: '用户状态',
      dataIndex: 'user_status',
      key: 'user_status',
      render: (status: string) => (
        <Tag color={status === 'active' ? 'green' : 'default'}>
          {status === 'active' ? '活跃' : status}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (text: string) => new Date(text).toLocaleString('zh-CN'),
    },
  ];

  // 订阅表格列
  const subscriptionColumns: ColumnType<UserSubscription>[] = [
    {
      title: '订阅计划',
      dataIndex: 'plan_name',
      key: 'plan_name',
      render: (text: string) => <strong>{text}</strong>,
    },
    {
      title: '所属组织',
      dataIndex: 'organization_name',
      key: 'organization_name',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const statusConfig: Record<string, { color: string; text: string }> = {
          active: { color: 'green', text: '活跃' },
          trial: { color: 'blue', text: '试用' },
          expired: { color: 'red', text: '已过期' },
          cancelled: { color: 'default', text: '已取消' },
        };
        const config = statusConfig[status] || { color: 'default', text: status };
        return <Tag color={config.color}>{config.text}</Tag>;
      },
    },
    {
      title: '计费周期',
      dataIndex: 'billing_cycle',
      key: 'billing_cycle',
      render: (cycle: string) => {
        const cycleMap: Record<string, string> = {
          monthly: '月付',
          quarterly: '季付',
          yearly: '年付',
        };
        return cycleMap[cycle] || cycle;
      },
    },
    {
      title: '开始时间',
      dataIndex: 'start_date',
      key: 'start_date',
      render: (text: string) => new Date(text).toLocaleDateString('zh-CN'),
    },
    {
      title: '结束时间',
      dataIndex: 'end_date',
      key: 'end_date',
      render: (text: string | null) => text ? new Date(text).toLocaleDateString('zh-CN') : '-',
    },
  ];

  if (loading) {
    return (
      <div style={{ padding: 24, textAlign: 'center' }}>
        <Spin size="large" />
      </div>
    );
  }

  if (!user) {
    return (
      <div style={{ padding: 24 }}>
        <Card>
          <p>用户不存在</p>
          <Button onClick={() => navigate('/users')}>返回用户列表</Button>
        </Card>
      </div>
    );
  }

  return (
    <div style={{ padding: 24, background: 'var(--theme-bg-primary)', minHeight: '100vh' }}>
      {/* 返回按钮 */}
      <div style={{ marginBottom: 16 }}>
        <Button
          icon={<ArrowLeftOutlined />}
          onClick={() => navigate('/users')}
        >
          返回用户列表
        </Button>
      </div>

      {/* 用户基本信息 */}
      <Card title="用户基本信息" style={{ marginBottom: 24 }}>
        <Row gutter={16}>
          <Col span={4}>
            <div style={{ textAlign: 'center' }}>
              <Avatar
                size={100}
                src={user.avatar_url}
                icon={!user.avatar_url ? <UserOutlined /> : undefined}
              />
            </div>
          </Col>
          <Col span={20}>
            <Descriptions column={2}>
              <Descriptions.Item label="用户名">{user.username}</Descriptions.Item>
              <Descriptions.Item label="昵称">{user.nickname || '-'}</Descriptions.Item>
              <Descriptions.Item label="邮箱">
                <Space>
                  {user.email || '-'}
                  {user.email_verified && (
                    <Tag icon={<CheckCircleOutlined />} color="success">
                      已验证
                    </Tag>
                  )}
                </Space>
              </Descriptions.Item>
              <Descriptions.Item label="认证方式">
                <Tag color={user.oauth_provider ? 'blue' : 'default'}>
                  {user.oauth_provider || '本地'}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="Clerk用户ID">
                {user.clerk_user_id || '-'}
              </Descriptions.Item>
              <Descriptions.Item label="创建时间">
                {new Date(user.created_at).toLocaleString('zh-CN')}
              </Descriptions.Item>
            </Descriptions>
          </Col>
        </Row>
      </Card>

      {/* 统计卡片 */}
      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={8}>
          <Card>
            <Statistic
              title="所属组织"
              value={user.organization_count}
              prefix={<TeamOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card>
            <Statistic
              title="工作空间"
              value={user.workspace_count}
              prefix={<CloudServerOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card>
            <Statistic
              title="活跃订阅"
              value={user.subscription_count}
              prefix={<DollarOutlined />}
              valueStyle={{ color: '#fa8c16' }}
            />
          </Card>
        </Col>
      </Row>

      {/* 详细信息标签页 */}
      <Card>
        <Tabs
          items={[
            {
              key: 'organizations',
              label: (
                <span>
                  <TeamOutlined />
                  所属组织 ({organizationsTotal})
                </span>
              ),
              children: (
                <Table
                  dataSource={organizations}
                  columns={organizationColumns}
                  rowKey="organization_id"
                  pagination={false}
                />
              ),
            },
            {
              key: 'workspaces',
              label: (
                <span>
                  <CloudServerOutlined />
                  工作空间 ({workspacesTotal})
                </span>
              ),
              children: (
                <Table
                  dataSource={workspaces}
                  columns={workspaceColumns}
                  rowKey="workspace_id"
                  pagination={false}
                />
              ),
            },
            {
              key: 'subscriptions',
              label: (
                <span>
                  <DollarOutlined />
                  订阅记录 ({subscriptionsTotal})
                </span>
              ),
              children: (
                <Table
                  dataSource={subscriptions}
                  columns={subscriptionColumns}
                  rowKey="subscription_id"
                  pagination={false}
                />
              ),
            },
          ]}
        />
      </Card>
    </div>
  );
};

export default UserDetail;
