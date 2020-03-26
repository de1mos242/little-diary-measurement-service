package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/gormigrate.v1"
	"little-diary-measurement-service/src/config"
	"little-diary-measurement-service/src/migrations"
	"little-diary-measurement-service/src/router"
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

	r := router.GetMainEngine()

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
