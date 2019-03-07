package internal

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
	err := cdc.UnmarshalBinaryBare(data, rst)
	if err != nil {
		panic("error in json decode userMeta" + err.Error())
	}
	return rst
}

func encodeUserMeta(dt *userMeta) []byte {
	if dt == nil {
		return nil
	}
	rst, err := cdc.MarshalBinaryBare(dt)
	if err != nil {
		panic("error in encoding: " + err.Error())
	}
	return []byte(rst)
}

func decodePostMeta(data []byte) *postMeta {
	if data == nil {
		return nil
	}
	rst := &postMeta{}
	err := cdc.UnmarshalBinaryBare(data, rst)
	if err != nil {
		panic("error in json decode postMeta" + err.Error())
	}
	return rst
}

func encodePostMeta(dt *postMeta) []byte {
	if dt == nil {
		return nil
	}
	rst, err := cdc.MarshalBinaryBare(dt)
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
	err := cdc.UnmarshalBinaryBare(data, rst)
	if err != nil {
		panic("error in json decode roundMeta" + err.Error())
	}
	return rst
}

func encodeRoundMeta(dt *roundMeta) []byte {
	if dt == nil {
		return nil
	}
	rst, err := cdc.MarshalBinaryBare(dt)
	if err != nil {
		panic("error in encoding: " + err.Error())
	}
	return []byte(rst)
}

func decodeUserPostMeta(data []byte) *userPostMeta {
	if data == nil {
		return nil
	}
	rst := &userPostMeta{}
	err := cdc.UnmarshalBinaryBare(data, rst)
	if err != nil {
		panic("error in json decode userPostMeta" + err.Error())
	}
	return rst
}

func encodeUserPostMeta(dt *userPostMeta) []byte {
	if dt == nil {
		return nil
	}
	rst, err := cdc.MarshalBinaryBare(dt)
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
	err := cdc.UnmarshalBinaryBare(data, rst)
	if err != nil {
		panic("error in json decode roundPostMeta" + err.Error())
	}
	return rst
}

func encodeRoundPostMeta(dt *roundPostMeta) []byte {
	if dt == nil {
		return nil
	}
	rst, err := cdc.MarshalBinaryBare(dt)
	if err != nil {
		panic("error in encoding: " + err.Error())
	}
	return []byte(rst)
}

func decodeRoundUserPostMeta(data []byte) *roundUserPostMeta {
	if data == nil {
		return nil
	}
	rst := &roundUserPostMeta{}
	err := cdc.UnmarshalBinaryBare(data, rst)
	if err != nil {
		panic("error in json decode roundUserPostMeta" + err.Error())
	}
	return rst
}

func encodeRoundUserPostMeta(dt *roundUserPostMeta) []byte {
	if dt == nil {
		return nil
	}
	rst, err := cdc.MarshalBinaryBare(dt)
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
	err := cdc.UnmarshalBinaryBare(data, rst)
	if err != nil {
		panic("error in json decode gameMeta" + err.Error())
	}
	return rst
}

func encodeGameMeta(dt *gameMeta) []byte {
	if dt == nil {
		return nil
	}
	rst, err := cdc.MarshalBinaryBare(dt)
	if err != nil {
		panic("error in encoding: " + err.Error())
	}
	return []byte(rst)
}
