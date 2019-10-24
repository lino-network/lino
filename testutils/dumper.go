package testutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ValueInterface interface {
}

type JSONKV struct {
	Prefix string         `json:"prefix"`
	Key    string         `json:"key"`
	Val    ValueInterface `json:"val"`
}

type JSONState = []JSONKV

type prefixMatcher struct {
	prefix      []byte
	mkContainer func() interface{}
	raw         bool
}

type Dumper struct {
	prefixes  map[string]prefixMatcher
	storeCdc  *wire.Codec
	dumperCdc *wire.Codec
	key       sdk.StoreKey
}

type OptionCodec func(cdc *wire.Codec)

func prefixToStr(prefix []byte) string {
	return string([]byte{byte(int(prefix[0]) + int('0'))})
}

func strToPrefx(str string) []byte {
	return []byte{byte(int([]byte(str)[0]) - int('0'))}
}

func NewDumper(key sdk.StoreKey, storeCdc *wire.Codec, options ...OptionCodec) *Dumper {
	str := ""
	dumperCdc := wire.New()
	wire.RegisterCrypto(dumperCdc)
	dumperCdc.RegisterInterface((*ValueInterface)(nil), nil)
	dumperCdc.RegisterConcrete(str, "str", nil)
	for _, option := range options {
		option(dumperCdc)
	}
	return &Dumper{
		prefixes:  make(map[string]prefixMatcher),
		storeCdc:  storeCdc,
		dumperCdc: dumperCdc,
		key:       key,
	}
}

func (d *Dumper) RegisterType(t interface{}, name string, subStore []byte) {
	d.prefixes[string(subStore)] = prefixMatcher{
		subStore,
		func() interface{} {
			return reflect.New(reflect.ValueOf(t).Elem().Type()).Interface()
		},
		false}
	d.dumperCdc.RegisterConcrete(t, name, nil)
}

func (d *Dumper) RegisterRawString(subStore []byte) {
	d.prefixes[string(subStore)] = prefixMatcher{
		prefix: subStore,
		raw:    true,
	}
}

func (d *Dumper) ToJSON(ctx sdk.Context) []byte {
	store := ctx.KVStore(d.key)
	state := make([]JSONKV, 0)
	itr := sdk.KVStorePrefixIterator(store, nil)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		key := itr.Key()
		val := itr.Value()
		if len(key) == 0 {
			panic("zero length key")
		}
		prefix, ok := d.prefixes[string(key[0])]
		if !ok {
			panic(fmt.Sprintf("unknown substoreprefix: %d", int(key[0])))
		}
		pre := prefixToStr(key)
		var kv JSONKV
		if !prefix.raw {
			container := prefix.mkContainer()
			d.storeCdc.MustUnmarshalBinaryLengthPrefixed(val, container)
			kv = JSONKV{
				Prefix: pre,
				Key:    string(key[1:]),
				Val:    container,
			}
		} else {
			kv = JSONKV{
				Prefix: pre,
				Key:    string(key[1:]),
				Val:    string(val),
			}
		}
		state = append(state, kv)
	}
	bz, err := d.dumperCdc.MarshalJSONIndent(state, "", "  ")
	if err != nil {
		panic(err)
	}
	return bz
}

func (d *Dumper) DumpToFile(ctx sdk.Context, filepath string) {
	bz := d.ToJSON(ctx)

	f, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err = f.Write(bz); err != nil {
		panic(err)
	}
	if err = f.Sync(); err != nil {
		panic(err)
	}
}

func (d *Dumper) LoadFromFile(ctx sdk.Context, filepath string) {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	bz, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	state := make(JSONState, 0)
	d.dumperCdc.MustUnmarshalJSON(bz, &state)

	store := ctx.KVStore(d.key)
	for _, v := range state {
		pre := strToPrefx(v.Prefix)
		prefix, ok := d.prefixes[string(pre)]
		if !ok {
			panic(fmt.Sprintf("unknown prefix: %v", v.Prefix))
		}
		if prefix.raw {
			store.Set(append(pre, []byte(v.Key)...), []byte(v.Val.(string)))
		} else {
			store.Set(append(pre, []byte(v.Key)...),
				d.storeCdc.MustMarshalBinaryLengthPrefixed(v.Val))
		}
	}
}
