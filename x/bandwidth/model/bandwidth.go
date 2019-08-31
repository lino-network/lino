package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	linotypes "github.com/lino-network/lino/types"
)

// BandwidthInfo - stores info about the moving average of mps and max mps
type BandwidthInfo struct {
	GeneralMsgEMA sdk.Dec `json:"general_msg_ema"`
	AppMsgEMA     sdk.Dec `json:"app_msg_ema"`
	MaxMPS        sdk.Dec `json:"max_mps"`
}

// BlockInfo - stores info about number of tx in last block
type BlockInfo struct {
	TotalMsgSignedByApp  int64          `json:"total_tx_signed_by_app"`
	TotalMsgSignedByUser int64          `json:"total_tx_signed_by_user"`
	CurMsgFee            linotypes.Coin `json:"cur_msg_fee"`
}

// AppBandwidthInfo - stores info about each app's bandwidth
type AppBandwidthInfo struct {
	MaxBandwidthCredit sdk.Dec `json:"max_bandwidth_credit"`
	CurBandwidthCredit sdk.Dec `json:"cur_bandwidth_credit"`
	MessagesInCurBlock int64   `json:"messages_in_cur_block"`
	ExpectedMPS        sdk.Dec `json:"expected_mps"`
	LastRefilledAt     int64   `json:"last_refilled_at"`
}
