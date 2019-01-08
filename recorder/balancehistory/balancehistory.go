package balancehistory

import (
	"time"

	"github.com/lino-network/lino/types"
)

type BalanceHistory struct {
	ID         int64                    `json:"id"`
	Username   string                   `json:"username"`
	FromUser   string                   `json:"from_user"`
	ToUser     string                   `json:"to_user"`
	Amount     int64                    `json:"amount"`
	Balance    int64                    `json:"balance"`
	DetailType types.TransferDetailType `json:"detail_type"`
	CreatedAt  time.Time                `json:"created_at"`
	Memo       string                   `json:"memo"`
}
