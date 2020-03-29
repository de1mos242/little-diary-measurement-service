package common

import "little-diary-measurement-service/src/integrations"

type ServiceLocator struct {
	PublicKeyGetter            integrations.AuthServerJwtPublicKeyGetter
	UserHasAccessToBabyChecker integrations.UserHasAccessToBabyChecker
}
