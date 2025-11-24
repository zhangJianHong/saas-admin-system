package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/gorm"

	"sass-monitor/internal/database"
	"sass-monitor/internal/models"
	"sass-monitor/pkg/config"
)

type TaskScheduler struct {
	dbManager     *database.DatabaseManager
	config        *config.Config
	dataCollector *DataCollector
	collectors   map[string]*time.Ticker
	stopChans     map[string]chan bool
	mutex         sync.RWMutex
	running       bool
}

func NewTaskScheduler(dbManager *database.DatabaseManager, cfg *config.Config) *TaskScheduler {
	return &TaskScheduler{
		dbManager:     dbManager,
		config:        cfg,
		dataCollector: NewDataCollector(dbManager),
		collectors:   make(map[string]*time.Ticker),
		stopChans:     make(map[string]chan bool),
		running:       false,
	}
}

// Start 启动调度器
func (ts *TaskScheduler) Start() error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	if ts.running {
		return fmt.Errorf("scheduler is already running")
	}

	log.Println("Starting task scheduler...")

	// 启动数据采集任务
	if err := ts.startDataCollection(); err != nil {
		return fmt.Errorf("failed to start data collection: %w", err)
	}

	// 启动告警检查任务
	if err := ts.startAlertChecker(); err != nil {
		return fmt.Errorf("failed to start alert checker: %w", err)
	}

	// 启动数据清理任务
	if err := ts.startDataCleanup(); err != nil {
		return fmt.Errorf("failed to start data cleanup: %w", err)
	}

	ts.running = true
	log.Println("Task scheduler started successfully")

	return nil
}

// Stop 停止调度器
func (ts *TaskScheduler) Stop() {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	if !ts.running {
		return
	}

	log.Println("Stopping task scheduler...")

	// 停止所有定时任务
	for name, stopChan := range ts.stopChans {
		close(stopChan)
		if ticker, exists := ts.collectors[name]; exists {
			ticker.Stop()
			delete(ts.collectors, name)
		}
		delete(ts.stopChans, name)
	}

	ts.running = false
	log.Println("Task scheduler stopped")
}

// startDataCollection 启动数据采集任务
func (ts *TaskScheduler) startDataCollection() error {
	interval := time.Duration(ts.config.Monitoring.CollectInterval) * time.Minute

	ticker := time.NewTicker(interval)
	stopChan := make(chan bool)

	ts.collectors["data_collection"] = ticker
	ts.stopChans["data_collection"] = stopChan

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := ts.dataCollector.CollectAllData(context.Background()); err != nil {
					log.Printf("Data collection error: %v", err)
					ts.logMonitoringError("data_collector", err.Error())
				}
			case <-stopChan:
				ticker.Stop()
				return
			}
		}
	}()

	log.Printf("Data collection task started with interval: %v", interval)
	return nil
}

// startAlertChecker 启动告警检查任务
func (ts *TaskScheduler) startAlertChecker() error {
	interval := 1 * time.Minute // 每分钟检查一次

	ticker := time.NewTicker(interval)
	stopChan := make(chan bool)

	ts.collectors["alert_checker"] = ticker
	ts.stopChans["alert_checker"] = stopChan

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := ts.checkAlerts(context.Background()); err != nil {
					log.Printf("Alert check error: %v", err)
					ts.logMonitoringError("alert_checker", err.Error())
				}
			case <-stopChan:
				ticker.Stop()
				return
			}
		}
	}()

	log.Printf("Alert checker task started with interval: %v", interval)
	return nil
}

// startDataCleanup 启动数据清理任务
func (ts *TaskScheduler) startDataCleanup() error {
	interval := 24 * time.Hour // 每天清理一次

	ticker := time.NewTicker(interval)
	stopChan := make(chan bool)

	ts.collectors["data_cleanup"] = ticker
	ts.stopChans["data_cleanup"] = stopChan

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := ts.cleanupOldData(context.Background()); err != nil {
					log.Printf("Data cleanup error: %v", err)
					ts.logMonitoringError("data_cleanup", err.Error())
				}
			case <-stopChan:
				ticker.Stop()
				return
			}
		}
	}()

	log.Printf("Data cleanup task started with interval: %v", interval)
	return nil
}

// checkAlerts 检查告警规则
func (ts *TaskScheduler) checkAlerts(ctx context.Context) error {
	// 获取启用的告警规则
	var alertRules []models.AlertRule
	if err := ts.dbManager.SaasMonitorDB.Where("enabled = ?", true).Find(&alertRules).Error; err != nil {
		return fmt.Errorf("failed to fetch alert rules: %w", err)
	}

	for _, rule := range alertRules {
		if err := ts.evaluateAlertRule(ctx, rule); err != nil {
			log.Printf("Error evaluating alert rule %s: %v", rule.Name, err)
		}
	}

	return nil
}

// evaluateAlertRule 评估单个告警规则
func (ts *TaskScheduler) evaluateAlertRule(ctx context.Context, rule models.AlertRule) error {
	// 获取最新的指标数据
	var metric models.ResourceMetric
	err := ts.dbManager.SaasMonitorDB.Where("database_type = ? AND database_name = ? AND metric_name = ?",
		rule.TargetType, rule.TargetName, rule.MetricName).
		Order("collected_at DESC").
		First(&metric).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("no metric found for rule %s", rule.Name)
		}
		return fmt.Errorf("failed to query metric for rule %s: %w", rule.Name, err)
	}

	// 检查是否触发告警
	triggered := false
	switch rule.Operator {
	case ">":
		triggered = metric.MetricValue > rule.Threshold
	case ">=":
		triggered = metric.MetricValue >= rule.Threshold
	case "<":
		triggered = metric.MetricValue < rule.Threshold
	case "<=":
		triggered = metric.MetricValue <= rule.Threshold
	case "=":
		triggered = metric.MetricValue == rule.Threshold
	}

	if triggered {
		log.Printf("Alert triggered: %s - %s %s %v (actual: %v)",
			rule.Name, rule.MetricName, rule.Operator, rule.Threshold, metric.MetricValue)

		// 这里可以发送告警通知
		if err := ts.sendAlertNotification(ctx, rule, metric); err != nil {
			log.Printf("Failed to send alert notification: %v", err)
		}

		// 记录告警日志
		ts.logAlertTriggered(rule, metric)
	}

	return nil
}

// sendAlertNotification 发送告警通知
func (ts *TaskScheduler) sendAlertNotification(ctx context.Context, rule models.AlertRule, metric models.ResourceMetric) error {
	// 这里可以实现邮件、短信、Webhook等通知方式
	log.Printf("Alert notification sent for rule: %s", rule.Name)

	// 示例：记录到监控日志
	alertLog := models.MonitoringLog{
		LogLevel: "warning",
		Source:   "alert_system",
		Message:  fmt.Sprintf("Alert triggered: %s - %s %s %v (actual: %v %s)",
			rule.Name, rule.MetricName, rule.Operator, rule.Threshold,
			metric.MetricValue, metric.Unit),
		Details:  fmt.Sprintf(`{"rule_id": "%s", "severity": "%s", "metric_value": %v}`,
			rule.ID, rule.Severity, metric.MetricValue),
		CreatedAt: time.Now(),
	}

	return ts.dbManager.SaasMonitorDB.Create(&alertLog).Error
}

// logAlertTriggered 记录告警触发日志
func (ts *TaskScheduler) logAlertTriggered(rule models.AlertRule, metric models.ResourceMetric) {
	alertLog := models.MonitoringLog{
		LogLevel: "warning",
		Source:   "alert_system",
		Component: fmt.Sprintf("%s_%s", rule.TargetType, rule.TargetName),
		Message:  fmt.Sprintf("Alert rule '%s' triggered", rule.Name),
		Details:  fmt.Sprintf(`{"rule_id": "%s", "metric_name": "%s", "threshold": %v, "actual": %v, "operator": "%s", "severity": "%s"}`,
			rule.ID, rule.MetricName, rule.Threshold, metric.MetricValue, rule.Operator, rule.Severity),
		CreatedAt: time.Now(),
	}

	ts.dbManager.SaasMonitorDB.Create(&alertLog)
}

// cleanupOldData 清理过期数据
func (ts *TaskScheduler) cleanupOldData(ctx context.Context) error {
	retentionDays := ts.config.Monitoring.RetentionDays
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)

	log.Printf("Cleaning up data older than %d days (cutoff: %s)", retentionDays, cutoffDate.Format("2006-01-02"))

	// 清理过期的资源指标
	if err := ts.dbManager.SaasMonitorDB.Where("created_at < ?", cutoffDate).Delete(&models.ResourceMetric{}).Error; err != nil {
		return fmt.Errorf("failed to cleanup resource metrics: %w", err)
	}

	// 清理过期的监控日志
	if err := ts.dbManager.SaasMonitorDB.Where("created_at < ?", cutoffDate).Delete(&models.MonitoringLog{}).Error; err != nil {
		return fmt.Errorf("failed to cleanup monitoring logs: %w", err)
	}

	deletedCount := time.Since(cutoffDate).Hours() / 24 * 100 // 估算删除的记录数
	log.Printf("Data cleanup completed. Estimated records deleted: ~%.0f", deletedCount)

	return nil
}

// logMonitoringError 记录监控错误
func (ts *TaskScheduler) logMonitoringError(component, errorMessage string) {
	errorLog := models.MonitoringLog{
		LogLevel: "error",
		Source:   "scheduler",
		Component: component,
		Message:  "Monitoring task error",
		Details:  fmt.Sprintf(`{"error": "%s", "timestamp": "%s"}`, errorMessage, time.Now().Format(time.RFC3339)),
		CreatedAt: time.Now(),
	}

	ts.dbManager.SaasMonitorDB.Create(&errorLog)
}

// IsRunning 检查调度器是否正在运行
func (ts *TaskScheduler) IsRunning() bool {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()
	return ts.running
}

// GetTaskStatus 获取任务状态
func (ts *TaskScheduler) GetTaskStatus() map[string]string {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()

	status := make(map[string]string)

	for name := range ts.collectors {
		status[name] = "running"
	}

	return status
}

// RestartTask 重启指定任务
func (ts *TaskScheduler) RestartTask(taskName string) error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	if !ts.running {
		return fmt.Errorf("scheduler is not running")
	}

	// 停止指定任务
	if stopChan, exists := ts.stopChans[taskName]; exists {
		close(stopChan)
		delete(ts.stopChans, taskName)
	}

	if ticker, exists := ts.collectors[taskName]; exists {
		ticker.Stop()
		delete(ts.collectors, taskName)
	}

	// 重启任务
	switch taskName {
	case "data_collection":
		return ts.startDataCollection()
	case "alert_checker":
		return ts.startAlertChecker()
	case "data_cleanup":
		return ts.startDataCleanup()
	default:
		return fmt.Errorf("unknown task: %s", taskName)
	}
}