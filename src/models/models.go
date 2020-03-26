package models

import "time"

type MeasurementId uint
type MeasurementType string
type MeasurementUUID string
type TargetUUID string

const (
	MAESUREMENT_TYPE_HEIGHT MeasurementType = "HEIGHT"
	MAESUREMENT_TYPE_WEIGHT MeasurementType = "WEIGHT"
)

type Measurement struct {
	ID         MeasurementId   `gorm:"primary_key;column:id" json:"-"`
	CreatedAt  time.Time       `gorm:"column:created_at" json:"-"`
	UpdatedAt  time.Time       `gorm:"column:updated_at" json:"-"`
	Type       MeasurementType `gorm:"column:measurement_type;not null" json:"type"`
	Timestamp  time.Time       `gorm:"column:measurement_date;not null" json:"ts" swaggertype:"string" format:"datetime"`
	Value      float32         `gorm:"column:measurement_value;not null" json:"value"`
	Uuid       MeasurementUUID `gorm:"column:measurement_uuid;unique;not null;type:uuid" json:"uuid" swaggertype:"string" format:"uuid"`
	TargetUuid TargetUUID      `gorm:"column:target_uuid;not null;index;type:uuid" json:"target_uuid" swaggertype:"string" format:"uuid"`
}
