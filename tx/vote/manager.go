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

func (vm VoteManager) IsLegalWithdraw(ctx sdk.Context, username acc.AccountKey, coin types.Coin) bool {
	voter, getErr := vm.GetVoter(ctx, username)
	if getErr != nil {
		return false
	}
	//reject if the remaining coins are less than register fee
	res := voter.Deposit.Minus(coin)
	if !res.IsGTE(valRegisterFee) {
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

func (vm VoteManager) DeleteVoter(ctx sdk.Context, username acc.AccountKey) sdk.Error {
	store := ctx.KVStore(vm.key)
	store.Delete(GetVoterKey(username))
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

func (vm VoteManager) Withdraw(ctx sdk.Context, username acc.AccountKey, coin types.Coin) sdk.Error {
	return nil
}

func (vm VoteManager) WithdrawAll(ctx sdk.Context, username acc.AccountKey) sdk.Error {
	voter, getErr := vm.GetVoter(ctx, username)
	if getErr != nil {
		return getErr
	}
	if err := vm.Withdraw(ctx, username, voter.Deposit); err != nil {
		return err
	}
	return nil
}

func (vm VoteManager) GetDelegation(ctx sdk.Context, voter acc.AccountKey, delegator acc.AccountKey) (*Delegation, sdk.Error) {
	store := ctx.KVStore(vm.key)
	delegationByte := store.Get(GetDelegationKey(voter, delegator))
	if delegationByte == nil {
		return nil, ErrGetDelegation()
	}
	delegation := new(Delegation)
	if err := vm.cdc.UnmarshalJSON(delegationByte, delegation); err != nil {
		return nil, ErrDelegationUnmarshalError(err)
	}
	return delegation, nil
}

func (vm VoteManager) SetDelegation(ctx sdk.Context, voter acc.AccountKey, delegator acc.AccountKey, delegation *Delegation) sdk.Error {
	store := ctx.KVStore(vm.key)
	delegationByte, err := vm.cdc.MarshalJSON(*delegation)
	if err != nil {
		return ErrDelegationMarshalError(err)
	}
	store.Set(GetDelegationKey(voter, delegator), delegationByte)
	return nil
}

func (vm VoteManager) AddDelegation(ctx sdk.Context, voterName acc.AccountKey, delegatorName acc.AccountKey, coin types.Coin) sdk.Error {
	delegation, getErr := vm.GetDelegation(ctx, voterName, delegatorName)
	if getErr != nil {
		return getErr
	}

	voter, getErr := vm.GetVoter(ctx, voterName)
	if getErr != nil {
		return getErr
	}

	voter.DelegatedPower = voter.DelegatedPower.Plus(coin)
	delegation.Amount = delegation.Amount.Plus(coin)

	if err := vm.SetDelegation(ctx, voterName, delegatorName, delegation); err != nil {
		return err
	}
	if err := vm.SetVoter(ctx, voterName, voter); err != nil {
		return err
	}
	return nil
}

func (vm VoteManager) DeleteDelegation(ctx sdk.Context, voter acc.AccountKey, delegator acc.AccountKey) sdk.Error {
	store := ctx.KVStore(vm.key)
	store.Delete(GetDelegationKey(voter, delegator))
	return nil
}

func (vm VoteManager) GetAllDelegators(ctx sdk.Context, username acc.AccountKey) ([]acc.AccountKey, sdk.Error) {
	return nil, nil
}

func (vm VoteManager) ReturnCoinToDelegator(ctx sdk.Context, voterName acc.AccountKey, delegatorName acc.AccountKey) sdk.Error {
	voter, getErr := vm.GetVoter(ctx, voterName)
	if getErr != nil {
		return getErr
	}
	delegation, getErr := vm.GetDelegation(ctx, voterName, delegatorName)
	if getErr != nil {
		return getErr
	}

	voter.DelegatedPower = voter.DelegatedPower.Minus(delegation.Amount)
	// TODO return coin
	if err := vm.SetVoter(ctx, voterName, voter); err != nil {
		return err
	}
	if err := vm.DeleteDelegation(ctx, voterName, delegatorName); err != nil {
		return err
	}
	return nil
}

func GetDelegatorPrefix(me acc.AccountKey) []byte {
	return append(append(DelegatorSubstore, me...), types.KeySeparator...)
}

// "delegator substore" + "me(voter)" + "my delegator"
func GetDelegationKey(me acc.AccountKey, myDelegator acc.AccountKey) []byte {
	return append(GetDelegatorPrefix(me), myDelegator...)
}

func GetVoterKey(me acc.AccountKey) []byte {
	return append(VoterSubstore, me...)
}
