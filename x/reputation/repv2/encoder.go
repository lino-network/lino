package repv2

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
)

var cdc = wire.New()

// ------ following codes are generated from codegen/genGobCode.py --------
// ------------------------- DO NOT CHANGE --------------------------------
func decodeUserMeta(data []byte) *userMeta {
	if data == nil {
		return nil
	}
	rst := &userMeta{}
	err := cdc.UnmarshalJSON(data, rst)
	if err != nil {
		panic("error in json decode userMeta" + err.Error())
	}
	return rst
}

func encodeUserMeta(dt *userMeta) []byte {
	if dt == nil {
		return nil
	}
	rst, err := cdc.MarshalJSON(dt)
	if err != nil {
		panic("error in encoding: " + err.Error())
	}
	return []byte(rst)
}

func decodeRoundMeta(data []byte) *roundMeta {
	if data == nil {
		return nil
	}
	rst := &roundMeta{}
	err := cdc.UnmarshalJSON(data, rst)
	if err != nil {
		panic("error in json decode roundMeta" + err.Error())
	}
	return rst
}

func encodeRoundMeta(dt *roundMeta) []byte {
	if dt == nil {
		return nil
	}
	rst, err := cdc.MarshalJSON(dt)
	if err != nil {
		panic("error in encoding: " + err.Error())
	}
	return []byte(rst)
}

func decodeRoundPostMeta(data []byte) *roundPostMeta {
	if data == nil {
		return nil
	}
	rst := &roundPostMeta{}
	err := cdc.UnmarshalJSON(data, rst)
	if err != nil {
		panic("error in json decode roundPostMeta" + err.Error())
	}
	return rst
}

func encodeRoundPostMeta(dt *roundPostMeta) []byte {
	if dt == nil {
		return nil
	}
	rst, err := cdc.MarshalJSON(dt)
	if err != nil {
		panic("error in encoding: " + err.Error())
	}
	return []byte(rst)
}

func decodeGameMeta(data []byte) *gameMeta {
	if data == nil {
		return nil
	}
	rst := &gameMeta{}
	err := cdc.UnmarshalJSON(data, rst)
	if err != nil {
		panic("error in json decode gameMeta" + err.Error())
	}
	return rst
}

func encodeGameMeta(dt *gameMeta) []byte {
	if dt == nil {
		return nil
	}
	rst, err := cdc.MarshalJSON(dt)
	if err != nil {
		panic("error in encoding: " + err.Error())
	}
	return []byte(rst)
}
