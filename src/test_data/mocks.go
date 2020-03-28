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
