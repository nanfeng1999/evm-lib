/**
 * @Author: yzy
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/9/17 16:46
 * @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */
package mongodb

import (
	"errors"
	"fmt"
	"github.com/JodeZer/mgop"
	"github.com/nanfeng1999/evm-lib/common"
	"github.com/nanfeng1999/evm-lib/kernel"
	"github.com/nanfeng1999/evm-lib/rlp"
	"gopkg.in/mgo.v2/bson"
)

const (
	ACCOUNT_TABLE = "account"
)

var (
	accountExistErr    = errors.New("the account exist")
	accountNotExistErr = errors.New("the account not exist")
)

type Account struct {
	Name     string
	Password string
	Adderss  string
}

var DefaultDBPool mgop.SessionPool

func init() {
	DefaultDBPool = newMongoDBPool()
}

func newMongoDBPool() mgop.SessionPool {
	pool, err := mgop.DialStrongPool("mongodb://mis:mis20201001@localhost/blockchain", 5)
	if err != nil {
		panic(err)
	}
	return pool
}

func CreateAccount(name string, password string) (string, error) {
	sess := DefaultDBPool.AcquireSession()
	defer sess.Release()

	var acc = Account{
		Name:     name,
		Password: password,
	}
	if acc := GetAccount(acc.Name); acc != nil {
		return "", accountExistErr
	}

	acc.Adderss = createAddress(&acc).Hex()
	return acc.Adderss, sess.DB("blockchain").C(ACCOUNT_TABLE).Insert(&acc)
}

func GetAccountAddr(name string) (string, error) {
	sess := DefaultDBPool.AcquireSession()
	defer sess.Release()

	if acc := GetAccount(name); acc == nil {
		return "", accountExistErr
	} else {
		return acc.Adderss, nil
	}
}

func createAddress(acc *Account) common.Address {
	data, _ := rlp.EncodeToBytes(acc)
	return common.BytesToAddress(kernel.Keccak256(data)[12:])
}

func GetAccount(name string) *Account {
	sess := DefaultDBPool.AcquireSession()
	defer sess.Release()

	var accounts []Account
	err := sess.DB("blockchain").C(ACCOUNT_TABLE).Find(bson.M{"name": name}).All(&accounts)
	if err != nil {
		fmt.Println("account %s not exist", name)
		return nil
	}

	if len(accounts) > 0 {
		return &accounts[0]
	}

	return nil
}
