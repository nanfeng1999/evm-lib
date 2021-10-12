/**
 * @Author: yzy
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/9/24 16:58
 * @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */
package kernel

import (
	"fmt"
)

//go:generate msgp
type StateObjectJson struct {
	ABI      []byte `msg:"abi"`
	Address  []byte `msg:"address"`
	AddrHash []byte `msg:"addrhash"`
	Data     []byte `msg:"data"`
	Code     []byte `msg:"code"`
	Origin   []byte `msg:"origin"`
}

func (s *StateObjectJson) ToByteArray() ([]byte, error) {
	data, err := s.MarshalMsg(nil)
	if err != nil {
		fmt.Println("StateObjectJson to byte err=", err)
	}
	return data, nil
}

func (s *StateObjectJson) FromByteArray(data []byte) error {
	_, err := s.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("StateObjectJson from byte err=", err)
	}
	return err
}
