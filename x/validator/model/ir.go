package model

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/lino-network/lino/types"
)

// ABCIPubKeyIR.
type ABCIPubKeyIR struct {
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func NewABCIPubKeyIRFromTM(pubkey crypto.PubKey) ABCIPubKeyIR {
	key := tmtypes.TM2PB.PubKey(pubkey)
	return ABCIPubKeyIR{
		Type: key.Type,
		Data: key.Data,
	}
}

func (pubkey ABCIPubKeyIR) ToTM() crypto.PubKey {
	rst, err := tmtypes.PB2TM.PubKey(abci.PubKey{
		Type: pubkey.Type,
		Data: pubkey.Data,
	})
	if err != nil {
		panic(err)
	}
	return rst
}

// ABCIValidatorIR.
type ABCIValidatorIR struct {
	Address []byte `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Power   int64  `protobuf:"varint,3,opt,name=power,proto3" json:"power,omitempty"`
}

// ValidatorIR
type ValidatorIR struct {
	ABCIValidator   ABCIValidatorIR  `json:"abci_validator"`
	PubKey          ABCIPubKeyIR     `json:"pub_key"`
	Username        types.AccountKey `json:"username"`
	ReceivedVotes   types.Coin       `json:"received_votes"`
	HasRevoked      bool             `json:"has_revoked"`
	AbsentCommit    int64            `json:"absent_commit"`
	ByzantineCommit int64            `json:"byzantine_commit"`
	ProducedBlocks  int64            `json:"produced_blocks"`
	Link            string           `json:"link"`
}

type ElectionVoteIR struct {
	ValidatorName types.AccountKey `json:"validator_name"`
	Vote          types.Coin       `json:"votes"`
}

type ElectionVoteListIR struct {
	Username      types.AccountKey `json:"username"`
	ElectionVotes []ElectionVoteIR `json:"election_votes"`
}

// ValidatorList
type ValidatorListIR struct {
	Oncall             []types.AccountKey `json:"oncall"`
	Standby            []types.AccountKey `json:"standby"`
	Candidates         []types.AccountKey `json:"candidates"`
	Jail               []types.AccountKey `json:"jail"`
	PreBlockValidators []types.AccountKey `json:"pre_block_validators"`
	LowestOncallVotes  types.Coin         `json:"lowest_oncall_votes"`
	LowestOncall       types.AccountKey   `json:"lowest_oncall"`
	LowestStandbyVotes types.Coin         `json:"lowest_standby_votes"`
	LowestStandby      types.AccountKey   `json:"lowest_standby"`
}

// ValidatorTablesIR - Validators changed.
type ValidatorTablesIR struct {
	Version    int                  `json:"version"`
	Validators []ValidatorIR        `json:"validators"`
	Votes      []ElectionVoteListIR `json:"votes"`
	List       ValidatorListIR      `json:"list"`
}
