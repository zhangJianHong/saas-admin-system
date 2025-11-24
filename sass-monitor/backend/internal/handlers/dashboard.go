package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"sass-monitor/internal/database"
	"sass-monitor/internal/models"
	"sass-monitor/pkg/config"
)

type DashboardHandler struct {
	dbManager *database.DatabaseManager
	config    *config.Config
}

func NewDashboardHandler(dbManager *database.DatabaseManager, cfg *config.Config) *DashboardHandler {
	return &DashboardHandler{
		dbManager: dbManager,
		config:    cfg,
	}
}

// OverviewResponse 概览响应结构
type OverviewResponse struct {
	TotalOrganizations int64             `json:"total_organizations"`
	TotalUsers         int64             `json:"total_users"`
	TotalSubscriptions int64             `json:"total_subscriptions"`
	DatabaseStatus     map[string]string `json:"database_status"`
	SystemHealth       map[string]string `json:"system_health"`
	RecentMetrics      []RecentMetric    `json:"recent_metrics"`
}

// RecentMetric 最近指标结构
type RecentMetric struct {
	DatabaseName string  `json:"database_name"`
	MetricType   string  `json:"metric_type"`
	MetricValue  float64 `json:"metric_value"`
	Unit         string  `json:"unit"`
	Timestamp    int64   `json:"timestamp"`
}

// OrganizationInfo 组织信息结构
type OrganizationInfo struct {
	ID                      string  `json:"id" gorm:"column:id"`
	Name                    string  `json:"name" gorm:"column:name"`
	OwnerID                 string  `json:"owner_id" gorm:"column:owner_id"`
	Description             *string `json:"description" gorm:"column:description"`
	UserCount               int64   `json:"user_count" gorm:"column:user_count"`
	SubscriptionCount       int64   `json:"subscription_count" gorm:"column:subscription_count"`
	ActiveSubscriptionCount int64   `json:"active_subscription_count" gorm:"column:active_subscription_count"`
	StorageUsage            float64 `json:"storage_usage" gorm:"column:storage_usage"`
	CreatedAt               string  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt               *string `json:"updated_at" gorm:"column:updated_at"`
	// 订阅到期相关字段
	SubscriptionStatus  string  `json:"subscription_status" gorm:"-"`
	SubscriptionEndDate *string `json:"subscription_end_date" gorm:"column:subscription_end_date"`
	DaysUntilExpiration *int    `json:"days_until_expiration" gorm:"-"`
}

// DatabaseStatusResponse 数据库状态响应
type DatabaseStatusResponse struct {
	PostgreSQL DatabaseInfo                 `json:"postgresql"`
	ClickHouse map[string]ClickHouseInfo    `json:"clickhouse"`
	Redis      RedisInfo                    `json:"redis"`
}

// DatabaseInfo 数据库基础信息
type DatabaseInfo struct {
	Status         string  `json:"status"`
	Connections    int     `json:"connections"`
	MaxConnections int     `json:"max_connections"`
	DatabaseSize   float64 `json:"database_size"`
	ResponseTime   int     `json:"response_time"`
}

// ClickHouseInfo ClickHouse信息
type ClickHouseInfo struct {
	Status       string  `json:"status"`
	DatabaseSize float64 `json:"database_size"`
	TableCount   int     `json:"table_count"`
	RowCount     int64   `json:"row_count"`
	ResponseTime int     `json:"response_time"`
}

// RedisInfo Redis信息
type RedisInfo struct {
	Status        string  `json:"status"`
	UsedMemory    float64 `json:"used_memory"`
	MaxMemory     float64 `json:"max_memory"`
	ConnectedClients int  `json:"connected_clients"`
	ResponseTime  int     `json:"response_time"`
}

// GetOverview 获取仪表板概览
func (h *DashboardHandler) GetOverview(c *gin.Context) {
	overview := OverviewResponse{
		DatabaseStatus: make(map[string]string),
		SystemHealth:   make(map[string]string),
	}

	// 获取组织总数
	var totalOrgs int64
	if err := h.dbManager.LightAdminDB.Table("auth_organizations").Count(&totalOrgs).Error; err != nil {
		overview.DatabaseStatus["light_admin"] = "error"
	} else {
		overview.DatabaseStatus["light_admin"] = "connected"
		overview.TotalOrganizations = totalOrgs
	}

	// 获取用户总数
	var totalUsers int64
	if err := h.dbManager.LightAdminDB.Table("auth_users").Count(&totalUsers).Error; err != nil {
		overview.DatabaseStatus["auth_users"] = "error"
	} else {
		overview.TotalUsers = totalUsers
	}

	// 获取订阅总数
	var totalSubs int64
	if err := h.dbManager.LightAdminDB.Table("subscription_users").Where("status = ?", "active").Count(&totalSubs).Error; err != nil {
		overview.DatabaseStatus["subscription_users"] = "error"
	} else {
		overview.TotalSubscriptions = totalSubs
	}

	// 检查系统健康状态
	healthStatus := h.dbManager.HealthCheck()
	for component, status := range healthStatus {
		if status == nil {
			overview.SystemHealth[component] = "healthy"
		} else {
			overview.SystemHealth[component] = "unhealthy"
		}
	}

	// 获取最近指标（示例数据）
	overview.RecentMetrics = []RecentMetric{
		{
			DatabaseName: "light_admin",
			MetricType:   "storage_usage",
			MetricValue:  1024.5,
			Unit:         "MB",
			Timestamp:    1640995200,
		},
		{
			DatabaseName: "traces",
			MetricType:   "table_count",
			MetricValue:  25,
			Unit:         "count",
			Timestamp:    1640995200,
		},
	}

	c.JSON(http.StatusOK, overview)
}

// GetOrganizations 获取组织列表
func (h *DashboardHandler) GetOrganizations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "100"))
	search := c.Query("search")

	offset := (page - 1) * pageSize

	query := h.dbManager.LightAdminDB.Table("auth_organizations").
		Select(`
			auth_organizations.id,
			auth_organizations.name,
			auth_organizations.owner_id,
			auth_organizations.description,
			auth_organizations.created_at::text as created_at,
			auth_organizations.updated_at::text as updated_at,
			COUNT(DISTINCT au.id) as user_count,
			COUNT(DISTINCT su.id) as subscription_count,
			COUNT(DISTINCT CASE WHEN su.status = 'active' THEN su.id END) as active_subscription_count,
			0 as storage_usage,
			MIN(CASE WHEN su.status IN ('active','trial') THEN su.end_date::text END) as subscription_end_date
		`).
		Joins("LEFT JOIN auth_user_organization auo ON auth_organizations.id = auo.organization_id").
		Joins("LEFT JOIN auth_users au ON auo.user_id = au.id").
		Joins("LEFT JOIN subscription_users su ON su.organization_id = auth_organizations.id::text")

	if search != "" {
		query = query.Where("auth_organizations.name ILIKE ? OR auth_organizations.description ILIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	var organizations []OrganizationInfo
	query.Group("auth_organizations.id").
		Order("auth_organizations.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Scan(&organizations)

	// 计算订阅状态和到期天数
	now := time.Now()
	for i := range organizations {
		org := &organizations[i]
		if org.SubscriptionEndDate == nil || *org.SubscriptionEndDate == "" {
			org.SubscriptionStatus = "none"
		} else {
			// 尝试多种日期格式
			var endDate time.Time
			var parseErr error

			// 尝试格式1: "2006-01-02 15:04:05.999999"
			endDate, parseErr = time.Parse("2006-01-02 15:04:05.999999", *org.SubscriptionEndDate)
			if parseErr != nil {
				// 尝试格式2: "2006-01-02 15:04:05"
				endDate, parseErr = time.Parse("2006-01-02 15:04:05", *org.SubscriptionEndDate)
			}
			if parseErr != nil {
				// 尝试格式3: "2006-01-02T15:04:05Z"
				endDate, parseErr = time.Parse(time.RFC3339, *org.SubscriptionEndDate)
			}

			if parseErr == nil {
				days := int(endDate.Sub(now).Hours() / 24)
				org.DaysUntilExpiration = &days

				if days < 0 {
					org.SubscriptionStatus = "expired"
				} else if days <= 7 {
					org.SubscriptionStatus = "expiring_soon"
				} else {
					org.SubscriptionStatus = "active"
				}
			} else {
				org.SubscriptionStatus = "none"
			}
		}
	}

	c.JSON(http.StatusOK, organizations)
}

// GetOrganizationMetrics 获取组织指标
func (h *DashboardHandler) GetOrganizationMetrics(c *gin.Context) {
	orgID := c.Param("id")

	// 验证UUID格式
	if _, err := uuid.Parse(orgID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid organization ID format",
		})
		return
	}

	// 获取组织基础信息
	var org models.AuthOrganization
	if err := h.dbManager.LightAdminDB.Where("id = ?", orgID).First(&org).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Organization not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	// 获取用户统计
	var userStats struct {
		TotalUsers   int64 `json:"total_users"`
		ActiveUsers  int64 `json:"active_users"`
	}

	h.dbManager.LightAdminDB.Table("auth_users au").
		Joins("INNER JOIN auth_user_organization auo ON au.id = auo.user_id").
		Where("auo.organization_id = ?", orgID).
		Count(&userStats.TotalUsers)

	h.dbManager.LightAdminDB.Table("auth_users au").
		Joins("INNER JOIN auth_user_organization auo ON au.id = auo.user_id").
		Joins("INNER JOIN auth_user_workspace auw ON au.id = auw.user_id").
		Where("auo.organization_id = ? AND auw.user_status = ?", orgID, "active").
		Count(&userStats.ActiveUsers)

	// 获取订阅统计
	var subStats struct {
		TotalSubs      int64   `json:"total_subscriptions"`
		ActiveSubs     int64   `json:"active_subscriptions"`
		MonthlyRevenue float64 `json:"monthly_revenue"`
	}

	h.dbManager.LightAdminDB.Table("subscription_users").
		Where("organization_id = ?", orgID).
		Count(&subStats.TotalSubs)

	h.dbManager.LightAdminDB.Table("subscription_users").
		Where("organization_id = ? AND status = ?", orgID, "active").
		Count(&subStats.ActiveSubs)

	// 计算月收入（这里需要根据实际的计费逻辑）
	h.dbManager.LightAdminDB.Table("subscription_users su").
		Joins("INNER JOIN subscription_plans sp ON su.plan_id = sp.id").
		Where("su.organization_id = ? AND su.status = ? AND su.billing_cycle = ?", orgID, "active", "monthly").
		Select("COALESCE(SUM(sp.pricing_monthly), 0)").
		Scan(&subStats.MonthlyRevenue)

	// 获取工作空间统计
	var workspaceCount int64
	h.dbManager.LightAdminDB.Table("auth_workspaces").
		Where("organization_id = ?", orgID).
		Count(&workspaceCount)

	// 获取资源使用情况
	var resourceUsage struct {
		StorageUsage float64 `json:"storage_usage_mb"`
		QueryCount   int64   `json:"query_count_today"`
	}

	// 这里需要实现实际的资源统计逻辑
	// 目前使用示例数据
	resourceUsage.StorageUsage = 512.5
	resourceUsage.QueryCount = 1250

	metrics := gin.H{
		"organization": org,
		"user_stats":   userStats,
		"sub_stats":    subStats,
		"workspace_count": workspaceCount,
		"resource_usage": resourceUsage,
	}

	c.JSON(http.StatusOK, metrics)
}

// GetDatabaseStatus 获取数据库状态
func (h *DashboardHandler) GetDatabaseStatus(c *gin.Context) {
	response := DatabaseStatusResponse{
		ClickHouse: make(map[string]ClickHouseInfo),
	}

	// PostgreSQL状态
	if sqlDB, err := h.dbManager.LightAdminDB.DB(); err == nil {
		stats := sqlDB.Stats()
		response.PostgreSQL = DatabaseInfo{
			Status:         "healthy",
			Connections:    stats.OpenConnections,
			MaxConnections: h.dbManager.Config.Databases.LightAdmin.MaxOpenConns,
			ResponseTime:   10, // 示例响应时间
		}

		// 获取数据库大小
		var dbSize float64
		h.dbManager.LightAdminDB.Raw(`
			SELECT pg_size_pretty(pg_database_size(?)) as size
		`, "light_admin").Scan(&dbSize)
		response.PostgreSQL.DatabaseSize = dbSize
	} else {
		response.PostgreSQL = DatabaseInfo{
			Status: "unhealthy",
		}
	}

	// ClickHouse状态
	for name, conn := range h.dbManager.ClickHouse {
		info := ClickHouseInfo{
			Status: "healthy",
		}

		// 获取数据库大小和表数量
		rows, err := conn.Query(c.Request.Context(), `
			SELECT
				database,
				COUNT(*) as table_count,
				SUM(bytes) as total_bytes
			FROM system.parts
			WHERE database = ?
			GROUP BY database
		`, name)

		if err != nil {
			info.Status = "unhealthy"
		} else {
			if rows.Next() {
				var dbName string
				var tableCount int
				var totalBytes int64
				rows.Scan(&dbName, &tableCount, &totalBytes)
				info.TableCount = tableCount
				info.DatabaseSize = float64(totalBytes) / (1024 * 1024) // MB
			}
			rows.Close()
		}

		info.ResponseTime = 15 // 示例响应时间
		response.ClickHouse[name] = info
	}

	// Redis状态
	if h.dbManager.RedisClient != nil {
		info := h.dbManager.RedisClient.Info(c.Request.Context())
		if info.Err() == nil {
			response.Redis = RedisInfo{
				Status:          "healthy",
				ConnectedClients: 5, // 示例数据
				ResponseTime:    2,  // 示例响应时间
			}

			// 获取内存使用情况
			if memInfo, err := h.dbManager.RedisClient.Info(c.Request.Context(), "memory").Result(); err == nil {
				// 解析内存使用信息
				if strings.Contains(memInfo, "used_memory:") {
					// 简单解析used_memory字段
					lines := strings.Split(memInfo, "\r\n")
					for _, line := range lines {
						if strings.HasPrefix(line, "used_memory:") {
							parts := strings.Split(line, ":")
							if len(parts) == 2 {
								if usedMem, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64); err == nil {
									response.Redis.UsedMemory = usedMem / (1024 * 1024) // MB
								}
							}
							break
						}
					}
				}
			}
		} else {
			response.Redis = RedisInfo{
				Status: "unhealthy",
			}
		}
	}

	c.JSON(http.StatusOK, response)
}