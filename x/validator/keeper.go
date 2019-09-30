package validator

//go:generate mockery -name ValidatorKeeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"

	linotypes "github.com/lino-network/lino/types"
	votemn "github.com/lino-network/lino/x/validator/manager"
	"github.com/lino-network/lino/x/validator/model"
)

type ValidatorKeeper interface {
	InitGenesis(ctx sdk.Context)
	OnBeginBlock(ctx sdk.Context, req abci.RequestBeginBlock)
	RegisterValidator(ctx sdk.Context, username linotypes.AccountKey, valPubKey crypto.PubKey, link string) sdk.Error
	RevokeValidator(ctx sdk.Context, username linotypes.AccountKey) sdk.Error
	VoteValidator(ctx sdk.Context, username linotypes.AccountKey, votedValidators []linotypes.AccountKey) sdk.Error
	DistributeInflationToValidator(ctx sdk.Context) sdk.Error
	Hooks() votemn.Hooks

	// getters
	GetInitValidators(ctx sdk.Context) ([]abci.ValidatorUpdate, sdk.Error)
	GetValidatorUpdates(ctx sdk.Context) ([]abci.ValidatorUpdate, sdk.Error)
	IsLegalValidator(ctx sdk.Context, accKey linotypes.AccountKey) bool
	GetValidator(ctx sdk.Context, username linotypes.AccountKey) (*model.Validator, sdk.Error)
	GetValidatorList(ctx sdk.Context) *model.ValidatorList
	GetElectionVoteList(ctx sdk.Context, accKey linotypes.AccountKey) *model.ElectionVoteList
	GetCommittingValidators(ctx sdk.Context) []linotypes.AccountKey
	GetCommittingValidatorVoteStatus(ctx sdk.Context) ([]model.ReceivedVotesStatus, sdk.Error)
}

var _ ValidatorKeeper = votemn.ValidatorManager{}
