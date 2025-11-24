import React, { useState, useEffect } from 'react';
import { Routes, Route, Navigate, Link } from 'react-router-dom';
import {
  Typography,
  Card,
  Row,
  Col,
  Table,
  Tag,
  Button,
  Space,
  Alert,
  Spin,
} from 'antd';
import {
  DatabaseOutlined,
  LineChartOutlined,
  BellOutlined,
  ArrowRightOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import MonitoringService from '../services/monitoring';
import { ResourceMetric, AlertRule } from '../types';

const { Title, Text } = Typography;

const Monitoring: React.FC = () => {
  const [metrics, setMetrics] = useState<ResourceMetric[]>([]);
  const [alerts, setAlerts] = useState<AlertRule[]>([]);
  const [metricsLoading, setMetricsLoading] = useState(false);
  const [alertsLoading, setAlertsLoading] = useState(false);

  // 获取监控指标
  const fetchMetrics = async () => {
    try {
      setMetricsLoading(true);
      const data = await MonitoringService.getMetrics({});
      setMetrics(Array.isArray(data.metrics) ? data.metrics : []);
    } catch (error: any) {
      console.error('Fetch metrics error:', error);
    } finally {
      setMetricsLoading(false);
    }
  };

  // 获取告警规则
  const fetchAlerts = async () => {
    try {
      setAlertsLoading(true);
      const data = await MonitoringService.getAlerts({});
      setAlerts(Array.isArray(data.alerts) ? data.alerts : []);
    } catch (error: any) {
      console.error('Fetch alerts error:', error);
    } finally {
      setAlertsLoading(false);
    }
  };

  useEffect(() => {
    fetchMetrics();
    fetchAlerts();
  }, []);

  // 指标表格列定义
  const metricsColumns = [
    {
      title: '数据库类型',
      dataIndex: 'database_type',
      key: 'database_type',
      render: (type: string) => (
        <Tag color={type === 'postgresql' ? 'blue' : type === 'clickhouse' ? 'green' : 'red'}>
          {type.toUpperCase()}
        </Tag>
      ),
    },
    {
      title: '数据库名称',
      dataIndex: 'database_name',
      key: 'database_name',
      render: (name: string) => (
        <Space>
          <DatabaseOutlined />
          {name}
        </Space>
      ),
    },
    {
      title: '指标类型',
      dataIndex: 'metric_type',
      key: 'metric_type',
    },
    {
      title: '指标名称',
      dataIndex: 'metric_name',
      key: 'metric_name',
    },
    {
      title: '当前值',
      dataIndex: 'metric_value',
      key: 'metric_value',
      render: (value: number, record: ResourceMetric) => (
        <Text strong style={{ fontSize: 16 }}>
          {value} {record.unit || ''}
        </Text>
      ),
    },
    {
      title: '采集时间',
      dataIndex: 'collected_at',
      key: 'collected_at',
      render: (time: string) => new Date(time).toLocaleString(),
    },
  ];

  // 告警表格列定义
  const alertsColumns = [
    {
      title: '规则名称',
      dataIndex: 'name',
      key: 'name',
      render: (name: string) => (
        <Space>
          <BellOutlined />
          <Text strong>{name}</Text>
        </Space>
      ),
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: '目标类型',
      dataIndex: 'target_type',
      key: 'target_type',
      render: (type: string) => (
        <Tag color={type === 'postgresql' ? 'blue' : type === 'clickhouse' ? 'green' : 'red'}>
          {type.toUpperCase()}
        </Tag>
      ),
    },
    {
      title: '指标',
      dataIndex: 'metric_name',
      key: 'metric_name',
    },
    {
      title: '阈值',
      dataIndex: 'threshold',
      key: 'threshold',
      render: (threshold: number, record: AlertRule) => (
        <Text>{record.operator} {threshold}</Text>
      ),
    },
    {
      title: '严重级别',
      dataIndex: 'severity',
      key: 'severity',
      render: (severity: string) => {
        const color = severity === 'critical' ? 'red' : severity === 'warning' ? 'orange' : 'blue';
        return <Tag color={color}>{severity.toUpperCase()}</Tag>;
      },
    },
    {
      title: '状态',
      dataIndex: 'enabled',
      key: 'enabled',
      render: (enabled: boolean) => (
        <Tag color={enabled ? 'green' : 'red'}>
          {enabled ? '启用' : '禁用'}
        </Tag>
      ),
    },
  ];

  return (
    <div style={{ padding: 24 }}>
      {/* 页面标题 */}
      <div style={{ marginBottom: 24 }}>
        <Title level={2} style={{ margin: 0 }}>监控数据</Title>
        <Text type="secondary">实时监控数据库性能和系统告警</Text>
      </div>

      {/* 导航卡片 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={12}>
          <Card
            hoverable
            bodyStyle={{ padding: 24 }}
            onClick={() => window.location.href = '/monitoring/metrics'}
          >
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
              <div>
                <div style={{ display: 'flex', alignItems: 'center', marginBottom: 8 }}>
                  <LineChartOutlined style={{ fontSize: 24, color: '#1890ff', marginRight: 12 }} />
                  <Title level={4} style={{ margin: 0 }}>实时指标</Title>
                </div>
                <Text type="secondary">查看数据库性能指标</Text>
              </div>
              <ArrowRightOutlined style={{ fontSize: 16, color: '#1890ff' }} />
            </div>
          </Card>
        </Col>
        <Col xs={24} sm={12}>
          <Card
            hoverable
            bodyStyle={{ padding: 24 }}
            onClick={() => window.location.href = '/monitoring/alerts'}
          >
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
              <div>
                <div style={{ display: 'flex', alignItems: 'center', marginBottom: 8 }}>
                  <BellOutlined style={{ fontSize: 24, color: '#fa8c16', marginRight: 12 }} />
                  <Title level={4} style={{ margin: 0 }}>告警管理</Title>
                </div>
                <Text type="secondary">配置和管理告警规则</Text>
              </div>
              <ArrowRightOutlined style={{ fontSize: 16, color: '#fa8c16' }} />
            </div>
          </Card>
        </Col>
      </Row>

      {/* 数据概览 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col span={12}>
          <Card
            title="最新指标"
            extra={
              <Button
                icon={<ReloadOutlined />}
                onClick={fetchMetrics}
                loading={metricsLoading}
              >
                刷新
              </Button>
            }
          >
            {metricsLoading ? (
              <div style={{ textAlign: 'center', padding: 20 }}>
                <Spin />
              </div>
            ) : metrics.length > 0 ? (
              <Table
                columns={metricsColumns}
                dataSource={metrics.slice(0, 5)}
                rowKey="id"
                pagination={false}
                size="small"
              />
            ) : (
              <Alert message="暂无指标数据" type="info" showIcon />
            )}
            {metrics.length > 5 && (
              <div style={{ textAlign: 'right', marginTop: 12 }}>
                <Link to="/monitoring/metrics">
                  <Button type="link">查看全部指标 →</Button>
                </Link>
              </div>
            )}
          </Card>
        </Col>
        <Col span={12}>
          <Card
            title="告警规则"
            extra={
              <Button
                icon={<ReloadOutlined />}
                onClick={fetchAlerts}
                loading={alertsLoading}
              >
                刷新
              </Button>
            }
          >
            {alertsLoading ? (
              <div style={{ textAlign: 'center', padding: 20 }}>
                <Spin />
              </div>
            ) : alerts.length > 0 ? (
              <Table
                columns={alertsColumns}
                dataSource={alerts.slice(0, 5)}
                rowKey="id"
                pagination={false}
                size="small"
              />
            ) : (
              <Alert message="暂无告警规则" type="info" showIcon />
            )}
            {alerts.length > 5 && (
              <div style={{ textAlign: 'right', marginTop: 12 }}>
                <Link to="/monitoring/alerts">
                  <Button type="link">查看全部告警 →</Button>
                </Link>
              </div>
            )}
          </Card>
        </Col>
      </Row>

      {/* 子路由 */}
      <Routes>
        <Route path="/" element={<Navigate to="/monitoring/metrics" replace />} />
        <Route
          path="/metrics"
          element={
            <Card title="实时指标">
              <Table
                columns={metricsColumns}
                dataSource={metrics}
                rowKey="id"
                loading={metricsLoading}
                pagination={{
                  pageSize: 10,
                  showSizeChanger: true,
                  showQuickJumper: true,
                  showTotal: (total, range) =>
                    `第 ${range[0]}-${range[1]} 条，共 ${total} 条记录`,
                }}
              />
            </Card>
          }
        />
        <Route
          path="/alerts"
          element={
            <Card title="告警规则">
              <Table
                columns={alertsColumns}
                dataSource={alerts}
                rowKey="id"
                loading={alertsLoading}
                pagination={{
                  pageSize: 10,
                  showSizeChanger: true,
                  showQuickJumper: true,
                  showTotal: (total, range) =>
                    `第 ${range[0]}-${range[1]} 条，共 ${total} 条记录`,
                }}
              />
            </Card>
          }
        />
      </Routes>
    </div>
  );
};

export default Monitoring;