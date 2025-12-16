import React, { useState, useEffect, useCallback } from 'react';
import {
  Row,
  Col,
  Card,
  Table,
  Button,
  Modal,
  Form,
  Input,
  InputNumber,
  Switch,
  Select,
  message,
  Popconfirm,
  Tag,
  Tooltip,
  Space,
  Typography,
  Statistic,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ReloadOutlined,
  EyeOutlined,
  DollarOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
} from '@ant-design/icons';
import { SubscriptionPlan, CreateSubscriptionPlanRequest, UpdateSubscriptionPlanRequest } from '../types';
import SubscriptionPlanService from '../services/subscriptionPlanService';

const { Title, Text } = Typography;

// 流量套餐类型枚举
const FLOW_PACKAGE_OPTIONS = [
  { label: '小型套餐 (50M)', value: 'small' },
  { label: '专业套餐 (200M)', value: 'pro' },
  { label: '大型套餐 (500M)', value: 'large' },
  { label: '自定义套餐', value: 'custom' },
];

const SubscriptionPlansPage: React.FC = () => {
  const [plans, setPlans] = useState<SubscriptionPlan[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [searchText, setSearchText] = useState('');
  const [isCreateModalVisible, setIsCreateModalVisible] = useState(false);
  const [isEditModalVisible, setIsEditModalVisible] = useState(false);
  const [editingPlan, setEditingPlan] = useState<SubscriptionPlan | null>(null);

  const [form] = Form.useForm();
  const [createLoading, setCreateLoading] = useState(false);
  const [updateLoading, setUpdateLoading] = useState(false);

  // 获取订阅计划列表
  const fetchPlans = useCallback(async () => {
    try {
      setLoading(true);
      const response = await SubscriptionPlanService.getSubscriptionPlans({
        page: currentPage,
        page_size: pageSize,
        search: searchText,
      });
      setPlans(response.data);
      setTotal(response.total);
    } catch (error) {
      console.error('获取订阅计划失败:', error);
      message.error('获取订阅计划列表失败');
    } finally {
      setLoading(false);
    }
  }, [currentPage, pageSize, searchText]);

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

  // 创建订阅计划
  const handleCreate = async (values: CreateSubscriptionPlanRequest) => {
    try {
      setCreateLoading(true);
      await SubscriptionPlanService.createSubscriptionPlan(values);
      message.success('订阅计划创建成功');
      setIsCreateModalVisible(false);
      form.resetFields();
      fetchPlans();
    } catch (error) {
      console.error('创建订阅计划失败:', error);
      message.error('创建订阅计划失败');
    } finally {
      setCreateLoading(false);
    }
  };

  // 更新订阅计划
  const handleUpdate = async (values: UpdateSubscriptionPlanRequest) => {
    if (!editingPlan) return;

    try {
      setUpdateLoading(true);
      await SubscriptionPlanService.updateSubscriptionPlan(editingPlan.id, values);
      message.success('订阅计划更新成功');
      setIsEditModalVisible(false);
      setEditingPlan(null);
      fetchPlans();
    } catch (error) {
      console.error('更新订阅计划失败:', error);
      message.error('更新订阅计划失败');
    } finally {
      setUpdateLoading(false);
    }
  };

  // 删除订阅计划
  const handleDelete = async (id: string, tierName: string) => {
    try {
      await SubscriptionPlanService.deleteSubscriptionPlan(id);
      message.success(`订阅计划 "${tierName}" 删除成功`);
      fetchPlans();
    } catch (error) {
      console.error('删除订阅计划失败:', error);
      message.error('删除订阅计划失败');
    }
  };

  // 打开编辑模态框
  const openEditModal = (plan: SubscriptionPlan) => {
    setEditingPlan(plan);
    form.setFieldsValue({
      tier_name: plan.tier_name,
      pricing_monthly: plan.pricing_monthly,
      pricing_quarterly: plan.pricing_quarterly,
      pricing_yearly: plan.pricing_yearly,
      limits: plan.limits,
      features: plan.features,
      target_users: plan.target_users,
      upgrade_path: plan.upgrade_path,
      is_custom: plan.is_custom,
      default_flow_package: plan.default_flow_package,
      is_active: plan.is_active,
      stripe_price_id_monthly: plan.stripe_price_id_monthly,
      stripe_price_id_quarterly: plan.stripe_price_id_quarterly,
      stripe_price_id_yearly: plan.stripe_price_id_yearly,
    });
    setIsEditModalVisible(true);
  };

  // 刷新数据
  const handleRefresh = () => {
    fetchPlans();
  };

  // 初始加载
  useEffect(() => {
    fetchPlans();
  }, [currentPage, pageSize, searchText, fetchPlans]);

  // 统计信息
  const activeCount = plans.filter(plan => plan.is_active).length;

  const columns = [
    {
      title: '计划名称',
      dataIndex: 'tier_name',
      key: 'tier_name',
      sorter: (a: SubscriptionPlan, b: SubscriptionPlan) => a.tier_name.localeCompare(b.tier_name),
      render: (text: string, record: SubscriptionPlan) => (
        <Space>
          <span style={{ fontWeight: 500 }}>{text}</span>
          {record.is_active && <Tag color="success">活跃</Tag>}
        </Space>
      ),
    },
    {
      title: '月费',
      dataIndex: 'pricing_monthly',
      key: 'pricing_monthly',
      sorter: true,
      render: (value: number) => (
        <Statistic
          value={value}
          prefix="¥"
          precision={2}
          valueStyle={{ fontSize: 14 }}
        />
      ),
    },
    {
      title: '季费',
      dataIndex: 'pricing_quarterly',
      key: 'pricing_quarterly',
      sorter: true,
      render: (value: number) => (
        <Statistic
          value={value}
          prefix="¥"
          precision={2}
          valueStyle={{ fontSize: 14 }}
        />
      ),
    },
    {
      title: '年费',
      dataIndex: 'pricing_yearly',
      key: 'pricing_yearly',
      sorter: true,
      render: (value: number) => (
        <Statistic
          value={value}
          prefix="¥"
          precision={2}
          valueStyle={{ fontSize: 14 }}
        />
      ),
    },
    {
      title: '目标用户',
      dataIndex: 'target_users',
      key: 'target_users',
      render: (text: string) => (
        <Text ellipsis style={{ maxWidth: 150 }}>{text || '-'}</Text>
      ),
    },
    {
      title: '默认流量套餐',
      dataIndex: 'default_flow_package',
      key: 'default_flow_package',
      align: 'center' as const,
      render: (value: string) => {
        const packageInfo = FLOW_PACKAGE_OPTIONS.find(opt => opt.value === value);
        const colorMap: Record<string, string> = {
          small: 'blue',
          pro: 'green',
          large: 'orange',
          custom: 'purple',
        };
        return (
          <Tag color={colorMap[value] || 'default'}>
            {packageInfo?.label || value || '-'}
          </Tag>
        );
      },
    },
    {
      title: '状态',
      dataIndex: 'is_active',
      key: 'is_active',
      align: 'center' as const,
      render: (isActive: boolean) => (
        <Tag color={isActive ? 'success' : 'default'} icon={isActive ? <CheckCircleOutlined /> : <CloseCircleOutlined />}>
          {isActive ? '活跃' : '停用'}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      sorter: true,
      render: (text: string) => (
        <Text>{new Date(text).toLocaleDateString()}</Text>
      ),
    },
    {
      title: '操作',
      key: 'actions',
      align: 'center' as const,
      render: (_: any, record: SubscriptionPlan) => (
        <Space>
          <Tooltip title="查看详情">
            <Button
              type="text"
              icon={<EyeOutlined />}
              size="small"
              onClick={() => {
                Modal.info({
                  title: '订阅计划详情',
                  width: 600,
                  content: (
                    <div>
                      <p><strong>计划名称:</strong> {record.tier_name}</p>
                      <p><strong>月费:</strong> ¥{record.pricing_monthly.toFixed(2)}</p>
                      {record.stripe_price_id_monthly && (
                        <p style={{ paddingLeft: 20, color: '#666' }}>
                          <strong>Stripe 月度价格 ID:</strong> {record.stripe_price_id_monthly}
                        </p>
                      )}
                      <p><strong>季费:</strong> ¥{record.pricing_quarterly.toFixed(2)}</p>
                      {record.stripe_price_id_quarterly && (
                        <p style={{ paddingLeft: 20, color: '#666' }}>
                          <strong>Stripe 季度价格 ID:</strong> {record.stripe_price_id_quarterly}
                        </p>
                      )}
                      <p><strong>年费:</strong> ¥{record.pricing_yearly.toFixed(2)}</p>
                      {record.stripe_price_id_yearly && (
                        <p style={{ paddingLeft: 20, color: '#666' }}>
                          <strong>Stripe 年度价格 ID:</strong> {record.stripe_price_id_yearly}
                        </p>
                      )}
                      <p><strong>默认流量套餐:</strong> {FLOW_PACKAGE_OPTIONS.find(opt => opt.value === record.default_flow_package)?.label || record.default_flow_package || '-'}</p>
                      <p><strong>目标用户:</strong> {record.target_users || '-'}</p>
                      <p><strong>限制:</strong> {record.limits || '-'}</p>
                      <p><strong>功能:</strong> {record.features || '-'}</p>
                    </div>
                  ),
                });
              }}
            />
          </Tooltip>
          <Tooltip title="编辑">
            <Button
              type="text"
              icon={<EditOutlined />}
              size="small"
              onClick={() => openEditModal(record)}
            />
          </Tooltip>
          <Popconfirm
            title={`确定要删除订阅计划 "${record.tier_name}" 吗？`}
            description="删除后不可恢复，请谨慎操作"
            onConfirm={() => handleDelete(record.id, record.tier_name)}
            okText="确定"
            cancelText="取消"
          >
            <Tooltip title="删除">
              <Button
                type="text"
                danger
                icon={<DeleteOutlined />}
                size="small"
              />
            </Tooltip>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div style={{ padding: 24, background: 'var(--theme-bg-primary)', minHeight: '100vh' }}>
      {/* 页面标题和操作栏 */}
      <div style={{ marginBottom: 24, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Title level={2} style={{ margin: 0, color: 'var(--theme-text-primary)' }}>
          订阅计划管理
        </Title>
        <Space>
          <Input.Search
            placeholder="搜索订阅计划..."
            allowClear
            style={{ width: 250 }}
            onSearch={handleSearch}
            onChange={(e) => e.target.value && handleSearch(e.target.value)}
          />
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => {
              form.resetFields();
              setIsCreateModalVisible(true);
            }}
          >
            新增计划
          </Button>
          <Button icon={<ReloadOutlined />} onClick={handleRefresh} loading={loading}>
            刷新
          </Button>
        </Space>
      </div>

      {/* 统计卡片 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={12}>
          <Card>
            <Statistic
              title="总计划数"
              value={total}
              prefix={<DollarOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12}>
          <Card>
            <Statistic
              title="活跃计划"
              value={activeCount}
              prefix={<CheckCircleOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
      </Row>

      {/* 订阅计划列表 */}
      <Card>
        <Table
          columns={columns}
          dataSource={plans}
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
          scroll={{ x: 1200 }}
        />
      </Card>

      {/* 创建订阅计划模态框 */}
      <Modal
        title="新增订阅计划"
        open={isCreateModalVisible}
        onCancel={() => {
          setIsCreateModalVisible(false);
          form.resetFields();
        }}
        footer={null}
        width={800}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleCreate}
          initialValues={{
            pricing_monthly: 0,
            pricing_quarterly: 0,
            pricing_yearly: 0,
            default_flow_package: 'small',
            is_active: true,
            is_custom: false,
          }}
        >
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                label="计划名称"
                name="tier_name"
                rules={[{ required: true, message: '请输入计划名称' }]}
              >
                <Input placeholder="请输入计划名称" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                label="目标用户"
                name="target_users"
              >
                <Input placeholder="请输入目标用户描述" />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={8}>
              <Form.Item
                label="月费"
                name="pricing_monthly"
                rules={[{ required: true, message: '请输入月费' }]}
              >
                <InputNumber
                  style={{ width: '100%' }}
                  min={0}
                  precision={2}
                  prefix="¥"
                  placeholder="0.00"
                />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item
                label="季费"
                name="pricing_quarterly"
              >
                <InputNumber
                  style={{ width: '100%' }}
                  min={0}
                  precision={2}
                  prefix="¥"
                  placeholder="0.00"
                />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item
                label="年费"
                name="pricing_yearly"
              >
                <InputNumber
                  style={{ width: '100%' }}
                  min={0}
                  precision={2}
                  prefix="¥"
                  placeholder="0.00"
                />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={8}>
              <Form.Item
                label="Stripe 月度价格 ID"
                name="stripe_price_id_monthly"
                tooltip="Stripe 月度订阅价格 ID，例如：price_xxxxx"
              >
                <Input placeholder="price_xxxxx" />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item
                label="Stripe 季度价格 ID"
                name="stripe_price_id_quarterly"
                tooltip="Stripe 季度订阅价格 ID，例如：price_xxxxx"
              >
                <Input placeholder="price_xxxxx" />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item
                label="Stripe 年度价格 ID"
                name="stripe_price_id_yearly"
                tooltip="Stripe 年度订阅价格 ID，例如：price_xxxxx"
              >
                <Input placeholder="price_xxxxx" />
              </Form.Item>
            </Col>
          </Row>

          <Form.Item
            label="资源限制 (JSON)"
            name="limits"
            rules={[{ required: true, message: '请输入资源限制配置' }]}
          >
            <Input.TextArea
              rows={3}
              placeholder='{"users": 10, "storage": "10GB", "api_calls": 1000}'
            />
          </Form.Item>

          <Form.Item
            label="功能特性 (JSON)"
            name="features"
          >
            <Input.TextArea
              rows={3}
              placeholder='["基础API调用", "邮件支持", "数据分析"]'
            />
          </Form.Item>

          <Form.Item
            label="升级路径"
            name="upgrade_path"
          >
            <Input placeholder="请输入升级路径说明" />
          </Form.Item>

          <Form.Item
            label="默认流量套餐"
            name="default_flow_package"
            rules={[{ required: true, message: '请选择默认流量套餐' }]}
          >
            <Select
              placeholder="请选择流量套餐类型"
              options={FLOW_PACKAGE_OPTIONS}
            />
          </Form.Item>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                label="是否自定义计划"
                name="is_custom"
                valuePropName="checked"
              >
                <Switch />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                label="是否启用"
                name="is_active"
                valuePropName="checked"
              >
                <Switch defaultChecked />
              </Form.Item>
            </Col>
          </Row>

          <Form.Item>
            <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
              <Button onClick={() => setIsCreateModalVisible(false)}>取消</Button>
              <Button type="primary" htmlType="submit" loading={createLoading}>
                创建
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>

      {/* 编辑订阅计划模态框 */}
      <Modal
        title="编辑订阅计划"
        open={isEditModalVisible}
        onCancel={() => {
          setIsEditModalVisible(false);
          setEditingPlan(null);
          form.resetFields();
        }}
        footer={null}
        width={800}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleUpdate}
        >
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                label="计划名称"
                name="tier_name"
                rules={[{ required: true, message: '请输入计划名称' }]}
              >
                <Input placeholder="请输入计划名称" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                label="目标用户"
                name="target_users"
              >
                <Input placeholder="请输入目标用户描述" />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={8}>
              <Form.Item
                label="月费"
                name="pricing_monthly"
                rules={[{ required: true, message: '请输入月费' }]}
              >
                <InputNumber
                  style={{ width: '100%' }}
                  min={0}
                  precision={2}
                  prefix="¥"
                  placeholder="0.00"
                />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item
                label="季费"
                name="pricing_quarterly"
              >
                <InputNumber
                  style={{ width: '100%' }}
                  min={0}
                  precision={2}
                  prefix="¥"
                  placeholder="0.00"
                />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item
                label="年费"
                name="pricing_yearly"
              >
                <InputNumber
                  style={{ width: '100%' }}
                  min={0}
                  precision={2}
                  prefix="¥"
                  placeholder="0.00"
                />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={8}>
              <Form.Item
                label="Stripe 月度价格 ID"
                name="stripe_price_id_monthly"
                tooltip="Stripe 月度订阅价格 ID，例如：price_xxxxx"
              >
                <Input placeholder="price_xxxxx" />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item
                label="Stripe 季度价格 ID"
                name="stripe_price_id_quarterly"
                tooltip="Stripe 季度订阅价格 ID，例如：price_xxxxx"
              >
                <Input placeholder="price_xxxxx" />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item
                label="Stripe 年度价格 ID"
                name="stripe_price_id_yearly"
                tooltip="Stripe 年度订阅价格 ID，例如：price_xxxxx"
              >
                <Input placeholder="price_xxxxx" />
              </Form.Item>
            </Col>
          </Row>

          <Form.Item
            label="资源限制 (JSON)"
            name="limits"
            rules={[{ required: true, message: '请输入资源限制配置' }]}
          >
            <Input.TextArea
              rows={3}
              placeholder='{"users": 10, "storage": "10GB", "api_calls": 1000}'
            />
          </Form.Item>

          <Form.Item
            label="功能特性 (JSON)"
            name="features"
          >
            <Input.TextArea
              rows={3}
              placeholder='["基础API调用", "邮件支持", "数据分析"]'
            />
          </Form.Item>

          <Form.Item
            label="升级路径"
            name="upgrade_path"
          >
            <Input placeholder="请输入升级路径说明" />
          </Form.Item>

          <Form.Item
            label="默认流量套餐"
            name="default_flow_package"
            rules={[{ required: true, message: '请选择默认流量套餐' }]}
          >
            <Select
              placeholder="请选择流量套餐类型"
              options={FLOW_PACKAGE_OPTIONS}
            />
          </Form.Item>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                label="是否自定义计划"
                name="is_custom"
                valuePropName="checked"
              >
                <Switch />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                label="是否启用"
                name="is_active"
                valuePropName="checked"
              >
                <Switch />
              </Form.Item>
            </Col>
          </Row>

          <Form.Item>
            <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
              <Button onClick={() => setIsEditModalVisible(false)}>取消</Button>
              <Button type="primary" htmlType="submit" loading={updateLoading}>
                更新
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default SubscriptionPlansPage;