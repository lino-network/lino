package vote

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/global"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
	oldwire "github.com/tendermint/go-wire"
)

var (
	DelegatorSubstore = []byte{0x00}
	VoterSubstore     = []byte{0x01}
)

const returnCoinEvent = 0x1

var _ = oldwire.RegisterInterface(
	struct{ global.Event }{},
	oldwire.ConcreteType{ReturnCoinEvent{}, returnCoinEvent},
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

func (vm VoteManager) IsDelegationExist(ctx sdk.Context, voter acc.AccountKey, delegator acc.AccountKey) bool {
	store := ctx.KVStore(vm.key)
	if infoByte := store.Get(GetDelegationKey(voter, delegator)); infoByte == nil {
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
	return res.IsGTE(valRegisterFee)
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

// this method won't check if it is a legal withdraw, caller should check by itself
func (vm VoteManager) Withdraw(ctx sdk.Context, username acc.AccountKey, coin types.Coin, gm global.GlobalProxy) sdk.Error {
	voter, getErr := vm.GetVoter(ctx, username)
	if getErr != nil {
		return getErr
	}
	voter.Deposit = voter.Deposit.Minus(coin)

	if err := vm.SetVoter(ctx, username, voter); err != nil {
		return err
	}
	if err := vm.CreateReturnCoinEvent(ctx, username, coin, gm); err != nil {
		return nil
	}
	return nil
}

func (vm VoteManager) WithdrawAll(ctx sdk.Context, username acc.AccountKey, gm global.GlobalProxy) sdk.Error {
	voter, getErr := vm.GetVoter(ctx, username)
	if getErr != nil {
		return getErr
	}
	if err := vm.Withdraw(ctx, username, voter.Deposit, gm); err != nil {
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
	var delegation *Delegation
	var getErr sdk.Error

	if !vm.IsDelegationExist(ctx, voterName, delegatorName) {
		delegation = &Delegation{
			Delegator: delegatorName,
		}
	} else {
		delegation, getErr = vm.GetDelegation(ctx, voterName, delegatorName)
		if getErr != nil {
			return getErr
		}
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

func (vm VoteManager) GetAllDelegators(ctx sdk.Context, voterName acc.AccountKey) ([]acc.AccountKey, sdk.Error) {
	store := ctx.KVStore(vm.key)
	iterator := store.Iterator(subspace(GetDelegatorPrefix(voterName)))

	var delegators []acc.AccountKey

	for ; iterator.Valid(); iterator.Next() {
		delegationBytes := iterator.Value()
		var delegation Delegation
		err := vm.cdc.UnmarshalJSON(delegationBytes, &delegation)
		if err != nil {
			return nil, ErrDelegationUnmarshalError(err)
		}
		delegators = append(delegators, delegation.Delegator)
	}
	iterator.Close()
	return delegators, nil
}

func (vm VoteManager) ReturnCoinToDelegator(ctx sdk.Context, voterName acc.AccountKey, delegatorName acc.AccountKey, gm global.GlobalProxy) sdk.Error {
	voter, getErr := vm.GetVoter(ctx, voterName)
	if getErr != nil {
		return getErr
	}
	delegation, getErr := vm.GetDelegation(ctx, voterName, delegatorName)
	if getErr != nil {
		return getErr
	}

	voter.DelegatedPower = voter.DelegatedPower.Minus(delegation.Amount)
	if err := vm.CreateReturnCoinEvent(ctx, delegatorName, delegation.Amount, gm); err != nil {
		return err
	}
	if err := vm.SetVoter(ctx, voterName, voter); err != nil {
		return err
	}
	if err := vm.DeleteDelegation(ctx, voterName, delegatorName); err != nil {
		return err
	}
	return nil
}

// return coin to an user periodically
func (vm VoteManager) CreateReturnCoinEvent(ctx sdk.Context, username acc.AccountKey, amount types.Coin, gm global.GlobalProxy) sdk.Error {
	event := ReturnCoinEvent{
		Username: username,
		Amount:   amount,
	}
	if err := gm.RegisterEventAtHeight(ctx, types.Height(ctx.BlockHeight()+1000), event); err != nil {
		return err
	}
	return nil
}

func (vm VoteManager) GetVotingPower(ctx sdk.Context, voterName acc.AccountKey) (types.Coin, sdk.Error) {
	voter, getErr := vm.GetVoter(ctx, voterName)
	if getErr != nil {
		return types.Coin{}, getErr
	}
	res := voter.Deposit.Plus(voter.DelegatedPower)
	return res, nil
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

func subspace(prefix []byte) (start, end []byte) {
	end = make([]byte, len(prefix))
	copy(end, prefix)
	end[len(end)-1]++
	return prefix, end
}
