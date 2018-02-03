package commands

import (
	"fmt"
	"github.com/pkg/errors"
	bc "github.com/lino-network/lino/types"
	crypto "github.com/tendermint/go-crypto"
	keys "github.com/tendermint/go-crypto/keys"
	wire "github.com/tendermint/go-wire"
)

type PostTx struct {
	chainID string
	signers []crypto.PubKey
	Tx      *bc.PostTx
}

var _ keys.Signable = &PostTx{}

// SignBytes returned the unsigned bytes, needing a signature
func (p *PostTx) SignBytes() []byte {
	fmt.Println("Sign the byte")
	return p.Tx.SignBytes(p.chainID)
}

// AddSigner sets address and pubkey info on the tx based on the key that
// will be used for signing
func (p *PostTx) AddSigner(pk crypto.PubKey) {
	if p.Tx.Sequence == 1 {
		p.Tx.PubKey = pk
	}
}

// Sign will add a signature and pubkey.
//
// Depending on the Signable, one may be able to call this multiple times for multisig
// Returns error if called with invalid data or too many times
func (p *PostTx) Sign(pubkey crypto.PubKey, sig crypto.Signature) error {
	fmt.Println("Sign the post")
	addr := pubkey.Address()
	set := p.Tx.SetSignature(sig)
	if !set {
		return errors.Errorf("Cannot add signature for address %X", addr)
	}
	return nil
}
// Signers will return the public key(s) that signed if the signature
// is valid, or an error if there is any issue with the signature,
// including if there are no signatures
func (p *PostTx) Signers() ([]crypto.PubKey, error) {
	if len(p.signers) == 0 {
		return nil, errors.New("No signatures on SendTx")
	}
	return p.signers, nil
}

// TxBytes returns the transaction data as well as all signatures
// It should return an error if Sign was never called
func (p *PostTx) TxBytes() ([]byte, error) {
	// TODO: verify it is signed

	// Code and comment from: basecoin/cmd/basecoin/commands/tx.go
	// Don't you hate having to do this?
	// How many times have I lost an hour over this trick?!
	txBytes := wire.BinaryBytes(struct {
		bc.Tx `json:"unwrap"`
	}{p.Tx})
	return txBytes, nil
}

// TODO: this should really be in the basecoin.types SendTx,
// but that code is too ugly now, needs refactor..
func (p *PostTx) ValidateBasic() error {
	if p.chainID == "" {
		return errors.New("No chain-id specified")
	}
	if len(p.Tx.Address) != 20 {
		return errors.Errorf("Invalid address length: %d", len(p.Tx.Address))
	}
	if p.Tx.Sequence <= 0 {
		return errors.New("Sequence must be greater than 0")
	}

	return nil
}
