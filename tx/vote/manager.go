package vote

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
)

var (
	DelegatorSubstore = []byte{0x00}
	VoterSubstore     = []byte{0x01}
)

type VoteManager struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// The wire codec for binary encoding/decoding of accounts.
	cdc *wire.Codec
}

// NewValidatorManager returns a new ValidatorManager
func NewVoteMananger(key sdk.StoreKey) VoteManager {
	cdc := wire.NewCodec()
	vm := VoteManager{
		key: key,
		cdc: cdc,
	}

	return vm
}

func (vm VoteManager) IsVoterExist(ctx sdk.Context, accKey acc.AccountKey) bool {
	store := ctx.KVStore(vm.key)
	if infoByte := store.Get(GetVoterKey(accKey)); infoByte == nil {
		return false
	}
	return true
}

func (vm VoteManager) GetVoter(ctx sdk.Context, accKey acc.AccountKey) (*Voter, sdk.Error) {
	store := ctx.KVStore(vm.key)
	voterByte := store.Get(GetVoterKey(accKey))
	if voterByte == nil {
		return nil, ErrGetVoter()
	}
	voter := new(Voter)
	if err := vm.cdc.UnmarshalJSON(voterByte, voter); err != nil {
		return nil, ErrVoterUnmarshalError(err)
	}
	return voter, nil
}

func (vm VoteManager) SetVoter(ctx sdk.Context, accKey acc.AccountKey, voter *Voter) sdk.Error {
	store := ctx.KVStore(vm.key)
	voterByte, err := vm.cdc.MarshalJSON(*voter)
	if err != nil {
		return ErrVoterMarshalError(err)
	}
	store.Set(GetVoterKey(accKey), voterByte)
	return nil
}

func (vm VoteManager) RegisterVoter(ctx sdk.Context, username acc.AccountKey, coin types.Coin) sdk.Error {
	voter := &Voter{
		Username: username,
		Deposit:  coin,
	}
	// check minimum requirements for registering as a voter
	if !coin.IsGTE(valRegisterFee) {
		return ErrRegisterFeeNotEnough()
	}

	if setErr := vm.SetVoter(ctx, username, voter); setErr != nil {
		return setErr
	}
	return nil
}

func (vm VoteManager) Deposit(ctx sdk.Context, username acc.AccountKey, coin types.Coin) sdk.Error {
	voter, err := vm.GetVoter(ctx, username)
	if err != nil {
		return err
	}
	voter.Deposit = voter.Deposit.Plus(coin)
	if setErr := vm.SetVoter(ctx, username, voter); setErr != nil {
		return setErr
	}
	return nil
}

func GetDelegatorPrefix(me acc.AccountKey) []byte {
	return append(append(DelegatorSubstore, me...), types.KeySeparator...)
}

// "delegator substore" + "me(voter)" + "my delegator"
func GetDelegatorKey(me acc.AccountKey, myDelegator acc.AccountKey) []byte {
	return append(GetDelegatorPrefix(me), myDelegator...)
}

func GetVoterKey(me acc.AccountKey) []byte {
	return append(VoterSubstore, me...)
}
