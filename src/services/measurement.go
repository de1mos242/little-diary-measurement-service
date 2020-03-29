package services

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"little-diary-measurement-service/src/common"
	"little-diary-measurement-service/src/dto"
	"little-diary-measurement-service/src/errors"
	"little-diary-measurement-service/src/models"
)

type measurementDAO interface {
	GetByMeasurementUuid(measurementUuid models.MeasurementUUID) (*models.Measurement, error)
	SaveMeasurement(measurement *models.Measurement) error
	GetMeasurementsByTargetUuid(targetUuid models.TargetUUID) ([]*models.Measurement, error)
}

type MeasurementService struct {
	dao            measurementDAO
	serviceLocator *common.ServiceLocator
}

func NewMeasurementService(dao measurementDAO, locator *common.ServiceLocator) *MeasurementService {
	return &MeasurementService{dao, locator}
}

func (s *MeasurementService) GetByMeasurementUuid(measurementUuid string, userUuid string) (*models.Measurement, error) {
	measurement, err := s.dao.GetByMeasurementUuid(models.MeasurementUUID(measurementUuid))
	if err != nil {
		return nil, err
	}
	allowed, err := s.serviceLocator.UserHasAccessToBabyChecker.CheckUserHasAccessToBaby(userUuid, string(measurement.TargetUuid))
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, &errors.ForbiddenError{S: "operation not allowed"}
	}
	return measurement, err
}

func (s *MeasurementService) Save(uuid string, request dto.MeasurementRequest, userUuid string) (*models.Measurement, error) {
	err := s.validateMeasurement(request)
	if err != nil {
		return nil, err
	}

	allowed, err := s.serviceLocator.UserHasAccessToBabyChecker.CheckUserHasAccessToBaby(userUuid, request.TargetUuid)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, &errors.ForbiddenError{S: "operation not allowed"}
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

func (s *MeasurementService) GetByTargetUuid(targetUuid string, userUuid string) ([]*models.Measurement, error) {
	allowed, err := s.serviceLocator.UserHasAccessToBabyChecker.CheckUserHasAccessToBaby(userUuid, targetUuid)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, &errors.ForbiddenError{S: "operation not allowed"}
	}
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
