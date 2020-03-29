package integrations

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"little-diary-measurement-service/src/test_data"
	"net/http"
	"testing"
)

func TestFamilyIntegration_CheckUserHasAccessToBaby(t *testing.T) {
	type fields struct {
		Client      HttpClient
		Config      FamilyServerConfig
		AuthService AuthService
	}
	type args struct {
		userUuid   string
		targetUuid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Test check good",
			fields: fields{
				Client: func() HttpClient {
					mockObj := new(test_data.MockHttpClient)
					mockObj.On("Do", mock.Anything).
						Return(test_data.MakeHttpResponse(200, &CheckAccessResponseDto{HasAccess: true}), nil)
					return mockObj
				}(),
				Config: func() FamilyServerConfig {
					mockObj := test_data.MockFamilyServerConfig{}
					mockObj.On("GetFamilyServerUrl").Return("https://family.little-diary.net")
					return &mockObj
				}(),
				AuthService: func() AuthService {
					mockObj := test_data.MockAuthService{}
					mockObj.On("GetAccessToken", mock.Anything).Return("fake_token", nil)
					return &mockObj
				}(),
			},
			args: args{
				userUuid:   "74822532-56e4-4ff0-8890-3247a17433cc",
				targetUuid: "11111111-3333-412d-ade7-47be43827d68",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Test check failure",
			fields: fields{
				Client: func() HttpClient {
					mockObj := new(test_data.MockHttpClient)
					mockObj.On("Do", mock.Anything).
						Return(test_data.MakeHttpResponse(200, &CheckAccessResponseDto{HasAccess: false}), nil)
					return mockObj
				}(),
				Config: func() FamilyServerConfig {
					mockObj := test_data.MockFamilyServerConfig{}
					mockObj.On("GetFamilyServerUrl").Return("https://family.little-diary.net")
					return &mockObj
				}(),
				AuthService: func() AuthService {
					mockObj := test_data.MockAuthService{}
					mockObj.On("GetAccessToken", mock.Anything).Return("fake_token", nil)
					return &mockObj
				}(),
			},
			args: args{
				userUuid:   "74822532-56e4-4ff0-8890-3247a17433cc",
				targetUuid: "11111111-3333-412d-ade7-47be43827d68",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Test target not found",
			fields: fields{
				Client: func() HttpClient {
					mockObj := new(test_data.MockHttpClient)
					mockObj.On("Do", mock.Anything).
						Return(test_data.MakeHttpResponse(404, &CheckAccessResponseDto{HasAccess: false}), nil)
					return mockObj
				}(),
				Config: func() FamilyServerConfig {
					mockObj := test_data.MockFamilyServerConfig{}
					mockObj.On("GetFamilyServerUrl").Return("https://family.little-diary.net")
					return &mockObj
				}(),
				AuthService: func() AuthService {
					mockObj := test_data.MockAuthService{}
					mockObj.On("GetAccessToken", mock.Anything).Return("fake_token", nil)
					return &mockObj
				}(),
			},
			args: args{
				userUuid:   "74822532-56e4-4ff0-8890-3247a17433cc",
				targetUuid: "11111111-3333-412d-ade7-47be43827d68",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FamilyIntegration{
				Client:      tt.fields.Client,
				Config:      tt.fields.Config,
				AuthService: tt.fields.AuthService,
			}
			got, err := f.CheckUserHasAccessToBaby(tt.args.userUuid, tt.args.targetUuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckUserHasAccessToBaby() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckUserHasAccessToBaby() got = %v, want %v", got, tt.want)
			}

			actualReq := tt.fields.Client.(*test_data.MockHttpClient).Calls[0].Arguments.Get(0).(*http.Request)
			assert.Equal(t, actualReq.URL.String(), fmt.Sprintf("%s/v1/access/%s/baby/%s",
				tt.fields.Config.GetFamilyServerUrl(),
				tt.args.userUuid,
				tt.args.targetUuid))
			assert.Equal(t, actualReq.Header.Get("Authorization"), fmt.Sprintf("Bearer %s", "fake_token"))
		})
	}
}
