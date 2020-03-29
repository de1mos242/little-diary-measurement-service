package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"gopkg.in/gormigrate.v1"
	"little-diary-measurement-service/src/common"
	"little-diary-measurement-service/src/config"
	_ "little-diary-measurement-service/src/docs"
	"little-diary-measurement-service/src/integrations"
	"little-diary-measurement-service/src/migrations"
	"little-diary-measurement-service/src/router"
	"net/http"
)

// @title Measurement service API
// @version 1.0

// @contact.name Denis Yakovlev
// @contact.email de1m0s242@gmail.com

// @BasePath /api/v1
func main() {

	if err := config.LoadConfig("./config"); err != nil {
		panic(fmt.Errorf("invalid application configuration: %s", err))
	}

	authIntegration := integrations.AuthIntegration{
		Client: &http.Client{},
		Config: &config.Config,
	}
	serviceLocator := common.ServiceLocator{
		PublicKeyGetter: &config.Config,
		UserHasAccessToBabyChecker: &integrations.FamilyIntegration{
			Client:      &http.Client{},
			Config:      &config.Config,
			AuthService: &authIntegration,
		},
	}

	r := router.GetMainEngine(&serviceLocator)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	config.Config.DB, config.Config.DBErr = gorm.Open("postgres", config.Config.DSN)
	if config.Config.DBErr != nil {
		panic(config.Config.DBErr)
	}
	defer config.Config.DB.Close()

	config.Config.DB.LogMode(true)

	m := gormigrate.New(config.Config.DB, gormigrate.DefaultOptions, migrations.GetMigrations())

	if err := m.Migrate(); err != nil {
		panic(fmt.Errorf("could not migrate: %v", err))
	}

	r.Run(fmt.Sprintf(":%v", config.Config.ServerPort))
}
