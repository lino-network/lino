package account

import (
	"encoding/hex"
	"strconv"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/account/model"
	"github.com/lino-network/lino/x/account/types"
)

// creates a querier for account REST endpoints
func NewQuerier(am AccountKeeper) sdk.Querier {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryAccountInfo:
			return queryAccountInfo(ctx, cdc, path[1:], req, am)
		case types.QueryAccountBank:
			return queryAccountBank(ctx, cdc, path[1:], req, am)
		case types.QueryAccountBankByAddress:
			return queryAccountBankByAddress(ctx, cdc, path[1:], req, am)
		case types.QueryAccountMeta:
			return queryAccountMeta(ctx, cdc, path[1:], req, am)
		case types.QueryAccountGrantPubKeys:
			return queryAccountGrantPubKeys(ctx, cdc, path[1:], req, am)
		case types.QueryAccountAllGrantPubKeys:
			return queryAccountAllGrantPubKeys(ctx, cdc, path[1:], req, am)
		case types.QueryTxAndAccountSequence:
			return queryTxAndSequenceNumber(ctx, cdc, path[1:], req, am)
		default:
			return nil, sdk.ErrUnknownRequest("unknown account query endpoint")
		}
	}
}

func queryAccountInfo(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, am AccountKeeper) ([]byte, sdk.Error) {
	if err := linotypes.CheckPathContentAndMinLength(path, 1); err != nil {
		return nil, err
	}
	accountInfo, err := am.GetInfo(ctx, linotypes.AccountKey(path[0]))
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(accountInfo)
	if marshalErr != nil {
		return nil, types.ErrQueryFailed()
	}
	return res, nil
}

func queryAccountBank(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, am AccountKeeper) ([]byte, sdk.Error) {
	if err := linotypes.CheckPathContentAndMinLength(path, 1); err != nil {
		return nil, err
	}
	bank, err := am.GetBank(ctx, linotypes.AccountKey(path[0]))
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(bank)
	if marshalErr != nil {
		return nil, types.ErrQueryFailed()
	}
	return res, nil
}

func queryAccountBankByAddress(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, am AccountKeeper) ([]byte, sdk.Error) {
	if err := linotypes.CheckPathContentAndMinLength(path, 1); err != nil {
		return nil, err
	}
	addr, e := sdk.AccAddressFromBech32(path[0])
	if e != nil {
		return nil, types.ErrQueryFailed()
	}
	bank, err := am.GetBankByAddress(ctx, addr)
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(bank)
	if marshalErr != nil {
		return nil, types.ErrQueryFailed()
	}
	return res, nil
}

func queryAccountMeta(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, am AccountKeeper) ([]byte, sdk.Error) {
	if err := linotypes.CheckPathContentAndMinLength(path, 1); err != nil {
		return nil, err
	}
	accountMeta, err := am.GetMeta(ctx, linotypes.AccountKey(path[0]))
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(accountMeta)
	if marshalErr != nil {
		return nil, types.ErrQueryFailed()
	}
	return res, nil
}

func queryTxAndSequenceNumber(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, am AccountKeeper) ([]byte, sdk.Error) {
	if err := linotypes.CheckPathContentAndMinLength(path, 2); err != nil {
		return nil, err
	}
	isAddr := false
	var e error
	if len(path) == 3 {
		isAddr, e = strconv.ParseBool(path[2])
		if e != nil {
			return nil, types.ErrQueryFailed()
		}
	}

	var bank *model.AccountBank
	var err sdk.Error
	if isAddr {
		addr, e := sdk.AccAddressFromBech32(path[0])
		if e != nil {
			return nil, types.ErrQueryFailed()
		}
		bank, err = am.GetBankByAddress(ctx, addr)
		if err != nil {
			return nil, err
		}
	} else {
		bank, err = am.GetBank(ctx, linotypes.AccountKey(path[0]))
		if err != nil {
			return nil, err
		}
	}

	txAndSeq := model.TxAndSequenceNumber{
		Username: path[0],
		Sequence: bank.Sequence,
	}

	txHash, decodeFail := hex.DecodeString(path[1])
	if decodeFail != nil {
		return nil, types.ErrQueryFailed()
	}

	rpc := rpcclient.NewHTTP("http://localhost:26657", "/websocket")
	tx, _ := rpc.Tx(txHash, false)
	if tx != nil {
		txAndSeq.Tx = &model.Transaction{
			Hash:   hex.EncodeToString(tx.Hash),
			Height: tx.Height,
			Tx:     tx.Tx,
			Code:   tx.TxResult.Code,
			Log:    tx.TxResult.Log,
		}
	}
	res, marshalErr := cdc.MarshalJSON(txAndSeq)
	if marshalErr != nil {
		return nil, types.ErrQueryFailed()
	}
	return res, nil
}

func queryAccountGrantPubKeys(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, am AccountKeeper) ([]byte, sdk.Error) {
	if err := linotypes.CheckPathContentAndMinLength(path, 2); err != nil {
		return nil, err
	}
	grantPubKeys, err := am.GetGrantPubKeys(ctx, linotypes.AccountKey(path[0]), linotypes.AccountKey(path[1]))
	if err != nil {
		return nil, types.ErrQueryFailed()
	}
	res, marshalErr := cdc.MarshalJSON(grantPubKeys)
	if marshalErr != nil {
		return nil, types.ErrQueryFailed()
	}
	return res, nil
}

func queryAccountAllGrantPubKeys(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, am AccountKeeper) ([]byte, sdk.Error) {
	if err := linotypes.CheckPathContentAndMinLength(path, 1); err != nil {
		return nil, err
	}
	pubKeys, err := am.GetAllGrantPubKeys(ctx, linotypes.AccountKey(path[0]))
	if err != nil {
		return nil, types.ErrQueryFailed()
	}
	res, marshalErr := cdc.MarshalJSON(pubKeys)
	if marshalErr != nil {
		return nil, types.ErrQueryFailed()
	}
	return res, nil
}
