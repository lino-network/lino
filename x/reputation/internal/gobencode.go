package internal

import (
	"bytes"
	"encoding/gob"
)

// ------ following codes are generated from codegen/genGobCode.py --------
// ------------------------- DO NOT CHANGE --------------------------------
func decodeUserMeta(data []byte) *userMeta {
	if data == nil {
		return nil
	}
	rst := &userMeta{}
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	err := dec.Decode(rst)
	if err != nil {
		panic("error in gob decode userMeta" + err.Error())
	}
	return rst
}

func encodeUserMeta(dt *userMeta) []byte {
	if dt == nil {
		return nil
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(dt)
	if err != nil {
		panic("error in encoding: " + err.Error())
	}
	return buf.Bytes()
}

func decodePostMeta(data []byte) *postMeta {
	if data == nil {
		return nil
	}
	rst := &postMeta{}
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	err := dec.Decode(rst)
	if err != nil {
		panic("error in gob decode postMeta" + err.Error())
	}
	return rst
}

func encodePostMeta(dt *postMeta) []byte {
	if dt == nil {
		return nil
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(dt)
	if err != nil {
		panic("error in encoding: " + err.Error())
	}
	return buf.Bytes()
}

func decodeRoundMeta(data []byte) *roundMeta {
	if data == nil {
		return nil
	}
	rst := &roundMeta{}
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	err := dec.Decode(rst)
	if err != nil {
		panic("error in gob decode roundMeta" + err.Error())
	}
	return rst
}

func encodeRoundMeta(dt *roundMeta) []byte {
	if dt == nil {
		return nil
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(dt)
	if err != nil {
		panic("error in encoding: " + err.Error())
	}
	return buf.Bytes()
}

func decodeUserPostMeta(data []byte) *userPostMeta {
	if data == nil {
		return nil
	}
	rst := &userPostMeta{}
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	err := dec.Decode(rst)
	if err != nil {
		panic("error in gob decode userPostMeta" + err.Error())
	}
	return rst
}

func encodeUserPostMeta(dt *userPostMeta) []byte {
	if dt == nil {
		return nil
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(dt)
	if err != nil {
		panic("error in encoding: " + err.Error())
	}
	return buf.Bytes()
}

func decodeRoundPostMeta(data []byte) *roundPostMeta {
	if data == nil {
		return nil
	}
	rst := &roundPostMeta{}
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	err := dec.Decode(rst)
	if err != nil {
		panic("error in gob decode roundPostMeta" + err.Error())
	}
	return rst
}

func encodeRoundPostMeta(dt *roundPostMeta) []byte {
	if dt == nil {
		return nil
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(dt)
	if err != nil {
		panic("error in encoding: " + err.Error())
	}
	return buf.Bytes()
}

func decodeRoundUserPostMeta(data []byte) *roundUserPostMeta {
	if data == nil {
		return nil
	}
	rst := &roundUserPostMeta{}
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	err := dec.Decode(rst)
	if err != nil {
		panic("error in gob decode roundUserPostMeta" + err.Error())
	}
	return rst
}

func encodeRoundUserPostMeta(dt *roundUserPostMeta) []byte {
	if dt == nil {
		return nil
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(dt)
	if err != nil {
		panic("error in encoding: " + err.Error())
	}
	return buf.Bytes()
}

func decodeGameMeta(data []byte) *gameMeta {
	if data == nil {
		return nil
	}
	rst := &gameMeta{}
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	err := dec.Decode(rst)
	if err != nil {
		panic("error in gob decode gameMeta" + err.Error())
	}
	return rst
}

func encodeGameMeta(dt *gameMeta) []byte {
	if dt == nil {
		return nil
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(dt)
	if err != nil {
		panic("error in encoding: " + err.Error())
	}
	return buf.Bytes()
}
