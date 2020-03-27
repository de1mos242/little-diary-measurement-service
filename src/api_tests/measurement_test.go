package api_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"little-diary-measurement-service/src/config"
	"little-diary-measurement-service/src/dto"
	"little-diary-measurement-service/src/models"
	"little-diary-measurement-service/src/router"
	"little-diary-measurement-service/src/test_data"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetMeasurement(t *testing.T) {
	tx := test_data.OnBeforeDBTest()
	defer test_data.OnAfterDBTest(tx)

	stored := test_data.MeasurementStoredFactory.MustCreate().(*models.Measurement)

	r := router.GetMainEngine()
	ts := httptest.NewServer(r)
	defer ts.Close()

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/measurement/%s", stored.Uuid), nil)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if assert.Equal(t, w.Code, http.StatusOK) {

		p, err := ioutil.ReadAll(w.Body)
		assert.Nil(t, err)

		var m = new(dto.MeasurementResponse)
		err = json.Unmarshal(p, &m)
		if assert.Nil(t, err) {

			assert.Equal(t, m.Uuid, string(stored.Uuid))
			assert.Equal(t, m.Value, stored.Value)
			assert.Equal(t, m.Type, string(stored.Type))
			assert.Equal(t, m.TargetUuid, string(stored.TargetUuid))
			assert.Equal(t, m.Timestamp.Format(time.RFC3339), stored.Timestamp.In(time.UTC).Format(time.RFC3339))
		}
	}
}

func TestGetMeasurementsByTarget(t *testing.T) {
	tx := test_data.OnBeforeDBTest()
	defer test_data.OnAfterDBTest(tx)

	t1 := models.TargetUUID(fmt.Sprintf("%s", uuid.New()))

	m1 := test_data.MeasurementStoredFactory.MustCreateWithOption(map[string]interface{}{"TargetUuid": t1}).(*models.Measurement)
	test_data.MeasurementStoredFactory.MustCreate()
	m3 := test_data.MeasurementStoredFactory.MustCreateWithOption(map[string]interface{}{"TargetUuid": t1}).(*models.Measurement)

	r := router.GetMainEngine()
	ts := httptest.NewServer(r)
	defer ts.Close()

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/measurements?target-uuid=%s", t1), nil)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if assert.Equal(t, w.Code, http.StatusOK) {

		p, err := ioutil.ReadAll(w.Body)
		if assert.Nil(t, err) {

			var m []*dto.MeasurementResponse
			err = json.Unmarshal(p, &m)
			if assert.Nil(t, err) {

				if assert.Len(t, m, 2) {
					assert.Equal(t, m[0].Uuid, string(m1.Uuid))
					assert.Equal(t, m[1].Uuid, string(m3.Uuid))
				}
			}
		}
	}
}

func TestAddMeasurement(t *testing.T) {
	tx := test_data.OnBeforeDBTest()
	defer test_data.OnAfterDBTest(tx)

	r := router.GetMainEngine()
	ts := httptest.NewServer(r)
	defer ts.Close()

	mUuid := fmt.Sprintf("%s", uuid.New())

	requestDto := &dto.MeasurementRequest{
		Type:       "WEIGHT",
		Timestamp:  time.Now(),
		Value:      9552,
		TargetUuid: fmt.Sprintf("%s", uuid.New()),
	}
	body, _ := json.Marshal(requestDto)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/measurement/%s", mUuid), bytes.NewReader(body))

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)

	p, err := ioutil.ReadAll(w.Body)
	assert.Nil(t, err)

	var responseDto = new(dto.MeasurementResponse)
	err = json.Unmarshal(p, &responseDto)
	assert.Nil(t, err)

	assert.Equal(t, responseDto.Uuid, mUuid)
	assert.Equal(t, responseDto.Value, requestDto.Value)
	assert.Equal(t, responseDto.Type, requestDto.Type)
	assert.Equal(t, responseDto.TargetUuid, requestDto.TargetUuid)
	assert.Equal(t, responseDto.Timestamp.In(time.UTC).Format(time.RFC3339), requestDto.Timestamp.In(time.UTC).Format(time.RFC3339))

	var stored models.Measurement

	err = config.Config.DB.
		Where("measurement_uuid = ?", mUuid).
		First(&stored).
		Error

	assert.Nil(t, err)
	assert.Equal(t, string(stored.Uuid), mUuid)
	assert.Equal(t, stored.Value, requestDto.Value)
	assert.Equal(t, string(stored.Type), requestDto.Type)
	assert.Equal(t, string(stored.TargetUuid), requestDto.TargetUuid)
	assert.Equal(t, stored.Timestamp.Format(time.RFC3339), requestDto.Timestamp.In(time.UTC).Format(time.RFC3339))
}
