package daos

import (
	"little-diary-measurement-service/src/config"
	"little-diary-measurement-service/src/models"
)

type MeasurementDAO struct{}

func NewMeasurementDAO() *MeasurementDAO {
	return &MeasurementDAO{}
}

func (dao *MeasurementDAO) GetByMeasurementUuid(measurementUuid models.MeasurementUUID) (*models.Measurement, error) {
	var measurement models.Measurement

	err := config.Config.DB.
		Where("measurement_uuid = ?", measurementUuid).
		First(&measurement).
		Error

	return &measurement, err
}

func (dao *MeasurementDAO) SaveMeasurement(measurement *models.Measurement) error {
	return config.Config.DB.Save(measurement).Error
}

func (dao *MeasurementDAO) GetMeasurementsByTargetUuid(targetUuid models.TargetUUID) ([]*models.Measurement, error) {
	var measurements []*models.Measurement
	err := config.Config.DB.
		Where("target_uuid = ?", targetUuid).
		Order("measurement_date ASC").
		Find(&measurements).
		Error
	return measurements, err
}
