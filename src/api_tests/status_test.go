package api_tests

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"little-diary-measurement-service/src/common"
	"little-diary-measurement-service/src/config"
	"little-diary-measurement-service/src/integrations"
	"little-diary-measurement-service/src/router"
	"little-diary-measurement-service/src/test_data"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckHealth(t *testing.T) {
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

	req := httptest.NewRequest("GET", "/status/health", nil)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)
}
