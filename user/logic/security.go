package logic

import (
	"context"

	"github.com/junqirao/gocomponents/security"
)

const (
	securityProviderNameTransport = "c_user_transport"
)

var (
	Security = &sec{}
)

type sec struct{}

func (l *sec) GetTransportPublicKeyPem(ctx context.Context) (p string, err error) {
	provider, err := security.GetProvider(ctx, securityProviderNameTransport, security.StorageTypeMysql)
	if err != nil {
		return
	}
	p = provider.GetPublicKeyPem()
	return
}
