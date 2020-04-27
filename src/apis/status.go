package apis

import (
	"github.com/gin-gonic/gin"
	"little-diary-measurement-service/src/config"
	"little-diary-measurement-service/src/models"
	"net/http"
)

// GetHealth godoc
// @Summary Health check endpoint
// @Success 204
// @Router /status/health [get]
func GetHealth(c *gin.Context) {
	var measurement models.Measurement

	err := config.Config.DB.FirstOrInit(&measurement).Error
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.Status(http.StatusOK)
}
