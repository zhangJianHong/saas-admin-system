package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sass-monitor/internal/database"
	"sass-monitor/internal/models"
)

// OrganizationService 组织管理服务（只读模式）
type OrganizationService struct {
	dbManager *database.DatabaseManager
}

func NewOrganizationService(dbManager *database.DatabaseManager) *OrganizationService {
	return &OrganizationService{
		dbManager: dbManager,
	}
}

// OrganizationDetail 组织详细信息结构
type OrganizationDetail struct {
	ID                string     `json:"id" gorm:"column:id"`
	Name              string     `json:"name" gorm:"column:name"`
	OwnerID           string     `json:"owner_id" gorm:"column:owner_id"`
	Description       *string    `json:"description" gorm:"column:description"`
	UserCount         int64      `json:"user_count" gorm:"column:user_count"`
	ActiveUsers       int64      `json:"active_users" gorm:"column:active_users"`
	SubscriptionCount int64      `json:"subscription_count" gorm:"column:subscription_count"`
	WorkspaceCount    int64      `json:"workspace_count" gorm:"column:workspace_count"`
	StorageUsage      float64    `json:"storage_usage" gorm:"column:storage_usage"`
	// 订阅到期相关
	SubscriptionStatus      string     `json:"subscription_status" gorm:"-"`    // active/expiring_soon/expired/none (计算字段)
	SubscriptionEndDate     *time.Time `json:"subscription_end_date" gorm:"column:subscription_end_date"`  // 最近的订阅到期时间
	DaysUntilExpiration     *int       `json:"days_until_expiration" gorm:"-"`  // 距离到期天数 (计算字段)
	ActiveSubscriptionCount int64      `json:"active_subscription_count" gorm:"column:active_subscription_count"` // 活跃订阅数
	CreatedAt               time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt               *time.Time `json:"updated_at" gorm:"column:updated_at"`
}

// OrganizationUser 组织用户信息
type OrganizationUser struct {
	ID        string     `json:"id"`
	Username  string     `json:"username"`
	Nickname  *string    `json:"nickname"`
	Email     *string    `json:"email"`
	AvatarURL string     `json:"avatar_url"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// OrganizationWorkspace 组织工作空间信息
type OrganizationWorkspace struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Status      string     `json:"status"`
	Description *string    `json:"description"`
	UserCount   int64      `json:"user_count"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

// OrganizationSubscription 组织订阅信息
type OrganizationSubscription struct {
	ID               string     `json:"id"`
	UserID           string     `json:"user_id"`
	Username         string     `json:"username"`          // 订阅用户名
	UserEmail        *string    `json:"user_email"`        // 订阅用户邮箱
	PlanID           string     `json:"plan_id"`
	PlanName         string     `json:"plan_name"`
	PlanPricing      float64    `json:"plan_pricing"`      // 套餐价格
	Status           string     `json:"status"`
	BillingCycle     string     `json:"billing_cycle"`
	StartDate        time.Time  `json:"start_date"`
	EndDate          *time.Time `json:"end_date"`
	DaysUntilExpiry  *int       `json:"days_until_expiry"` // 距离到期天数
	PaymentMethod    *string    `json:"payment_method"`
	LastBilledAt     *time.Time `json:"last_billed_at"`
	TrialDaysUsed    *int       `json:"trial_days_used"`
	CreatedAt        time.Time  `json:"created_at"`
}

// WorkspaceUser 工作空间用户信息
type WorkspaceUser struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	Nickname     *string    `json:"nickname"`
	Email        *string    `json:"email"`
	Role         string     `json:"role"`
	Status       string     `json:"status"`
	JoinedAt     time.Time  `json:"joined_at"`
	LastActiveAt *time.Time `json:"last_active_at"`
}

// PaginatedResponse 分页响应结构
type PaginatedResponse[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// GetOrganizations 获取组织列表（只读）
func (s *OrganizationService) GetOrganizations(page, pageSize int, search string) (*PaginatedResponse[OrganizationDetail], error) {
	offset := (page - 1) * pageSize

	query := s.dbManager.LightAdminDB.Table("auth_organizations").
		Select(`auth_organizations.id, auth_organizations.name, auth_organizations.owner_id,
			auth_organizations.description, auth_organizations.created_at, auth_organizations.updated_at,
			COUNT(DISTINCT auo.user_id) as user_count,
			COUNT(DISTINCT aw.id) as workspace_count,
			COUNT(DISTINCT CASE WHEN su.status IN ('active','trial') THEN su.id END) as subscription_count,
			COUNT(DISTINCT CASE WHEN su.status = 'active' THEN su.id END) as active_subscription_count,
			MIN(CASE WHEN su.status IN ('active','trial') THEN su.end_date END) as subscription_end_date`).
		Joins("LEFT JOIN auth_user_organization auo ON auth_organizations.id = auo.organization_id").
		Joins("LEFT JOIN auth_workspaces aw ON auth_organizations.id = aw.organization_id").
		Joins("LEFT JOIN subscription_users su ON su.organization_id = auth_organizations.id::text").
		Group("auth_organizations.id")
	if search != "" {
		query = query.Where("auth_organizations.name ILIKE ? OR auth_organizations.description ILIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	var organizations []OrganizationDetail
	err := query.Order("auth_organizations.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Scan(&organizations).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get organizations: %w", err)
	}

	// 计算订阅到期状态和天数
	now := time.Now()
	for i := range organizations {
		s.calculateSubscriptionStatus(&organizations[i], now)
	}

	// 获取总数
	var total int64
	countQuery := s.dbManager.LightAdminDB.Table("auth_organizations")
	if search != "" {
		countQuery = countQuery.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	countQuery.Count(&total)

	return &PaginatedResponse[OrganizationDetail]{
		Data:       organizations,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}, nil
}

// GetOrganizationByID 根据ID获取组织详细信息（只读）
func (s *OrganizationService) GetOrganizationByID(orgID string) (*OrganizationDetail, error) {
	// 验证UUID格式
	if _, err := uuid.Parse(orgID); err != nil {
		return nil, fmt.Errorf("invalid organization ID format")
	}

	var org OrganizationDetail
	err := s.dbManager.LightAdminDB.Table("auth_organizations").
		Select(`auth_organizations.id, auth_organizations.name, auth_organizations.owner_id,
			auth_organizations.description, auth_organizations.created_at, auth_organizations.updated_at,
			COUNT(DISTINCT auo.user_id) as user_count,
			COUNT(DISTINCT aw.id) as workspace_count,
			COUNT(DISTINCT CASE WHEN su.status IN ('active','trial') THEN su.id END) as subscription_count,
			COUNT(DISTINCT CASE WHEN su.status = 'active' THEN su.id END) as active_subscription_count,
			MIN(CASE WHEN su.status IN ('active','trial') THEN su.end_date END) as subscription_end_date`).
		Joins("LEFT JOIN auth_user_organization auo ON auth_organizations.id = auo.organization_id").
		Joins("LEFT JOIN auth_workspaces aw ON auth_organizations.id = aw.organization_id").
		Joins("LEFT JOIN subscription_users su ON su.organization_id = auth_organizations.id::text").
		Where("auth_organizations.id = ?", orgID).
		Group("auth_organizations.id").
		First(&org).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("organization not found")
		}
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	// 获取活跃用户数（30天内有活动）
	thirtyDaysAgo := time.Now().AddDate(-30, 0, 0)
	s.dbManager.LightAdminDB.Table("auth_user_organization").
		Where("organization_id = ? AND created_at > ?",
			orgID, thirtyDaysAgo).Count(&org.ActiveUsers)

	// 设置存储使用情况（示例数据）
	org.StorageUsage = 0

	// 计算订阅到期状态和天数
	s.calculateSubscriptionStatus(&org, time.Now())

	return &org, nil
}

// GetOrganizationUsers 获取组织用户列表（只读）
func (s *OrganizationService) GetOrganizationUsers(orgID string, page, pageSize int, search string) (*PaginatedResponse[OrganizationUser], error) {
	offset := (page - 1) * pageSize

	query := s.dbManager.LightAdminDB.Table("auth_user_organization auo").
		Select(`au.id, au.username, au.nickname, au.email, au.avatar_url,
			auo.created_at, auo.updated_at`).
		Joins("INNER JOIN auth_users au ON auo.user_id = au.id").
		Where("auo.organization_id = ?", orgID)

	if search != "" {
		query = query.Where("au.username ILIKE ? OR au.nickname ILIKE ? OR au.email ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	var users []OrganizationUser
	err := query.Order("auo.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Scan(&users).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get organization users: %w", err)
	}

	// 获取总数
	var total int64
	countQuery := s.dbManager.LightAdminDB.Table("auth_user_organization auo").
		Joins("INNER JOIN auth_users au ON auo.user_id = au.id").
		Where("auo.organization_id = ?", orgID)
	if search != "" {
		countQuery = countQuery.Where("au.username ILIKE ? OR au.nickname ILIKE ? OR au.email ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	countQuery.Count(&total)

	return &PaginatedResponse[OrganizationUser]{
		Data:       users,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}, nil
}

// GetOrganizationWorkspaces 获取组织工作空间列表（只读）
func (s *OrganizationService) GetOrganizationWorkspaces(orgID string, page, pageSize int, search string) (*PaginatedResponse[OrganizationWorkspace], error) {
	offset := (page - 1) * pageSize

	query := s.dbManager.LightAdminDB.Table("auth_workspaces").
		Select(`id, name, status, description, created_at, updated_at,
			(SELECT COUNT(*) FROM auth_user_workspace WHERE workspace_id = auth_workspaces.id) as user_count`).
		Where("organization_id = ?", orgID)

	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	var workspaces []OrganizationWorkspace
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Scan(&workspaces).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get organization workspaces: %w", err)
	}

	// 获取总数
	var total int64
	countQuery := s.dbManager.LightAdminDB.Table("auth_workspaces").
		Where("organization_id = ?", orgID)
	if search != "" {
		countQuery = countQuery.Where("name ILIKE ? OR description ILIKE ?",
			"%"+search+"%", "%"+search+"%")
	}
	countQuery.Count(&total)

	return &PaginatedResponse[OrganizationWorkspace]{
		Data:       workspaces,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}, nil
}

// GetOrganizationSubscriptions 获取组织订阅信息（只读）
func (s *OrganizationService) GetOrganizationSubscriptions(orgID string, page, pageSize int) (*PaginatedResponse[OrganizationSubscription], error) {
	offset := (page - 1) * pageSize

	query := s.dbManager.LightAdminDB.Table("subscription_users su").
		Select(`su.id, su.user_id, au.username, au.email as user_email,
			su.plan_id, sp.tier_name as plan_name,
			CASE
				WHEN su.billing_cycle = 'monthly' THEN sp.pricing_monthly
				WHEN su.billing_cycle = 'quarterly' THEN sp.pricing_quarterly
				WHEN su.billing_cycle = 'yearly' THEN sp.pricing_yearly
				ELSE 0
			END as plan_pricing,
			su.status, su.billing_cycle, su.start_date, su.end_date,
			su.payment_method, su.last_billed_at, su.trial_days_used, su.created_at`).
		Joins("LEFT JOIN subscription_plans sp ON su.plan_id = sp.id").
		Joins("LEFT JOIN auth_users au ON su.user_id = au.id").
		Where("su.organization_id = ?", orgID)

	var subscriptions []OrganizationSubscription
	err := query.Order("su.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Scan(&subscriptions).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get organization subscriptions: %w", err)
	}

	// 计算每个订阅的到期天数
	now := time.Now()
	for i := range subscriptions {
		if subscriptions[i].EndDate != nil {
			days := int(subscriptions[i].EndDate.Sub(now).Hours() / 24)
			subscriptions[i].DaysUntilExpiry = &days
		}
	}

	// 获取总数
	var total int64
	s.dbManager.LightAdminDB.Table("subscription_users").
		Where("organization_id = ?", orgID).Count(&total)

	return &PaginatedResponse[OrganizationSubscription]{
		Data:       subscriptions,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}, nil
}

// GetOrganizationMetrics 获取组织指标统计（只读）
func (s *OrganizationService) GetOrganizationMetrics(orgID string) (*OrganizationDetail, error) {
	return s.GetOrganizationByID(orgID)
}

// GetWorkspaceUsers 获取工作空间用户列表（只读）
func (s *OrganizationService) GetWorkspaceUsers(workspaceID string, page, pageSize int) (*PaginatedResponse[WorkspaceUser], error) {
	offset := (page - 1) * pageSize

	query := s.dbManager.LightAdminDB.Table("auth_user_workspace auw").
		Select(`au.id, au.username, au.nickname, au.email, au.avatar_url,
			auw.role, auw.status, auw.joined_at, auw.last_active_at`).
		Joins("INNER JOIN auth_users au ON auw.user_id = au.id").
		Where("auw.workspace_id = ?", workspaceID)

	var users []WorkspaceUser
	err := query.Order("auw.joined_at DESC").
		Offset(offset).
		Limit(pageSize).
		Scan(&users).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get workspace users: %w", err)
	}

	// 获取总数
	var total int64
	s.dbManager.LightAdminDB.Table("auth_user_workspace").
		Where("workspace_id = ?", workspaceID).Count(&total)

	return &PaginatedResponse[WorkspaceUser]{
		Data:       users,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}, nil
}

// CreateOrganization - 不支持，因为light_admin表为只读
func (s *OrganizationService) CreateOrganization(name, ownerID string, description *string, createdBy string) (*models.AuthOrganization, error) {
	return nil, fmt.Errorf("create operation not supported for light_admin tables (read-only mode)")
}

// UpdateOrganization - 不支持，因为light_admin表为只读
func (s *OrganizationService) UpdateOrganization(orgID string, name, description *string, updatedBy string) (*models.AuthOrganization, error) {
	return nil, fmt.Errorf("update operation not supported for light_admin tables (read-only mode)")
}

// DeleteOrganization - 不支持，因为light_admin表为只读
func (s *OrganizationService) DeleteOrganization(orgID string) error {
	return fmt.Errorf("delete operation not supported for light_admin tables (read-only mode)")
}

// AddUserToOrganization - 不支持，因为light_admin表为只读
func (s *OrganizationService) AddUserToOrganization(orgID, userID, role, addedBy string) error {
	return fmt.Errorf("add user operation not supported for light_admin tables (read-only mode)")
}

// UpdateUserInOrganization - 不支持，因为light_admin表为只读
func (s *OrganizationService) UpdateUserInOrganization(orgID, userID, role, updatedBy string) error {
	return fmt.Errorf("update user operation not supported for light_admin tables (read-only mode)")
}

// RemoveUserFromOrganization - 不支持，因为light_admin表为只读
func (s *OrganizationService) RemoveUserFromOrganization(orgID, userID string) error {
	return fmt.Errorf("remove user operation not supported for light_admin tables (read-only mode)")
}

// calculateSubscriptionStatus 计算组织订阅状态和到期天数
func (s *OrganizationService) calculateSubscriptionStatus(org *OrganizationDetail, now time.Time) {
	if org.SubscriptionEndDate == nil || org.ActiveSubscriptionCount == 0 {
		org.SubscriptionStatus = "none"
		org.DaysUntilExpiration = nil
		return
	}

	// 计算距离到期天数
	days := int(org.SubscriptionEndDate.Sub(now).Hours() / 24)
	org.DaysUntilExpiration = &days

	// 判断订阅状态
	if days < 0 {
		org.SubscriptionStatus = "expired"
	} else if days <= 7 {
		org.SubscriptionStatus = "expiring_soon" // 7天内到期
	} else {
		org.SubscriptionStatus = "active"
	}
}

// SendExpiryReminder 发送订阅到期提醒邮件（预留接口）
func (s *OrganizationService) SendExpiryReminder(orgID string) error {
	// 验证组织存在
	org, err := s.GetOrganizationByID(orgID)
	if err != nil {
		return err
	}

	// TODO: 实现邮件发送逻辑
	// 1. 获取组织所有者和管理员邮箱
	// 2. 构造邮件内容（订阅到期时间、剩余天数等）
	// 3. 调用邮件服务发送

	// 当前仅记录日志
	fmt.Printf("Reminder email would be sent to organization: %s (ID: %s)\n", org.Name, org.ID)
	fmt.Printf("Subscription status: %s, Days until expiration: %v\n",
		org.SubscriptionStatus, org.DaysUntilExpiration)

	return nil
}
