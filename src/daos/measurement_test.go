package daos

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"little-diary-measurement-service/src/config"
	"little-diary-measurement-service/src/models"
	"little-diary-measurement-service/src/test_data"
	"testing"
	"time"
)

func TestMeasurementDAO_Get(t *testing.T) {
	tx := test_data.OnBeforeDBTest()
	defer test_data.OnAfterDBTest(tx)

	type args struct {
		measurementUuid models.MeasurementUUID
	}

	stored := test_data.MeasurementStoredFactory.MustCreate().(*models.Measurement)

	tests := []struct {
		name    string
		args    args
		want    *models.Measurement
		wantErr bool
	}{
		{name: "find existed measurement", args: args{stored.Uuid}, want: stored, wantErr: false},
		{name: "find not existed measurement",
			args: args{models.MeasurementUUID(fmt.Sprintf("%s", uuid.New()))},
			want: stored, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := &MeasurementDAO{}
			got, err := dao.GetByMeasurementUuid(tt.args.measurementUuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByMeasurementUuid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			assert.Equal(t, got.ID, tt.want.ID)
			assert.Equal(t, got.Uuid, tt.want.Uuid)
			assert.Equal(t, got.CreatedAt.Format(time.RFC3339), tt.want.CreatedAt.In(time.UTC).Format(time.RFC3339))
			assert.Equal(t, got.TargetUuid, tt.want.TargetUuid)
			assert.Equal(t, got.Timestamp.Format(time.RFC3339), tt.want.Timestamp.In(time.UTC).Format(time.RFC3339))
			assert.Equal(t, got.Type, tt.want.Type)
			assert.Equal(t, got.UpdatedAt.Format(time.RFC3339), tt.want.UpdatedAt.In(time.UTC).Format(time.RFC3339))
			assert.Equal(t, got.Value, tt.want.Value)

		})
	}
}

func TestMeasurementDAO_SaveMeasurement(t *testing.T) {
	tx := test_data.OnBeforeDBTest()
	defer test_data.OnAfterDBTest(tx)

	m := test_data.MeasurementFactory.MustCreate().(*models.Measurement)

	type args struct {
		measurement *models.Measurement
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "store", args: args{m}},
		{name: "store fail", args: args{&models.Measurement{}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := &MeasurementDAO{}
			if err := dao.SaveMeasurement(tt.args.measurement); err != nil {
				assert.True(t, tt.wantErr)
				return
			}

			assert.NotNil(t, m.ID)
			var stored models.Measurement

			err := config.Config.DB.
				Where("measurement_uuid = ?", m.Uuid).
				First(&stored).
				Error
			assert.Nil(t, err)

			assert.NotNil(t, stored)
			assert.NotNil(t, stored.ID)
		})
	}
}

func TestMeasurementDAO_GetMeasurementsByTargetUuid(t *testing.T) {
	tx := test_data.OnBeforeDBTest()
	defer test_data.OnAfterDBTest(tx)

	t1 := models.TargetUUID(fmt.Sprintf("%s", uuid.New()))
	t2 := models.TargetUUID(fmt.Sprintf("%s", uuid.New()))
	t3 := models.TargetUUID(fmt.Sprintf("%s", uuid.New()))

	m1 := test_data.MeasurementStoredFactory.MustCreateWithOption(map[string]interface{}{"TargetUuid": t1}).(*models.Measurement)
	m2 := test_data.MeasurementStoredFactory.MustCreateWithOption(map[string]interface{}{"TargetUuid": t2}).(*models.Measurement)
	m3 := test_data.MeasurementStoredFactory.MustCreateWithOption(map[string]interface{}{"TargetUuid": t1}).(*models.Measurement)

	type args struct {
		targetUuid models.TargetUUID
	}
	tests := []struct {
		name    string
		args    args
		want    []*models.Measurement
		wantErr bool
	}{
		{
			name: "Find measurements by t1",
			args: args{t1},
			want: []*models.Measurement{m1, m3},
		},
		{
			name: "Find measurements by t2",
			args: args{t2},
			want: []*models.Measurement{m2},
		},
		{
			name: "Find measurements by empty t3",
			args: args{t3},
			want: []*models.Measurement{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := &MeasurementDAO{}
			got, err := dao.GetMeasurementsByTargetUuid(tt.args.targetUuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMeasurementsByTargetUuid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i, v := range tt.want {
				assert.Equal(t, got[i].ID, v.ID)
			}
		})
	}
}
