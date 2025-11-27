package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"paper_ai/internal/infrastructure/security"
	apperrors "paper_ai/pkg/errors"
	"paper_ai/pkg/response"
)

// AuthRequired JWT认证中间件
func AuthRequired(jwtManager *security.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从Header获取Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, apperrors.NewUnauthorizedError("缺少认证令牌"))
			c.Abort()
			return
		}

		// 2. 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, apperrors.NewUnauthorizedError("令牌格式错误"))
			c.Abort()
			return
		}

		token := parts[1]

		// 3. 验证Token
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			response.Error(c, apperrors.NewUnauthorizedError("无效的令牌"))
			c.Abort()
			return
		}

		// 4. 保存用户信息到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// OptionalAuth 可选认证中间件（不强制要求认证，但如果有token则验证）
func OptionalAuth(jwtManager *security.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		token := parts[1]
		claims, err := jwtManager.ValidateToken(token)
		if err == nil {
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
		}

		c.Next()
	}
}
