import React, { useState, useEffect, useCallback } from 'react';
import {
  Card,
  Table,
  Button,
  message,
  Tag,
  Space,
  Input,
  Row,
  Col,
  Statistic,
  Tooltip,
  Avatar,
} from 'antd';
import {
  ReloadOutlined,
  EyeOutlined,
  UserOutlined,
  TeamOutlined,
  CloudServerOutlined,
  DollarOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import UserService from '../services/userService';
import { LightAdminUser } from '../types';
import type { ColumnType } from 'antd/es/table';

const Users: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [users, setUsers] = useState<LightAdminUser[]>([]);
  const [total, setTotal] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [searchText, setSearchText] = useState('');

  // 获取用户列表
  const fetchUsers = useCallback(async () => {
    try {
      setLoading(true);
      const response = await UserService.getUsers({
        page: currentPage,
        page_size: pageSize,
        search: searchText,
      });
      setUsers(response.data || []);
      setTotal(response.total || 0);
    } catch (error: any) {
      message.error('获取用户列表失败');
      console.error('Fetch users error:', error);
    } finally {
      setLoading(false);
    }
  }, [currentPage, pageSize, searchText]);

  useEffect(() => {
    fetchUsers();
  }, [currentPage, pageSize, searchText, fetchUsers]);

  // 处理搜索
  const handleSearch = (value: string) => {
    setSearchText(value);
    setCurrentPage(1);
  };

  // 查看用户详情
  const handleViewDetail = (userId: string) => {
    navigate(`/users/${userId}`);
  };

  // 表格列定义
  const columns: ColumnType<LightAdminUser>[] = [
    {
      title: '头像',
      dataIndex: 'avatar_url',
      key: 'avatar',
      width: 80,
      render: (avatarUrl: string, record: LightAdminUser) => (
        <Avatar
          src={avatarUrl}
          icon={!avatarUrl ? <UserOutlined /> : undefined}
          alt={record.username}
        />
      ),
    },
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
      width: 150,
      render: (text: string) => <strong>{text}</strong>,
    },
    {
      title: '昵称',
      dataIndex: 'nickname',
      key: 'nickname',
      width: 150,
      render: (text: string | null) => text || '-',
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      key: 'email',
      width: 200,
      ellipsis: true,
      render: (text: string | null, record: LightAdminUser) => (
        <Space>
          {text || '-'}
          {record.email_verified && (
            <Tooltip title="邮箱已验证">
              <CheckCircleOutlined style={{ color: '#52c41a' }} />
            </Tooltip>
          )}
        </Space>
      ),
    },
    {
      title: '认证方式',
      dataIndex: 'oauth_provider',
      key: 'oauth_provider',
      width: 120,
      render: (text: string) => (
        <Tag color={text ? 'blue' : 'default'}>
          {text || '本地'}
        </Tag>
      ),
    },
    {
      title: '组织数',
      dataIndex: 'organization_count',
      key: 'organization_count',
      width: 100,
      align: 'center' as const,
      render: (count: number) => (
        <Tag icon={<TeamOutlined />} color="blue">
          {count}
        </Tag>
      ),
    },
    {
      title: '工作空间数',
      dataIndex: 'workspace_count',
      key: 'workspace_count',
      width: 120,
      align: 'center' as const,
      render: (count: number) => (
        <Tag icon={<CloudServerOutlined />} color="cyan">
          {count}
        </Tag>
      ),
    },
    {
      title: '订阅数',
      dataIndex: 'subscription_count',
      key: 'subscription_count',
      width: 100,
      align: 'center' as const,
      render: (count: number) => (
        <Tag icon={<DollarOutlined />} color="gold">
          {count}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
      render: (text: string) => new Date(text).toLocaleString('zh-CN'),
    },
    {
      title: '操作',
      key: 'action',
      width: 120,
      fixed: 'right' as const,
      render: (_: any, record: LightAdminUser) => (
        <Space>
          <Tooltip title="查看详情">
            <Button
              type="link"
              size="small"
              icon={<EyeOutlined />}
              onClick={() => handleViewDetail(record.id)}
            >
              详情
            </Button>
          </Tooltip>
        </Space>
      ),
    },
  ];

  return (
    <div style={{ padding: 24, background: 'var(--theme-bg-primary)', minHeight: '100vh' }}>
      {/* 页面标题和操作栏 */}
      <div style={{ marginBottom: 24, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h2 style={{ margin: 0, color: 'var(--theme-text-primary)' }}>
          <UserOutlined style={{ marginRight: 8 }} />
          用户管理
        </h2>
        <Space>
          <Space.Compact>
            <Input.Search
              placeholder="搜索用户名、昵称或邮箱"
              allowClear
              style={{ width: 300 }}
              onSearch={handleSearch}
            />
          </Space.Compact>
          <Button
            icon={<ReloadOutlined />}
            onClick={fetchUsers}
            loading={loading}
          >
            刷新
          </Button>
        </Space>
      </div>

      {/* 统计卡片 */}
      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={6}>
          <Card>
            <Statistic
              title="总用户数"
              value={total}
              prefix={<UserOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={6}>
          <Card>
            <Statistic
              title="总组织关联"
              value={users?.reduce((sum, user) => sum + user.organization_count, 0) || 0}
              prefix={<TeamOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={6}>
          <Card>
            <Statistic
              title="总工作空间"
              value={users?.reduce((sum, user) => sum + user.workspace_count, 0) || 0}
              prefix={<CloudServerOutlined />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={6}>
          <Card>
            <Statistic
              title="活跃订阅"
              value={users?.reduce((sum, user) => sum + user.subscription_count, 0) || 0}
              prefix={<DollarOutlined />}
              valueStyle={{ color: '#fa8c16' }}
            />
          </Card>
        </Col>
      </Row>

      {/* 用户列表表格 */}
      <Card>
        <Table
          dataSource={users}
          columns={columns}
          rowKey="id"
          loading={loading}
          pagination={{
            current: currentPage,
            pageSize: pageSize,
            total: total,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total) => `共 ${total} 条记录`,
            onChange: (page, size) => {
              setCurrentPage(page);
              setPageSize(size || 10);
            },
          }}
          scroll={{ x: 1400 }}
        />
      </Card>
    </div>
  );
};

export default Users;
