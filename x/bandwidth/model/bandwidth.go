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

// LastBlockInfo - stores info about number of tx in last block
type LastBlockInfo struct {
	TotalMsgSignedByApp  uint32 `json:"total_tx_signed_by_app"`
	TotalMsgSignedByUser uint32 `json:"total_tx_signed_by_user"`
}

// BlockStatsCache - message stats in-memory cache of current block
type BlockStatsCache struct {
	CurMsgFee            sdk.Dec `json:"cur_msg_fee"`
	TotalMsgSignedByApp  uint32  `json:"total_tx_signed_by_app"`
	TotalMsgSignedByUser uint32  `json:"total_tx_signed_by_user"`
}
