package apis

import (
	"github.com/gin-gonic/gin"
	"little-diary-measurement-service/src/daos"
	"little-diary-measurement-service/src/services"
	"log"
	"net/http"
)

func GetMeasurement(c *gin.Context) {
	s := services.NewMeasurementService(daos.NewMeasurementDAO())
	uuid := c.Param("uuid")
	if measurement, err := s.GetByMeasurementUuid(uuid); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		log.Println(err)
	} else {
		c.JSON(http.StatusOK, measurement)
	}
}
