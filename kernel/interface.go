package kernel

import (
	"math/big"
)

// StateDB is an EVM database for full state querying.
type StateDB interface {
	CreateAccount(Address)

	SubBalance(Address, *big.Int)
	AddBalance(Address, *big.Int)
	GetBalance(Address) *big.Int

	GetNonce(Address) uint64
	SetNonce(Address, uint64)

	GetCodeHash(Address) Hash
	GetCode(Address) []byte
	SetCode(Address, []byte)
	GetCodeSize(Address) int

	AddRefund(uint64)
	GetRefund() uint64

	GetState(Address, Hash) Hash
	SetState(Address, Hash, Hash)

	Suicide(Address) bool
	HasSuicided(Address) bool

	Exist(Address) bool
	Empty(Address) bool

	RevertToSnapshot(int)
	Snapshot() int

	HaveSufficientBalance(Address, *big.Int) bool
	TransferBalance(Address, Address, *big.Int)

	AddLog(*Log)
	AddPreimage(Hash, []byte)
	ForEachStorage(Address, func(Hash, Hash) bool)
}

type AddressHandler interface {
	CreateAddress(b Address, nonce uint64) Address
}
