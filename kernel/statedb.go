package kernel

import (
	"errors"
	"math/big"
)

const (
	maxSnapNum = 10
)

var (
	StateObjectNotFoundErr = errors.New("stateObject not exist")
)

type MStateDB struct {
	db           DB                       // 落库的数据库引用
	stateObjects map[Address]*stateObject // 缓存
	version      int                      // 快照版本号
}

var _ StateDB = (*MStateDB)(nil)

func MakeNewStateDB(db DB) StateDB {
	statedb := new(MStateDB)
	statedb.db = db
	statedb.stateObjects = make(map[Address]*stateObject)
	statedb.version = -1
	return statedb
}

func (s *MStateDB) createObject(addr Address) *stateObject {
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

func (s *MStateDB) CreateAccount(addr Address) {
	s.createObject(addr)
}

func (s *MStateDB) getStateObject(addr Address) *stateObject {
	return s.stateObjects[addr]
}

func (s *MStateDB) SubBalance(addr Address, amount *big.Int) {
	obj := s.getStateObject(addr)
	if obj != nil {
		obj.SubBalance(amount)

	}
}
func (s *MStateDB) AddBalance(addr Address, amount *big.Int) {
	obj := s.getStateObject(addr)
	if obj != nil {
		obj.AddBalance(amount)
	}
}
func (s *MStateDB) GetBalance(addr Address) *big.Int {
	obj := s.getStateObject(addr)
	if obj != nil {
		return obj.data.GetBalance()
	}
	return new(big.Int)
}
func (s *MStateDB) GetNonce(addr Address) uint64 {
	obj := s.getStateObject(addr)
	if obj != nil {
		return obj.Nonce()
	}
	return 0
}
func (s *MStateDB) SetNonce(addr Address, nonce uint64) {
	obj := s.getStateObject(addr)
	if obj != nil {
		obj.SetNonce(nonce)
	}
}

func (s *MStateDB) GetCodeHash(addr Address) Hash {
	obj := s.getStateObject(addr)
	if obj != nil {
		return obj.CodeHash()
	}
	return Hash{}
}

func (s *MStateDB) GetCode(addr Address) []byte {
	obj := s.getStateObject(addr)
	if obj != nil {
		return obj.code
	}
	return []byte{}
}

func (s *MStateDB) SetCode(addr Address, data []byte) {
	obj := s.getStateObject(addr)
	if obj != nil {
		obj.code = data
	}
}

func (s *MStateDB) GetCodeSize(addr Address) int {
	obj := s.getStateObject(addr)
	if obj != nil {
		return len(obj.code)
	}
	return 0
}

// AddRefund 没有用到 暂不实现
func (s *MStateDB) AddRefund(uint64)  {}
func (s *MStateDB) GetRefund() uint64 { return 0 }

func (s *MStateDB) GetState(addr Address, key Hash) Hash {
	obj := s.getStateObject(addr)
	if obj != nil {
		return obj.GetState(key)
	}
	return Hash{}
}
func (s *MStateDB) SetState(addr Address, key Hash, value Hash) {
	obj := s.getStateObject(addr)
	if obj != nil {
		obj.SetState(key, value)
	}
}

// Suicide 没有用到 暂不实现
func (s *MStateDB) Suicide(Address) bool     { return false }
func (s *MStateDB) HasSuicided(Address) bool { return false }

func (s *MStateDB) Exist(addr Address) bool {
	return s.getStateObject(addr) != nil
}

func (s *MStateDB) Empty(addr Address) bool {
	return s.getStateObject(addr) == nil
}

func (s *MStateDB) RevertToSnapshot(pre int) {
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

func (s *MStateDB) AddPreimage(Hash, []byte) {

}

func (s *MStateDB) ForEachStorage(Address, func(Hash, Hash) bool) {

}

func (s *MStateDB) HaveSufficientBalance(Address, *big.Int) bool {
	return true
}

func (s *MStateDB) TransferBalance(Address, Address, *big.Int) {

}
