package test_data

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/gormigrate.v1"
	"little-diary-measurement-service/src/config"
	"little-diary-measurement-service/src/migrations"
)

func init() {
	if err := config.LoadConfig("../../config"); err != nil {
		panic(fmt.Errorf("invalid test application configuration: %s", err))
	}

	config.Config.DB, config.Config.DBErr = gorm.Open("postgres", config.Config.DSN)
	if config.Config.DBErr != nil {
		panic(config.Config.DBErr)
	}

	config.Config.DB.LogMode(true)

	migrate()
}

func migrate() {
	m := gormigrate.New(config.Config.DB, gormigrate.DefaultOptions, migrations.GetMigrations())

	if err := m.Migrate(); err != nil {
		//panic(fmt.Errorf("could not migrate: %v", err))
		// don't panic, it's another thread run migration or we fall in tests
	}
}

func OnBeforeDBTest() *gorm.DB {
	old := config.Config.DB
	config.Config.DB = old.Begin()
	return old
}

func OnAfterDBTest(tx *gorm.DB) {
	config.Config.DB.Rollback()
	config.Config.DB = tx
}
