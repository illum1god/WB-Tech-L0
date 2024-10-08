package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ErrorRecoveryMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.Errorf("Паника восстановлена: %v", recovered)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}
