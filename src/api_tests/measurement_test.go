package api_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"little-diary-measurement-service/src/common"
	"little-diary-measurement-service/src/config"
	"little-diary-measurement-service/src/dto"
	"little-diary-measurement-service/src/integrations"
	"little-diary-measurement-service/src/models"
	"little-diary-measurement-service/src/router"
	"little-diary-measurement-service/src/test_data"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetMeasurement(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tx := test_data.OnBeforeDBTest()
	defer test_data.OnAfterDBTest(tx)

	stored := test_data.MeasurementStoredFactory.MustCreate().(*models.Measurement)

	r := router.GetMainEngine(&common.ServiceLocator{
		PublicKeyGetter: &config.Config,
		UserHasAccessToBabyChecker: func() integrations.UserHasAccessToBabyChecker {
			mockObj := new(test_data.MockUserHasAccessToBabyChecker)
			mockObj.On("CheckUserHasAccessToBaby", mock.Anything, mock.Anything).Return(true, nil)
			return mockObj
		}(),
	})
	ts := httptest.NewServer(r)
	defer ts.Close()

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/measurement/%s", stored.Uuid), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJpYXQiOjE1ODU0ODk5NTgsIm5iZiI6MTU4NTQ4OTk1OCwianRpIjoiNWE4OGNlYTctMzQ5NC00MThiLTlmMWYtMGQ5NjI0NDkwNjU1IiwiZXhwIjoyNTg1NDkwODU4LCJpZGVudGl0eSI6MSwiZnJlc2giOmZhbHNlLCJ0eXBlIjoiYWNjZXNzIiwidXNlcl9jbGFpbXMiOnsicm9sZSI6ImFkbWluIiwidXVpZCI6IjNhZWY5YjVjLTg4M2YtNDEyZC1hZGU3LTQ3YmU0MzgyN2Q2OCIsInJlc291cmNlcyI6W119fQ.m_tC6rROZWPU1quds3lf07r8toOX5Gwg-FRKL8UIo6We_KT5dVYeU23330619BLaaNgSOObpTNysl31hxLBBXgyxngnyqRroXXWIz0lcieY76X18Z7xfq8twKavyaAjN4XyUmTVtTzBuN1VKknBIzGGN9QpKfJrjpFWGFVNXj3bOmJRqHyDO49uGQ5CgJdidnQ_OFVRXugXuLnkP94fe4bEAQxW-vE53zkrgdNFqElCvj6t5IA7yKX216SJxUs6vj-nkHsjJ-_sw1QdDCQ-ZKvb-kRAU8x9UHYiMaeQifIg3SVQt2jHqLzZpF1S1uOmyfJ1GOVIZt1MMy9wjYrtl-A"))

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

func TestGetMeasurementUnAuthorize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tx := test_data.OnBeforeDBTest()
	defer test_data.OnAfterDBTest(tx)

	stored := test_data.MeasurementStoredFactory.MustCreate().(*models.Measurement)

	r := router.GetMainEngine(&common.ServiceLocator{
		PublicKeyGetter: &config.Config,
		UserHasAccessToBabyChecker: func() integrations.UserHasAccessToBabyChecker {
			mockObj := new(test_data.MockUserHasAccessToBabyChecker)
			mockObj.On("CheckUserHasAccessToBaby", mock.Anything, mock.Anything).Return(true, nil)
			return mockObj
		}(),
	})
	ts := httptest.NewServer(r)
	defer ts.Close()

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/measurement/%s", stored.Uuid), nil)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusUnauthorized)
}

func TestGetMeasurementForbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tx := test_data.OnBeforeDBTest()
	defer test_data.OnAfterDBTest(tx)

	stored := test_data.MeasurementStoredFactory.MustCreate().(*models.Measurement)

	r := router.GetMainEngine(&common.ServiceLocator{
		PublicKeyGetter: &config.Config,
		UserHasAccessToBabyChecker: func() integrations.UserHasAccessToBabyChecker {
			mockObj := new(test_data.MockUserHasAccessToBabyChecker)
			mockObj.On("CheckUserHasAccessToBaby", mock.Anything, mock.Anything).Return(false, nil)
			return mockObj
		}(),
	})
	ts := httptest.NewServer(r)
	defer ts.Close()

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/measurement/%s", stored.Uuid), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJpYXQiOjE1ODU0ODk5NTgsIm5iZiI6MTU4NTQ4OTk1OCwianRpIjoiNWE4OGNlYTctMzQ5NC00MThiLTlmMWYtMGQ5NjI0NDkwNjU1IiwiZXhwIjoyNTg1NDkwODU4LCJpZGVudGl0eSI6MSwiZnJlc2giOmZhbHNlLCJ0eXBlIjoiYWNjZXNzIiwidXNlcl9jbGFpbXMiOnsicm9sZSI6ImFkbWluIiwidXVpZCI6IjNhZWY5YjVjLTg4M2YtNDEyZC1hZGU3LTQ3YmU0MzgyN2Q2OCIsInJlc291cmNlcyI6W119fQ.m_tC6rROZWPU1quds3lf07r8toOX5Gwg-FRKL8UIo6We_KT5dVYeU23330619BLaaNgSOObpTNysl31hxLBBXgyxngnyqRroXXWIz0lcieY76X18Z7xfq8twKavyaAjN4XyUmTVtTzBuN1VKknBIzGGN9QpKfJrjpFWGFVNXj3bOmJRqHyDO49uGQ5CgJdidnQ_OFVRXugXuLnkP94fe4bEAQxW-vE53zkrgdNFqElCvj6t5IA7yKX216SJxUs6vj-nkHsjJ-_sw1QdDCQ-ZKvb-kRAU8x9UHYiMaeQifIg3SVQt2jHqLzZpF1S1uOmyfJ1GOVIZt1MMy9wjYrtl-A"))

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusForbidden)
}

func TestGetMeasurementsByTarget(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tx := test_data.OnBeforeDBTest()
	defer test_data.OnAfterDBTest(tx)

	t1 := models.TargetUUID(fmt.Sprintf("%s", uuid.New()))

	m1 := test_data.MeasurementStoredFactory.MustCreateWithOption(map[string]interface{}{"TargetUuid": t1}).(*models.Measurement)
	test_data.MeasurementStoredFactory.MustCreate()
	m3 := test_data.MeasurementStoredFactory.MustCreateWithOption(map[string]interface{}{"TargetUuid": t1}).(*models.Measurement)

	r := router.GetMainEngine(&common.ServiceLocator{
		PublicKeyGetter: &config.Config,
		UserHasAccessToBabyChecker: func() integrations.UserHasAccessToBabyChecker {
			mockObj := new(test_data.MockUserHasAccessToBabyChecker)
			mockObj.On("CheckUserHasAccessToBaby", mock.Anything, mock.Anything).Return(true, nil)
			return mockObj
		}(),
	})
	ts := httptest.NewServer(r)
	defer ts.Close()

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/measurements?target-uuid=%s", t1), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJpYXQiOjE1ODU0ODk5NTgsIm5iZiI6MTU4NTQ4OTk1OCwianRpIjoiNWE4OGNlYTctMzQ5NC00MThiLTlmMWYtMGQ5NjI0NDkwNjU1IiwiZXhwIjoyNTg1NDkwODU4LCJpZGVudGl0eSI6MSwiZnJlc2giOmZhbHNlLCJ0eXBlIjoiYWNjZXNzIiwidXNlcl9jbGFpbXMiOnsicm9sZSI6ImFkbWluIiwidXVpZCI6IjNhZWY5YjVjLTg4M2YtNDEyZC1hZGU3LTQ3YmU0MzgyN2Q2OCIsInJlc291cmNlcyI6W119fQ.m_tC6rROZWPU1quds3lf07r8toOX5Gwg-FRKL8UIo6We_KT5dVYeU23330619BLaaNgSOObpTNysl31hxLBBXgyxngnyqRroXXWIz0lcieY76X18Z7xfq8twKavyaAjN4XyUmTVtTzBuN1VKknBIzGGN9QpKfJrjpFWGFVNXj3bOmJRqHyDO49uGQ5CgJdidnQ_OFVRXugXuLnkP94fe4bEAQxW-vE53zkrgdNFqElCvj6t5IA7yKX216SJxUs6vj-nkHsjJ-_sw1QdDCQ-ZKvb-kRAU8x9UHYiMaeQifIg3SVQt2jHqLzZpF1S1uOmyfJ1GOVIZt1MMy9wjYrtl-A"))

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

func TestGetMeasurementsByTargetForbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tx := test_data.OnBeforeDBTest()
	defer test_data.OnAfterDBTest(tx)

	t1 := models.TargetUUID(fmt.Sprintf("%s", uuid.New()))

	r := router.GetMainEngine(&common.ServiceLocator{
		PublicKeyGetter: &config.Config,
		UserHasAccessToBabyChecker: func() integrations.UserHasAccessToBabyChecker {
			mockObj := new(test_data.MockUserHasAccessToBabyChecker)
			mockObj.On("CheckUserHasAccessToBaby", mock.Anything, mock.Anything).Return(false, nil)
			return mockObj
		}(),
	})
	ts := httptest.NewServer(r)
	defer ts.Close()

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/measurements?target-uuid=%s", t1), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJpYXQiOjE1ODU0ODk5NTgsIm5iZiI6MTU4NTQ4OTk1OCwianRpIjoiNWE4OGNlYTctMzQ5NC00MThiLTlmMWYtMGQ5NjI0NDkwNjU1IiwiZXhwIjoyNTg1NDkwODU4LCJpZGVudGl0eSI6MSwiZnJlc2giOmZhbHNlLCJ0eXBlIjoiYWNjZXNzIiwidXNlcl9jbGFpbXMiOnsicm9sZSI6ImFkbWluIiwidXVpZCI6IjNhZWY5YjVjLTg4M2YtNDEyZC1hZGU3LTQ3YmU0MzgyN2Q2OCIsInJlc291cmNlcyI6W119fQ.m_tC6rROZWPU1quds3lf07r8toOX5Gwg-FRKL8UIo6We_KT5dVYeU23330619BLaaNgSOObpTNysl31hxLBBXgyxngnyqRroXXWIz0lcieY76X18Z7xfq8twKavyaAjN4XyUmTVtTzBuN1VKknBIzGGN9QpKfJrjpFWGFVNXj3bOmJRqHyDO49uGQ5CgJdidnQ_OFVRXugXuLnkP94fe4bEAQxW-vE53zkrgdNFqElCvj6t5IA7yKX216SJxUs6vj-nkHsjJ-_sw1QdDCQ-ZKvb-kRAU8x9UHYiMaeQifIg3SVQt2jHqLzZpF1S1uOmyfJ1GOVIZt1MMy9wjYrtl-A"))

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusForbidden)
}

func TestAddMeasurement(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tx := test_data.OnBeforeDBTest()
	defer test_data.OnAfterDBTest(tx)

	r := router.GetMainEngine(&common.ServiceLocator{
		PublicKeyGetter: &config.Config,
		UserHasAccessToBabyChecker: func() integrations.UserHasAccessToBabyChecker {
			mockObj := new(test_data.MockUserHasAccessToBabyChecker)
			mockObj.On("CheckUserHasAccessToBaby", mock.Anything, mock.Anything).Return(true, nil)
			return mockObj
		}(),
	})
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJpYXQiOjE1ODU0ODk5NTgsIm5iZiI6MTU4NTQ4OTk1OCwianRpIjoiNWE4OGNlYTctMzQ5NC00MThiLTlmMWYtMGQ5NjI0NDkwNjU1IiwiZXhwIjoyNTg1NDkwODU4LCJpZGVudGl0eSI6MSwiZnJlc2giOmZhbHNlLCJ0eXBlIjoiYWNjZXNzIiwidXNlcl9jbGFpbXMiOnsicm9sZSI6ImFkbWluIiwidXVpZCI6IjNhZWY5YjVjLTg4M2YtNDEyZC1hZGU3LTQ3YmU0MzgyN2Q2OCIsInJlc291cmNlcyI6W119fQ.m_tC6rROZWPU1quds3lf07r8toOX5Gwg-FRKL8UIo6We_KT5dVYeU23330619BLaaNgSOObpTNysl31hxLBBXgyxngnyqRroXXWIz0lcieY76X18Z7xfq8twKavyaAjN4XyUmTVtTzBuN1VKknBIzGGN9QpKfJrjpFWGFVNXj3bOmJRqHyDO49uGQ5CgJdidnQ_OFVRXugXuLnkP94fe4bEAQxW-vE53zkrgdNFqElCvj6t5IA7yKX216SJxUs6vj-nkHsjJ-_sw1QdDCQ-ZKvb-kRAU8x9UHYiMaeQifIg3SVQt2jHqLzZpF1S1uOmyfJ1GOVIZt1MMy9wjYrtl-A"))

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

func TestAddMeasurementForbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tx := test_data.OnBeforeDBTest()
	defer test_data.OnAfterDBTest(tx)

	r := router.GetMainEngine(&common.ServiceLocator{
		PublicKeyGetter: &config.Config,
		UserHasAccessToBabyChecker: func() integrations.UserHasAccessToBabyChecker {
			mockObj := new(test_data.MockUserHasAccessToBabyChecker)
			mockObj.On("CheckUserHasAccessToBaby", mock.Anything, mock.Anything).Return(false, nil)
			return mockObj
		}(),
	})
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJpYXQiOjE1ODU0ODk5NTgsIm5iZiI6MTU4NTQ4OTk1OCwianRpIjoiNWE4OGNlYTctMzQ5NC00MThiLTlmMWYtMGQ5NjI0NDkwNjU1IiwiZXhwIjoyNTg1NDkwODU4LCJpZGVudGl0eSI6MSwiZnJlc2giOmZhbHNlLCJ0eXBlIjoiYWNjZXNzIiwidXNlcl9jbGFpbXMiOnsicm9sZSI6ImFkbWluIiwidXVpZCI6IjNhZWY5YjVjLTg4M2YtNDEyZC1hZGU3LTQ3YmU0MzgyN2Q2OCIsInJlc291cmNlcyI6W119fQ.m_tC6rROZWPU1quds3lf07r8toOX5Gwg-FRKL8UIo6We_KT5dVYeU23330619BLaaNgSOObpTNysl31hxLBBXgyxngnyqRroXXWIz0lcieY76X18Z7xfq8twKavyaAjN4XyUmTVtTzBuN1VKknBIzGGN9QpKfJrjpFWGFVNXj3bOmJRqHyDO49uGQ5CgJdidnQ_OFVRXugXuLnkP94fe4bEAQxW-vE53zkrgdNFqElCvj6t5IA7yKX216SJxUs6vj-nkHsjJ-_sw1QdDCQ-ZKvb-kRAU8x9UHYiMaeQifIg3SVQt2jHqLzZpF1S1uOmyfJ1GOVIZt1MMy9wjYrtl-A"))

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusForbidden)
}
