package models

import (
	"time"
	"github.com/google/uuid"
)

// AuthOrganization 对应light_admin.auth_organizations表
type AuthOrganization struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CreatedBy   string     `gorm:"size:50" json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedBy   *string    `gorm:"size:50" json:"updated_by"`
	UpdatedAt   *time.Time `json:"updated_at"`
	Name        string     `gorm:"size:100;not null" json:"name"`
	OwnerID     uuid.UUID  `gorm:"type:uuid;not null" json:"owner_id"`
	Description *string    `gorm:"size:255" json:"description"`
}

// AuthWorkspace 对应light_admin.auth_workspaces表
type AuthWorkspace struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CreatedBy    string     `gorm:"size:50" json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedBy    *string    `gorm:"size:50" json:"updated_by"`
	UpdatedAt    *time.Time `json:"updated_at"`
	Name         string     `gorm:"size:100;not null" json:"name"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null" json:"organization_id"`
	OwnerID      uuid.UUID  `gorm:"type:uuid;not null" json:"owner_id"`
	Status       string     `gorm:"size:20;not null" json:"status"`
	Description  *string    `gorm:"size:255" json:"description"`
	FeatureMenu  *string    `gorm:"type:jsonb" json:"feature_menu"`
}

// AuthUser 对应light_admin.auth_users表
type AuthUser struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	CreatedBy     string    `gorm:"size:50" json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedBy     *string   `gorm:"size:50" json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	Username      string    `gorm:"size:50;not null;uniqueIndex" json:"username"`
	PasswordHash  string    `gorm:"size:255;not null" json:"-"`
	Nickname      *string   `gorm:"size:50" json:"nickname"`
	Email         *string   `gorm:"size:255" json:"email"`
	ClerkUserID   string    `gorm:"size:255;default:''" json:"clerk_user_id"`
	OAuthProvider string    `gorm:"size:50;default:''" json:"oauth_provider"`
	AvatarURL     string    `gorm:"type:text;default:''" json:"avatar_url"`
	EmailVerified *bool     `json:"email_verified"`
}

// AuthUserOrganization 对应light_admin.auth_user_organization表（用户-组织关联表）
type AuthUserOrganization struct {
	UserID         uuid.UUID  `gorm:"type:uuid;primaryKey;not null" json:"user_id"`
	OrganizationID uuid.UUID  `gorm:"type:uuid;primaryKey;not null" json:"organization_id"`
	CreatedBy      string     `gorm:"size:50;not null" json:"created_by"`
	CreatedAt      time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedBy      *string    `gorm:"size:50" json:"updated_by"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

// AuthUserWorkspace 对应light_admin.auth_user_workspace表（用户-工作空间关联表）
type AuthUserWorkspace struct {
	UserID      uuid.UUID  `gorm:"type:uuid;primaryKey;not null" json:"user_id"`
	WorkspaceID uuid.UUID  `gorm:"type:uuid;primaryKey;not null" json:"workspace_id"`
	UserStatus  string     `gorm:"size:20;default:'active';not null" json:"user_status"`
	CreatedBy   string     `gorm:"size:50;not null" json:"created_by"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedBy   *string    `gorm:"size:50" json:"updated_by"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

// SubscriptionPlan 对应light_admin.subscription_plans表
type SubscriptionPlan struct {
	ID                         uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	TierName                   string     `gorm:"size:50;not null;uniqueIndex" json:"tier_name"`
	PricingMonthly             float64    `gorm:"type:numeric(10,2);default:0.00" json:"pricing_monthly"`
	PricingQuarterly           float64    `gorm:"type:numeric(10,2);default:0.00" json:"pricing_quarterly"`
	PricingYearly              float64    `gorm:"type:numeric(10,2);default:0.00" json:"pricing_yearly"`
	Limits                     string     `gorm:"type:jsonb;default:'{}'" json:"limits"`
	Features                   *string    `gorm:"type:jsonb" json:"features"`
	TargetUsers                *string    `gorm:"type:text" json:"target_users"`
	UpgradePath                *string    `gorm:"type:text" json:"upgrade_path"`
	IsCustom                   *bool      `json:"is_custom"`
	DefaultFlowPackage         *string    `gorm:"type:public.flow_package_type" json:"default_flow_package"`
	IsActive                   bool       `gorm:"default:true" json:"is_active"`
	CreatedAt                  time.Time  `json:"created_at"`
	UpdatedAt                  time.Time  `json:"updated_at"`
	StripePriceIDMonthly       *string    `gorm:"size:255" json:"stripe_price_id_monthly"`
	StripePriceIDQuarterly     *string    `gorm:"size:255" json:"stripe_price_id_quarterly"`
	StripePriceIDYearly        *string    `gorm:"size:255" json:"stripe_price_id_yearly"`
}

// SubscriptionUser 对应light_admin.subscription_users表
type SubscriptionUser struct {
	ID                  string     `gorm:"primary_key;default:gen_random_uuid()" json:"id"`
	UserID              uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	PlanID              uuid.UUID  `gorm:"type:uuid;not null" json:"plan_id"`
	Status              string     `gorm:"size:20;default:'active';not null" json:"status"`
	BillingCycle        string     `gorm:"size:20;default:'monthly';not null" json:"billing_cycle"`
	StartDate           time.Time  `gorm:"default:CURRENT_TIMESTAMP;not null" json:"start_date"`
	EndDate             *time.Time `json:"end_date"`
	PaymentMethod       *string    `gorm:"size:50" json:"payment_method"`
	LastBilledAt        *time.Time `json:"last_billed_at"`
	TrialDaysUsed       *int       `gorm:"default:0" json:"trial_days_used"`
	OrganizationID      string     `gorm:"not null" json:"organization_id"`
	Notes               *string    `gorm:"type:text" json:"notes"`
	CreatedAt           time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt           time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	StripeSessionID     *string    `gorm:"size:255" json:"stripe_session_id"`
	StripeSessionData   *string    `gorm:"type:jsonb" json:"stripe_session_data"`
}

// OrgUsage 对应light_admin.org_usage表
type OrgUsage struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID string    `gorm:"not null" json:"organization_id"`
	WorkspaceID    string    `gorm:"not null" json:"workspace_id"`
	Month          string    `gorm:"not null" json:"month"`
	Usage          string    `gorm:"type:jsonb;not null" json:"usage"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

// OrgBilling 对应light_admin.org_billing表
type OrgBilling struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrgID      *uuid.UUID `gorm:"type:uuid" json:"org_id"`
	Month      *string    `gorm:"size:7" json:"month"`
	UsageCount *int64     `json:"usage_count"`
	FreeQuota  *int64     `json:"free_quota"`
	Overage    *int64     `json:"overage"`
	Amount     *float64   `gorm:"type:numeric(10,2)" json:"amount"`
	Status     *string    `gorm:"size:20" json:"status"`
}

// Payment 对应light_admin.payments表
type Payment struct {
	ID            string     `gorm:"primary_key" json:"id"`
	CustomerID    string     `gorm:"not null" json:"customer_id"`
	SubscriptionID *string   `json:"subscription_id"`
	Amount        int64      `gorm:"not null" json:"amount"`
	Currency      string     `gorm:"size:10;not null" json:"currency"`
	Status        string     `gorm:"size:50;not null" json:"status"`
	StripeEventID *string    `gorm:"size:255" json:"stripe_event_id"`
	Metadata      *string    `gorm:"type:text" json:"metadata"`
	CreatedAt     time.Time  `gorm:"default:now()" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"default:now()" json:"updated_at"`
}

// 指定表名
func (AuthOrganization) TableName() string {
	return "auth_organizations"
}

func (AuthWorkspace) TableName() string {
	return "auth_workspaces"
}

func (AuthUser) TableName() string {
	return "auth_users"
}

func (AuthUserOrganization) TableName() string {
	return "auth_user_organization"
}

func (AuthUserWorkspace) TableName() string {
	return "auth_user_workspace"
}

func (SubscriptionPlan) TableName() string {
	return "subscription_plans"
}

func (SubscriptionUser) TableName() string {
	return "subscription_users"
}

func (OrgUsage) TableName() string {
	return "org_usage"
}

func (OrgBilling) TableName() string {
	return "org_billing"
}

func (Payment) TableName() string {
	return "payments"
}