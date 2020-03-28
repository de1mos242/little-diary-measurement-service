package integrations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type HttpClient interface {
	Do(request *http.Request) (*http.Response, error)
}

type AuthServerConfig interface {
	GetAuthServerUrl() string
	GetAuthServerLoginPath() string
	GetAuthServerUsername() string
	GetAuthServerPassword() string
	GetAuthServerJwtPublicKey() string
}

type AuthIntegration struct {
	Client HttpClient
	Config AuthServerConfig
}

var (
	accessToken atomic.Value
	o           sync.Once
)

func NewAuthIntegration(client HttpClient, config AuthServerConfig) *AuthIntegration {
	return &AuthIntegration{Client: client, Config: config}
}

func (a *AuthIntegration) GetAccessToken() (token string, err error) {
	o.Do(func() {
		accessToken.Store("")
	})

	token = accessToken.Load().(string)
	if token == "" || a.isTokenExpired(token) {
		token, err = a.requestToken()
		if err != nil {
			return "", err
		}
		accessToken.Store(token)
	}
	return
}

func (a *AuthIntegration) isTokenExpired(token string) bool {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		publicKeyBytes := []byte(a.Config.GetAuthServerJwtPublicKey())
		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
		if err != nil {
			return nil, err
		}
		return verifyKey, nil
	})
	if err != nil {
		return true
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		// Check that token will be valid for at least next hour
		if !claims.VerifyExpiresAt(time.Now().Add(time.Hour).Unix(), true) {
			return true
		}
	} else {
		return true
	}
	return false
}

type LoginRequestDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginResponseDto struct {
	AccessToken string `json:"access_token"`
}

func (a *AuthIntegration) requestToken() (string, error) {
	requestDto := LoginRequestDto{
		Username: a.Config.GetAuthServerUsername(),
		Password: a.Config.GetAuthServerPassword(),
	}
	jsonBytes, err := json.Marshal(&requestDto)
	if err != nil {
		return "", err
	}
	loginUrl, err := url.Parse(a.Config.GetAuthServerUrl())
	if err != nil {
		return "", err
	}
	loginUrl.Path = a.Config.GetAuthServerLoginPath()
	request, err := http.NewRequest(http.MethodPost, loginUrl.String(), bytes.NewReader(jsonBytes))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := a.Client.Do(request)
	if err != nil {
		return "", err
	}
	if response.StatusCode != http.StatusOK {
		textData, _ := ioutil.ReadAll(response.Body)
		return "", fmt.Errorf("loing error from auth server %d: %s", response.StatusCode, textData)
	}

	var responseDto LoginResponseDto
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&responseDto)
	if err != nil {
		return "", err
	}

	return responseDto.AccessToken, nil
}
