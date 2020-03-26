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
