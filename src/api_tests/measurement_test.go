package api_tests

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
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

	assert.Equal(t, w.Code, http.StatusOK)

	p, err := ioutil.ReadAll(w.Body)
	assert.Nil(t, err)

	var m = new(models.Measurement)
	err = json.Unmarshal(p, &m)
	assert.Nil(t, err)

	assert.Equal(t, m.Uuid, stored.Uuid)
	assert.Equal(t, m.Value, stored.Value)
	assert.Equal(t, m.Type, stored.Type)
	assert.Equal(t, m.TargetUuid, stored.TargetUuid)
	assert.Equal(t, m.Timestamp.Format(time.RFC3339), stored.Timestamp.In(time.UTC).Format(time.RFC3339))

}
