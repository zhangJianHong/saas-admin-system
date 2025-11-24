package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"sass-monitor/internal/services"
	"sass-monitor/pkg/config"
)

type SubscriptionPlanHandler struct {
	subscriptionPlanService *services.SubscriptionPlanService
	config                  *config.Config
}

func NewSubscriptionPlanHandler(subscriptionPlanService *services.SubscriptionPlanService, cfg *config.Config) *SubscriptionPlanHandler {
	return &SubscriptionPlanHandler{
		subscriptionPlanService: subscriptionPlanService,
		config:                  cfg,
	}
}

// GetSubscriptionPlans 获取订阅计划列表
func (h *SubscriptionPlanHandler) GetSubscriptionPlans(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")

	result, err := h.subscriptionPlanService.GetSubscriptionPlans(page, pageSize, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get subscription plans: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetSubscriptionPlanByID 根据ID获取订阅计划
func (h *SubscriptionPlanHandler) GetSubscriptionPlanByID(c *gin.Context) {
	planID := c.Param("id")

	plan, err := h.subscriptionPlanService.GetSubscriptionPlanByID(planID)
	if err != nil {
		if err.Error() == "invalid plan ID format" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid plan ID format",
			})
			return
		}
		if err.Error() == "subscription plan not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Subscription plan not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get subscription plan: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, plan)
}

// CreateSubscriptionPlan 创建订阅计划
func (h *SubscriptionPlanHandler) CreateSubscriptionPlan(c *gin.Context) {
	var req services.CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format: " + err.Error(),
		})
		return

	}

	plan, err := h.subscriptionPlanService.CreateSubscriptionPlan(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create subscription plan: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, plan)
}

// UpdateSubscriptionPlan 更新订阅计划
func (h *SubscriptionPlanHandler) UpdateSubscriptionPlan(c *gin.Context) {
	planID := c.Param("id")

	// 验证UUID格式
	if _, err := uuid.Parse(planID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid plan ID format",
		})
		return
	}

	var req services.UpdatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}

	plan, err := h.subscriptionPlanService.UpdateSubscriptionPlan(planID, &req)
	if err != nil {
		if err.Error() == "subscription plan not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Subscription plan not found",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update subscription plan: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, plan)
}

// DeleteSubscriptionPlan 删除订阅计划
func (h *SubscriptionPlanHandler) DeleteSubscriptionPlan(c *gin.Context) {
	planID := c.Param("id")

	// 验证UUID格式
	if _, err := uuid.Parse(planID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid plan ID format",
		})
		return
	}

	err := h.subscriptionPlanService.DeleteSubscriptionPlan(planID)
	if err != nil {
		if err.Error() == "invalid plan ID format" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid plan ID format",
			})
			return
		}
		if err.Error() == "subscription plan not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Subscription plan not found",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete subscription plan: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Subscription plan deleted successfully",
	})
}

// GetActiveSubscriptionPlans 获取活跃的订阅计划
func (h *SubscriptionPlanHandler) GetActiveSubscriptionPlans(c *gin.Context) {
	plans, err := h.subscriptionPlanService.GetActiveSubscriptionPlans()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get active subscription plans: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"plans": plans,
	})
}

// GetSubscriptionPlansByPricingRange 根据价格范围获取订阅计划
func (h *SubscriptionPlanHandler) GetSubscriptionPlansByPricingRange(c *gin.Context) {
	minPrice, err := strconv.ParseFloat(c.DefaultQuery("min_price", "0"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid min_price parameter",
		})
		return
	}

	maxPrice, err := strconv.ParseFloat(c.DefaultQuery("max_price", "999999"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid max_price parameter",
		})
		return
	}

	plans, err := h.subscriptionPlanService.GetSubscriptionPlansByPricingRange(minPrice, maxPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get subscription plans by pricing range: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"plans": plans,
	})
}
