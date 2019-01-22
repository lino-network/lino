package grantpermission

import (
	"time"

	"github.com/lino-network/lino/types"
)

type GrantPubKey struct {
	Username   string           `json:"username"`
	AuthTo     string           `json:"auth_to"`
	Permission types.Permission `json:"permission"`
	CreatedAt  time.Time        `json:"created_at"`
	ExpiresAt  time.Time        `json:"expires_at"`
	Amount     string           `json:"amount"`
}
