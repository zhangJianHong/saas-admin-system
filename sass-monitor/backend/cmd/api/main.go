package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sass-monitor/internal/database"
	"sass-monitor/internal/handlers"
	"sass-monitor/internal/middleware"
	"sass-monitor/internal/models"
	"sass-monitor/internal/services"
	"sass-monitor/pkg/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// @title Sass后台监控系统API
// @version 1.0
// @description 企业级多数据库监控后台系统API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库连接
	dbManager := database.GetDatabaseManager(cfg)
	if err := dbManager.Initialize(); err != nil {
		log.Fatalf("Failed to initialize database connections: %v", err)
	}

	// 自动迁移数据库表结构
	if err := autoMigrate(dbManager); err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}

	// 初始化任务调度器
	scheduler := services.NewTaskScheduler(dbManager, cfg)
	if err := scheduler.Start(); err != nil {
		log.Printf("Failed to start scheduler: %v", err)
		// 调度器启动失败不应该阻止服务器启动
	}

	// 确保在程序退出时关闭资源
	defer func() {
		log.Println("Stopping scheduler...")
		scheduler.Stop()

		if err := dbManager.Close(); err != nil {
			log.Printf("Error closing database connections: %v", err)
		}
	}()

	// 创建Gin引擎
	router := gin.New()

	// 添加中间件
	setupMiddleware(router, cfg)

	// 设置路由
	setupRoutes(router, dbManager, cfg)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// 启动服务器
	go func() {
		log.Printf("Starting server on port %d", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 设置关闭超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// setupMiddleware 设置中间件
func setupMiddleware(router *gin.Engine, cfg *config.Config) {
	// 恢复中间件
	router.Use(gin.Recovery())

	// 日志中间件
	router.Use(gin.Logger())

	// CORS中间件
	corsConfig := cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     cfg.CORS.AllowedMethods,
		AllowHeaders:     cfg.CORS.AllowedHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))

	// 限流中间件
	router.Use(middleware.RateLimitMiddleware())

	// 安全中间件
	router.Use(middleware.SecurityMiddleware())
}

// setupRoutes 设置路由
func setupRoutes(router *gin.Engine, dbManager *database.DatabaseManager, cfg *config.Config) {
	// 健康检查端点
	router.GET("/health", func(c *gin.Context) {
		healthStatus := dbManager.HealthCheck()

		allHealthy := true
		for _, status := range healthStatus {
			if status != nil {
				allHealthy = false
				break
			}
		}

		if allHealthy {
			c.JSON(http.StatusOK, gin.H{
				"status":    "healthy",
				"timestamp": time.Now(),
				"version":   "1.0.0",
			})
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":    "unhealthy",
				"timestamp": time.Now(),
				"details":   healthStatus,
			})
		}
	})

	// API版本分组
	v1 := router.Group("/api/v1")
	{
		// 创建处理器
		authHandler := handlers.NewAuthHandler(dbManager.SaasMonitorDB, cfg)
		dashboardHandler := handlers.NewDashboardHandler(dbManager, cfg)
		monitoringHandler := handlers.NewMonitoringHandler(dbManager, cfg)
		organizationService := services.NewOrganizationService(dbManager)
		organizationHandler := handlers.NewOrganizationHandler(organizationService, cfg)
		subscriptionPlanService := services.NewSubscriptionPlanService(dbManager)
		subscriptionPlanHandler := handlers.NewSubscriptionPlanHandler(subscriptionPlanService, cfg)
		userService := services.NewUserService(dbManager)
		userHandler := handlers.NewUserHandler(userService)

		// 认证路由（无需JWT）
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/logout", authHandler.Logout)
			authGroup.GET("/refresh", authHandler.RefreshToken)
		}

		// 需要认证的路由
		protectedGroup := v1.Group("")
		protectedGroup.Use(middleware.JWTMiddleware(cfg.Server.JWTSecret))
		{
			// 用户管理
			protectedGroup.GET("/profile", authHandler.GetProfile)
			protectedGroup.PUT("/profile", authHandler.UpdateProfile)
			protectedGroup.POST("/change-password", authHandler.ChangePassword)

			// 仪表板
			dashboardGroup := protectedGroup.Group("/dashboard")
			{
				dashboardGroup.GET("/overview", dashboardHandler.GetOverview)
				dashboardGroup.GET("/organizations", dashboardHandler.GetOrganizations)
				dashboardGroup.GET("/organizations/:id/metrics", dashboardHandler.GetOrganizationMetrics)
				dashboardGroup.GET("/database-status", dashboardHandler.GetDatabaseStatus)
			}

			// 监控数据
			monitoringGroup := protectedGroup.Group("/monitoring")
			{
				monitoringGroup.GET("/metrics", monitoringHandler.GetMetrics)
				monitoringGroup.GET("/metrics/history", monitoringHandler.GetMetricsHistory)
				monitoringGroup.GET("/organizations", monitoringHandler.GetOrganizations)
				monitoringGroup.GET("/organizations/overview", monitoringHandler.GetOrganizationOverview)
				monitoringGroup.GET("/organizations/:id/usage", monitoringHandler.GetOrganizationUsage)
				monitoringGroup.GET("/databases", monitoringHandler.GetDatabaseInfo)
				monitoringGroup.GET("/alerts", monitoringHandler.GetAlerts)
				monitoringGroup.POST("/alerts", monitoringHandler.CreateAlert)
				monitoringGroup.PUT("/alerts/:id", monitoringHandler.UpdateAlert)
				monitoringGroup.DELETE("/alerts/:id", monitoringHandler.DeleteAlert)
			}

			// 组织管理（只读模式）
			organizationGroup := protectedGroup.Group("/organizations")
			{
				// 组织查询（只读）
				organizationGroup.GET("", organizationHandler.GetOrganizations)
				organizationGroup.GET("/:id", organizationHandler.GetOrganizationByID)
				organizationGroup.GET("/:id/metrics", organizationHandler.GetOrganizationMetrics)

				// 组织用户查询（只读）
				organizationGroup.GET("/:id/users", organizationHandler.GetOrganizationUsers)

				// 组织工作空间查询（只读）
				organizationGroup.GET("/:id/workspaces", organizationHandler.GetOrganizationWorkspaces)

				// 组织订阅查询（只读）
				organizationGroup.GET("/:id/subscriptions", organizationHandler.GetOrganizationSubscriptions)

				// 发送订阅到期提醒邮件（预留功能）
				organizationGroup.POST("/:id/send-expiry-reminder", organizationHandler.SendExpiryReminder)

				// 不支持的操作（返回405 Method Not Allowed）
				organizationGroup.POST("", organizationHandler.CreateOrganization)
				organizationGroup.PUT("/:id", organizationHandler.UpdateOrganization)
				organizationGroup.DELETE("/:id", organizationHandler.DeleteOrganization)
			}

			// 工作空间管理
			workspaceGroup := protectedGroup.Group("/workspaces")
			{
				workspaceGroup.GET("/:workspaceId/users", organizationHandler.GetWorkspaceUsers)
			}

			// 订阅计划管理（完整CRUD）
			subscriptionGroup := protectedGroup.Group("/subscription-plans")
			{
				subscriptionGroup.GET("", subscriptionPlanHandler.GetSubscriptionPlans)
				subscriptionGroup.POST("", subscriptionPlanHandler.CreateSubscriptionPlan)
				subscriptionGroup.GET("/:id", subscriptionPlanHandler.GetSubscriptionPlanByID)
				subscriptionGroup.PUT("/:id", subscriptionPlanHandler.UpdateSubscriptionPlan)
				subscriptionGroup.DELETE("/:id", subscriptionPlanHandler.DeleteSubscriptionPlan)
				subscriptionGroup.GET("/active", subscriptionPlanHandler.GetActiveSubscriptionPlans)
				subscriptionGroup.GET("/search", subscriptionPlanHandler.GetSubscriptionPlansByPricingRange)
			}

			// 系统管理
			systemGroup := protectedGroup.Group("/system")
			{
				systemGroup.GET("/health", monitoringHandler.GetSystemHealth)
				systemGroup.GET("/logs", monitoringHandler.GetSystemLogs)
				systemGroup.GET("/configs", monitoringHandler.GetSystemConfigs)
				systemGroup.PUT("/configs", monitoringHandler.UpdateSystemConfigs)
			}

			// 用户管理（只读模式）
			userGroup := protectedGroup.Group("/users")
			{
				// 用户查询（只读）
				userGroup.GET("", userHandler.GetUsers)
				userGroup.GET("/:id", userHandler.GetUserByID)

				// 用户关联信息查询（只读）
				userGroup.GET("/:id/organizations", userHandler.GetUserOrganizations)
				userGroup.GET("/:id/workspaces", userHandler.GetUserWorkspaces)
				userGroup.GET("/:id/subscriptions", userHandler.GetUserSubscriptions)

				// 不支持的操作（返回405 Method Not Allowed）
				userGroup.POST("", userHandler.CreateUser)
				userGroup.PUT("/:id", userHandler.UpdateUser)
				userGroup.DELETE("/:id", userHandler.DeleteUser)
			}
		}
	}

	// Swagger文档
	router.Static("/swagger", "./docs/swagger")
}

// autoMigrate 自动迁移数据库表结构
func autoMigrate(dbManager *database.DatabaseManager) error {
	log.Println("Starting database migration...")

	// 检查数据库连接是否为空
	if dbManager == nil {
		return fmt.Errorf("database manager is nil")
	}

	if dbManager.SaasMonitorDB == nil {
		return fmt.Errorf("saas monitor database connection is nil")
	}

	log.Println("Database connections validated, starting migration...")

	// 迁移Sass监控数据库表
	if err := dbManager.SaasMonitorDB.AutoMigrate(
		&models.AdminUser{},
		&models.MonitoringConfig{},
		&models.AlertRule{},
		&models.ResourceMetric{},
		&models.MonitoringLog{},
		&models.SystemHealth{},
	); err != nil {
		return fmt.Errorf("failed to migrate saas_monitor database: %w", err)
	}

	log.Println("Sass monitor database migration completed")

	// 执行初始化SQL脚本
	if err := executeInitSQL(dbManager); err != nil {
		log.Printf("Warning: Failed to execute init SQL: %v", err)
	}

	return nil
}

// executeInitSQL 执行初始化SQL脚本
func executeInitSQL(dbManager *database.DatabaseManager) error {
	// 读取并执行初始化脚本
	initSQL := `
-- 插入默认监控配置（如果不存在）
INSERT INTO monitoring_configs (config_key, config_value, description) VALUES
('collect_interval', '5', '数据采集间隔（分钟）'),
('retention_days', '30', '数据保留天数'),
('alert_enabled', 'true', '是否启用告警'),
('cpu_threshold', '80', 'CPU使用率告警阈值'),
('memory_threshold', '85', '内存使用率告警阈值'),
('disk_threshold', '90', '磁盘使用率告警阈值'),
('connection_threshold', '100', '数据库连接数告警阈值')
ON CONFLICT (config_key) DO NOTHING;

-- 插入默认管理员用户（如果不存在）
INSERT INTO admin_users (id, username, password_hash, email, full_name, role, status, created_at, updated_at) VALUES
('00000000-0000-0000-0000-000000000001', 'admin', '$2a$10$7bmQZ.QALZRfjzf82Y1NfeGoW.8ojHG6R6G3g37F6bprHNwnKEDGq', 'admin@sass-monitor.com', '系统管理员', 'super_admin', 'active', NOW(), NOW())
ON CONFLICT (username) DO NOTHING;

-- 插入默认告警规则（如果不存在）
INSERT INTO alert_rules (name, description, rule_type, target_type, target_name, metric_name, operator, threshold, severity, enabled, created_by) VALUES
('PostgreSQL连接数告警', '当PostgreSQL连接数超过阈值时触发告警', 'database', 'postgresql', 'light_admin', 'active_connections', '>', 80, 'warning', true, '00000000-0000-0000-0000-000000000001'),
('ClickHouse存储告警', '当ClickHouse数据库存储使用率超过阈值时触发', 'database', 'clickhouse', 'traces', 'database_size_mb', '>', 1024, 'warning', true, '00000000-0000-0000-0000-000000000001'),
('Redis内存告警', '当Redis内存使用率超过阈值时触发', 'database', 'redis', 'default', 'used_memory_bytes', '>', 1073741824, 'critical', true, '00000000-0000-0000-0000-000000000001')
ON CONFLICT DO NOTHING;
`

	// 在事务中执行
	return dbManager.SaasMonitorDB.Exec(initSQL).Error
}
