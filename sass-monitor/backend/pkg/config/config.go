package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Databases DatabaseConfig  `mapstructure:"databases"`
	ClickHouse []ClickHouseConfig `mapstructure:"clickhouse"`
	Redis     RedisConfig     `mapstructure:"redis"`
	Monitoring MonitoringConfig `mapstructure:"monitoring"`
	Logging   LoggingConfig   `mapstructure:"logging"`
	CORS      CORSConfig      `mapstructure:"cors"`
}

type ServerConfig struct {
	Port           int    `mapstructure:"port"`
	Mode           string `mapstructure:"mode"`
	JWTSecret      string `mapstructure:"jwt_secret"`
	JWTExpireHours int    `mapstructure:"jwt_expire_hours"`
}

type DatabaseConfig struct {
	SaasMonitor DatabaseConnectionConfig `mapstructure:"saas_monitor"`
	LightAdmin  DatabaseConnectionConfig `mapstructure:"light_admin"`
}

type DatabaseConnectionConfig struct {
	Type          string `mapstructure:"type"`
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	User          string `mapstructure:"user"`
	Password      string `mapstructure:"password"`
	Database      string `mapstructure:"database"`
	SSLMode       string `mapstructure:"ssl_mode"`
	MaxOpenConns  int    `mapstructure:"max_open_conns"`
	MaxIdleConns  int    `mapstructure:"max_idle_conns"`
	ReadOnly      bool   `mapstructure:"readonly"`
}

type ClickHouseConfig struct {
	Name     string `mapstructure:"name"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
	PoolSize int    `mapstructure:"pool_size"`
}

type MonitoringConfig struct {
	CollectInterval int          `mapstructure:"collect_interval"`
	RetentionDays   int          `mapstructure:"retention_days"`
	Alerts          AlertConfig  `mapstructure:"alerts"`
}

type AlertConfig struct {
	Enabled             bool    `mapstructure:"enabled"`
	CPUThreshold        int     `mapstructure:"cpu_threshold"`
	MemoryThreshold     int     `mapstructure:"memory_threshold"`
	DiskThreshold       int     `mapstructure:"disk_threshold"`
	ConnectionThreshold int     `mapstructure:"connection_threshold"`
}

type LoggingConfig struct {
	Level    string `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	Output   string `mapstructure:"output"`
	FilePath string `mapstructure:"file_path"`
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 设置默认值
	setDefaults()

	// 允许环境变量覆盖配置
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// setDefaults 设置默认配置值
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.jwt_expire_hours", 24)

	// Database defaults
	viper.SetDefault("databases.saas_monitor.type", "postgres")
	viper.SetDefault("databases.saas_monitor.ssl_mode", "disable")
	viper.SetDefault("databases.saas_monitor.max_open_conns", 10)
	viper.SetDefault("databases.saas_monitor.max_idle_conns", 5)

	viper.SetDefault("databases.light_admin.type", "postgres")
	viper.SetDefault("databases.light_admin.ssl_mode", "disable")
	viper.SetDefault("databases.light_admin.max_open_conns", 10)
	viper.SetDefault("databases.light_admin.max_idle_conns", 5)
	viper.SetDefault("databases.light_admin.readonly", true)

	// ClickHouse defaults
	viper.SetDefault("clickhouse.0.host", "localhost")
	viper.SetDefault("clickhouse.0.port", 9000)
	viper.SetDefault("clickhouse.0.user", "default")
	viper.SetDefault("clickhouse.0.password", "")

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.database", 0)
	viper.SetDefault("redis.pool_size", 10)

	// Monitoring defaults
	viper.SetDefault("monitoring.collect_interval", 5)
	viper.SetDefault("monitoring.retention_days", 30)
	viper.SetDefault("monitoring.alerts.enabled", true)

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")
}

// GetDSN 获取数据库连接字符串
func (d *DatabaseConnectionConfig) GetDSN() string {
	if d.Type == "postgres" {
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			d.Host, d.Port, d.User, d.Password, d.Database, d.SSLMode)
	}
	return ""
}

// GetClickHouseDSN 获取ClickHouse连接字符串
func (c *ClickHouseConfig) GetClickHouseDSN() string {
	return fmt.Sprintf("tcp://%s:%d?database=%s&username=%s&password=%s",
		c.Host, c.Port, c.Database, c.User, c.Password)
}

// GetRedisAddr 获取Redis地址
func (r *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}