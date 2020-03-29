package integrations

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type FamilyServerConfig interface {
	GetFamilyServerUrl() string
}

type AuthService interface {
	GetAccessToken() (token string, err error)
}

type FamilyIntegration struct {
	Client      HttpClient
	Config      FamilyServerConfig
	AuthService AuthService
}

type CheckAccessResponseDto struct {
	HasAccess bool `json:"has_access"`
}

func (f *FamilyIntegration) CheckUserHasAccessToBaby(userUuid string, targetUuid string) (bool, error) {
	accessToken, err := f.AuthService.GetAccessToken()
	if err != nil {
		return false, nil
	}

	accessUrl, err := url.Parse(f.Config.GetFamilyServerUrl())
	if err != nil {
		return false, err
	}
	accessUrl.Path = fmt.Sprintf("/v1/access/%s/baby/%s", userUuid, targetUuid)
	request, err := http.NewRequest(http.MethodGet, accessUrl.String(), nil)
	if err != nil {
		return false, err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	response, err := f.Client.Do(request)
	if err != nil {
		return false, err
	}
	if response.StatusCode != http.StatusOK {
		textData, _ := ioutil.ReadAll(response.Body)
		return false, fmt.Errorf("check access error from family server %d: %s", response.StatusCode, textData)
	}

	var responseDto CheckAccessResponseDto
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&responseDto)
	if err != nil {
		return false, err
	}

	return responseDto.HasAccess, nil
}
