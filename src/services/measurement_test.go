package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"little-diary-measurement-service/src/models"
	"little-diary-measurement-service/src/test_data"
	"reflect"
	"testing"
)

var records = []*models.Measurement{
	test_data.MeasurementFactory.MustCreate().(*models.Measurement),
	test_data.MeasurementFactory.MustCreate().(*models.Measurement),
}

type mockMeasurementDAO struct {
	records []*models.Measurement
}

func (m *mockMeasurementDAO) GetByMeasurementUuid(measurementUuid models.MeasurementUUID) (*models.Measurement, error) {
	for _, record := range m.records {
		if record.Uuid == measurementUuid {
			return record, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Record %s not found", measurementUuid))
}

func newMockMeasurementDAO() measurementDAO {
	return &mockMeasurementDAO{
		records: records,
	}
}

func TestMeasurementService_Get(t *testing.T) {
	type fields struct {
		dao measurementDAO
	}
	type args struct {
		measurementUuid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Measurement
		wantErr bool
	}{
		{name: "find existed record", fields: fields{newMockMeasurementDAO()},
			args: args{string(records[0].Uuid)}, want: records[0]},
		{name: "find existed record 2", fields: fields{newMockMeasurementDAO()},
			args: args{string(records[1].Uuid)}, want: records[1]},
		{name: "find not existed record", fields: fields{newMockMeasurementDAO()},
			args: args{fmt.Sprintf("%s", uuid.New())}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MeasurementService{
				dao: tt.fields.dao,
			}
			got, err := s.GetByMeasurementUuid(tt.args.measurementUuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByMeasurementUuid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByMeasurementUuid() got = %v, want %v", got, tt.want)
			}
		})
	}
}
