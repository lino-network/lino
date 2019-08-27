package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BandwidthInfo - stores info about the moving average of mps and max mps
type BandwidthInfo struct {
	GeneralMsgEMA sdk.Dec `json:"general_msg_ema"`
	AppMsgEMA     sdk.Dec `json:"app_msg_ema"`
	MaxMPS        sdk.Dec `json:"max_mps"`
}

// CurBlockInfo - stores info about number of tx in current block
type CurBlockInfo struct {
	TotalMsgSignedByApp  uint32  `json:"total_tx_signed_by_app"`
	TotalMsgSignedByUser uint32  `json:"total_tx_signed_by_user"`
	CurMsgFee            sdk.Dec `json:"cur_msg_fee"`
}
