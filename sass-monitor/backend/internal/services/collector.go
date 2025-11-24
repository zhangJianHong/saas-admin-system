package services

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"gorm.io/gorm"

	"sass-monitor/internal/database"
	"sass-monitor/internal/models"
)

type DataCollector struct {
	dbManager *database.DatabaseManager
}

func NewDataCollector(dbManager *database.DatabaseManager) *DataCollector {
	return &DataCollector{
		dbManager: dbManager,
	}
}

// CollectAllData 采集所有监控数据
func (dc *DataCollector) CollectAllData(ctx context.Context) error {
	log.Println("Starting data collection...")

	// 采集PostgreSQL数据
	if err := dc.collectPostgreSQLData(ctx); err != nil {
		log.Printf("Error collecting PostgreSQL data: %v", err)
	}

	// 采集ClickHouse数据
	if err := dc.collectClickHouseData(ctx); err != nil {
		log.Printf("Error collecting ClickHouse data: %v", err)
	}

	// 采集Redis数据
	if err := dc.collectRedisData(ctx); err != nil {
		log.Printf("Error collecting Redis data: %v", err)
	}

	// 采集系统健康状态
	if err := dc.collectSystemHealth(ctx); err != nil {
		log.Printf("Error collecting system health: %v", err)
	}

	log.Println("Data collection completed")
	return nil
}

// collectPostgreSQLData 采集PostgreSQL监控数据
func (dc *DataCollector) collectPostgreSQLData(ctx context.Context) error {
	// 获取数据库连接统计
	if err := dc.collectPostgreSQLConnections(ctx); err != nil {
		return fmt.Errorf("failed to collect PostgreSQL connections: %w", err)
	}

	// 获取数据库大小
	if err := dc.collectPostgreSQLDatabaseSize(ctx); err != nil {
		return fmt.Errorf("failed to collect PostgreSQL database size: %w", err)
	}

	// 获取表大小统计
	if err := dc.collectPostgreSQLTableSize(ctx); err != nil {
		return fmt.Errorf("failed to collect PostgreSQL table size: %w", err)
	}

	// 获取用户和订阅统计
	if err := dc.collectPostgreSQLUserStats(ctx); err != nil {
		return fmt.Errorf("failed to collect PostgreSQL user stats: %w", err)
	}

	// 获取组织维度统计
	if err := dc.collectPostgreSQLOrganizationStats(ctx); err != nil {
		return fmt.Errorf("failed to collect PostgreSQL organization stats: %w", err)
	}

	return nil
}

// collectPostgreSQLConnections 采集PostgreSQL连接数据
func (dc *DataCollector) collectPostgreSQLConnections(ctx context.Context) error {
	sqlDB, err := dc.dbManager.LightAdminDB.DB()
	if err != nil {
		return err
	}

	stats := sqlDB.Stats()

	// 记录连接数指标
	metric := models.ResourceMetric{
		DatabaseType:  "postgresql",
		DatabaseName:  "light_admin",
		MetricType:    "connection",
		MetricName:    "active_connections",
		MetricValue:   float64(stats.OpenConnections),
		Unit:          "count",
		CollectedAt:   time.Now(),
	}

	// 添加标签信息
	tags := map[string]interface{}{
		"max_connections": dc.dbManager.Config.Databases.LightAdmin.MaxOpenConns,
		"idle_connections": stats.Idle,
	}
	metric.Tags = dc.formatTags(tags)

	return dc.dbManager.SaasMonitorDB.Create(&metric).Error
}

// collectPostgreSQLDatabaseSize 采集PostgreSQL数据库大小
func (dc *DataCollector) collectPostgreSQLDatabaseSize(ctx context.Context) error {
	var dbSize float64

	err := dc.dbManager.LightAdminDB.Raw(`
		SELECT pg_database_size(current_database()) as size_bytes
	`).Scan(&dbSize).Error

	if err != nil {
		return err
	}

	// 转换为MB
	dbSizeMB := dbSize / (1024 * 1024)

	metric := models.ResourceMetric{
		DatabaseType: "postgresql",
		DatabaseName: "light_admin",
		MetricType:   "storage",
		MetricName:   "database_size_mb",
		MetricValue:  dbSizeMB,
		Unit:         "MB",
		CollectedAt:  time.Now(),
	}

	return dc.dbManager.SaasMonitorDB.Create(&metric).Error
}

// collectPostgreSQLTableSize 采集PostgreSQL表大小统计
func (dc *DataCollector) collectPostgreSQLTableSize(ctx context.Context) error {
	var tableStats []struct {
		TableName string  `gorm:"column:tablename"`
		Size      float64 `gorm:"column:size_mb"`
		RowCount  int64   `gorm:"column:row_count"`
	}

	err := dc.dbManager.LightAdminDB.Raw(`
		SELECT
			schemaname||'.'||tablename as tablename,
			pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size,
			pg_total_relation_size(schemaname||'.'||tablename) / (1024*1024) as size_mb,
			COALESCE(n_tup_ins, 0) as row_count
		FROM pg_tables t
		LEFT JOIN pg_stat_user_tables s ON t.tablename = s.relname
		WHERE schemaname = 'public'
		ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC
		LIMIT 20
	`).Scan(&tableStats).Error

	if err != nil {
		return err
	}

	// 记录每个表的指标
	for _, stat := range tableStats {
		// 表大小指标
		sizeMetric := models.ResourceMetric{
			DatabaseType: "postgresql",
			DatabaseName: "light_admin",
			MetricType:   "storage",
			MetricName:   fmt.Sprintf("table_size_%s", stat.TableName),
			MetricValue:  stat.Size,
			Unit:         "MB",
			CollectedAt:  time.Now(),
			Tags:         dc.formatTags(map[string]interface{}{"table": stat.TableName}),
		}

		// 行数指标
		rowMetric := models.ResourceMetric{
			DatabaseType: "postgresql",
			DatabaseName: "light_admin",
			MetricType:   "row_count",
			MetricName:   fmt.Sprintf("table_rows_%s", stat.TableName),
			MetricValue:  float64(stat.RowCount),
			Unit:         "count",
			CollectedAt:  time.Now(),
			Tags:         dc.formatTags(map[string]interface{}{"table": stat.TableName}),
		}

		// 批量创建指标
		if err := dc.dbManager.SaasMonitorDB.Create([]models.ResourceMetric{sizeMetric, rowMetric}).Error; err != nil {
			log.Printf("Error creating table metrics for %s: %v", stat.TableName, err)
		}
	}

	return nil
}

// collectPostgreSQLUserStats 采集用户和订阅统计
func (dc *DataCollector) collectPostgreSQLUserStats(ctx context.Context) error {
	// 获取组织统计
	var orgStats struct {
		TotalOrgs int64 `json:"total_organizations"`
	}

	if err := dc.dbManager.LightAdminDB.Table("auth_organizations").Count(&orgStats.TotalOrgs).Error; err != nil {
		return err
	}

	// 获取用户统计
	var userStats struct {
		TotalUsers  int64 `json:"total_users"`
		ActiveUsers int64 `json:"active_users"`
	}

	dc.dbManager.LightAdminDB.Table("auth_users").Count(&userStats.TotalUsers)

	// 活跃用户定义：最近30天有登录记录
	dc.dbManager.LightAdminDB.Table("auth_users").
		Where("last_login_at > ?", time.Now().AddDate(0, 0, -30)).
		Count(&userStats.ActiveUsers)

	// 获取订阅统计
	var subStats struct {
		TotalSubs     int64   `json:"total_subscriptions"`
		ActiveSubs    int64   `json:"active_subscriptions"`
		MonthlyRevenue float64 `json:"monthly_revenue"`
	}

	dc.dbManager.LightAdminDB.Table("subscription_users").
		Count(&subStats.TotalSubs)

	dc.dbManager.LightAdminDB.Table("subscription_users").
		Where("status = ?", "active").
		Count(&subStats.ActiveSubs)

	// 计算月收入
	dc.dbManager.LightAdminDB.Raw(`
		SELECT COALESCE(SUM(pricing_monthly), 0)
		FROM subscription_users su
		JOIN subscription_plans sp ON su.plan_id = sp.id
		WHERE su.status = 'active' AND su.billing_cycle = 'monthly'
	`).Scan(&subStats.MonthlyRevenue)

	// 创建组织统计指标
	orgMetric := models.ResourceMetric{
		DatabaseType: "postgresql",
		DatabaseName: "light_admin",
		MetricType:   "organization_count",
		MetricName:   "total_organizations",
		MetricValue:  float64(orgStats.TotalOrgs),
		Unit:         "count",
		CollectedAt:  time.Now(),
	}

	// 创建用户统计指标
	userMetric := models.ResourceMetric{
		DatabaseType: "postgresql",
		DatabaseName: "light_admin",
		MetricType:   "user_count",
		MetricName:   "total_users",
		MetricValue:  float64(userStats.TotalUsers),
		Unit:         "count",
		CollectedAt:  time.Now(),
		Tags:         dc.formatTags(map[string]interface{}{"active_users": userStats.ActiveUsers}),
	}

	// 创建订阅统计指标
	subMetric := models.ResourceMetric{
		DatabaseType: "postgresql",
		DatabaseName: "light_admin",
		MetricType:   "subscription_count",
		MetricName:   "active_subscriptions",
		MetricValue:  float64(subStats.ActiveSubs),
		Unit:         "count",
		CollectedAt:  time.Now(),
		Tags:         dc.formatTags(map[string]interface{}{"monthly_revenue": subStats.MonthlyRevenue}),
	}

	// 创建收入指标
	revenueMetric := models.ResourceMetric{
		DatabaseType: "postgresql",
		DatabaseName: "light_admin",
		MetricType:   "revenue",
		MetricName:   "monthly_revenue",
		MetricValue:  subStats.MonthlyRevenue,
		Unit:         "USD",
		CollectedAt:  time.Now(),
	}

	return dc.dbManager.SaasMonitorDB.Create([]models.ResourceMetric{
		orgMetric, userMetric, subMetric, revenueMetric,
	}).Error
}

// collectClickHouseData 采集ClickHouse监控数据
func (dc *DataCollector) collectClickHouseData(ctx context.Context) error {
	for dbName, conn := range dc.dbManager.ClickHouse {
		if err := dc.collectClickHouseDatabaseData(ctx, dbName, conn); err != nil {
			log.Printf("Error collecting ClickHouse data for %s: %v", dbName, err)
		}
	}
	return nil
}

// collectClickHouseDatabaseData 采集单个ClickHouse数据库的数据
func (dc *DataCollector) collectClickHouseDatabaseData(ctx context.Context, dbName string, conn clickhouse.Conn) error {
	// 获取数据库大小
	if err := dc.collectClickHouseDatabaseSize(ctx, dbName, conn); err != nil {
		return err
	}

	// 获取表统计
	if err := dc.collectClickHouseTableStats(ctx, dbName, conn); err != nil {
		return err
	}

	// 获取查询性能指标
	if err := dc.collectClickHouseQueryStats(ctx, dbName, conn); err != nil {
		return err
	}

	return nil
}

// collectClickHouseDatabaseSize 采集ClickHouse数据库大小
func (dc *DataCollector) collectClickHouseDatabaseSize(ctx context.Context, dbName string, conn clickhouse.Conn) error {
	var dbSize struct {
		TotalBytes int64 `ch:"total_bytes"`
		TableCount  int64  `ch:"table_count"`
		RowCount    int64  `ch:"row_count"`
	}

	err := conn.QueryRow(ctx, `
		SELECT
			SUM(bytes) as total_bytes,
			COUNT(DISTINCT table) as table_count,
			SUM(rows) as row_count
		FROM system.parts
		WHERE database = ? AND active = 1
	`, dbName).Scan(&dbSize.TotalBytes, &dbSize.TableCount, &dbSize.RowCount)

	if err != nil {
		return err
	}

	// 转换为MB
	sizeMB := float64(dbSize.TotalBytes) / (1024 * 1024)

	// 数据库大小指标
	sizeMetric := models.ResourceMetric{
		DatabaseType: "clickhouse",
		DatabaseName: dbName,
		MetricType:   "storage",
		MetricName:   "database_size_mb",
		MetricValue:  sizeMB,
		Unit:         "MB",
		CollectedAt:  time.Now(),
	}

	// 表数量指标
	tableMetric := models.ResourceMetric{
		DatabaseType: "clickhouse",
		DatabaseName: dbName,
		MetricType:   "table_count",
		MetricName:   "total_tables",
		MetricValue:  float64(dbSize.TableCount),
		Unit:         "count",
		CollectedAt:  time.Now(),
	}

	// 行数指标
	rowMetric := models.ResourceMetric{
		DatabaseType: "clickhouse",
		DatabaseName: dbName,
		MetricType:   "row_count",
		MetricName:   "total_rows",
		MetricValue:  float64(dbSize.RowCount),
		Unit:         "count",
		CollectedAt:  time.Now(),
	}

	return dc.dbManager.SaasMonitorDB.Create([]models.ResourceMetric{
		sizeMetric, tableMetric, rowMetric,
	}).Error
}

// collectClickHouseTableStats 采集ClickHouse表统计
func (dc *DataCollector) collectClickHouseTableStats(ctx context.Context, dbName string, conn clickhouse.Conn) error {
	rows, err := conn.Query(ctx, `
		SELECT
			table,
			SUM(bytes) as total_bytes,
			SUM(rows) as total_rows
		FROM system.parts
		WHERE database = ? AND active = 1
		GROUP BY table
		ORDER BY total_bytes DESC
		LIMIT 20
	`, dbName)

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		var totalBytes int64
		var totalRows int64

		if err := rows.Scan(&tableName, &totalBytes, &totalRows); err != nil {
			log.Printf("Error scanning ClickHouse table stats: %v", err)
			continue
		}

		// 表大小指标
		sizeMetric := models.ResourceMetric{
			DatabaseType: "clickhouse",
			DatabaseName: dbName,
			MetricType:   "storage",
			MetricName:   fmt.Sprintf("table_size_%s", tableName),
			MetricValue:  float64(totalBytes) / (1024 * 1024), // MB
			Unit:         "MB",
			CollectedAt:  time.Now(),
			Tags:         dc.formatTags(map[string]interface{}{"table": tableName}),
		}

		// 行数指标
		rowMetric := models.ResourceMetric{
			DatabaseType: "clickhouse",
			DatabaseName: dbName,
			MetricType:   "row_count",
			MetricName:   fmt.Sprintf("table_rows_%s", tableName),
			MetricValue:  float64(totalRows),
			Unit:         "count",
			CollectedAt:  time.Now(),
			Tags:         dc.formatTags(map[string]interface{}{"table": tableName}),
		}

		if err := dc.dbManager.SaasMonitorDB.Create([]models.ResourceMetric{sizeMetric, rowMetric}).Error; err != nil {
			log.Printf("Error creating ClickHouse table metrics for %s.%s: %v", dbName, tableName, err)
		}
	}

	return nil
}

// collectClickHouseQueryStats 采集ClickHouse查询性能指标
func (dc *DataCollector) collectClickHouseQueryStats(ctx context.Context, dbName string, conn clickhouse.Conn) error {
	// 获取慢查询统计
	var slowQueries struct {
		Count    int64   `ch:"count"`
		AvgTime  float64 `ch:"avg_time"`
		MaxTime  float64 `ch:"max_time"`
	}

	err := conn.QueryRow(ctx, `
		SELECT
			COUNT() as count,
			avg(query_duration_ms) as avg_time,
			max(query_duration_ms) as max_time
		FROM system.query_log
		WHERE database = ?
			AND type = 'QueryFinish'
			AND event_time > now() - INTERVAL 1 HOUR
			AND query_duration_ms > 1000
	`, dbName).Scan(&slowQueries.Count, &slowQueries.AvgTime, &slowQueries.MaxTime)

	if err != nil {
		// 如果没有查询日志表，跳过
		return nil
	}

	// 慢查询数量指标
	slowCountMetric := models.ResourceMetric{
		DatabaseType: "clickhouse",
		DatabaseName: dbName,
		MetricType:   "query_performance",
		MetricName:   "slow_queries_count_1h",
		MetricValue:  float64(slowQueries.Count),
		Unit:         "count",
		CollectedAt:  time.Now(),
	}

	// 平均查询时间指标
	avgTimeMetric := models.ResourceMetric{
		DatabaseType: "clickhouse",
		DatabaseName: dbName,
		MetricType:   "query_performance",
		MetricName:   "avg_query_time_1h",
		MetricValue:  slowQueries.AvgTime,
		Unit:         "ms",
		CollectedAt:  time.Now(),
	}

	return dc.dbManager.SaasMonitorDB.Create([]models.ResourceMetric{
		slowCountMetric, avgTimeMetric,
	}).Error
}

// collectRedisData 采集Redis监控数据
func (dc *DataCollector) collectRedisData(ctx context.Context) error {
	if dc.dbManager.RedisClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	// 获取Redis信息
	info := dc.dbManager.RedisClient.Info(ctx)
	if info.Err() != nil {
		return info.Err()
	}

	// 解析Redis信息
	redisInfo := info.Val()

	// 提取关键指标
	memoryUsed := dc.extractRedisMetric(redisInfo, "used_memory:")
	_ = dc.extractRedisMetric(redisInfo, "maxmemory:") // 最大内存（暂不使用）
	connectedClients := dc.extractRedisMetric(redisInfo, "connected_clients:")
	keyspaceHits := dc.extractRedisMetric(redisInfo, "keyspace_hits:")
	keyspaceMisses := dc.extractRedisMetric(redisInfo, "keyspace_misses:")

	// 计算命中率
	totalRequests := keyspaceHits + keyspaceMisses
	hitRate := float64(0)
	if totalRequests > 0 {
		hitRate = float64(keyspaceHits) / float64(totalRequests) * 100
	}

	// 内存使用指标
	memoryMetric := models.ResourceMetric{
		DatabaseType: "redis",
		DatabaseName: "default",
		MetricType:   "memory",
		MetricName:   "used_memory_bytes",
		MetricValue:  float64(memoryUsed),
		Unit:         "bytes",
		CollectedAt:  time.Now(),
	}

	// 连接数指标
	connectionsMetric := models.ResourceMetric{
		DatabaseType: "redis",
		DatabaseName: "default",
		MetricType:   "connection",
		MetricName:   "connected_clients",
		MetricValue:  float64(connectedClients),
		Unit:         "count",
		CollectedAt:  time.Now(),
	}

	// 命中率指标
	hitRateMetric := models.ResourceMetric{
		DatabaseType: "redis",
		DatabaseName: "default",
		MetricType:   "performance",
		MetricName:   "hit_rate_percent",
		MetricValue:  hitRate,
		Unit:         "percent",
		CollectedAt:  time.Now(),
	}

	return dc.dbManager.SaasMonitorDB.Create([]models.ResourceMetric{
		memoryMetric, connectionsMetric, hitRateMetric,
	}).Error
}

// collectSystemHealth 采集系统健康状态
func (dc *DataCollector) collectSystemHealth(ctx context.Context) error {
	healthStatus := dc.dbManager.HealthCheck()

	for component, err := range healthStatus {
		status := "healthy"
		responseTime := int(0)
		errorMessage := ""

		if err != nil {
			status = "unhealthy"
			errorMessage = err.Error()
		} else {
			// 模拟响应时间测量
			start := time.Now()
			switch component {
			case "saas_monitor":
				dc.dbManager.SaasMonitorDB.Exec("SELECT 1")
			case "light_admin":
				dc.dbManager.LightAdminDB.Exec("SELECT 1")
			case "redis":
				dc.dbManager.RedisClient.Ping(ctx)
			}
			responseTime = int(time.Since(start).Milliseconds())
		}

		// 查找或创建健康记录
		var healthRecord models.SystemHealth
		result := dc.dbManager.SaasMonitorDB.Where("component_name = ?", component).First(&healthRecord)

		now := time.Now()
		if result.Error == gorm.ErrRecordNotFound {
			// 创建新记录
			healthRecord = models.SystemHealth{
				ComponentName: component,
				ComponentType: dc.getComponentType(component),
				Status:        status,
				ResponseTime:  &responseTime,
				ErrorMessage:  &errorMessage,
				LastCheckedAt: now,
				CreatedAt:     now,
				UpdatedAt:     now,
			}
			dc.dbManager.SaasMonitorDB.Create(&healthRecord)
		} else {
			// 更新现有记录
			healthRecord.Status = status
			healthRecord.ResponseTime = &responseTime
			healthRecord.ErrorMessage = &errorMessage
			healthRecord.LastCheckedAt = now
			healthRecord.UpdatedAt = now
			dc.dbManager.SaasMonitorDB.Save(&healthRecord)
		}
	}

	return nil
}

// formatTags 格式化标签为JSON字符串
func (dc *DataCollector) formatTags(tags map[string]interface{}) string {
	// 简单的JSON格式化
	result := "{"
	first := true
	for key, value := range tags {
		if !first {
			result += ","
		}
		result += fmt.Sprintf(`"%s":"%v"`, key, value)
		first = false
	}
	result += "}"
	return result
}

// getComponentType 根据组件名称获取组件类型
func (dc *DataCollector) getComponentType(componentName string) string {
	if componentName == "saas_monitor" || componentName == "light_admin" {
		return "database"
	}
	if componentName == "redis" {
		return "cache"
	}
	if strings.Contains(componentName, "clickhouse") {
		return "database"
	}
	return "unknown"
}

// collectPostgreSQLOrganizationStats 采集组织维度统计数据
func (dc *DataCollector) collectPostgreSQLOrganizationStats(ctx context.Context) error {
	// 获取所有组织
	var organizations []models.AuthOrganization
	if err := dc.dbManager.LightAdminDB.Find(&organizations).Error; err != nil {
		return fmt.Errorf("failed to fetch organizations: %w", err)
	}

	var metrics []models.ResourceMetric
	now := time.Now()

	// 按组织统计资源使用情况
	for _, org := range organizations {
		orgIDStr := org.ID.String()

		// 统计该组织的用户数（简化统计，假设订阅用户数）
		var userCount int64
		dc.dbManager.LightAdminDB.Table("subscription_users").
			Where("organization_id = ?", orgIDStr).
			Count(&userCount)

		// 统计该组织的工作空间数
		var workspaceCount int64
		dc.dbManager.LightAdminDB.Table("auth_workspaces").
			Where("organization_id = ?", org.ID).
			Count(&workspaceCount)

		// 统计该组织的活跃订阅数
		var activeSubs int64
		dc.dbManager.LightAdminDB.Raw(`
			SELECT COUNT(*)
			FROM subscription_users su
			JOIN auth_organizations ao ON su.organization_id = ao.id::text
			WHERE su.organization_id = ? AND su.status = 'active'
		`, orgIDStr).Scan(&activeSubs)

		// 统计该组织的数据使用量（从org_usage表）
		var usageRecords []models.OrgUsage
		var totalUsage float64
		dc.dbManager.LightAdminDB.Table("org_usage").
			Where("organization_id = ? AND month = ?", orgIDStr, time.Now().Format("2006-01")).
			Find(&usageRecords)

		// 累加使用量（简化处理）
		for range usageRecords {
			// 这里可以根据实际的使用量JSON字段进行解析
			totalUsage += 1.0 // 暂时每个记录计为1单位使用量
		}

		// 创建组织用户数指标
		userMetric := models.ResourceMetric{
			OrganizationID: &orgIDStr,
			DatabaseType:   "postgresql",
			DatabaseName:   "light_admin",
			MetricType:     "organization_users",
			MetricName:     "user_count",
			MetricValue:    float64(userCount),
			Unit:           "count",
			CollectedAt:    now,
			Tags:           dc.formatTags(map[string]interface{}{"organization_name": org.Name}),
		}
		metrics = append(metrics, userMetric)

		// 创建组织工作空间数指标
		workspaceMetric := models.ResourceMetric{
			OrganizationID: &orgIDStr,
			DatabaseType:   "postgresql",
			DatabaseName:   "light_admin",
			MetricType:     "organization_workspaces",
			MetricName:     "workspace_count",
			MetricValue:    float64(workspaceCount),
			Unit:           "count",
			CollectedAt:    now,
			Tags:           dc.formatTags(map[string]interface{}{"organization_name": org.Name}),
		}
		metrics = append(metrics, workspaceMetric)

		// 创建组织订阅数指标
		subMetric := models.ResourceMetric{
			OrganizationID: &orgIDStr,
			DatabaseType:   "postgresql",
			DatabaseName:   "light_admin",
			MetricType:     "organization_subscriptions",
			MetricName:     "active_subscriptions",
			MetricValue:    float64(activeSubs),
			Unit:           "count",
			CollectedAt:    now,
			Tags:           dc.formatTags(map[string]interface{}{"organization_name": org.Name}),
		}
		metrics = append(metrics, subMetric)

		// 创建组织使用量指标
		usageMetric := models.ResourceMetric{
			OrganizationID: &orgIDStr,
			DatabaseType:   "postgresql",
			DatabaseName:   "light_admin",
			MetricType:     "organization_usage",
			MetricName:     "monthly_usage",
			MetricValue:    totalUsage,
			Unit:           "units",
			CollectedAt:    now,
			Tags:           dc.formatTags(map[string]interface{}{"organization_name": org.Name}),
		}
		metrics = append(metrics, usageMetric)
	}

	// 批量插入指标数据
	if len(metrics) > 0 {
		return dc.dbManager.SaasMonitorDB.CreateInBatches(metrics, 100).Error
	}

	return nil
}

// extractRedisMetric 从Redis信息中提取指定指标的值
func (dc *DataCollector) extractRedisMetric(info, metric string) int64 {
	lines := strings.Split(info, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, metric) {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				if value, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
					return value
				}
			}
		}
	}
	return 0
}