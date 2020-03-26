package daos

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
