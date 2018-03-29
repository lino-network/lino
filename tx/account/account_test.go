package account

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	user := AccountKey("user")
	lam := newLinoAccountManager()
	acc := NewProxyAccount(user, &lam)
	assert.Nil(t, acc.accountInfo)
	assert.Nil(t, acc.accountBank)
	assert.Nil(t, acc.accountMeta)
}
