package logic

import (
	"context"

	"github.com/junqirao/gocomponents/security"
)

const (
	securityProviderNameTransport = "c_user_transport"
)

func GetSecurityPublicKeyPem(ctx context.Context) (p string, err error) {
	provider, err := security.GetProvider(ctx, securityProviderNameTransport, security.StorageTypeMysql)
	if err != nil {
		return
	}
	p = provider.GetPublicKeyPem()
	return
}
