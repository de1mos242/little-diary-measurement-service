package models

import "time"

type MeasurementId uint
type MeasurementType string
type MeasurementUUID string
type TargetUUID string

const (
	MeasurementTypeHeight MeasurementType = "HEIGHT"
	MeasurementTypeWeight MeasurementType = "WEIGHT"
)

type Measurement struct {
	ID         MeasurementId   `gorm:"primary_key;column:id"`
	CreatedAt  time.Time       `gorm:"column:created_at"`
	UpdatedAt  time.Time       `gorm:"column:updated_at"`
	Type       MeasurementType `gorm:"column:measurement_type;not null"`
	Timestamp  time.Time       `gorm:"column:measurement_date;not null"`
	Value      float32         `gorm:"column:measurement_value;not null"`
	Uuid       MeasurementUUID `gorm:"column:measurement_uuid;unique;not null;type:uuid"`
	TargetUuid TargetUUID      `gorm:"column:target_uuid;not null;index;type:uuid"`
}
