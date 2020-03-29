package apis

import (
	"github.com/gin-gonic/gin"
	"little-diary-measurement-service/src/common"
	"little-diary-measurement-service/src/daos"
	"little-diary-measurement-service/src/dto"
	"little-diary-measurement-service/src/errors"
	"little-diary-measurement-service/src/services"
	"log"
	"net/http"
)

// GetMeasurement godoc
// @Summary Retrieves measurement based on given UUID
// @Security ApiKeyAuth
// @Produce json
// @Param uuid path string true "Measurement UUID" format(uuid)
// @Success 200 {object} dto.MeasurementResponse
// @Router /measurement/{uuid} [get]
func GetMeasurement(c *gin.Context, locator *common.ServiceLocator) {
	s := services.NewMeasurementService(daos.NewMeasurementDAO(), locator)
	uuid := c.Param("uuid")
	userUuid := c.GetString("UserUuid")
	if measurement, err := s.GetByMeasurementUuid(uuid, userUuid); err != nil {
		if _, ok := err.(*errors.ForbiddenError); ok {
			c.AbortWithStatus(http.StatusForbidden)
		} else {
			c.AbortWithStatus(http.StatusNotFound)
			log.Println(err)
		}
	} else {
		c.JSON(http.StatusOK, dto.MeasurementResponseFromModel(measurement))
	}
}

// GetMeasurementsByTarget godoc
// @Summary Retrieves measurements based on given target UUID
// @Security ApiKeyAuth
// @Produce json
// @Param target-uuid query string true "Target UUID" format(uuid)
// @Success 200 {array} dto.MeasurementResponse
// @Router /measurements [get]
func GetMeasurementsByTarget(c *gin.Context, locator *common.ServiceLocator) {
	s := services.NewMeasurementService(daos.NewMeasurementDAO(), locator)
	uuid := c.Query("target-uuid")
	userUuid := c.GetString("UserUuid")
	if measurements, err := s.GetByTargetUuid(uuid, userUuid); err != nil {
		if _, ok := err.(*errors.ForbiddenError); ok {
			c.AbortWithStatus(http.StatusForbidden)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
			log.Println(err)
		}
	} else {
		var dtos []*dto.MeasurementResponse
		for _, m := range measurements {
			dtos = append(dtos, dto.MeasurementResponseFromModel(m))
		}
		c.JSON(http.StatusOK, dtos)
	}
}

// SaveMeasurement godoc
// @Summary Create or update measurement
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param uuid path string true "Measurement UUID" format(uuid)
// @Param measurement body dto.MeasurementRequest true "Measurement data" format(uuid)
// @Success 200 {object} dto.MeasurementResponse
// @Router /measurement/{uuid} [put]
func SaveMeasurement(c *gin.Context, locator *common.ServiceLocator) {
	s := services.NewMeasurementService(daos.NewMeasurementDAO(), locator)
	uuid := c.Param("uuid")
	userUuid := c.GetString("UserUuid")
	var requestDto dto.MeasurementRequest
	err := c.BindJSON(&requestDto)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Println(err)
	}
	if measurement, err := s.Save(uuid, requestDto, userUuid); err != nil {
		if _, ok := err.(*errors.ForbiddenError); ok {
			c.AbortWithStatus(http.StatusForbidden)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
			log.Println(err)
		}
	} else {
		responseDto := dto.MeasurementResponseFromModel(measurement)
		c.JSON(http.StatusOK, responseDto)
	}
}
