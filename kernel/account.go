/**
 * @Author: yzy
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/8/19 15:26
 * @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */
package kernel

import "math/big"

type Account struct {
	Nonce    uint64   // 账户产生交易次数
	Balance  *big.Int // 账户余额
	CodeHash []byte   // 智能合约代码哈希数组
}

func NewEmptyAccount() Account {
	return Account{}
}

func (act *Account) GetBalance() *big.Int {
	return act.Balance
}

func (act *Account) SetBalance(amount *big.Int) {
	act.Balance = amount
}

func (act *Account) DeepCopy() Account {
	newAcc := Account{}
	newAcc.Nonce = act.Nonce
	newAcc.Balance = new(big.Int).Add(act.Balance, new(big.Int))
	newAcc.CodeHash = make([]byte, len(act.CodeHash))
	copy(newAcc.CodeHash, act.CodeHash)
	return newAcc
}
