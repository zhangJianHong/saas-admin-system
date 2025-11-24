package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"sass-monitor/internal/services"
	"sass-monitor/pkg/config"
)

type OrganizationHandler struct {
	orgService *services.OrganizationService
	config     *config.Config
}

func NewOrganizationHandler(orgService *services.OrganizationService, cfg *config.Config) *OrganizationHandler {
	return &OrganizationHandler{
		orgService: orgService,
		config:     cfg,
	}
}

// GetOrganizations 获取组织列表
func (h *OrganizationHandler) GetOrganizations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")

	result, err := h.orgService.GetOrganizations(page, pageSize, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get organizations: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetOrganizationByID 根据ID获取组织详细信息
func (h *OrganizationHandler) GetOrganizationByID(c *gin.Context) {
	orgID := c.Param("id")

	org, err := h.orgService.GetOrganizationByID(orgID)
	if err != nil {
		if err.Error() == "invalid organization ID format" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid organization ID format",
			})
			return
		}
		if err.Error() == "organization not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Organization not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get organization: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, org)
}

// CreateOrganizationRequest 创建组织请求结构
type CreateOrganizationRequest struct {
	Name        string  `json:"name" binding:"required"`
	OwnerID     string  `json:"owner_id" binding:"required"`
	Description *string `json:"description"`
}

// CreateOrganization 创建新组织（不支持）
func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{
		"error": "Create operation not supported for light_admin tables (read-only mode)",
	})
}

// UpdateOrganizationRequest 更新组织请求结构
type UpdateOrganizationRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

// UpdateOrganization 更新组织信息（不支持）
func (h *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{
		"error": "Update operation not supported for light_admin tables (read-only mode)",
	})
}

// DeleteOrganization 删除组织（不支持）
func (h *OrganizationHandler) DeleteOrganization(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{
		"error": "Delete operation not supported for light_admin tables (read-only mode)",
	})
}

// GetOrganizationUsers 获取组织用户列表
func (h *OrganizationHandler) GetOrganizationUsers(c *gin.Context) {
	orgID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")

	// 验证UUID格式
	if _, err := uuid.Parse(orgID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid organization ID format",
		})
		return
	}

	result, err := h.orgService.GetOrganizationUsers(orgID, page, pageSize, search)
	if err != nil {
		if err.Error() == "organization not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Organization not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get organization users: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// AddUserToOrganizationRequest 添加用户到组织请求结构
type AddUserToOrganizationRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required"`
}

// AddUserToOrganization 添加用户到组织（不支持）
func (h *OrganizationHandler) AddUserToOrganization(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{
		"error": "Add user operation not supported for light_admin tables (read-only mode)",
	})
}

// UpdateUserInOrganizationRequest 更新组织用户角色请求结构
type UpdateUserInOrganizationRequest struct {
	Role string `json:"role" binding:"required"`
}

// UpdateUserInOrganization 更新组织中的用户角色（不支持）
func (h *OrganizationHandler) UpdateUserInOrganization(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{
		"error": "Update user operation not supported for light_admin tables (read-only mode)",
	})
}

// RemoveUserFromOrganization 从组织中移除用户（不支持）
func (h *OrganizationHandler) RemoveUserFromOrganization(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{
		"error": "Remove user operation not supported for light_admin tables (read-only mode)",
	})
}

// GetOrganizationWorkspaces 获取组织工作空间列表
func (h *OrganizationHandler) GetOrganizationWorkspaces(c *gin.Context) {
	orgID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")

	// 验证UUID格式
	if _, err := uuid.Parse(orgID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid organization ID format",
		})
		return
	}

	result, err := h.orgService.GetOrganizationWorkspaces(orgID, page, pageSize, search)
	if err != nil {
		if err.Error() == "organization not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Organization not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get organization workspaces: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetOrganizationSubscriptions 获取组织订阅信息
func (h *OrganizationHandler) GetOrganizationSubscriptions(c *gin.Context) {
	orgID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 验证UUID格式
	if _, err := uuid.Parse(orgID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid organization ID format",
		})
		return
	}

	result, err := h.orgService.GetOrganizationSubscriptions(orgID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get organization subscriptions: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetOrganizationMetrics 获取组织指标统计
func (h *OrganizationHandler) GetOrganizationMetrics(c *gin.Context) {
	orgID := c.Param("id")

	// 验证UUID格式
	if _, err := uuid.Parse(orgID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid organization ID format",
		})
		return
	}

	metrics, err := h.orgService.GetOrganizationMetrics(orgID)
	if err != nil {
		if err.Error() == "organization not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Organization not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get organization metrics: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// SendExpiryReminder 发送订阅到期提醒邮件（预留）
func (h *OrganizationHandler) SendExpiryReminder(c *gin.Context) {
	orgID := c.Param("id")

	// 验证UUID格式
	if _, err := uuid.Parse(orgID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid organization ID format",
		})
		return
	}

	err := h.orgService.SendExpiryReminder(orgID)
	if err != nil {
		if err.Error() == "organization not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Organization not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send expiry reminder: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Expiry reminder sent successfully (placeholder)",
	})
}

// GetWorkspaceUsers 获取工作空间用户列表
func (h *OrganizationHandler) GetWorkspaceUsers(c *gin.Context) {
	workspaceID := c.Param("workspaceId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 验证UUID格式
	if _, err := uuid.Parse(workspaceID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid workspace ID format",
		})
		return
	}

	result, err := h.orgService.GetWorkspaceUsers(workspaceID, page, pageSize)
	if err != nil {
		if err.Error() == "workspace not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Workspace not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get workspace users: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}