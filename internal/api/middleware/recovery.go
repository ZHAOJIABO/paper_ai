package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10005,
			"message": "internal server error",
		})
		c.Abort()
	})
}
