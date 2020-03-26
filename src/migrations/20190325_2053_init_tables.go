package migrations

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
	"time"
)

func Getmigration201903252053InitTables() *gormigrate.Migration {
	m := gormigrate.Migration{ID: "20190325_2053_init_tables",
		Migrate: func(tx *gorm.DB) error {
			type MeasurementId uint
			type MeasurementType string
			type MeasurementUUID string
			type TargetUUID string

			type Measurement struct {
				ID         MeasurementId   `gorm:"primary_key;column:id" json:"-"`
				CreatedAt  time.Time       `gorm:"column:created_at" json:"-"`
				UpdatedAt  time.Time       `gorm:"column:updated_at" json:"-"`
				Type       MeasurementType `gorm:"column:measurement_type;not null" json:"type"`
				Timestamp  time.Time       `gorm:"column:measurement_date;not null" json:"ts"`
				Value      float32         `gorm:"column:measurement_value;not null" json:"value"`
				Uuid       MeasurementUUID `gorm:"column:measurement_uuid;unique;not null;type:uuid" json:"uuid"`
				TargetUuid TargetUUID      `gorm:"column:target_uuid;not null;index;type:uuid" json:"target_uuid"`
			}

			return tx.AutoMigrate(&Measurement{}).Error
		}}
	return &m
}
