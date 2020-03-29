package test_data

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
)

type MockHttpClient struct {
	mock.Mock
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func MakeHttpResponse(status int, returnObject interface{}) *http.Response {
	bytesJson, _ := json.Marshal(returnObject)
	r := ioutil.NopCloser(bytes.NewReader(bytesJson))
	return &http.Response{
		StatusCode: status,
		Body:       r,
	}
}

type MockAuthServerConfig struct {
	mock.Mock
}

func (m *MockAuthServerConfig) GetAuthServerJwtPublicKey() string {
	return m.Called().String(0)
}

func (m *MockAuthServerConfig) GetAuthServerUrl() string {
	return m.Called().String(0)
}

func (m *MockAuthServerConfig) GetAuthServerLoginPath() string {
	return m.Called().String(0)
}

func (m *MockAuthServerConfig) GetAuthServerUsername() string {
	return m.Called().String(0)
}

func (m *MockAuthServerConfig) GetAuthServerPassword() string {
	return m.Called().String(0)
}

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) GetAccessToken() (token string, err error) {
	arguments := m.Called()
	return arguments.String(0), arguments.Error(1)
}

type MockFamilyServerConfig struct {
	mock.Mock
}

func (m *MockFamilyServerConfig) GetFamilyServerUrl() string {
	return m.Called().String(0)
}

type MockUserHasAccessToBabyChecker struct {
	mock.Mock
}

func (m MockUserHasAccessToBabyChecker) CheckUserHasAccessToBaby(userUuid string, targetUuid string) (bool, error) {
	args := m.Called(userUuid, targetUuid)
	return args.Bool(0), args.Error(1)
}
