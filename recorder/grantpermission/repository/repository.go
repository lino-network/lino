package repository

import (
	errors "github.com/lino-network/lino/recorder/errors"
	"github.com/lino-network/lino/recorder/grantpermission"
)

type GrantPermissionRepository interface {
	Add(grantPubKey *grantpermission.GrantPubKey) errors.Error
	Get(username, authTo string) (*grantpermission.GrantPubKey, errors.Error)
	SetAmount(username, authTo string, amount string) errors.Error
	Update(grantPubKey *grantpermission.GrantPubKey) errors.Error
	Delete(username, authTo string) errors.Error
	IsEnable() bool
}
