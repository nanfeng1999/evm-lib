/**
 * @Author: yzy
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/8/19 15:26
 * @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */
package kernel

import (
	"evm/crypto"
	"math/big"
)

type StateObject interface {
	GetState(key Hash) Hash
	SetState(key, value Hash)
	AddBalance(amount *big.Int)
	SubBalance(amount *big.Int)
	SetBalance(amount *big.Int)
	Address() Address
	Nonce() uint64
	SetNonce(nonce uint64)
	CodeHash() Hash
	SnapShot()
	RevertToSnap(pre int)
	RevertToInit()
}

type stateObject struct {
	address  Address // 账户地址
	addrHash Hash    // 账户地址哈希
	data     Account // 账户结构体
	code     []byte  // 智能合约代码字节数组
	//Trie          *mpt.Trie           // MPT 用来存储状态集合 用于最后落库
	originStorage Storage             // 存储某个账户执行合约过程中的临时状态信息
	dirtyStorage  [maxSnapNum]Storage // 存储某个账户执行合约过程中之前的状态信息
	version       int                 // 回退版本号
}

// 检测结构体是否已经实现接口
var _ StateObject = (*stateObject)(nil)

var emptyHash = crypto.Keccak256(nil)

func newStateObject(addr Address, account Account) *stateObject {
	if account.Balance == nil {
		account.SetBalance(new(big.Int))
	}
	if account.CodeHash == nil {
		account.CodeHash = emptyHash
	}
	return &stateObject{
		address:       addr,                                // 账户地址
		addrHash:      Hash(crypto.Keccak256Hash(addr[:])), // 账户地址哈希
		data:          account,                             // 账户结构体
		originStorage: make(map[Hash]Hash),                 //存储某个账户执行合约过程中的临时状态信息
		version:       -1,
	}
}

func (obj *stateObject) isExist(key Hash) bool {
	_, ok := obj.originStorage[key]
	return ok
}

func (obj *stateObject) GetState(key Hash) Hash {
	if obj.isExist(key) {
		return obj.originStorage[key]
	}
	return Hash{}
}

func (obj *stateObject) SetState(key, value Hash) {
	obj.originStorage[key] = value
}

func (obj *stateObject) AddBalance(amount *big.Int) {
	if amount.Sign() == 0 {
		return
	}
	obj.data.SetBalance(new(big.Int).Add(obj.data.GetBalance(), amount))
}

func (obj *stateObject) SubBalance(amount *big.Int) {
	if amount.Sign() == 0 {
		return
	}
	obj.data.SetBalance(new(big.Int).Sub(obj.data.GetBalance(), amount))
}

func (obj *stateObject) SetBalance(amount *big.Int) {
	obj.data.Balance = amount
}

func (obj *stateObject) Address() Address {
	return obj.address
}

func (obj *stateObject) Nonce() uint64 {
	return obj.data.Nonce
}

func (obj *stateObject) SetNonce(nonce uint64) {
	obj.data.Nonce = nonce
}

func (obj *stateObject) CodeHash() Hash {
	return BytesToHash(obj.data.CodeHash)
}

func (obj *stateObject) SnapShot() {
	obj.version++
	if checkVersion(obj.version) {
		obj.dirtyStorage[obj.version] = make(map[Hash]Hash)
		obj.dirtyStorage[obj.version] = obj.originStorage.Copy()
	}
}

func (obj *stateObject) RevertToSnap(pre int) {
	if checkVersion(pre) {
		obj.originStorage = obj.dirtyStorage[pre].Copy()
		for i := pre + 1; i < obj.version; i++ {
			obj.dirtyStorage[i] = nil
		}
		obj.version = pre
	}
}

func (obj *stateObject) RevertToInit() {
	obj.RevertToSnap(0)
}

func checkVersion(version int) bool {
	if version < maxSnapNum {
		return true
	}
	return false
}
