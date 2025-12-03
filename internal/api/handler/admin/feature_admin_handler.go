package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"paper_ai/internal/domain/repository"
	apperrors "paper_ai/pkg/errors"
	"paper_ai/pkg/response"
)

// FeatureAdminHandler 功能管理处理器
type FeatureAdminHandler struct {
	userRepo repository.UserRepository
}

// NewFeatureAdminHandler 创建功能管理处理器
func NewFeatureAdminHandler(userRepo repository.UserRepository) *FeatureAdminHandler {
	return &FeatureAdminHandler{
		userRepo: userRepo,
	}
}

// EnableMultiVersionRequest 启用多版本请求
type EnableMultiVersionRequest struct {
	Quota int `json:"quota"` // 配额，0表示无限
}

// EnableMultiVersionForUser 为用户开通多版本功能
// @Summary 开通多版本功能
// @Description 为指定用户开通多版本润色功能
// @Tags admin
// @Accept json
// @Produce json
// @Param user_id path int true "用户ID"
// @Param request body EnableMultiVersionRequest true "配额信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/users/{user_id}/multi-version/enable [post]
func (h *FeatureAdminHandler) EnableMultiVersionForUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.NewInvalidParameterError("无效的用户ID"))
		return
	}

	var req EnableMultiVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// 获取用户
	user, err := h.userRepo.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 更新权限
	user.EnableMultiVersion = true
	user.MultiVersionQuota = req.Quota

	if err := h.userRepo.Update(c.Request.Context(), user); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "已为用户开通多版本润色功能",
		"user_id": userID,
		"quota":   req.Quota,
	})
}

// DisableMultiVersionForUser 为用户关闭多版本功能
// @Summary 关闭多版本功能
// @Description 为指定用户关闭多版本润色功能
// @Tags admin
// @Produce json
// @Param user_id path int true "用户ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/users/{user_id}/multi-version/disable [post]
func (h *FeatureAdminHandler) DisableMultiVersionForUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.NewInvalidParameterError("无效的用户ID"))
		return
	}

	// 获取用户
	user, err := h.userRepo.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 关闭权限
	user.EnableMultiVersion = false

	if err := h.userRepo.Update(c.Request.Context(), user); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "已关闭用户的多版本润色功能",
		"user_id": userID,
	})
}

// UpdateQuota 更新用户配额
// @Summary 更新用户配额
// @Description 更新用户的多版本润色配额
// @Tags admin
// @Accept json
// @Produce json
// @Param user_id path int true "用户ID"
// @Param request body EnableMultiVersionRequest true "配额信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/users/{user_id}/multi-version/quota [put]
func (h *FeatureAdminHandler) UpdateQuota(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.NewInvalidParameterError("无效的用户ID"))
		return
	}

	var req EnableMultiVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// 获取用户
	user, err := h.userRepo.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 更新配额
	user.MultiVersionQuota = req.Quota

	if err := h.userRepo.Update(c.Request.Context(), user); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "已更新用户配额",
		"user_id": userID,
		"quota":   req.Quota,
	})
}

// GetUserMultiVersionStatus 获取用户多版本功能状态
// @Summary 获取用户状态
// @Description 获取用户的多版本润色功能状态和配额使用情况
// @Tags admin
// @Produce json
// @Param user_id path int true "用户ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/users/{user_id}/multi-version/status [get]
func (h *FeatureAdminHandler) GetUserMultiVersionStatus(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.NewInvalidParameterError("无效的用户ID"))
		return
	}

	// 获取用户
	user, err := h.userRepo.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	// TODO: 查询用户已使用的配额
	// usedQuota, err := h.polishRepo.CountUserMultiVersionRecords(c.Request.Context(), userID)

	response.Success(c, gin.H{
		"user_id":              userID,
		"enable_multi_version": user.EnableMultiVersion,
		"quota":                user.MultiVersionQuota,
		"used_quota":           0, // TODO: 实际使用量
		"has_unlimited_quota":  user.HasUnlimitedQuota(),
	})
}
