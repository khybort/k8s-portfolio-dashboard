package middleware

import (
	"net/http"
	"go.uber.org/zap"
	"github.com/gin-gonic/gin"
)

func Recovery(zapLogger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		zapLogger.Error("Panic recovered", zap.Any("error", recovered))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		c.Abort()
	})
}

