package types

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cauth "github.com/cosmos/cosmos-sdk/x/auth"
)

// Theses two functions must be paired, that if txEncoder use JSON format,
// the txDecoder must use JSON format as well.
// All libraries use use the TxEncoder below.

// TxDecoder - default tx decoder, decode tx before authenticate handler
func TxDecoder(cdc *wire.Codec) sdk.TxDecoder {
	return func(txBytes []byte) (tx sdk.Tx, err sdk.Error) {
		defer func() {
			if r := recover(); r != nil {
				err = sdk.ErrTxDecode("tx decode panic")
			}
		}()
		tx = cauth.StdTx{}

		if len(txBytes) == 0 {
			return nil, sdk.ErrTxDecode("txBytes are empty")
		}

		// StdTx.Msg is an interface. The concrete types
		// are registered by MakeTxCodec
		unmarshalErr := cdc.UnmarshalJSON(txBytes, &tx)
		if unmarshalErr != nil {
			return nil, sdk.ErrTxDecode("")
		}
		return tx, nil
	}
}

// TxDecoder - default tx decoder, decode tx before authenticate handler
func TxEncoder(cdc *wire.Codec) sdk.TxEncoder {
	return func(tx sdk.Tx) ([]byte, error) {
		return cdc.MarshalJSON(tx)
	}
}
