package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sass-monitor/internal/database"
	"sass-monitor/internal/models"
)

// UserService 用户管理服务（只读模式）
type UserService struct {
	dbManager *database.DatabaseManager
}

func NewUserService(dbManager *database.DatabaseManager) *UserService {
	return &UserService{
		dbManager: dbManager,
	}
}

// UserDetail 用户详细信息
type UserDetail struct {
	ID              string     `json:"id"`
	Username        string     `json:"username"`
	Nickname        *string    `json:"nickname"`
	Email           *string    `json:"email"`
	AvatarURL       string     `json:"avatar_url"`
	EmailVerified   *bool      `json:"email_verified"`
	OAuthProvider   string     `json:"oauth_provider"`
	ClerkUserID     string     `json:"clerk_user_id"`
	OrganizationCount int64    `json:"organization_count"`
	WorkspaceCount  int64      `json:"workspace_count"`
	SubscriptionCount int64    `json:"subscription_count"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

// UserOrganization 用户所属组织信息
type UserOrganization struct {
	OrganizationID   string     `json:"organization_id"`
	OrganizationName string     `json:"organization_name"`
	Description      *string    `json:"description"`
	JoinedAt         time.Time  `json:"joined_at"`
	CreatedAt        time.Time  `json:"created_at"`
}

// UserWorkspace 用户所属工作空间信息
type UserWorkspace struct {
	WorkspaceID      string     `json:"workspace_id"`
	WorkspaceName    string     `json:"workspace_name"`
	OrganizationID   string     `json:"organization_id"`
	OrganizationName string     `json:"organization_name"`
	UserStatus       string     `json:"user_status"`
	CreatedAt        time.Time  `json:"created_at"`
}

// UserSubscription 用户订阅信息
type UserSubscription struct {
	SubscriptionID   string     `json:"subscription_id"`
	PlanID           string     `json:"plan_id"`
	PlanName         string     `json:"plan_name"`
	OrganizationID   string     `json:"organization_id"`
	OrganizationName string     `json:"organization_name"`
	Status           string     `json:"status"`
	BillingCycle     string     `json:"billing_cycle"`
	StartDate        time.Time  `json:"start_date"`
	EndDate          *time.Time `json:"end_date"`
	TrialDaysUsed    *int       `json:"trial_days_used"`
}

// GetUsers 获取用户列表（只读）
func (s *UserService) GetUsers(page, pageSize int, search string) (*PaginatedResponse[UserDetail], error) {
	offset := (page - 1) * pageSize

	query := s.dbManager.LightAdminDB.Table("auth_users").
		Select(`auth_users.id, auth_users.username, auth_users.nickname, auth_users.email,
			auth_users.avatar_url, auth_users.email_verified, auth_users.oauth_provider,
			auth_users.clerk_user_id, auth_users.created_at, auth_users.updated_at,
			COUNT(DISTINCT auo.organization_id) as organization_count,
			COUNT(DISTINCT auw.workspace_id) as workspace_count,
			COUNT(DISTINCT CASE WHEN su.status IN ('active','trial') THEN su.id END) as subscription_count`).
		Joins("LEFT JOIN auth_user_organization auo ON auth_users.id = auo.user_id").
		Joins("LEFT JOIN auth_user_workspace auw ON auth_users.id = auw.user_id").
		Joins("LEFT JOIN subscription_users su ON su.user_id = auth_users.id").
		Group("auth_users.id")

	if search != "" {
		query = query.Where("auth_users.username ILIKE ? OR auth_users.nickname ILIKE ? OR auth_users.email ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	var users []UserDetail
	err := query.Order("auth_users.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Scan(&users).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	// 获取总数
	var total int64
	countQuery := s.dbManager.LightAdminDB.Table("auth_users")
	if search != "" {
		countQuery = countQuery.Where("username ILIKE ? OR nickname ILIKE ? OR email ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	countQuery.Count(&total)

	return &PaginatedResponse[UserDetail]{
		Data:       users,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}, nil
}

// GetUserByID 根据ID获取用户详细信息（只读）
func (s *UserService) GetUserByID(userID string) (*UserDetail, error) {
	// 验证UUID格式
	if _, err := uuid.Parse(userID); err != nil {
		return nil, fmt.Errorf("invalid user ID format")
	}

	var user UserDetail
	err := s.dbManager.LightAdminDB.Table("auth_users").
		Select(`auth_users.id, auth_users.username, auth_users.nickname, auth_users.email,
			auth_users.avatar_url, auth_users.email_verified, auth_users.oauth_provider,
			auth_users.clerk_user_id, auth_users.created_at, auth_users.updated_at,
			COUNT(DISTINCT auo.organization_id) as organization_count,
			COUNT(DISTINCT auw.workspace_id) as workspace_count,
			COUNT(DISTINCT CASE WHEN su.status IN ('active','trial') THEN su.id END) as subscription_count`).
		Joins("LEFT JOIN auth_user_organization auo ON auth_users.id = auo.user_id").
		Joins("LEFT JOIN auth_user_workspace auw ON auth_users.id = auw.user_id").
		Joins("LEFT JOIN subscription_users su ON su.user_id = auth_users.id").
		Where("auth_users.id = ?", userID).
		Group("auth_users.id").
		First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetUserOrganizations 获取用户所属组织列表（只读）
func (s *UserService) GetUserOrganizations(userID string, page, pageSize int) (*PaginatedResponse[UserOrganization], error) {
	offset := (page - 1) * pageSize

	query := s.dbManager.LightAdminDB.Table("auth_user_organization auo").
		Select(`ao.id as organization_id, ao.name as organization_name, ao.description,
			auo.created_at as joined_at, auo.created_at`).
		Joins("INNER JOIN auth_organizations ao ON auo.organization_id = ao.id").
		Where("auo.user_id = ?", userID)

	var organizations []UserOrganization
	err := query.Order("auo.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Scan(&organizations).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get user organizations: %w", err)
	}

	// 获取总数
	var total int64
	s.dbManager.LightAdminDB.Table("auth_user_organization").
		Where("user_id = ?", userID).Count(&total)

	return &PaginatedResponse[UserOrganization]{
		Data:       organizations,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}, nil
}

// GetUserWorkspaces 获取用户所属工作空间列表（只读）
func (s *UserService) GetUserWorkspaces(userID string, page, pageSize int) (*PaginatedResponse[UserWorkspace], error) {
	offset := (page - 1) * pageSize

	query := s.dbManager.LightAdminDB.Table("auth_user_workspace auw").
		Select(`aw.id as workspace_id, aw.name as workspace_name,
			ao.id as organization_id, ao.name as organization_name,
			auw.user_status, auw.created_at`).
		Joins("INNER JOIN auth_workspaces aw ON auw.workspace_id = aw.id").
		Joins("INNER JOIN auth_organizations ao ON aw.organization_id = ao.id").
		Where("auw.user_id = ?", userID)

	var workspaces []UserWorkspace
	err := query.Order("auw.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Scan(&workspaces).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get user workspaces: %w", err)
	}

	// 获取总数
	var total int64
	s.dbManager.LightAdminDB.Table("auth_user_workspace").
		Where("user_id = ?", userID).Count(&total)

	return &PaginatedResponse[UserWorkspace]{
		Data:       workspaces,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}, nil
}

// GetUserSubscriptions 获取用户订阅列表（只读）
func (s *UserService) GetUserSubscriptions(userID string, page, pageSize int) (*PaginatedResponse[UserSubscription], error) {
	offset := (page - 1) * pageSize

	query := s.dbManager.LightAdminDB.Table("subscription_users su").
		Select(`su.id as subscription_id, su.plan_id, sp.tier_name as plan_name,
			su.organization_id, ao.name as organization_name,
			su.status, su.billing_cycle, su.start_date, su.end_date, su.trial_days_used`).
		Joins("LEFT JOIN subscription_plans sp ON su.plan_id = sp.id").
		Joins("LEFT JOIN auth_organizations ao ON su.organization_id = ao.id::text").
		Where("su.user_id = ?", userID)

	var subscriptions []UserSubscription
	err := query.Order("su.start_date DESC").
		Offset(offset).
		Limit(pageSize).
		Scan(&subscriptions).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get user subscriptions: %w", err)
	}

	// 获取总数
	var total int64
	s.dbManager.LightAdminDB.Table("subscription_users").
		Where("user_id = ?", userID).Count(&total)

	return &PaginatedResponse[UserSubscription]{
		Data:       subscriptions,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}, nil
}

// CreateUser - 不支持，因为light_admin表为只读
func (s *UserService) CreateUser(username, email string) (*models.AuthUser, error) {
	return nil, fmt.Errorf("create operation not supported for light_admin tables (read-only mode)")
}

// UpdateUser - 不支持，因为light_admin表为只读
func (s *UserService) UpdateUser(userID string, nickname, email *string) (*models.AuthUser, error) {
	return nil, fmt.Errorf("update operation not supported for light_admin tables (read-only mode)")
}

// DeleteUser - 不支持，因为light_admin表为只读
func (s *UserService) DeleteUser(userID string) error {
	return fmt.Errorf("delete operation not supported for light_admin tables (read-only mode)")
}
