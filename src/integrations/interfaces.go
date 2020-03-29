package integrations

import "net/http"

type HttpClient interface {
	Do(request *http.Request) (*http.Response, error)
}

type UserHasAccessToBabyChecker interface {
	CheckUserHasAccessToBaby(userUuid string, targetUuid string) (bool, error)
}

type AuthServerJwtPublicKeyGetter interface {
	GetAuthServerJwtPublicKey() string
}
