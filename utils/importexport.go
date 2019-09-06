package utils

import (
	"io/ioutil"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
)

// Load import and unmarshal by amino.
func Load(filepath string, cdc *codec.Codec, factory func() interface{}) (interface{}, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	bytes, err := ioutil.ReadAll(f)

	table := factory()
	err = cdc.UnmarshalJSON(bytes, table)
	if err != nil {
		return nil, err
	}
	return table, nil
}
