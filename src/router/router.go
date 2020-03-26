package router

import (
	"github.com/gin-gonic/gin"
	"little-diary-measurement-service/src/apis"
)

func GetMainEngine() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	v1 := r.Group("/api/v1")
	{
		//v1.Use(auth())
		v1.GET("/measurement/:uuid", apis.GetMeasurement)
	}
	return r
}
