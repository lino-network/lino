package account

import (
	"encoding/hex"
	"strconv"
	"strings"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
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
			return utils.NewQueryResolver(1, func(args ...string) (interface{}, sdk.Error) {
				return am.GetInfo(ctx, linotypes.AccountKey(args[0]))
			})(ctx, cdc, path)
		case types.QueryAccountBank:
			return utils.NewQueryResolver(1, func(args ...string) (interface{}, sdk.Error) {
				return am.GetBank(ctx, linotypes.AccountKey(args[0]))
			})(ctx, cdc, path)
		case types.QueryAccountBankByAddress:
			return utils.NewQueryResolver(1, func(args ...string) (interface{}, sdk.Error) {
				addr, e := sdk.AccAddressFromBech32(args[0])
				if e != nil {
					return nil, types.ErrQueryFailed()
				}
				return am.GetBankByAddress(ctx, addr)
			})(ctx, cdc, path)
		case types.QueryAccountMeta:
			return utils.NewQueryResolver(1, func(args ...string) (interface{}, sdk.Error) {
				return am.GetMeta(ctx, linotypes.AccountKey(args[0]))
			})(ctx, cdc, path)
		case types.QueryTxAndAccountSequence:
			return queryTxAndSequenceNumber(ctx, cdc, path[1:], req, am)
		case types.QueryPool:
			return utils.NewQueryResolver(1, func(args ...string) (interface{}, sdk.Error) {
				poolname := strings.Join(args, "/")
				return am.GetPool(ctx, linotypes.PoolName(poolname))
			})(ctx, cdc, path)
		case types.QuerySupply:
			return utils.NewQueryResolver(0, func(args ...string) (interface{}, sdk.Error) {
				return am.GetSupply(ctx), nil
			})(ctx, cdc, path)
		default:
			return nil, sdk.ErrUnknownRequest("unknown account query endpoint")
		}
	}
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
