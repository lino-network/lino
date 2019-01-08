package user

import "time"

type User struct {
	Referrer          string    `json:"referrer"`
	Username          string    `json:"username"`
	CreatedAt         time.Time `json:"created_at"`
	ResetPubKey       string    `json:"reset_public_key"`
	TransactionPubKey string    `json:"transaction_public_key"`
	AppPubKey         string    `json:"app_public_key"`
	Saving            int64     `json:"saving"`
	Sequence          int64     `json:"sequence"`
}
