package services

import (
	"little-diary-measurement-service/src/models"
)

type measurementDAO interface {
	GetByMeasurementUuid(measurementUuid models.MeasurementUUID) (*models.Measurement, error)
}

type MeasurementService struct {
	dao measurementDAO
}

func NewMeasurementService(dao measurementDAO) *MeasurementService {
	return &MeasurementService{dao}
}

func (s *MeasurementService) GetByMeasurementUuid(measurementUuid string) (*models.Measurement, error) {
	return s.dao.GetByMeasurementUuid(models.MeasurementUUID(measurementUuid))
}
