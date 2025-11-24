package database

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"sass-monitor/pkg/config"
)

var once sync.Once
var dbManager *DatabaseManager

type DatabaseManager struct {
	Config *config.Config

	// GORM instances
	SaasMonitorDB *gorm.DB
	LightAdminDB  *gorm.DB

	// ClickHouse connections
	ClickHouse map[string]clickhouse.Conn

	// Redis
	RedisClient *redis.Client
}

// GetDatabaseManager 获取数据库管理器单例
func GetDatabaseManager(cfg *config.Config) *DatabaseManager {
	once.Do(func() {
		dbManager = &DatabaseManager{
			Config:     cfg,
			ClickHouse: make(map[string]clickhouse.Conn),
		}
	})
	return dbManager
}

// Initialize 初始化所有数据库连接
func (dm *DatabaseManager) Initialize() error {
	// 初始化PostgreSQL连接
	if err := dm.initPostgreSQL(); err != nil {
		return fmt.Errorf("failed to initialize PostgreSQL: %w", err)
	}

	// 初始化ClickHouse连接
	if err := dm.initClickHouse(); err != nil {
		return fmt.Errorf("failed to initialize ClickHouse: %w", err)
	}

	// 初始化Redis连接
	if err := dm.initRedis(); err != nil {
		return fmt.Errorf("failed to initialize Redis: %w", err)
	}

	log.Println("All database connections initialized successfully")
	return nil
}

// initPostgreSQL 初始化PostgreSQL连接
func (dm *DatabaseManager) initPostgreSQL() error {
	// 配置GORM日志
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// 连接Sass监控数据库
	saasMonitorDB, err := gorm.Open(postgres.Open(dm.Config.Databases.SaasMonitor.GetDSN()), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to saas_monitor database: %w", err)
	}

	// 获取底层数据库连接并配置连接池
	sqlDB, err := saasMonitorDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(dm.Config.Databases.SaasMonitor.MaxOpenConns)
	sqlDB.SetMaxIdleConns(dm.Config.Databases.SaasMonitor.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	dm.SaasMonitorDB = saasMonitorDB

	// 连接Light Admin数据库（只读）
	lightAdminDB, err := gorm.Open(postgres.Open(dm.Config.Databases.LightAdmin.GetDSN()), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to light_admin database: %w", err)
	}

	// 配置Light Admin数据库连接池
	lightAdminSQLDB, err := lightAdminDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get light_admin underlying sql.DB: %w", err)
	}

	lightAdminSQLDB.SetMaxOpenConns(dm.Config.Databases.LightAdmin.MaxOpenConns)
	lightAdminSQLDB.SetMaxIdleConns(dm.Config.Databases.LightAdmin.MaxIdleConns)
	lightAdminSQLDB.SetConnMaxLifetime(time.Hour)

	dm.LightAdminDB = lightAdminDB

	log.Println("PostgreSQL connections established")
	return nil
}

// initClickHouse 初始化ClickHouse连接
func (dm *DatabaseManager) initClickHouse() error {
	for _, chConfig := range dm.Config.ClickHouse {
		// 创建ClickHouse连接选项
		options := clickhouse.Options{
			Addr: []string{fmt.Sprintf("%s:%d", chConfig.Host, chConfig.Port)},
			Auth: clickhouse.Auth{
				Database: chConfig.Database,
				Username: chConfig.User,
				Password: chConfig.Password,
			},
			Settings: clickhouse.Settings{
				"max_execution_time": 60,
			},
			DialTimeout: 30 * time.Second,
		}

		// 建立连接
		conn, err := clickhouse.Open(&options)
		if err != nil {
			return fmt.Errorf("failed to connect to ClickHouse %s: %w", chConfig.Name, err)
		}

		// 验证连接
		if err := conn.Ping(context.Background()); err != nil {
			return fmt.Errorf("failed to ping ClickHouse %s: %w", chConfig.Name, err)
		}

		dm.ClickHouse[chConfig.Name] = conn
		log.Printf("ClickHouse connection established for %s", chConfig.Name)
	}

	return nil
}

// initRedis 初始化Redis连接
func (dm *DatabaseManager) initRedis() error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     dm.Config.Redis.GetRedisAddr(),
		Password: dm.Config.Redis.Password,
		DB:       dm.Config.Redis.Database,
		PoolSize: dm.Config.Redis.PoolSize,
	})

	// 验证连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	dm.RedisClient = rdb
	log.Println("Redis connection established")
	return nil
}

// Close 关闭所有数据库连接
func (dm *DatabaseManager) Close() error {
	var errors []error

	// 关闭PostgreSQL连接
	if dm.SaasMonitorDB != nil {
		if sqlDB, err := dm.SaasMonitorDB.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				errors = append(errors, fmt.Errorf("failed to close saas_monitor db: %w", err))
			}
		}
	}

	if dm.LightAdminDB != nil {
		if sqlDB, err := dm.LightAdminDB.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				errors = append(errors, fmt.Errorf("failed to close light_admin db: %w", err))
			}
		}
	}

	// 关闭ClickHouse连接
	for name, conn := range dm.ClickHouse {
		if err := conn.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close ClickHouse %s: %w", name, err))
		}
	}

	// 关闭Redis连接
	if dm.RedisClient != nil {
		if err := dm.RedisClient.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close Redis: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors occurred while closing databases: %v", errors)
	}

	log.Println("All database connections closed")
	return nil
}

// GetClickHouseConnection 获取指定名称的ClickHouse连接
func (dm *DatabaseManager) GetClickHouseConnection(name string) (clickhouse.Conn, error) {
	conn, exists := dm.ClickHouse[name]
	if !exists {
		return nil, fmt.Errorf("ClickHouse connection '%s' not found", name)
	}
	return conn, nil
}

// HealthCheck 检查所有数据库连接健康状态
func (dm *DatabaseManager) HealthCheck() map[string]error {
	status := make(map[string]error)

	// 检查Sass监控数据库
	if dm.SaasMonitorDB != nil {
		if sqlDB, err := dm.SaasMonitorDB.DB(); err == nil {
			if err := sqlDB.Ping(); err != nil {
				status["saas_monitor"] = err
			} else {
				status["saas_monitor"] = nil
			}
		} else {
			status["saas_monitor"] = err
		}
	}

	// 检查Light Admin数据库
	if dm.LightAdminDB != nil {
		if sqlDB, err := dm.LightAdminDB.DB(); err == nil {
			if err := sqlDB.Ping(); err != nil {
				status["light_admin"] = err
			} else {
				status["light_admin"] = nil
			}
		} else {
			status["light_admin"] = err
		}
	}

	// 检查ClickHouse连接
	for name, conn := range dm.ClickHouse {
		if err := conn.Ping(context.Background()); err != nil {
			status[fmt.Sprintf("clickhouse_%s", name)] = err
		} else {
			status[fmt.Sprintf("clickhouse_%s", name)] = nil
		}
	}

	// 检查Redis连接
	if dm.RedisClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := dm.RedisClient.Ping(ctx).Err(); err != nil {
			status["redis"] = err
		} else {
			status["redis"] = nil
		}
	}

	return status
}