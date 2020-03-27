package services

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"little-diary-measurement-service/src/dto"
	"little-diary-measurement-service/src/models"
)

type measurementDAO interface {
	GetByMeasurementUuid(measurementUuid models.MeasurementUUID) (*models.Measurement, error)
	SaveMeasurement(measurement *models.Measurement) error
	GetMeasurementsByTargetUuid(targetUuid models.TargetUUID) ([]*models.Measurement, error)
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

func (s *MeasurementService) Save(uuid string, request dto.MeasurementRequest) (*models.Measurement, error) {
	err := s.validateMeasurement(request)
	if err != nil {
		return nil, err
	}

	measurementUUID := models.MeasurementUUID(uuid)
	measurement, err := s.dao.GetByMeasurementUuid(measurementUUID)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			measurement = &models.Measurement{
				Uuid:       measurementUUID,
				TargetUuid: models.TargetUUID(request.TargetUuid),
				Type:       models.MeasurementType(request.Type),
			}
		} else {
			return nil, err
		}
	}
	measurement.Value = request.Value
	measurement.Timestamp = request.Timestamp
	err = s.dao.SaveMeasurement(measurement)
	return measurement, err
}

func (s *MeasurementService) GetByTargetUuid(targetUuid string) ([]*models.Measurement, error) {
	return s.dao.GetMeasurementsByTargetUuid(models.TargetUUID(targetUuid))
}

func (s *MeasurementService) validateMeasurement(request dto.MeasurementRequest) error {
	allowedTypes := map[string]bool{
		string(models.MeasurementTypeHeight): true,
		string(models.MeasurementTypeWeight): true,
	}
	if _, exists := allowedTypes[request.Type]; exists == false {
		return fmt.Errorf("measurement type %s does not exist", request.Type)
	}
	return nil
}
