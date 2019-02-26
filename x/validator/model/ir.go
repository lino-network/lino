package model

import (
	"github.com/lino-network/lino/types"
)

// ABCIPubKeyIR - type changed during upgrade.
type ABCIPubKeyIR struct {
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

// ABCIValidatorIR - type changed during upgrade.
type ABCIValidatorIR struct {
	Address []byte       `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	PubKey  ABCIPubKeyIR `protobuf:"bytes,2,opt,name=pub_key,json=pubKey" json:"pub_key"`
	Power   int64        `protobuf:"varint,3,opt,name=power,proto3" json:"power,omitempty"`
}

// ValidatorIR - ABCIValidator internally changed.
type ValidatorIR struct {
	ABCIValidator   ABCIValidatorIR
	Username        types.AccountKey `json:"username"`
	Deposit         types.Coin       `json:"deposit"`
	AbsentCommit    int64            `json:"absent_commit"`
	ByzantineCommit int64            `json:"byzantine_commit"`
	ProducedBlocks  int64            `json:"produced_blocks"`
	Link            string           `json:"link"`
}

// ValidatorRowIR - pk: (Username)
type ValidatorRowIR struct {
	Username types.AccountKey `json:"username"`
	// XXX(yumin): type changed.
	Validator ValidatorIR `json:"validator"`
}

// ValidatorTablesIR - Validators changed.
type ValidatorTablesIR struct {
	Validators    []ValidatorRowIR `json:"validators"`
	ValidatorList ValidatorListRow `json:"validator_list"`
}
