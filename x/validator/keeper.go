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
	InitGenesis(ctx sdk.Context) error
	RegisterValidator(ctx sdk.Context, username linotypes.AccountKey, valPubKey crypto.PubKey, link string) sdk.Error
	RevokeValidator(ctx sdk.Context, username linotypes.AccountKey) sdk.Error
	VoteValidator(ctx sdk.Context, username linotypes.AccountKey, votedValidators []linotypes.AccountKey) sdk.Error
	Hooks() votemn.Hooks
	GetInitValidators(ctx sdk.Context) ([]abci.ValidatorUpdate, sdk.Error)
	GetValidatorUpdates(ctx sdk.Context) ([]abci.ValidatorUpdate, sdk.Error)
	DistributeInflationToValidator(ctx sdk.Context) sdk.Error
	FireIncompetentValidator(ctx sdk.Context, byzantineValidators []abci.Evidence) sdk.Error
	UpdateSigningStats(ctx sdk.Context, voteInfos []abci.VoteInfo) sdk.Error
	// getter and setter
	GetValidator(ctx sdk.Context, username linotypes.AccountKey) (*model.Validator, sdk.Error)
	GetValidatorList(ctx sdk.Context) (*model.ValidatorList, sdk.Error)
	GetElectionVoteList(ctx sdk.Context, accKey linotypes.AccountKey) (*model.ElectionVoteList, sdk.Error)
	GetCommittingValidators(ctx sdk.Context) ([]linotypes.AccountKey, sdk.Error)
	SetValidatorList(ctx sdk.Context, lst *model.ValidatorList) sdk.Error
}

var _ ValidatorKeeper = votemn.ValidatorManager{}
