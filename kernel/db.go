/**
 * @Author: yzy
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/8/19 15:32
 * @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */
package kernel

import "github.com/nanfeng1999/evm-lib/common"

type DB interface {
	// 根据传入的hash 从数据库中取出rlp编码的stateObject 并进行解码
	OpenAccount(addr common.Address) []byte
	// 传入stateObject 对其进行rlp编码 然后插入数据库中去
	SaveToDB(common.Address, []byte) error
	// 数据库是否存在账户
	ExistAccount(common.Address) bool
	// 更新账户数据
	UpdateAccount(common.Address, []byte) error
}

type MockDB struct{}

func (*MockDB) OpenAccount(addr common.Address) []byte {
	return nil
}

func (*MockDB) SaveToDB(common.Address, []byte) error {
	return nil
}

func (*MockDB) ExistAccount(common.Address) bool {
	return true
}

func (MockDB) UpdateAccount(common.Address, []byte) error {
	return nil
}
