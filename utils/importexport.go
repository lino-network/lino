package utils

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Load import and unmarshal by cdc json unmarshal.
func Load(filepath string, cdc *codec.Codec, factory func() interface{}) (interface{}, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	table := factory()
	err = cdc.UnmarshalJSON(bytes, table)
	if err != nil {
		return nil, err
	}
	return table, nil
}

// Save save the state to file, using codec json marshal.
func Save(filepath string, cdc *codec.Codec, state interface{}) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	jsonbytes, err := cdc.MarshalJSON(state)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %s", err)
	}
	_, err = f.Write(jsonbytes)
	if err != nil {
		return err
	}
	err = f.Sync()
	if err != nil {
		return err
	}
	return nil
}

type ValueReactor = func(key []byte, val interface{}) bool
type ValueCreator = func() interface{}
type Unmarshaler = func(bz []byte, rst interface{})

type SubStore struct {
	Store      sdk.KVStore
	Prefix     []byte
	ValCreator ValueCreator
	Decoder    Unmarshaler
	NoValue    bool
}

func (s SubStore) Iterate(reactor ValueReactor) {
	itr := sdk.KVStorePrefixIterator(s.Store, s.Prefix)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		k := itr.Key()[len(s.Prefix):]
		v := itr.Value()
		var rst interface{}
		if s.ValCreator != nil {
			rst = s.ValCreator()
		}
		if !s.NoValue {
			s.Decoder(v, rst)
		}
		if reactor(k, rst) {
			break
		}
	}
}

type StoreMap map[string]SubStore

func NewStoreMap(ss []SubStore) StoreMap {
	rst := make(StoreMap)
	for _, v := range ss {
		rst[string(v.Prefix)] = v
	}
	return rst
}
