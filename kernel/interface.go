package kernel

import (
	"github.com/nanfeng1999/evm-lib/abi"
	"github.com/nanfeng1999/evm-lib/common"
	"math/big"
)

// StateDB is an EVM database for full state querying.
type StateDB interface {
	CreateAccount(common.Address)
	GetStateObject(common.Address) *stateObject
	ResetStateObject(common.Address)

	SetABI(common.Address, *abi.ABI, []byte)
	GetABI(common.Address) *abi.ABI

	SubBalance(common.Address, *big.Int)
	AddBalance(common.Address, *big.Int)
	GetBalance(common.Address) *big.Int

	GetNonce(common.Address) uint64
	SetNonce(common.Address, uint64)

	GetCodeHash(common.Address) common.Hash
	GetCode(common.Address) []byte
	SetCode(common.Address, []byte)
	GetCodeSize(common.Address) int

	AddRefund(uint64)
	GetRefund() uint64

	GetState(common.Address, common.Hash) common.Hash
	SetState(common.Address, common.Hash, common.Hash)

	Suicide(common.Address) bool
	HasSuicided(common.Address) bool

	Exist(common.Address) bool
	Empty(common.Address) bool

	RevertToSnapshot(int)
	Snapshot() int

	HaveSufficientBalance(common.Address, *big.Int) bool
	TransferBalance(common.Address, common.Address, *big.Int)

	AddLog(*Log)
	AddPreimage(common.Hash, []byte)
	ForEachStorage(common.Address, func(common.Hash, common.Hash) bool)
}

type AddressHandler interface {
	CreateAddress(b common.Address, nonce uint64) common.Address
}
