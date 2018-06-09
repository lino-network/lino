package types

var _ TransferObject = AccountKey("")
var _ TransferObject = PermLink("")
var _ TransferObject = InternalObject("")

type InternalObject string

// indicates all possible balance behavior types
type TransferInDetail int
type TransferOutDetail int

// transfer target maybe account or post or internal target
type TransferObject interface {
	IsTransferObject()
}

func (_ AccountKey) IsTransferObject()     {}
func (_ PermLink) IsTransferObject()       {}
func (_ InternalObject) IsTransferObject() {}
