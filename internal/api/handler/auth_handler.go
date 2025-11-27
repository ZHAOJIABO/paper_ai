package handler

import (
	"github.com/gin-gonic/gin"
	"paper_ai/internal/domain/model"
	"paper_ai/internal/service"
	apperrors "paper_ai/pkg/errors"
	"paper_ai/pkg/response"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.NewBadRequestError("参数错误"))
		return
	}

	userInfo, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, userInfo)
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.NewBadRequestError("参数错误"))
		return
	}

	// 获取客户端信息
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	loginResp, err := h.authService.Login(c.Request.Context(), &req, ip, userAgent)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, loginResp)
}

// RefreshToken 刷新访问令牌
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req model.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.NewBadRequestError("参数错误"))
		return
	}

	refreshResp, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, refreshResp)
}

// Logout 登出
func (h *AuthHandler) Logout(c *gin.Context) {
	var req model.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.NewBadRequestError("参数错误"))
		return
	}

	if err := h.authService.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "登出成功"})
}

// GetCurrentUser 获取当前用户信息
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// 从上下文获取用户ID（由中间件设置）
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, apperrors.NewUnauthorizedError("未授权"))
		return
	}

	userInfo, err := h.authService.GetCurrentUser(c.Request.Context(), userID.(int64))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, userInfo)
}
