package internal

import (
	"bytes"

	db "github.com/tendermint/tendermint/libs/db"
)

// A mock store for testing.
type MockStore struct {
	store map[string]([]byte)
}

func NewMockStore() *MockStore {
	return &MockStore{store: make(map[string]([]byte))}
}

func (s MockStore) Set(key []byte, val []byte) {
	s.store[string(key)] = val
}

func (s MockStore) Get(key []byte) []byte {
	if el, ok := s.store[string(key)]; ok {
		return el
	} else {
		return nil
	}
}

func (s MockStore) Has(key []byte) bool {
	_, ok := s.store[string(key)]
	return ok
}

func (s MockStore) Delete(key []byte) {
	delete(s.store, string(key))
}

func (s MockStore) Iterator(start, end []byte) db.Iterator {
	partial := make(map[string]([]byte))
	for k, v := range s.store {
		if inDomain([]byte(k), start, end) {
			partial[k] = v
		}
	}
	return NewMockStoreItr(partial)
}

type MockStoreItr struct {
	kvs   []([2]([]byte))
	cur   int
	start []byte
	end   []byte
}

func NewMockStoreItr(store map[string]([]byte)) *MockStoreItr {
	rst := &MockStoreItr{}
	for k, v := range store {
		rst.kvs = append(rst.kvs, [2][]byte{[]byte(k), v})
	}
	return rst
}

func (i *MockStoreItr) Domain() ([]byte, []byte) {
	return i.start, i.end
}

func (i *MockStoreItr) Valid() bool {
	return i.cur <= len(i.kvs)-1
}

func (i *MockStoreItr) Next() {
	if !i.Valid() {
		panic("next on invalidated key")
	}
	i.cur++
}

func (i *MockStoreItr) Key() []byte {
	return i.kvs[i.cur][0]
}

func (i *MockStoreItr) Value() []byte {
	return i.kvs[i.cur][1]
}

func (i *MockStoreItr) Close() {}

func inDomain(key, start, end []byte) bool {
	if bytes.Compare(key, start) < 0 {
		return false
	}
	if end != nil && bytes.Compare(end, key) <= 0 {
		return false
	}
	return true
}
