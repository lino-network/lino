package account

import (
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	user := types.AccountKey("user")
	lam := newLinoAccountManager()
	acc := NewLinoAccount(user, lam)
	assert.Nil(t, acc.accountInfo)
	assert.Nil(t, acc.accountBank)
	assert.Nil(t, acc.accountMeta)
	assert.Nil(t, acc.follower)
	assert.Nil(t, acc.following)
}
