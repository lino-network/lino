package validator

//go:generate mockery -name ValidatorKeeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/validator/manager"
	"github.com/lino-network/lino/x/validator/model"
)

type ValidatorKeeper interface {
	InitGenesis(ctx sdk.Context) error
	RegisterValidator(ctx sdk.Context, username linotypes.AccountKey, valPubKey crypto.PubKey, link string) sdk.Error
	RevokeValidator(ctx sdk.Context, username linotypes.AccountKey) sdk.Error
	VoteValidator(ctx sdk.Context, username linotypes.AccountKey, votedValidators []linotypes.AccountKey) sdk.Error
	// getter
	GetValidator(ctx sdk.Context, username linotypes.AccountKey) (*model.Validator, sdk.Error)
	GetValidatorList(ctx sdk.Context) (*model.ValidatorList, sdk.Error)
	GetElectionVoteList(ctx sdk.Context, accKey linotypes.AccountKey) (*model.ElectionVoteList, sdk.Error)
}

var _ ValidatorKeeper = manager.ValidatorManager{}
