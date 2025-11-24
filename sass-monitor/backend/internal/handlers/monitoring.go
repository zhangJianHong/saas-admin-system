package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"sass-monitor/internal/database"
	"sass-monitor/internal/models"
	"sass-monitor/pkg/config"
)

type MonitoringHandler struct {
	dbManager *database.DatabaseManager
	config    *config.Config
}

func NewMonitoringHandler(dbManager *database.DatabaseManager, cfg *config.Config) *MonitoringHandler {
	return &MonitoringHandler{
		dbManager: dbManager,
		config:    cfg,
	}
}

// MetricRequest 指标请求参数
type MetricRequest struct {
	DatabaseType  string    `form:"database_type"`
	DatabaseName  string    `form:"database_name"`
	MetricType    string    `form:"metric_type"`
	OrganizationID *string  `form:"organization_id"`
	StartTime     *time.Time `form:"start_time"`
	EndTime       *time.Time `form:"end_time"`
	Page          int       `form:"page,default=1"`
	PageSize      int       `form:"page_size,default=100"`
}

// AlertRequest 告警请求参数
type AlertRequest struct {
	Name               string  `json:"name" binding:"required"`
	Description        string  `json:"description"`
	RuleType           string  `json:"rule_type" binding:"required"`
	TargetType         string  `json:"target_type" binding:"required"`
	TargetName         string  `json:"target_name"`
	MetricName         string  `json:"metric_name" binding:"required"`
	Operator           string  `json:"operator" binding:"required"`
	Threshold          float64 `json:"threshold" binding:"required"`
	Duration           int     `json:"duration"`
	Severity           string  `json:"severity"`
	NotificationConfig string  `json:"notification_config"`
}

// GetMetrics 获取监控指标
func (h *MonitoringHandler) GetMetrics(c *gin.Context) {
	var req MetricRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request parameters",
			"details": err.Error(),
		})
		return
	}

	offset := (req.Page - 1) * req.PageSize

	query := h.dbManager.SaasMonitorDB.Model(&models.ResourceMetric{})

	if req.DatabaseType != "" {
		query = query.Where("database_type = ?", req.DatabaseType)
	}
	if req.DatabaseName != "" {
		query = query.Where("database_name = ?", req.DatabaseName)
	}
	if req.MetricType != "" {
		query = query.Where("metric_type = ?", req.MetricType)
	}
	if req.OrganizationID != nil {
		query = query.Where("organization_id = ?", *req.OrganizationID)
	}
	if req.StartTime != nil {
		query = query.Where("collected_at >= ?", *req.StartTime)
	}
	if req.EndTime != nil {
		query = query.Where("collected_at <= ?", *req.EndTime)
	}

	var metrics []models.ResourceMetric
	var total int64

	query.Count(&total)
	query.Order("collected_at DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&metrics)

	c.JSON(http.StatusOK, gin.H{
		"metrics": metrics,
		"total":   total,
		"page":    req.Page,
		"page_size": req.PageSize,
	})
}

// GetMetricsHistory 获取指标历史数据
func (h *MonitoringHandler) GetMetricsHistory(c *gin.Context) {
	databaseType := c.Query("database_type")
	databaseName := c.Query("database_name")
	metricType := c.Query("metric_type")
	orgID := c.Query("organization_id")

	// 解析时间范围
	hoursStr := c.DefaultQuery("hours", "24")
	hours, err := strconv.Atoi(hoursStr)
	if err != nil {
		hours = 24
	}

	endTime := time.Now()
	startTime := endTime.Add(-time.Duration(hours) * time.Hour)

	// 构建查询条件
	var results []gin.H

	switch databaseType {
	case "postgresql":
		// PostgreSQL指标历史
		results = h.getPostgreSQLMetrics(databaseName, metricType, startTime, endTime)
	case "clickhouse":
		// ClickHouse指标历史
		results = h.getClickHouseMetrics(databaseName, metricType, startTime, endTime)
	case "redis":
		// Redis指标历史
		results = h.getRedisMetrics(metricType, startTime, endTime)
	default:
		// 从监控数据库获取历史数据
		var metrics []models.ResourceMetric
		query := h.dbManager.SaasMonitorDB.Model(&models.ResourceMetric{}).
			Where("database_type = ? AND collected_at BETWEEN ? AND ?",
				databaseType, startTime, endTime)

		if databaseName != "" {
			query = query.Where("database_name = ?", databaseName)
		}
		if metricType != "" {
			query = query.Where("metric_type = ?", metricType)
		}
		if orgID != "" {
			query = query.Where("organization_id = ?", orgID)
		}

		query.Order("collected_at ASC").Find(&metrics)

		for _, metric := range metrics {
			results = append(results, gin.H{
				"timestamp": metric.CollectedAt.Unix(),
				"value":     metric.MetricValue,
				"metric_name": metric.MetricName,
				"unit":      metric.Unit,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"database_type": databaseType,
		"database_name": databaseName,
		"metric_type":   metricType,
		"start_time":    startTime.Unix(),
		"end_time":      endTime.Unix(),
		"data":          results,
	})
}

// GetOrganizations 获取组织列表（用于监控筛选）
func (h *MonitoringHandler) GetOrganizations(c *gin.Context) {
	var organizations []models.AuthOrganization

	query := h.dbManager.LightAdminDB.Model(&models.AuthOrganization{})

	search := c.Query("search")
	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	query.Order("created_at DESC").Find(&organizations)

	c.JSON(http.StatusOK, gin.H{
		"organizations": organizations,
	})
}

// GetOrganizationUsage 获取组织资源使用情况
func (h *MonitoringHandler) GetOrganizationUsage(c *gin.Context) {
	orgID := c.Param("id")

	// 验证UUID格式
	if _, err := uuid.Parse(orgID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid organization ID format",
		})
		return
	}

	// 获取时间范围
	hoursStr := c.DefaultQuery("hours", "24")
	hours, _ := strconv.Atoi(hoursStr)
	endTime := time.Now()
	startTime := endTime.Add(-time.Duration(hours) * time.Hour)

	usage := gin.H{
		"organization_id": orgID,
		"period": gin.H{
			"start_time": startTime.Unix(),
			"end_time":   endTime.Unix(),
			"hours":      hours,
		},
		"usage_by_database": gin.H{},
		"usage_by_type": gin.H{},
	}

	// 按数据库统计使用情况
	databaseTypes := []string{"postgresql", "clickhouse", "redis"}
	for _, dbType := range databaseTypes {
		var metrics []models.ResourceMetric
		h.dbManager.SaasMonitorDB.Model(&models.ResourceMetric{}).
			Where("organization_id = ? AND database_type = ? AND collected_at BETWEEN ? AND ?",
				orgID, dbType, startTime, endTime).
			Find(&metrics)

		usageByType := make(map[string][]gin.H)
		for _, metric := range metrics {
			usageByType[metric.MetricType] = append(usageByType[metric.MetricType], gin.H{
				"timestamp":    metric.CollectedAt.Unix(),
				"metric_name":  metric.MetricName,
				"metric_value": metric.MetricValue,
				"unit":         metric.Unit,
			})
		}

		usageByDatabase := usage["usage_by_database"].(gin.H)
		usageByDatabase[dbType] = usageByType
	}

	// 按使用类型统计
	var metrics []models.ResourceMetric
	h.dbManager.SaasMonitorDB.Model(&models.ResourceMetric{}).
		Where("organization_id = ? AND collected_at BETWEEN ? AND ?",
			orgID, startTime, endTime).
		Find(&metrics)

	usageByMetricType := make(map[string][]gin.H)
	for _, metric := range metrics {
		key := metric.DatabaseType + "_" + metric.MetricType
		usageByMetricType[key] = append(usageByMetricType[key], gin.H{
			"timestamp":    metric.CollectedAt.Unix(),
			"database_name": metric.DatabaseName,
			"metric_name":  metric.MetricName,
			"metric_value": metric.MetricValue,
			"unit":         metric.Unit,
		})
	}

	usage["usage_by_type"] = usageByMetricType

	c.JSON(http.StatusOK, usage)
}

// GetDatabaseInfo 获取数据库详细信息
func (h *MonitoringHandler) GetDatabaseInfo(c *gin.Context) {
	dbType := c.Query("type")

	response := gin.H{}

	switch dbType {
	case "postgresql":
		response = h.getPostgreSQLInfo()
	case "clickhouse":
		response = h.getClickHouseInfo()
	case "redis":
		response = h.getRedisInfo()
	default:
		response = gin.H{
			"postgresql": h.getPostgreSQLInfo(),
			"clickhouse": h.getClickHouseInfo(),
			"redis":      h.getRedisInfo(),
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetAlerts 获取告警列表
func (h *MonitoringHandler) GetAlerts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	enabled := c.Query("enabled")

	offset := (page - 1) * pageSize

	query := h.dbManager.SaasMonitorDB.Model(&models.AlertRule{})

	if enabled != "" {
		enabledBool := enabled == "true"
		query = query.Where("enabled = ?", enabledBool)
	}

	var alerts []models.AlertRule
	var total int64

	query.Count(&total)
	query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&alerts)

	c.JSON(http.StatusOK, gin.H{
		"alerts": alerts,
		"total":  total,
		"page":   page,
		"page_size": pageSize,
	})
}

// CreateAlert 创建告警规则
func (h *MonitoringHandler) CreateAlert(c *gin.Context) {
	userID := c.GetString("user_id")

	var req AlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	alert := models.AlertRule{
		Name:                req.Name,
		Description:         req.Description,
		RuleType:            req.RuleType,
		TargetType:          req.TargetType,
		TargetName:          req.TargetName,
		MetricName:          req.MetricName,
		Operator:            req.Operator,
		Threshold:           req.Threshold,
		Duration:            req.Duration,
		Severity:            req.Severity,
		Enabled:             true,
		NotificationConfig:  req.NotificationConfig,
		CreatedBy:           userUUID,
	}

	if alert.Duration == 0 {
		alert.Duration = 5 // 默认5分钟
	}
	if alert.Severity == "" {
		alert.Severity = "warning"
	}

	if err := h.dbManager.SaasMonitorDB.Create(&alert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create alert rule",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, alert)
}

// UpdateAlert 更新告警规则
func (h *MonitoringHandler) UpdateAlert(c *gin.Context) {
	alertID := c.Param("id")
	_ = c.GetString("user_id") // 用户ID，用于记录更新者

	alertUUID, err := uuid.Parse(alertID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid alert ID format",
		})
		return
	}

	var alert models.AlertRule
	if err := h.dbManager.SaasMonitorDB.Where("id = ?", alertUUID).First(&alert).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Alert rule not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	var req AlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// 更新字段
	alert.Name = req.Name
	alert.Description = req.Description
	alert.RuleType = req.RuleType
	alert.TargetType = req.TargetType
	alert.TargetName = req.TargetName
	alert.MetricName = req.MetricName
	alert.Operator = req.Operator
	alert.Threshold = req.Threshold
	alert.Duration = req.Duration
	alert.Severity = req.Severity
	alert.NotificationConfig = req.NotificationConfig

	if err := h.dbManager.SaasMonitorDB.Save(&alert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update alert rule",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, alert)
}

// DeleteAlert 删除告警规则
func (h *MonitoringHandler) DeleteAlert(c *gin.Context) {
	alertID := c.Param("id")

	alertUUID, err := uuid.Parse(alertID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid alert ID format",
		})
		return
	}

	if err := h.dbManager.SaasMonitorDB.Delete(&models.AlertRule{}, alertUUID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete alert rule",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Alert rule deleted successfully",
	})
}

// GetSystemHealth 获取系统健康状态
func (h *MonitoringHandler) GetSystemHealth(c *gin.Context) {
	var healthRecords []models.SystemHealth

	h.dbManager.SaasMonitorDB.Order("last_checked_at DESC").Find(&healthRecords)

	c.JSON(http.StatusOK, gin.H{
		"health_records": healthRecords,
	})
}

// GetSystemLogs 获取系统日志
func (h *MonitoringHandler) GetSystemLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	level := c.Query("level")
	source := c.Query("source")

	offset := (page - 1) * pageSize

	query := h.dbManager.SaasMonitorDB.Model(&models.MonitoringLog{})

	if level != "" {
		query = query.Where("log_level = ?", level)
	}
	if source != "" {
		query = query.Where("source = ?", source)
	}

	var logs []models.MonitoringLog
	var total int64

	query.Count(&total)
	query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}

// GetSystemConfigs 获取系统配置
func (h *MonitoringHandler) GetSystemConfigs(c *gin.Context) {
	var configs []models.MonitoringConfig

	h.dbManager.SaasMonitorDB.Find(&configs)

	configMap := make(map[string]string)
	for _, config := range configs {
		configMap[config.ConfigKey] = config.ConfigValue
	}

	c.JSON(http.StatusOK, gin.H{
		"configs": configMap,
	})
}

// UpdateSystemConfigs 更新系统配置
func (h *MonitoringHandler) UpdateSystemConfigs(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// 使用事务确保所有配置更新成功
	tx := h.dbManager.SaasMonitorDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for key, value := range req {
		var config models.MonitoringConfig
		if err := tx.Where("config_key = ?", key).First(&config).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// 创建新配置
				config = models.MonitoringConfig{
					ConfigKey:   key,
					ConfigValue: fmt.Sprintf("%v", value),
				}
				if err := tx.Create(&config).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "Failed to create configuration",
						"details": err.Error(),
						"config_key": key,
					})
					return
				}
			} else {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Database error",
					"details": err.Error(),
				})
				return
			}
		} else {
			// 更新现有配置
			config.ConfigValue = fmt.Sprintf("%v", value)
			if err := tx.Save(&config).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to update configuration",
					"details": err.Error(),
					"config_key": key,
				})
				return
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to commit transaction",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "System configurations updated successfully",
	})
}

// 辅助方法
func (h *MonitoringHandler) getPostgreSQLMetrics(dbName, metricType string, start, end time.Time) []gin.H {
	// 实现PostgreSQL指标获取逻辑
	return []gin.H{}
}

func (h *MonitoringHandler) getClickHouseMetrics(dbName, metricType string, start, end time.Time) []gin.H {
	// 实现ClickHouse指标获取逻辑
	return []gin.H{}
}

func (h *MonitoringHandler) getRedisMetrics(metricType string, start, end time.Time) []gin.H {
	// 实现Redis指标获取逻辑
	return []gin.H{}
}

func (h *MonitoringHandler) getPostgreSQLInfo() gin.H {
	// 返回PostgreSQL详细信息
	return gin.H{
		"status": "healthy",
		"databases": []string{"light_admin", "saas_monitor"},
	}
}

func (h *MonitoringHandler) getClickHouseInfo() gin.H {
	// 返回ClickHouse详细信息
	databases := []string{}
	for _, chConfig := range h.config.ClickHouse {
		databases = append(databases, chConfig.Database)
	}

	return gin.H{
		"status": "healthy",
		"databases": databases,
	}
}

func (h *MonitoringHandler) getRedisInfo() gin.H {
	// 返回Redis详细信息
	return gin.H{
		"status": "healthy",
		"database": h.config.Redis.Database,
	}
}

// GetOrganizationOverview 获取组织概览数据
func (h *MonitoringHandler) GetOrganizationOverview(c *gin.Context) {
	// 获取所有组织
	var organizations []models.AuthOrganization
	if err := h.dbManager.LightAdminDB.Find(&organizations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch organizations",
		})
		return
	}

	var overviewData []gin.H

	// 获取最近24小时的时间范围
	endTime := time.Now()
	startTime := endTime.Add(-24 * time.Hour)

	for _, org := range organizations {
		orgIDStr := org.ID.String()

		// 获取该组织的最新指标
		var userMetric models.ResourceMetric
		var workspaceMetric models.ResourceMetric
		var subscriptionMetric models.ResourceMetric

		// 获取用户数
		h.dbManager.SaasMonitorDB.Where(
			"organization_id = ? AND metric_name = ? AND collected_at BETWEEN ? AND ?",
			&orgIDStr, "user_count", startTime, endTime,
		).Order("collected_at DESC").First(&userMetric)

		// 获取工作空间数
		h.dbManager.SaasMonitorDB.Where(
			"organization_id = ? AND metric_name = ? AND collected_at BETWEEN ? AND ?",
			&orgIDStr, "workspace_count", startTime, endTime,
		).Order("collected_at DESC").First(&workspaceMetric)

		// 获取订阅数
		h.dbManager.SaasMonitorDB.Where(
			"organization_id = ? AND metric_name = ? AND collected_at BETWEEN ? AND ?",
			&orgIDStr, "active_subscriptions", startTime, endTime,
		).Order("collected_at DESC").First(&subscriptionMetric)

		// 构建组织概览数据
		orgData := gin.H{
			"organization_id":   orgIDStr,
			"organization_name": org.Name,
			"owner_id":         org.OwnerID.String(),
			"description":      org.Description,
			"created_at":       org.CreatedAt.Unix(),
			"updated_at":       org.UpdatedAt,
			"metrics": gin.H{
				"user_count": gin.H{
					"value":       userMetric.MetricValue,
					"unit":        userMetric.Unit,
					"last_updated": userMetric.CollectedAt.Unix(),
				},
				"workspace_count": gin.H{
					"value":       workspaceMetric.MetricValue,
					"unit":        workspaceMetric.Unit,
					"last_updated": workspaceMetric.CollectedAt.Unix(),
				},
				"active_subscriptions": gin.H{
					"value":       subscriptionMetric.MetricValue,
					"unit":        subscriptionMetric.Unit,
					"last_updated": subscriptionMetric.CollectedAt.Unix(),
				},
			},
		}

		overviewData = append(overviewData, orgData)
	}

	// 获取系统级统计
	var totalOrgs, totalUsers, totalWorkspaces, totalSubs int64

	h.dbManager.LightAdminDB.Table("auth_organizations").Count(&totalOrgs)
	h.dbManager.SaasMonitorDB.Model(&models.ResourceMetric{}).
		Where("metric_name = ? AND collected_at > ?", "user_count", startTime).
		Count(&totalUsers)
	h.dbManager.SaasMonitorDB.Model(&models.ResourceMetric{}).
		Where("metric_name = ? AND collected_at > ?", "workspace_count", startTime).
		Count(&totalWorkspaces)
	h.dbManager.SaasMonitorDB.Model(&models.ResourceMetric{}).
		Where("metric_name = ? AND collected_at > ?", "active_subscriptions", startTime).
		Count(&totalSubs)

	response := gin.H{
		"summary": gin.H{
			"total_organizations":     totalOrgs,
			"total_users":           totalUsers,
			"total_workspaces":      totalWorkspaces,
			"total_subscriptions":   totalSubs,
		},
		"organizations": overviewData,
		"period": gin.H{
			"start_time": startTime.Unix(),
			"end_time":   endTime.Unix(),
			"hours":      24,
		},
		"updated_at": endTime.Unix(),
	}

	c.JSON(http.StatusOK, response)
}