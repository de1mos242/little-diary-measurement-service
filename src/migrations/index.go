package migrations

import (
	"gopkg.in/gormigrate.v1"
)

func GetMigrations() []*gormigrate.Migration {
	migrations := []*gormigrate.Migration{
		Getmigration201903252053InitTables(),
	}
	return migrations
}
