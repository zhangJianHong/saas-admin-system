package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AdminUser 管理员用户模型
type AdminUser struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Username     string     `gorm:"uniqueIndex;not null;size:50" json:"username"`
	PasswordHash string     `gorm:"not null;size:255" json:"-"`
	Email        string     `gorm:"uniqueIndex;size:100" json:"email"`
	FullName     string     `gorm:"size:100" json:"full_name"`
	Role         string     `gorm:"default:'admin';size:20" json:"role"` // admin, super_admin
	Status       string     `gorm:"default:'active';size:20" json:"status"` // active, inactive, locked
	LastLoginAt  *time.Time `json:"last_login_at"`
	LoginAttempts int       `gorm:"default:0" json:"login_attempts"`
	LockedUntil   *time.Time `json:"locked_until"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// MonitoringConfig 监控配置模型
type MonitoringConfig struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ConfigKey   string    `gorm:"uniqueIndex;not null;size:100" json:"config_key"`
	ConfigValue string    `gorm:"type:text" json:"config_value"`
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AlertRule 告警规则模型
type AlertRule struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name            string    `gorm:"not null;size:100" json:"name"`
	Description     string    `gorm:"size:255" json:"description"`
	RuleType        string    `gorm:"not null;size:50" json:"rule_type"` // system, database, organization
	TargetType      string    `gorm:"not null;size:50" json:"target_type"` // postgresql, clickhouse, redis
	TargetName      string    `gorm:"size:100" json:"target_name"` // 具体的数据库名称
	MetricName      string    `gorm:"not null;size:100" json:"metric_name"` // cpu_usage, memory_usage, disk_usage
	Operator        string    `gorm:"not null;size:10" json:"operator"` // >, <, >=, <=, =
	Threshold       float64   `gorm:"not null" json:"threshold"`
	Duration        int       `gorm:"default:5" json:"duration"` // 持续时间(分钟)
	Severity        string    `gorm:"default:'warning';size:20" json:"severity"` // info, warning, critical
	Enabled         bool      `gorm:"default:true" json:"enabled"`
	NotificationConfig string `gorm:"type:jsonb" json:"notification_config"` // 通知配置JSON
	CreatedBy       uuid.UUID `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ResourceMetric 资源指标历史模型
type ResourceMetric struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID *string  `gorm:"size:255;index" json:"organization_id"` // 组织ID，可为空表示系统级指标
	DatabaseType  string    `gorm:"not null;size:20;index" json:"database_type"` // postgresql, clickhouse, redis
	DatabaseName  string    `gorm:"not null;size:100;index" json:"database_name"` // 具体数据库名称
	MetricType    string    `gorm:"not null;size:50;index" json:"metric_type"` // storage, connections, cpu, memory, query_performance
	MetricName    string    `gorm:"not null;size:100;index" json:"metric_name"` // 具体的指标名称
	MetricValue   float64   `gorm:"not null" json:"metric_value"`
	Unit          string    `gorm:"size:20" json:"unit"` // 单位：MB, GB, %, count, ms
	Tags          string    `gorm:"type:jsonb" json:"tags"` // 额外的标签JSON
	CollectedAt   time.Time `gorm:"not null;index" json:"collected_at"`
	CreatedAt     time.Time `json:"created_at"`
}

// MonitoringLog 监控日志模型
type MonitoringLog struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	LogLevel     string    `gorm:"not null;size:20;index" json:"log_level"` // info, warning, error, critical
	Source       string    `gorm:"not null;size:100;index" json:"source"` // 数据来源：database_collector, alert_system
	Component    string    `gorm:"size:100" json:"component"` // 组件名称：postgresql, clickhouse_traces等
	OrganizationID *string  `gorm:"size:255;index" json:"organization_id"`
	Message      string    `gorm:"type:text" json:"message"`
	Details      string    `gorm:"type:jsonb" json:"details"` // 详细信息JSON
	CreatedAt    time.Time `gorm:"index" json:"created_at"`
}

// SystemHealth 系统健康状态模型
type SystemHealth struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ComponentName string    `gorm:"not null;size:100;uniqueIndex" json:"component_name"` // postgresql, clickhouse_traces等
	ComponentType string    `gorm:"not null;size:50" json:"component_type"` // database, cache, message_queue
	Status        string    `gorm:"not null;size:20" json:"status"` // healthy, warning, critical, down
	ResponseTime  *int      `json:"response_time"` // 响应时间(毫秒)
	ErrorMessage  *string   `json:"error_message"`
	LastCheckedAt time.Time `gorm:"not null;index" json:"last_checked_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TableName 指定表名
func (AdminUser) TableName() string {
	return "admin_users"
}

func (MonitoringConfig) TableName() string {
	return "monitoring_configs"
}

func (AlertRule) TableName() string {
	return "alert_rules"
}

func (ResourceMetric) TableName() string {
	return "resource_metrics"
}

func (MonitoringLog) TableName() string {
	return "monitoring_logs"
}

func (SystemHealth) TableName() string {
	return "system_health"
}