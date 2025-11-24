import React, { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Card,
  Table,
  Button,
  Space,
  Tag,
  Input,
  message,
  Typography,
  Row,
  Col,
  Statistic,
  Tooltip,
} from 'antd';
import {
  ReloadOutlined,
  EyeOutlined,
  TeamOutlined,
  UserOutlined,
  CloudServerOutlined,
  ClockCircleOutlined,
  WarningOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  MailOutlined,
} from '@ant-design/icons';
import OrganizationService from '../services/organizationService';
import { Organization } from '../types';

const { Title } = Typography;

const Organizations: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [organizations, setOrganizations] = useState<Organization[]>([]);
  const [total, setTotal] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [searchText, setSearchText] = useState('');

  // 获取组织列表
  const fetchOrganizations = useCallback(async () => {
    try {
      setLoading(true);
      const response = await OrganizationService.getOrganizations({
        page: currentPage,
        page_size: pageSize,
        search: searchText,
      });
      setOrganizations(response.data || []);
      setTotal(response.total || 0);
    } catch (error: any) {
      message.error('获取组织列表失败');
      console.error('Fetch organizations error:', error);
    } finally {
      setLoading(false);
    }
  }, [currentPage, pageSize, searchText]);

  useEffect(() => {
    fetchOrganizations();
  }, [currentPage, pageSize, searchText, fetchOrganizations]);

  // 处理搜索
  const handleSearch = (value: string) => {
    setSearchText(value);
    setCurrentPage(1);
  };

  // 处理分页
  const handlePageChange = (page: number, size: number) => {
    setCurrentPage(page);
    setPageSize(size);
  };

  // 刷新数据
  const handleRefresh = () => {
    fetchOrganizations();
  };

  // 查看组织详情
  const handleViewDetail = (organization: Organization) => {
    navigate(`/organizations/${organization.id}`);
  };

  // 发送到期提醒
  const handleSendReminder = async (organizationId: string) => {
    try {
      await OrganizationService.sendExpiryReminder(organizationId);
      message.success('到期提醒已发送（预留功能）');
    } catch (error: any) {
      message.error('发送提醒失败: ' + (error.response?.data?.error || error.message));
    }
  };

  const columns = [
    {
      title: '组织名称',
      dataIndex: 'name',
      key: 'name',
      render: (name: string, record: Organization) => (
        <Space>
          <TeamOutlined />
          <div>
            <div style={{ fontWeight: 500 }}>{name}</div>
            <div style={{ fontSize: 12, color: '#666' }}>
              ID: {record.id.slice(0, 8)}...
            </div>
          </div>
        </Space>
      ),
    },
    {
      title: '用户数量',
      dataIndex: 'user_count',
      key: 'user_count',
      render: (count: number) => (
        <Statistic
          value={count}
          valueStyle={{ fontSize: 16 }}
          prefix={<UserOutlined style={{ color: '#1890ff' }} />}
        />
      ),
      sorter: (a: Organization, b: Organization) => a.user_count - b.user_count,
    },
    {
      title: '订阅数量',
      dataIndex: 'subscription_count',
      key: 'subscription_count',
      render: (count: number) => (
        <Statistic
          value={count}
          valueStyle={{ fontSize: 16 }}
          prefix={<CloudServerOutlined style={{ color: '#52c41a' }} />}
        />
      ),
      sorter: (a: Organization, b: Organization) => a.subscription_count - b.subscription_count,
    },
    {
      title: '订阅到期状态',
      key: 'subscription_expiry',
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

        return (
          <Space>
            <Tag color={config.color} icon={config.icon}>
              {config.text}
            </Tag>
            {status === 'expiring_soon' && (
              <Tooltip title="发送到期提醒">
                <Button
                  type="text"
                  size="small"
                  icon={<MailOutlined />}
                  onClick={() => handleSendReminder(record.id)}
                />
              </Tooltip>
            )}
          </Space>
        );
      },
    },
    {
      title: '存储使用',
      dataIndex: 'storage_usage',
      key: 'storage_usage',
      render: (usage: number) => {
        const usageMB = usage / (1024 * 1024);
        return `${usageMB.toFixed(1)} MB`;
      },
      sorter: (a: Organization, b: Organization) => a.storage_usage - b.storage_usage,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (dateStr: string) => new Date(dateStr).toLocaleDateString(),
      sorter: (a: Organization, b: Organization) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime(),
    },
    {
      title: '状态',
      key: 'status',
      render: (record: Organization) => (
        <Tag color={record.user_count > 0 ? 'green' : 'orange'}>
          {record.user_count > 0 ? '活跃' : '待激活'}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'actions',
      align: 'center' as const,
      render: (record: Organization) => (
        <Space>
          <Tooltip title="查看详情">
            <Button
              type="text"
              icon={<EyeOutlined />}
              size="small"
              onClick={() => handleViewDetail(record)}
            />
          </Tooltip>
        </Space>
      ),
    },
  ];

  // 计算统计数据（已在统计卡片中使用）

  return (
    <div style={{ padding: 24, background: 'var(--theme-bg-primary)', minHeight: '100vh' }}>
      {/* 页面标题和操作栏 */}
      <div style={{ marginBottom: 24, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Title level={2} style={{ margin: 0, color: 'var(--theme-text-primary)' }}>
          组织管理
        </Title>
        <Space.Compact>
          <Input.Search
            placeholder="搜索组织..."
            allowClear
            style={{ width: 250 }}
            onSearch={handleSearch}
            onChange={(e) => e.target.value && handleSearch(e.target.value)}
          />
          <Button icon={<ReloadOutlined />} onClick={handleRefresh} loading={loading}>
            刷新
          </Button>
        </Space.Compact>
      </div>

      {/* 统计卡片 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={8}>
          <Card>
            <Statistic
              title="总组织数"
              value={total}
              prefix={<TeamOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card>
            <Statistic
              title="总用户数"
              value={organizations?.reduce((sum, org) => sum + org.user_count, 0) || 0}
              prefix={<UserOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card>
            <Statistic
              title="活跃组织"
              value={organizations?.filter(org => org.user_count > 0).length || 0}
              prefix={<CloudServerOutlined />}
              valueStyle={{ color: '#fa8c16' }}
            />
          </Card>
        </Col>
      </Row>

      {/* 组织列表 */}
      <Card>
        <Table
          columns={columns}
          dataSource={organizations}
          rowKey="id"
          loading={loading}
          pagination={{
            current: currentPage,
            pageSize: pageSize,
            total: total,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
            onChange: handlePageChange,
          }}
          scroll={{ x: 1000 }}
        />
      </Card>
    </div>
  );
};

export default Organizations;