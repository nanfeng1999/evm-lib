/**
 * @Author: yzy
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/8/19 15:26
 * @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package kernel

import (
	"encoding/json"
	"github.com/nanfeng1999/evm-lib/abi"
	"github.com/nanfeng1999/evm-lib/common"
	"github.com/nanfeng1999/evm-lib/crypto"
	"math/big"
)

type StateObject interface {
	GetState(key common.Hash) common.Hash
	SetState(key, value common.Hash)
	AddBalance(amount *big.Int)
	SubBalance(amount *big.Int)
	SetBalance(amount *big.Int)
	Address() common.Address
	Nonce() uint64
	SetNonce(nonce uint64)
	CodeHash() common.Hash
	SnapShot()
	RevertToSnap(pre int)
	RevertToInit()
}

type stateObject struct {
	abi           *abi.ABI
	abiBytes      []byte
	address       common.Address
	addrHash      common.Hash
	data          Account
	code          []byte
	originStorage Storage
	dirtyStorage  [maxSnapNum]Storage
	version       int
}

// 检测结构体是否已经实现接口
var _ StateObject = (*stateObject)(nil)

var emptyHash = crypto.Keccak256(nil)

func newStateObject(addr common.Address, account Account) *stateObject {
	if account.Balance == nil {
		account.SetBalance(new(big.Int))
	}
	if account.CodeHash == nil {
		account.CodeHash = emptyHash
	}
	return &stateObject{
		address:       addr,                              // 账户地址
		addrHash:      crypto.Keccak256Hash(addr[:]),     // 账户地址哈希
		data:          account,                           // 账户结构体
		originStorage: make(map[common.Hash]common.Hash), //存储某个账户执行合约过程中的临时状态信息
		version:       -1,
	}
}

func (obj *stateObject) isExist(key common.Hash) bool {
	_, ok := obj.originStorage[key]
	return ok
}

func (obj *stateObject) GetState(key common.Hash) common.Hash {
	if obj.isExist(key) {
		return obj.originStorage[key]
	}
	return common.Hash{}
}

func (obj *stateObject) SetState(key, value common.Hash) {
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

func (obj *stateObject) Address() common.Address {
	return obj.address
}

func (obj *stateObject) Nonce() uint64 {
	return obj.data.Nonce
}

func (obj *stateObject) SetNonce(nonce uint64) {
	obj.data.Nonce = nonce
}

func (obj *stateObject) CodeHash() common.Hash {
	return common.BytesToHash(obj.data.CodeHash)
}

func (obj *stateObject) SnapShot() {
	obj.version++
	if checkVersion(obj.version) {
		obj.dirtyStorage[obj.version] = make(map[common.Hash]common.Hash)
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

func (obj *stateObject) ToByteArray() ([]byte, error) {

	dataBytes, _ := json.Marshal(&obj.data)

	originStorageBytes, _ := json.Marshal(&obj.originStorage)
	var stateObjectJson = &StateObjectJson{
		ABI:      obj.abiBytes,
		Address:  obj.address.Bytes(),
		AddrHash: obj.addrHash.Bytes(),
		Data:     dataBytes,
		Code:     obj.code,
		Origin:   originStorageBytes,
	}
	return stateObjectJson.ToByteArray()

}

func (obj *stateObject) FromByteArray(data []byte) error {
	var stateObjectJson = new(StateObjectJson)
	err := stateObjectJson.FromByteArray(data)
	if err != nil {
		return err
	}
	var acc Account
	err = json.Unmarshal(stateObjectJson.Data, &acc)
	if err != nil {
		return err
	}

	var abi abi.ABI
	if stateObjectJson.ABI != nil {
		err = json.Unmarshal(stateObjectJson.ABI, &abi)
		if err != nil {
			return err
		}
	}
	var origin Storage
	err = json.Unmarshal(stateObjectJson.Origin, &origin)
	if err != nil {
		return err
	}

	obj.data = acc
	obj.abi = &abi
	obj.abiBytes = stateObjectJson.ABI
	obj.address = common.BytesToAddress(stateObjectJson.Address)
	obj.addrHash = common.BytesToHash(stateObjectJson.AddrHash)
	obj.originStorage = origin
	obj.code = stateObjectJson.Code
	return nil
}

func checkVersion(version int) bool {
	if version < maxSnapNum {
		return true
	}
	return false
}
