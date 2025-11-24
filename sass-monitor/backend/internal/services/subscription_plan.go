package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sass-monitor/internal/database"
	"sass-monitor/internal/models"
)

// SubscriptionPlanService 订阅计划管理服务（支持增删查改）
type SubscriptionPlanService struct {
	dbManager *database.DatabaseManager
}

func NewSubscriptionPlanService(dbManager *database.DatabaseManager) *SubscriptionPlanService {
	return &SubscriptionPlanService{
		dbManager: dbManager,
	}
}

// CreatePlanRequest 创建订阅计划请求
type CreatePlanRequest struct {
	TierName               string  `json:"tier_name" binding:"required"`
	PricingMonthly         float64 `json:"pricing_monthly"`
	PricingQuarterly       float64 `json:"pricing_quarterly"`
	PricingYearly          float64 `json:"pricing_yearly"`
	Limits                 string  `json:"limits"`
	Features               *string `json:"features"`
	TargetUsers            *string `json:"target_users"`
	UpgradePath            *string `json:"upgrade_path"`
	IsCustom               *bool   `json:"is_custom"`
	DefaultFlowPackage     *string `json:"default_flow_package"`
	IsActive               bool    `json:"is_active"`
	StripePriceIDMonthly   *string `json:"stripe_price_id_monthly"`
	StripePriceIDQuarterly *string `json:"stripe_price_id_quarterly"`
	StripePriceIDYearly    *string `json:"stripe_price_id_yearly"`
}

// UpdatePlanRequest 更新订阅计划请求
type UpdatePlanRequest struct {
	TierName               string   `json:"tier_name"`
	PricingMonthly         *float64 `json:"pricing_monthly"`
	PricingQuarterly       *float64 `json:"pricing_quarterly"`
	PricingYearly          *float64 `json:"pricing_yearly"`
	Limits                 *string  `json:"limits"`
	Features               *string  `json:"features"`
	TargetUsers            *string  `json:"target_users"`
	UpgradePath            *string  `json:"upgrade_path"`
	IsCustom               *bool    `json:"is_custom"`
	DefaultFlowPackage     *string  `json:"default_flow_package"`
	IsActive               *bool    `json:"is_active"`
	StripePriceIDMonthly   *string  `json:"stripe_price_id_monthly"`
	StripePriceIDQuarterly *string  `json:"stripe_price_id_quarterly"`
	StripePriceIDYearly    *string  `json:"stripe_price_id_yearly"`
}

// GetSubscriptionPlans 获取订阅计划列表
func (s *SubscriptionPlanService) GetSubscriptionPlans(page, pageSize int, search string) (*PaginatedResponse[models.SubscriptionPlan], error) {
	offset := (page - 1) * pageSize

	query := s.dbManager.LightAdminDB.Model(&models.SubscriptionPlan{})

	if search != "" {
		query = query.Where("tier_name ILIKE ? OR target_users ILIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	var plans []models.SubscriptionPlan
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&plans).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get subscription plans: %w", err)
	}

	// 获取总数
	var total int64
	countQuery := s.dbManager.LightAdminDB.Model(&models.SubscriptionPlan{})
	if search != "" {
		countQuery = countQuery.Where("tier_name ILIKE ? OR target_users ILIKE ?",
			"%"+search+"%", "%"+search+"%")
	}
	countQuery.Count(&total)

	return &PaginatedResponse[models.SubscriptionPlan]{
		Data:       plans,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}, nil
}

// GetSubscriptionPlanByID 根据ID获取订阅计划
func (s *SubscriptionPlanService) GetSubscriptionPlanByID(planID string) (*models.SubscriptionPlan, error) {
	// 验证UUID格式
	if _, err := uuid.Parse(planID); err != nil {
		return nil, fmt.Errorf("invalid plan ID format")
	}

	var plan models.SubscriptionPlan
	err := s.dbManager.LightAdminDB.Where("id = ?", planID).First(&plan).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("subscription plan not found")
		}
		return nil, fmt.Errorf("failed to get subscription plan: %w", err)
	}

	return &plan, nil
}

// CreateSubscriptionPlan 创建订阅计划
func (s *SubscriptionPlanService) CreateSubscriptionPlan(req *CreatePlanRequest) (*models.SubscriptionPlan, error) {
	// 检查是否已存在相同名称的计划
	var existingPlan models.SubscriptionPlan
	err := s.dbManager.LightAdminDB.Where("tier_name = ?", req.TierName).First(&existingPlan).Error
	if err == nil {
		return nil, fmt.Errorf("subscription plan with tier_name '%s' already exists", req.TierName)
	}
	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check existing plan: %w", err)
	}

	if req.IsCustom == nil {
		req.IsCustom = new(bool)
		*req.IsCustom = false
	}

	plan := &models.SubscriptionPlan{
		ID:                     uuid.New(),
		TierName:               req.TierName,
		PricingMonthly:         req.PricingMonthly,
		PricingQuarterly:       req.PricingQuarterly,
		PricingYearly:          req.PricingYearly,
		Limits:                 req.Limits,
		Features:               req.Features,
		TargetUsers:            req.TargetUsers,
		UpgradePath:            req.UpgradePath,
		IsCustom:               req.IsCustom,
		DefaultFlowPackage:     req.DefaultFlowPackage,
		IsActive:               req.IsActive,
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
		StripePriceIDMonthly:   req.StripePriceIDMonthly,
		StripePriceIDQuarterly: req.StripePriceIDQuarterly,
		StripePriceIDYearly:    req.StripePriceIDYearly,
	}

	err = s.dbManager.LightAdminDB.Create(plan).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription plan: %w", err)
	}

	return plan, nil
}

// UpdateSubscriptionPlan 更新订阅计划
func (s *SubscriptionPlanService) UpdateSubscriptionPlan(planID string, req *UpdatePlanRequest) (*models.SubscriptionPlan, error) {
	// 验证UUID格式
	if _, err := uuid.Parse(planID); err != nil {
		return nil, fmt.Errorf("invalid plan ID format")
	}

	// 检查计划是否存在
	var plan models.SubscriptionPlan
	err := s.dbManager.LightAdminDB.Where("id = ?", planID).First(&plan).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("subscription plan not found")
		}
		return nil, fmt.Errorf("failed to get subscription plan: %w", err)
	}

	// 如果更新了tier_name，检查是否重复
	if req.TierName != "" && req.TierName != plan.TierName {
		var existingPlan models.SubscriptionPlan
		err = s.dbManager.LightAdminDB.Where("tier_name = ? AND id != ?", req.TierName, planID).First(&existingPlan).Error
		if err == nil {
			return nil, fmt.Errorf("subscription plan with tier_name '%s' already exists", req.TierName)
		}
		if err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to check existing plan: %w", err)
		}
		plan.TierName = req.TierName
	}

	// 更新其他字段
	updateData := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if req.PricingMonthly != nil {
		updateData["pricing_monthly"] = *req.PricingMonthly
	}
	if req.PricingQuarterly != nil {
		updateData["pricing_quarterly"] = *req.PricingQuarterly
	}
	if req.PricingYearly != nil {
		updateData["pricing_yearly"] = *req.PricingYearly
	}
	if req.Limits != nil {
		updateData["limits"] = *req.Limits
	}
	if req.Features != nil {
		updateData["features"] = *req.Features
	}
	if req.TargetUsers != nil {
		updateData["target_users"] = *req.TargetUsers
	}
	if req.UpgradePath != nil {
		updateData["upgrade_path"] = *req.UpgradePath
	}
	if req.IsCustom != nil {
		updateData["is_custom"] = *req.IsCustom
	}
	if req.DefaultFlowPackage != nil {
		updateData["default_flow_package"] = *req.DefaultFlowPackage
	}
	if req.IsActive != nil {
		updateData["is_active"] = *req.IsActive
	}
	if req.StripePriceIDMonthly != nil {
		updateData["stripe_price_id_monthly"] = *req.StripePriceIDMonthly
	}
	if req.StripePriceIDQuarterly != nil {
		updateData["stripe_price_id_quarterly"] = *req.StripePriceIDQuarterly
	}
	if req.StripePriceIDYearly != nil {
		updateData["stripe_price_id_yearly"] = *req.StripePriceIDYearly
	}

	err = s.dbManager.LightAdminDB.Model(&plan).Updates(updateData).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update subscription plan: %w", err)
	}

	// 重新获取更新后的计划
	err = s.dbManager.LightAdminDB.Where("id = ?", planID).First(&plan).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get updated subscription plan: %w", err)
	}

	return &plan, nil
}

// DeleteSubscriptionPlan 删除订阅计划
func (s *SubscriptionPlanService) DeleteSubscriptionPlan(planID string) error {
	// 验证UUID格式
	if _, err := uuid.Parse(planID); err != nil {
		return fmt.Errorf("invalid plan ID format")
	}

	// 检查是否有订阅用户使用此计划
	var subscriptionCount int64
	err := s.dbManager.LightAdminDB.Model(&models.SubscriptionUser{}).
		Where("plan_id = ?", planID).Count(&subscriptionCount).Error
	if err != nil {
		return fmt.Errorf("failed to check subscription users: %w", err)
	}

	if subscriptionCount > 0 {
		return fmt.Errorf("cannot delete subscription plan: %d users are currently subscribed", subscriptionCount)
	}

	// 删除计划
	result := s.dbManager.LightAdminDB.Where("id = ?", planID).Delete(&models.SubscriptionPlan{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete subscription plan: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("subscription plan not found")
	}

	return nil
}

// GetActiveSubscriptionPlans 获取活跃的订阅计划
func (s *SubscriptionPlanService) GetActiveSubscriptionPlans() ([]models.SubscriptionPlan, error) {
	var plans []models.SubscriptionPlan
	err := s.dbManager.LightAdminDB.Where("is_active = ?", true).
		Order("pricing_monthly ASC").Find(&plans).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get active subscription plans: %w", err)
	}
	return plans, nil
}

// GetSubscriptionPlansByPricingRange 根据价格范围获取订阅计划
func (s *SubscriptionPlanService) GetSubscriptionPlansByPricingRange(minPrice, maxPrice float64) ([]models.SubscriptionPlan, error) {
	var plans []models.SubscriptionPlan
	err := s.dbManager.LightAdminDB.Where("is_active = ? AND pricing_monthly BETWEEN ? AND ?", true, minPrice, maxPrice).
		Order("pricing_monthly ASC").Find(&plans).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription plans by pricing range: %w", err)
	}
	return plans, nil
}
