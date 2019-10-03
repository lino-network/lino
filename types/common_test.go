package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestIsUsername(t *testing.T) {
	testCases := map[string]struct {
		accountKeys []AccountKey
		expectRes   bool
	}{
		"username is too short": {
			accountKeys: []AccountKey{"us"},
			expectRes:   false,
		},
		"username is too long": {
			accountKeys: []AccountKey{"useruseruseruseruseru"},
			expectRes:   false,
		},
		"invalid characters": {
			accountKeys: []AccountKey{"register#", "_register", "-register", "reg@ister", "re--gister",
				"reg*ister", "register!", "register()", "reg$ister", "reg ister", " register", "re_-gister",
				"reg=ister", "register^", "register.", "reg$ister,", "Register", "r__egister", "reGister",
				"r_--gister", "re.-gister", ".re-gister", "re-gister.", "register_", "register-", "a.2.2.-.-..2",
				".register", "register..", "_.register", "123123", "reg--ster", "reg$ster", "re%gister", "regist\"er",
				"reg?ster", "reg:ster", "regi<ster", "regi>ster", "reg{ster", "regi}ster", "reg'ster", "reg`ster"},
			expectRes: false,
		},
		"address": {
			accountKeys: []AccountKey{AccountKey(secp256k1.GenPrivKey().PubKey().Address())},
			expectRes:   false,
		},
		"valid username": {
			accountKeys: []AccountKey{"register", "re.gister", "re-gister", "reg", "registerregisterregi"},
			expectRes:   true,
		},
	}

	for testName, tc := range testCases {
		for _, accKey := range tc.accountKeys {
			res := accKey.IsValid()
			assert.Equal(t, tc.expectRes, res, "%s", testName)
		}
	}
}
