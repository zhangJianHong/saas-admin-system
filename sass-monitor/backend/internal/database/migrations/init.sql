-- Sass监控系统数据库初始化脚本
-- 创建saas_monitor数据库表结构

-- 启用UUID扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 管理员用户表
CREATE TABLE IF NOT EXISTS admin_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE,
    full_name VARCHAR(100),
    role VARCHAR(20) DEFAULT 'admin',
    status VARCHAR(20) DEFAULT 'active',
    last_login_at TIMESTAMP,
    login_attempts INTEGER DEFAULT 0,
    locked_until TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_admin_users_username ON admin_users(username);
CREATE INDEX IF NOT EXISTS idx_admin_users_email ON admin_users(email);
CREATE INDEX IF NOT EXISTS idx_admin_users_status ON admin_users(status);
CREATE INDEX IF NOT EXISTS idx_admin_users_deleted_at ON admin_users(deleted_at);

-- 监控配置表
CREATE TABLE IF NOT EXISTS monitoring_configs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    config_key VARCHAR(100) UNIQUE NOT NULL,
    config_value TEXT,
    description VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_monitoring_configs_key ON monitoring_configs(config_key);

-- 告警规则表
CREATE TABLE IF NOT EXISTS alert_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description VARCHAR(255),
    rule_type VARCHAR(50) NOT NULL,
    target_type VARCHAR(50) NOT NULL,
    target_name VARCHAR(100),
    metric_name VARCHAR(100) NOT NULL,
    operator VARCHAR(10) NOT NULL,
    threshold DOUBLE PRECISION NOT NULL,
    duration INTEGER DEFAULT 5,
    severity VARCHAR(20) DEFAULT 'warning',
    enabled BOOLEAN DEFAULT true,
    notification_config JSONB,
    created_by UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_alert_rules_enabled ON alert_rules(enabled);
CREATE INDEX IF NOT EXISTS idx_alert_rules_type ON alert_rules(rule_type, target_type);
CREATE INDEX IF NOT EXISTS idx_alert_rules_created_by ON alert_rules(created_by);

-- 资源指标历史表
CREATE TABLE IF NOT EXISTS resource_metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id VARCHAR(255),
    database_type VARCHAR(20) NOT NULL,
    database_name VARCHAR(100) NOT NULL,
    metric_type VARCHAR(50) NOT NULL,
    metric_name VARCHAR(100) NOT NULL,
    metric_value DOUBLE PRECISION NOT NULL,
    unit VARCHAR(20),
    tags JSONB,
    collected_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_resource_metrics_org_db ON resource_metrics(organization_id, database_type);
CREATE INDEX IF NOT EXISTS idx_resource_metrics_type_name ON resource_metrics(metric_type, metric_name);
CREATE INDEX IF NOT EXISTS idx_resource_metrics_collected_at ON resource_metrics(collected_at);
CREATE INDEX IF NOT EXISTS idx_resource_metrics_composite ON resource_metrics(database_type, database_name, metric_type, collected_at);

-- 监控日志表
CREATE TABLE IF NOT EXISTS monitoring_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    log_level VARCHAR(20) NOT NULL,
    source VARCHAR(100) NOT NULL,
    component VARCHAR(100),
    organization_id VARCHAR(255),
    message TEXT,
    details JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_monitoring_logs_level ON monitoring_logs(log_level);
CREATE INDEX IF NOT EXISTS idx_monitoring_logs_source ON monitoring_logs(source);
CREATE INDEX IF NOT EXISTS idx_monitoring_logs_created_at ON monitoring_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_monitoring_logs_org ON monitoring_logs(organization_id);

-- 系统健康状态表
CREATE TABLE IF NOT EXISTS system_health (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    component_name VARCHAR(100) UNIQUE NOT NULL,
    component_type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    response_time INTEGER,
    error_message TEXT,
    last_checked_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_system_health_name ON system_health(component_name);
CREATE INDEX IF NOT EXISTS idx_system_health_type ON system_health(component_type);
CREATE INDEX IF NOT EXISTS idx_system_health_status ON system_health(status);
CREATE INDEX IF NOT EXISTS idx_system_health_checked_at ON system_health(last_checked_at);

-- 创建更新时间触发器函数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 为需要的表创建更新时间触发器
CREATE TRIGGER update_admin_users_updated_at BEFORE UPDATE ON admin_users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_monitoring_configs_updated_at BEFORE UPDATE ON monitoring_configs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_alert_rules_updated_at BEFORE UPDATE ON alert_rules
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_system_health_updated_at BEFORE UPDATE ON system_health
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 插入默认数据

-- 默认监控配置
INSERT INTO monitoring_configs (config_key, config_value, description) VALUES
('collect_interval', '5', '数据采集间隔（分钟）'),
('retention_days', '30', '数据保留天数'),
('alert_enabled', 'true', '是否启用告警'),
('cpu_threshold', '80', 'CPU使用率告警阈值'),
('memory_threshold', '85', '内存使用率告警阈值'),
('disk_threshold', '90', '磁盘使用率告警阈值'),
('connection_threshold', '100', '数据库连接数告警阈值')
ON CONFLICT (config_key) DO NOTHING;

-- 插入默认管理员用户（密码: admin123）
INSERT INTO admin_users (username, password_hash, email, full_name, role) VALUES
('admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin@sass-monitor.com', '系统管理员', 'super_admin')
ON CONFLICT (username) DO NOTHING;

-- 创建示例告警规则
INSERT INTO alert_rules (name, description, rule_type, target_type, target_name, metric_name, operator, threshold, severity, created_by) VALUES
('PostgreSQL连接数告警', '当PostgreSQL连接数超过阈值时触发告警', 'database', 'postgresql', 'light_admin', 'connection_count', '>', 80, 'warning', (SELECT id FROM admin_users WHERE username = 'admin' LIMIT 1)),
('ClickHouse存储告警', '当ClickHouse数据库存储使用率超过阈值时触发', 'database', 'clickhouse', 'traces', 'storage_usage_percent', '>', 85, 'warning', (SELECT id FROM admin_users WHERE username = 'admin' LIMIT 1)),
('Redis内存告警', '当Redis内存使用率超过阈值时触发', 'database', 'redis', 'default', 'memory_usage_percent', '>', 90, 'critical', (SELECT id FROM admin_users WHERE username = 'admin' LIMIT 1))
ON CONFLICT DO NOTHING;