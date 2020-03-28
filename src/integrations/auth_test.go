package integrations

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"little-diary-measurement-service/src/config"
	"little-diary-measurement-service/src/test_data"
	"net/http"
	"testing"
)

func TestAuthIntegration_RequestToken(t *testing.T) {
	type fields struct {
		Client HttpClient
		Config AuthServerConfig
	}

	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "Test login request",
			fields: fields{
				Client: func() HttpClient {
					mockObj := new(test_data.MockHttpClient)
					mockObj.On("Do", mock.Anything).
						Return(test_data.MakeHttpResponse(200, &LoginResponseDto{AccessToken: "jwt_token"}), nil)
					return mockObj
				}(),
				Config: func() AuthServerConfig {
					mockObj := new(test_data.MockAuthServerConfig)
					mockObj.On("GetAuthServerUrl", mock.Anything).Return("https://littlediary.net:8080")
					mockObj.On("GetAuthServerLoginPath", mock.Anything).Return("/auth/login")
					mockObj.On("GetAuthServerUsername", mock.Anything).Return("some name")
					mockObj.On("GetAuthServerPassword", mock.Anything).Return("some pass")
					return mockObj
				}()},
			want:    "jwt_token",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthIntegration{
				Client: tt.fields.Client,
				Config: tt.fields.Config,
			}
			got, err := a.requestToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("requestToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("requestToken() got = %v, want %v", got, tt.want)
			}

			tt.fields.Config.(*test_data.MockAuthServerConfig).AssertExpectations(t)

			actualReq := tt.fields.Client.(*test_data.MockHttpClient).Calls[0].Arguments.Get(0).(*http.Request)
			reqBody, _ := ioutil.ReadAll(actualReq.Body)
			expectedJson := fmt.Sprintf(`{"username":"%s","password":"%s"}`,
				tt.fields.Config.GetAuthServerUsername(),
				tt.fields.Config.GetAuthServerPassword())
			assert.Equal(t, reqBody, []byte(expectedJson))

			assert.Equal(t, actualReq.URL.String(), fmt.Sprintf("%s%s",
				tt.fields.Config.GetAuthServerUrl(),
				tt.fields.Config.GetAuthServerLoginPath()))
			assert.Equal(t, actualReq.Header.Get("Content-Type"), "application/json")
		})
	}
}

func TestAuthIntegration_isTokenExpired(t *testing.T) {
	type fields struct {
		Client HttpClient
		Config AuthServerConfig
	}
	type args struct {
		token string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "Test valid token",
			fields: fields{Config: &config.Config},
			args:   args{"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJpYXQiOjE1ODU0MjM2NTcsIm5iZiI6MTU4NTQyMzY1NywianRpIjoiMzViNjk0MWUtZmNhNC00MDhmLWI5MDItNjhmODZkOTQ0YThhIiwiZXhwIjo0MTAyNDQ0ODAwLCJpZGVudGl0eSI6MiwiZnJlc2giOmZhbHNlLCJ0eXBlIjoiYWNjZXNzIiwidXNlcl9jbGFpbXMiOnsicm9sZSI6InRlY2giLCJ1dWlkIjoiMmM4NmIwNmItMjY4MS00ZjA4LTk4YzAtZjhkYzcxMzg5ZTBhIiwicmVzb3VyY2VzIjpbImZhbWlseV9yZWFkIl19fQ.HbFiZtSQMfng6kjRXSQa27lRLFTMAlAJ2iDwFTSVzQlrPjm3vQi23ktdarIYcJzGssu1CZl6ajtdMi1sGI7uRaZ_ztYmciGqInvwTfhli3MmjLtvjZ9N5bwLS-P0qQaLMsng7vKaLIHSPC3DLPM7wbVZOcNvoP6BFDHFbNByffjFoKtKGAXpjnFJGEKIwGw_rZWVINb1OBn_WaYuzsJqjk_4IaDuN9CTEEFVN_hM2th44G-rYwX-FnQZFHueI_p5KfYuKGoYkjAgja9xHLnKdQC5eUy3FRl1pnZ09I91MXr4dNuLtjlj9iGlA_b6NrF3KrNE7r2fga99aaaTA_FqFA"},
			want:   false,
		},
		{
			name:   "Test expired token",
			fields: fields{Config: &config.Config},
			args:   args{"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJpYXQiOjE1ODU0MjM2NTcsIm5iZiI6MTU4NTQyMzY1NywianRpIjoiMzViNjk0MWUtZmNhNC00MDhmLWI5MDItNjhmODZkOTQ0YThhIiwiZXhwIjoxNTg1MzUzNjAwLCJpZGVudGl0eSI6MiwiZnJlc2giOmZhbHNlLCJ0eXBlIjoiYWNjZXNzIiwidXNlcl9jbGFpbXMiOnsicm9sZSI6InRlY2giLCJ1dWlkIjoiMmM4NmIwNmItMjY4MS00ZjA4LTk4YzAtZjhkYzcxMzg5ZTBhIiwicmVzb3VyY2VzIjpbImZhbWlseV9yZWFkIl19fQ.ADSRos61jPFV0xGZ8H2toY9jTo20hKmlzmb4ZJ__HybD5__Vger77pl9ZdO-ZGmE4V0ce_vvF5Nw1p1z9oHsi4mDFz65a9DcWg2vrm3gt_wpuU3vespKiJY87oZLUidUoZfWwp9gGBbJ1VMa6pMzMM_lFLBBVUAfUhXINGiTi6eiVyF0-7rYPaFCpTVfJA5KJKDpqXS4kI0ajre6DGwkiq2MtkNl8nrrTpD5tRcYGD-hZ_Iiqrhz0HAE4Mx_ZKqxLXN2-ilaX8-puesJyjw5s3_5R6d3sj2Bp2oqHZXKT-FHV06Hrq_y5wE786uRQo38PBLPnze8uQ0kTNYrSwtO4w"},
			want:   true,
		},
		{
			name:   "Test invalid signature",
			fields: fields{Config: &config.Config},
			args:   args{"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJpYXQiOjE1ODU0MjM2NTcsIm5iZiI6MTU4NTQyMzY1NywianRpIjoiMzViNjk0MWUtZmNhNC00MDhmLWI5MDItNjhmODZkOTQ0YThhIiwiZXhwIjo0MTAyNDQ0ODAwLCJpZGVudGl0eSI6MiwiZnJlc2giOmZhbHNlLCJ0eXBlIjoiYWNjZXNzIiwidXNlcl9jbGFpbXMiOnsicm9sZSI6InRlY2giLCJ1dWlkIjoiMmM4NmIwNmItMjY4MS00ZjA4LTk4YzAtZjhkYzcxMzg5ZTBhIiwicmVzb3VyY2VzIjpbImZhbWlseV9yZWFkIl19fQ.HbFiZtSQMfng6kjRXSQa27lRLFTMAlAJ2iDwFTSVzQlrPjm3vQi23ktdarIYcJzGssu1CZl6ajtdMi1sGI7uRaZ_ztYmciGqInvwTfhli3MmjLtvjZ9N5bwLS-P0qQaLMsng7vKaLIHSPC3DLPM7wbVZOcNvoP6BFDHFbNByffjFoKtKGAXpjnFJGEKIwGw_rZWVINb1OBn_WaYuzsJqjk_4IaDuN9CTEEFVN_hM2th44G-rYwX-FnQZFHueI_p5KfYuKGoYkjAgja9xHLnKdQC5eUy3FRl1pnZ09I91MXr4dNuLtjlj9iGlA_b6NrF3KrNE7r2fga99aaaTA_4444"},
			want:   true,
		},
		{
			name:   "Test corrupted token",
			fields: fields{Config: &config.Config},
			args:   args{"not a jwt"},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthIntegration{
				Client: tt.fields.Client,
				Config: tt.fields.Config,
			}
			if got := a.isTokenExpired(tt.args.token); got != tt.want {
				t.Errorf("isTokenExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthIntegration_GetAccessToken(t *testing.T) {
	expiredToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJpYXQiOjE1ODU0MjM2NTcsIm5iZiI6MTU4NTQyMzY1NywianRpIjoiMzViNjk0MWUtZmNhNC00MDhmLWI5MDItNjhmODZkOTQ0YThhIiwiZXhwIjoxNTg1MzUzNjAwLCJpZGVudGl0eSI6MiwiZnJlc2giOmZhbHNlLCJ0eXBlIjoiYWNjZXNzIiwidXNlcl9jbGFpbXMiOnsicm9sZSI6InRlY2giLCJ1dWlkIjoiMmM4NmIwNmItMjY4MS00ZjA4LTk4YzAtZjhkYzcxMzg5ZTBhIiwicmVzb3VyY2VzIjpbImZhbWlseV9yZWFkIl19fQ.ADSRos61jPFV0xGZ8H2toY9jTo20hKmlzmb4ZJ__HybD5__Vger77pl9ZdO-ZGmE4V0ce_vvF5Nw1p1z9oHsi4mDFz65a9DcWg2vrm3gt_wpuU3vespKiJY87oZLUidUoZfWwp9gGBbJ1VMa6pMzMM_lFLBBVUAfUhXINGiTi6eiVyF0-7rYPaFCpTVfJA5KJKDpqXS4kI0ajre6DGwkiq2MtkNl8nrrTpD5tRcYGD-hZ_Iiqrhz0HAE4Mx_ZKqxLXN2-ilaX8-puesJyjw5s3_5R6d3sj2Bp2oqHZXKT-FHV06Hrq_y5wE786uRQo38PBLPnze8uQ0kTNYrSwtO4w"
	token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJpYXQiOjE1ODU0MjM2NTcsIm5iZiI6MTU4NTQyMzY1NywianRpIjoiMzViNjk0MWUtZmNhNC00MDhmLWI5MDItNjhmODZkOTQ0YThhIiwiZXhwIjo0MTAyNDQ0ODAwLCJpZGVudGl0eSI6MiwiZnJlc2giOmZhbHNlLCJ0eXBlIjoiYWNjZXNzIiwidXNlcl9jbGFpbXMiOnsicm9sZSI6InRlY2giLCJ1dWlkIjoiMmM4NmIwNmItMjY4MS00ZjA4LTk4YzAtZjhkYzcxMzg5ZTBhIiwicmVzb3VyY2VzIjpbImZhbWlseV9yZWFkIl19fQ.HbFiZtSQMfng6kjRXSQa27lRLFTMAlAJ2iDwFTSVzQlrPjm3vQi23ktdarIYcJzGssu1CZl6ajtdMi1sGI7uRaZ_ztYmciGqInvwTfhli3MmjLtvjZ9N5bwLS-P0qQaLMsng7vKaLIHSPC3DLPM7wbVZOcNvoP6BFDHFbNByffjFoKtKGAXpjnFJGEKIwGw_rZWVINb1OBn_WaYuzsJqjk_4IaDuN9CTEEFVN_hM2th44G-rYwX-FnQZFHueI_p5KfYuKGoYkjAgja9xHLnKdQC5eUy3FRl1pnZ09I91MXr4dNuLtjlj9iGlA_b6NrF3KrNE7r2fga99aaaTA_FqFA"

	a := &AuthIntegration{
		Client: func() HttpClient {
			mockObj := new(test_data.MockHttpClient)
			mockObj.On("Do", mock.Anything).
				Return(test_data.MakeHttpResponse(200, &LoginResponseDto{AccessToken: expiredToken}), nil)
			return mockObj
		}(),
		Config: func() AuthServerConfig {
			mockObj := new(test_data.MockAuthServerConfig)
			mockObj.On("GetAuthServerUrl", mock.Anything).Return("https://littlediary.net:8080")
			mockObj.On("GetAuthServerLoginPath", mock.Anything).Return("/auth/login")
			mockObj.On("GetAuthServerUsername", mock.Anything).Return("some name")
			mockObj.On("GetAuthServerPassword", mock.Anything).Return("some pass")
			mockObj.On("GetAuthServerJwtPublicKey", mock.Anything).Return(config.Config.GetAuthServerJwtPublicKey())
			return mockObj
		}(),
	}
	got, err := a.GetAccessToken()
	assert.Nil(t, err)
	assert.Equal(t, got, expiredToken)

	a.Client = new(test_data.MockHttpClient)
	a.Client.(*test_data.MockHttpClient).On("Do", mock.Anything).
		Return(test_data.MakeHttpResponse(200, &LoginResponseDto{AccessToken: token}), nil)

	got, err = a.GetAccessToken()
	assert.Nil(t, err)
	assert.Equal(t, got, token)

	got, err = a.GetAccessToken()
	assert.Nil(t, err)
	assert.Equal(t, got, token)

	a.Client.(*test_data.MockHttpClient).AssertNumberOfCalls(t, "Do", 1)
}
