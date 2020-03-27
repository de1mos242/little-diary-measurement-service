package dto

import (
	"little-diary-measurement-service/src/models"
	"time"
)

type MeasurementRequest struct {
	Type       string    `json:"type" enums:"HEIGHT,WEIGHT"`
	Timestamp  time.Time `json:"ts" swaggertype:"string" format:"datetime"`
	Value      float32   `json:"value"`
	TargetUuid string    `json:"target_uuid" swaggertype:"string" format:"uuid"`
}

type MeasurementResponse struct {
	Type       string    `json:"type" enums:"HEIGHT,WEIGHT"`
	Timestamp  time.Time `json:"ts" swaggertype:"string" format:"datetime"`
	Value      float32   `json:"value"`
	Uuid       string    `json:"uuid" swaggertype:"string" format:"uuid"`
	TargetUuid string    `json:"target_uuid" swaggertype:"string" format:"uuid"`
}

func MeasurementResponseFromModel(source *models.Measurement) *MeasurementResponse {
	m := &MeasurementResponse{
		Type:       string(source.Type),
		Timestamp:  source.Timestamp,
		Value:      source.Value,
		Uuid:       string(source.Uuid),
		TargetUuid: string(source.TargetUuid),
	}
	return m
}
