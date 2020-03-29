package router

import (
	"github.com/gin-gonic/gin"
	"little-diary-measurement-service/src/apis"
	"little-diary-measurement-service/src/common"
	"little-diary-measurement-service/src/security"
	"net/http"
	"strings"
)

func wrapHandler(f func(c *gin.Context, locator *common.ServiceLocator), locator *common.ServiceLocator) gin.HandlerFunc {
	return func(context *gin.Context) {
		f(context, locator)
	}
}

func GetMainEngine(locator *common.ServiceLocator) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	v1 := r.Group("/api/v1")
	v1.Use(jwtMiddleware(locator))
	{
		//v1.Use(auth())
		v1.GET("/measurement/:uuid", wrapHandler(apis.GetMeasurement, locator))
		v1.PUT("/measurement/:uuid", wrapHandler(apis.SaveMeasurement, locator))

		v1.GET("/measurements", wrapHandler(apis.GetMeasurementsByTarget, locator))
	}
	return r
}

func jwtMiddleware(locator *common.ServiceLocator) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		token = strings.Replace(token, "Bearer ", "", 1)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token not provided"})
		}
		jwtTokenReader := security.JwtTokenReader{PublicKey: locator.PublicKeyGetter.GetAuthServerJwtPublicKey()}
		userUuid, err := jwtTokenReader.ReadUserUuid(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token not provided"})
		}

		c.Set("UserUuid", userUuid)
		c.Next()
	}
}
