package services

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"little-diary-measurement-service/src/dto"
	"little-diary-measurement-service/src/models"
	"little-diary-measurement-service/src/test_data"
	"reflect"
	"testing"
	"time"
)

var tUuid = fmt.Sprintf("%s", uuid.New())

var records = []*models.Measurement{
	test_data.MeasurementFactory.MustCreateWithOption(map[string]interface{}{"TargetUuid": models.TargetUUID(tUuid)}).(*models.Measurement),
	test_data.MeasurementFactory.MustCreate().(*models.Measurement),
	test_data.MeasurementFactory.MustCreateWithOption(map[string]interface{}{"TargetUuid": models.TargetUUID(tUuid)}).(*models.Measurement),
}

type mockMeasurementDAO struct {
	records []*models.Measurement
}

func (m *mockMeasurementDAO) GetMeasurementsByTargetUuid(targetUuid models.TargetUUID) ([]*models.Measurement, error) {
	var res []*models.Measurement
	for _, record := range m.records {
		if record.TargetUuid == targetUuid {
			res = append(res, record)
		}
	}
	return res, nil
}

func (m *mockMeasurementDAO) GetByMeasurementUuid(measurementUuid models.MeasurementUUID) (*models.Measurement, error) {
	for _, record := range m.records {
		if record.Uuid == measurementUuid {
			return record, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockMeasurementDAO) SaveMeasurement(measurement *models.Measurement) error {
	measurement.ID = models.MeasurementId(100500)
	return nil
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

func TestMeasurementService_validateMeasurement(t *testing.T) {
	type fields struct {
		dao measurementDAO
	}
	type args struct {
		request dto.MeasurementRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "test height", args: args{request: dto.MeasurementRequest{Type: "HEIGHT"}}, wantErr: false},
		{name: "test weight", args: args{request: dto.MeasurementRequest{Type: "WEIGHT"}}, wantErr: false},
		{name: "test height case fail", args: args{request: dto.MeasurementRequest{Type: "height"}}, wantErr: true},
		{name: "test random", args: args{request: dto.MeasurementRequest{Type: "dfa"}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MeasurementService{
				dao: tt.fields.dao,
			}
			if err := s.validateMeasurement(tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("validateMeasurement() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMeasurementService_Save(t *testing.T) {
	randomUuid := fmt.Sprintf("%s", uuid.New())
	targetUuid := fmt.Sprintf("%s", uuid.New())
	type fields struct {
		dao measurementDAO
	}
	type args struct {
		uuid    string
		request dto.MeasurementRequest
	}
	twoHoursBefore := time.Now().Add(time.Hour * -2)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Measurement
		wantErr bool
	}{
		{name: "test update measurement", fields: fields{dao: newMockMeasurementDAO()}, args: args{uuid: string(records[1].Uuid), request: dto.MeasurementRequest{
			Type:       "WEIGHT",
			Timestamp:  twoHoursBefore,
			Value:      9500,
			TargetUuid: targetUuid,
		}}, want: records[1]},
		{
			name:   "test create measurement",
			fields: fields{dao: newMockMeasurementDAO()},
			args: args{uuid: randomUuid, request: dto.MeasurementRequest{
				Type:       "HEIGHT",
				Timestamp:  twoHoursBefore,
				Value:      73,
				TargetUuid: targetUuid,
			}},
			want: &models.Measurement{
				ID:         models.MeasurementId(100500),
				Type:       models.MeasurementTypeHeight,
				Timestamp:  twoHoursBefore,
				Value:      73,
				Uuid:       models.MeasurementUUID(randomUuid),
				TargetUuid: models.TargetUUID(targetUuid),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MeasurementService{
				dao: tt.fields.dao,
			}
			got, err := s.Save(tt.args.uuid, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Save() got = %v, want %v", got, tt.want)
			}
			if got != nil {
				assert.Equal(t, got.Value, tt.args.request.Value)
				assert.Equal(t, got.Timestamp, tt.args.request.Timestamp)
			}
		})
	}
}

func TestMeasurementService_GetByTargetUuid(t *testing.T) {
	type fields struct {
		dao measurementDAO
	}
	type args struct {
		targetUuid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*models.Measurement
		wantErr bool
	}{
		{
			name:   "Find by target uuid",
			fields: fields{newMockMeasurementDAO()},
			args:   args{tUuid},
			want:   []*models.Measurement{records[0], records[2]},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MeasurementService{
				dao: tt.fields.dao,
			}
			got, err := s.GetByTargetUuid(tt.args.targetUuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByTargetUuid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByTargetUuid() got = %v, want %v", got, tt.want)
			}
		})
	}
}
