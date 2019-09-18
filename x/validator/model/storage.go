package model

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/validator/types"
)

var (
	validatorSubstore        = []byte{0x00}
	validatorListSubstore    = []byte{0x01}
	electionVoteListSubstore = []byte{0x02}
)

type ValidatorStorage struct {
	key sdk.StoreKey
	cdc *wire.Codec
}

func NewValidatorStorage(key sdk.StoreKey) ValidatorStorage {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	vs := ValidatorStorage{
		key: key,
		cdc: cdc,
	}
	return vs
}

func (vs ValidatorStorage) InitGenesis(ctx sdk.Context) error {
	lst := &ValidatorList{
		LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
		LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
	}

	if err := vs.SetValidatorList(ctx, lst); err != nil {
		return err
	}
	return nil
}

func (vs ValidatorStorage) DoesValidatorExist(ctx sdk.Context, accKey linotypes.AccountKey) bool {
	store := ctx.KVStore(vs.key)
	return store.Has(GetValidatorKey(accKey))
}

func (vs ValidatorStorage) GetValidator(ctx sdk.Context, accKey linotypes.AccountKey) (*Validator, sdk.Error) {
	store := ctx.KVStore(vs.key)
	validatorByte := store.Get(GetValidatorKey(accKey))
	if validatorByte == nil {
		return nil, types.ErrValidatorNotFound()
	}
	validator := new(Validator)
	vs.cdc.MustUnmarshalBinaryLengthPrefixed(validatorByte, validator)
	return validator, nil
}

func (vs ValidatorStorage) SetValidator(ctx sdk.Context, accKey linotypes.AccountKey, validator *Validator) sdk.Error {
	store := ctx.KVStore(vs.key)
	validatorByte := vs.cdc.MustMarshalBinaryLengthPrefixed(*validator)
	store.Set(GetValidatorKey(accKey), validatorByte)
	return nil
}

func (vs ValidatorStorage) DeleteValidator(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	store := ctx.KVStore(vs.key)
	store.Delete(GetValidatorKey(username))
	return nil
}

func (vs ValidatorStorage) GetValidatorList(ctx sdk.Context) (*ValidatorList, sdk.Error) {
	store := ctx.KVStore(vs.key)
	listByte := store.Get(GetValidatorListKey())
	if listByte == nil {
		return nil, types.ErrValidatorListNotFound()
	}
	lst := new(ValidatorList)
	vs.cdc.MustUnmarshalBinaryLengthPrefixed(listByte, lst)
	return lst, nil
}

func (vs ValidatorStorage) SetValidatorList(ctx sdk.Context, lst *ValidatorList) sdk.Error {
	store := ctx.KVStore(vs.key)
	listByte := vs.cdc.MustMarshalBinaryLengthPrefixed(*lst)
	store.Set(GetValidatorListKey(), listByte)
	return nil
}

func (vs ValidatorStorage) DoesElectionVoteListExist(ctx sdk.Context, accKey linotypes.AccountKey) bool {
	store := ctx.KVStore(vs.key)
	return store.Has(GetElectionVoteListKey(accKey))
}

func (vs ValidatorStorage) GetElectionVoteList(ctx sdk.Context, accKey linotypes.AccountKey) (*ElectionVoteList, sdk.Error) {
	store := ctx.KVStore(vs.key)
	lstByte := store.Get(GetElectionVoteListKey(accKey))
	if lstByte == nil {
		return nil, types.ErrElectionListNotFound()
	}
	lst := new(ElectionVoteList)
	vs.cdc.MustUnmarshalBinaryLengthPrefixed(lstByte, lst)
	return lst, nil
}

func (vs ValidatorStorage) SetElectionVoteList(ctx sdk.Context, accKey linotypes.AccountKey, lst *ElectionVoteList) sdk.Error {
	store := ctx.KVStore(vs.key)
	lstByte := vs.cdc.MustMarshalBinaryLengthPrefixed(*lst)
	store.Set(GetElectionVoteListKey(accKey), lstByte)
	return nil
}

// Export state of validators.
// func (vs ValidatorStorage) Export(ctx sdk.Context) *ValidatorTables {
// 	tables := &ValidatorTables{}
// 	store := ctx.KVStore(vs.key)
// 	// export table.validators
// 	func() {
// 		itr := sdk.KVStorePrefixIterator(store, validatorSubstore)
// 		defer itr.Close()
// 		for ; itr.Valid(); itr.Next() {
// 			k := itr.Key()
// 			username := linotypes.AccountKey(k[1:])
// 			val, err := vs.GetValidator(ctx, username)
// 			if err != nil {
// 				panic("failed to read validator: " + err.Error())
// 			}
// 			row := ValidatorRow{
// 				Username:  username,
// 				Validator: *val,
// 			}
// 			tables.Validators = append(tables.Validators, row)
// 		}
// 	}()
// 	// export table.validatorList
// 	list, err := vs.GetValidatorList(ctx)
// 	if err != nil {
// 		panic("failed to get validator list: " + err.Error())
// 	}
// 	tables.ValidatorList = ValidatorListRow{
// 		List: *list,
// 	}
// 	return tables
// }

// Import from tablesIR.
// func (vs ValidatorStorage) Import(ctx sdk.Context, tb *ValidatorTablesIR) {
// 	check := func(e error) {
// 		if e != nil {
// 			panic("[vs] Failed to import: " + e.Error())
// 		}
// 	}
// 	// import table.Validators
// 	for _, v := range tb.Validators {
// 		pubkey, err := tmtypes.PB2TM.PubKey(abci.PubKey{
// 			Type: v.Validator.ABCIValidator.PubKey.Type,
// 			Data: v.Validator.ABCIValidator.PubKey.Data,
// 		})
// 		check(err)
// 		err = vs.SetValidator(ctx, v.Username, &ValidatorV1{
// 			ABCIValidator: abci.Validator{
// 				Address: v.Validator.ABCIValidator.Address,
// 				Power:   v.Validator.ABCIValidator.Power,
// 			},
// 			PubKey:          pubkey,
// 			Username:        v.Validator.Username,
// 			Deposit:         v.Validator.Deposit,
// 			AbsentCommit:    v.Validator.AbsentCommit,
// 			ByzantineCommit: v.Validator.ByzantineCommit,
// 			ProducedBlocks:  v.Validator.ProducedBlocks,
// 			Link:            v.Validator.Link,
// 		})
// 		check(err)
// 	}
// 	// import ValidatorList
// 	err := vs.SetValidatorList(ctx, &tb.ValidatorList.List)
// 	check(err)
// }

func GetValidatorKey(accKey linotypes.AccountKey) []byte {
	return append(validatorSubstore, accKey...)
}

func GetElectionVoteListKey(accKey linotypes.AccountKey) []byte {
	return append(electionVoteListSubstore, accKey...)
}

func GetValidatorListKey() []byte {
	return validatorListSubstore
}
