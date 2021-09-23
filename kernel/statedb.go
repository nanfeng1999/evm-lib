package kernel

import (
	"errors"
	"fmt"
	"github.com/yzy-github/evm-lib/abi"
	"github.com/yzy-github/evm-lib/common"
	"math/big"
)

const (
	maxSnapNum = 10
)

var (
	StateObjectNotFoundErr = errors.New("stateObject not exist")
)

type MStateDB struct {
	db           DB                              // 落库的数据库引用
	stateObjects map[common.Address]*stateObject // 缓存
	version      int                             // 快照版本号
}

var _ StateDB = (*MStateDB)(nil)

func MakeNewStateDB(db DB) StateDB {
	statedb := new(MStateDB)
	statedb.db = db
	statedb.stateObjects = make(map[common.Address]*stateObject)
	statedb.version = -1
	return statedb
}

func (s *MStateDB) createObject(addr common.Address) *stateObject {
	// judge if the database has the account
	old := s.db.OpenAccount(addr)
	if old != nil {
		s.stateObjects[addr] = old
		return old
	}
	obj := newStateObject(addr, Account{})
	s.stateObjects[addr] = obj
	return obj
}

func (s *MStateDB) CreateAccount(addr common.Address) {
	s.createObject(addr)
}

func (s *MStateDB) SetABI(addr common.Address, abi *abi.ABI) {
	obj := s.getStateObject(addr)
	if obj != nil {
		s.setABI(addr, abi)
	}
}

func (s *MStateDB) GetABI(addr common.Address) *abi.ABI {
	obj := s.getStateObject(addr)
	if obj != nil {
		return s.stateObjects[addr].abi
	}
	return nil
}

func (s *MStateDB) setABI(addr common.Address, abi *abi.ABI) {
	s.stateObjects[addr].abi = abi
}

func (s *MStateDB) getStateObject(addr common.Address) *stateObject {
	return s.stateObjects[addr]
}

func (s *MStateDB) SubBalance(addr common.Address, amount *big.Int) {
	obj := s.getStateObject(addr)
	if obj != nil {
		obj.SubBalance(amount)
	}
}
func (s *MStateDB) AddBalance(addr common.Address, amount *big.Int) {
	obj := s.getStateObject(addr)
	if obj != nil {
		obj.AddBalance(amount)
	}
}
func (s *MStateDB) GetBalance(addr common.Address) *big.Int {
	obj := s.getStateObject(addr)
	if obj != nil {
		return obj.data.GetBalance()
	}
	return new(big.Int)
}
func (s *MStateDB) GetNonce(addr common.Address) uint64 {
	obj := s.getStateObject(addr)
	if obj != nil {
		return obj.Nonce()
	}
	return 0
}
func (s *MStateDB) SetNonce(addr common.Address, nonce uint64) {
	obj := s.getStateObject(addr)
	if obj != nil {
		obj.SetNonce(nonce)
	}
}

func (s *MStateDB) GetCodeHash(addr common.Address) common.Hash {
	obj := s.getStateObject(addr)
	if obj != nil {
		return obj.CodeHash()
	}
	return common.Hash{}
}

func (s *MStateDB) GetCode(addr common.Address) []byte {
	obj := s.getStateObject(addr)
	if obj != nil {
		return obj.code
	}
	return []byte{}
}

func (s *MStateDB) SetCode(addr common.Address, data []byte) {
	obj := s.getStateObject(addr)
	if obj != nil {
		obj.code = data
	}
}

func (s *MStateDB) GetCodeSize(addr common.Address) int {
	obj := s.getStateObject(addr)
	if obj != nil {
		return len(obj.code)
	}
	return 0
}

// AddRefund 没有用到 暂不实现
func (s *MStateDB) AddRefund(uint64)  {}
func (s *MStateDB) GetRefund() uint64 { return 0 }

func (s *MStateDB) GetState(addr common.Address, key common.Hash) common.Hash {
	obj := s.getStateObject(addr)
	if obj != nil {
		return obj.GetState(key)
	}
	return common.Hash{}
}
func (s *MStateDB) SetState(addr common.Address, key common.Hash, value common.Hash) {
	obj := s.getStateObject(addr)
	if obj != nil {
		obj.SetState(key, value)
	}
}

// Suicide 没有用到 暂不实现
func (s *MStateDB) Suicide(common.Address) bool     { return false }
func (s *MStateDB) HasSuicided(common.Address) bool { return false }

func (s *MStateDB) Exist(addr common.Address) bool {
	return s.getStateObject(addr) != nil
}

func (s *MStateDB) Empty(addr common.Address) bool {
	return s.getStateObject(addr) == nil
}

func (s *MStateDB) RevertToSnapshot(pre int) {
	fmt.Println("回退")
	for _, obj := range s.stateObjects {
		obj.RevertToSnap(pre)
	}
	s.version = pre
}

func (s *MStateDB) Snapshot() int {
	s.version++
	for _, obj := range s.stateObjects {
		obj.SnapShot()
	}
	return s.version
}

func (s *MStateDB) AddLog(*Log) {

}

func (s *MStateDB) AddPreimage(common.Hash, []byte) {

}

func (s *MStateDB) ForEachStorage(common.Address, func(common.Hash, common.Hash) bool) {

}

func (s *MStateDB) HaveSufficientBalance(common.Address, *big.Int) bool {
	return true
}

func (s *MStateDB) TransferBalance(common.Address, common.Address, *big.Int) {

}
