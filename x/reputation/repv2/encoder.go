package repv2

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
)

var cdc = wire.New()

func decodeUserMeta(data []byte) *userMeta {
	if data == nil {
		return nil
	}
	rst := &userMeta{}
	cdc.MustUnmarshalBinaryBare(data, rst)
	return rst
}

func encodeUserMeta(dt *userMeta) []byte {
	if dt == nil {
		return nil
	}
	rst := cdc.MustMarshalBinaryBare(dt)
	return []byte(rst)
}

func decodeRoundMeta(data []byte) *roundMeta {
	if data == nil {
		return nil
	}
	rst := &roundMeta{}
	cdc.MustUnmarshalBinaryBare(data, rst)
	return rst
}

func encodeRoundMeta(dt *roundMeta) []byte {
	if dt == nil {
		return nil
	}
	rst := cdc.MustMarshalBinaryBare(dt)
	return []byte(rst)
}

func decodeRoundPostMeta(data []byte) *roundPostMeta {
	if data == nil {
		return nil
	}
	rst := &roundPostMeta{}
	cdc.MustUnmarshalBinaryBare(data, rst)
	return rst
}

func encodeRoundPostMeta(dt *roundPostMeta) []byte {
	if dt == nil {
		return nil
	}
	rst := cdc.MustMarshalBinaryBare(dt)
	return []byte(rst)
}

func decodeGameMeta(data []byte) *gameMeta {
	if data == nil {
		return nil
	}
	rst := &gameMeta{}
	cdc.MustUnmarshalBinaryBare(data, rst)
	return rst
}

func encodeGameMeta(dt *gameMeta) []byte {
	if dt == nil {
		return nil
	}
	rst := cdc.MustMarshalBinaryBare(dt)
	return []byte(rst)
}
