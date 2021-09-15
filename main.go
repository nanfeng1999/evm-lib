package main

import (
	"encoding/json"
	"evm/kernel"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"io/ioutil"
	"math/big"
	"os"
)

const (
	defaultABIPath = "testdata/test.abi"
	defaultBINPath = "testdata/test.bin"
)

func JsonToABI(path string) (*abi.ABI, error) {
	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(fd)
	var abi abi.ABI
	err = decoder.Decode(&abi)
	if err != nil {
		return nil, err
	}
	return &abi, nil
}

func ReadBIN(path string) ([]byte, error) {
	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}
	return kernel.FromHex(string(bytes)), nil
}

func GetUintInput(ID []byte, input uint) []byte {
	newID := append(ID, common.LeftPadBytes([]byte{uint8(input)}, 32)...)
	return newID
}

func main() {
	testABI, _ := JsonToABI(defaultABIPath)
	CodeBytes, _ := ReadBIN(defaultBINPath)
	calleraddress := kernel.BytesToAddress([]byte("TestAddress"))
	evm := CreateExecuteRuntime(calleraddress)
	caller := kernel.AccountRef(evm.Origin)

	ret, contractAddr, _, err := evm.Create(caller, CodeBytes, evm.GasLimit, new(big.Int))

	input := GetUintInput(testABI.Methods["set"].ID, 100)

	ret, _, err = evm.Call(
		caller,
		contractAddr,
		input,
		evm.GasLimit,
		new(big.Int))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(kernel.Bytes2Hex(ret))
	}

	input = testABI.Methods["power"].ID
	ret, _, err = evm.Call(
		caller,
		contractAddr,
		input,
		evm.GasLimit,
		new(big.Int))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(kernel.Bytes2Hex(ret))
	}
}

//func main1() {
//
//	var (
//		//flAddress     = flag.String("address", "", "address")
//		flBlockNumber = flag.Int("number", 0, "block number")
//		//flCode        = flag.String("code", "", "bytecode")
//		flCoinbase = flag.String("coinbase", "", "coinbase")
//		//flData        = flag.String("data", "", "data")
//		//flDB          = flag.String("db", "", "database")
//		flDifficulty = flag.Int("difficulty", 0, "difficulty")
//		flGasLimit   = flag.Int("gaslimit", 100000, "gas limit")
//		flGasPrice   = flag.Int("gasprice", 1, "gas price")
//		flValue      = flag.Int64("value", 0, "value")
//	)
//
//	flag.Parse()
//	cfg := runtime.Config{}
//	cfg.BlockNumber = big.NewInt(int64(*flBlockNumber))
//	cfg.Coinbase = common.HexToAddress(*flCoinbase)
//	cfg.Difficulty = big.NewInt(int64(*flDifficulty))
//	cfg.GasLimit = uint64(*flGasLimit)
//	cfg.GasPrice = big.NewInt(int64(*flGasPrice))
//	cfg.Origin = common.HexToAddress("TestAddress")
//	cfg.Value = big.NewInt(*flValue)
//	cfg.EVMConfig.Debug = true
//	slg := vm.NewStructLogger(nil)
//	cfg.EVMConfig.Tracer = slg
//
//	ContactCode := "608060405234801561001057600080fd5b5061011e806100206000396000f3006080604052600436106053576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680634a4d59fa14605857806360fe47b11460805780636d4ce63c1460aa575b600080fd5b348015606357600080fd5b50606a60d2565b6040518082815260200191505060405180910390f35b348015608b57600080fd5b5060a86004803603810190808035906020019092919050505060df565b005b34801560b557600080fd5b5060bc60e9565b6040518082815260200191505060405180910390f35b6000805460005402905090565b8060008190555050565b600080549050905600a165627a7a723058209d9ea685d77d5adf0546a6abd70033a69e13ae9ecc4167b03c7488a18c9b1d240029"
//	CodeBytes := common.FromHex(ContactCode)
//
//	ret, add, _, err := runtime.Create(CodeBytes, &cfg)
//	if err != nil {
//		fmt.Println("create account fail,err:", err)
//		return
//	}
//	fmt.Println("Address =", add.String())
//	fmt.Println("Return  =", common.Bytes2Hex(ret))
//
//	intType, _ := abi.NewType("uint256", "", nil)
//
//	Inputs := abi.Arguments{abi.Argument{Name: "x", Type: intType}}
//	methods := abi.NewMethod("test", "set", abi.Function, "", false, false, Inputs, nil)
//	methods.ID = append(methods.ID, common.LeftPadBytes([]byte{100}, 32)...)
//	fmt.Println("Sig:", common.Bytes2Hex(methods.ID), methods.Sig)
//	//input1Str:=common.Bytes2Hex(methods.ID)
//
//	input1Str := "0x60fe47b10000000000000000000000000000000000000000000000000000000000000064"
//	ret, _, err = runtime.Call(add, common.FromHex(input1Str), &cfg)
//	if err != nil {
//		fmt.Println("create account fail,err:", err)
//		return
//	}
//	fmt.Println("Return  =", common.Bytes2Hex(ret))
//
//	input2Str := "0x6d4ce63c"
//	ret, _, err = runtime.Call(add, common.FromHex(input2Str), &cfg)
//	if err != nil {
//		fmt.Println("create account fail,err:", err)
//		return
//	}
//	fmt.Println("Return  =", common.Bytes2Hex(ret))
//
//	input3Str := "0x4a4d59fa"
//	ret, _, err = runtime.Call(add, common.FromHex(input3Str), &cfg)
//	if err != nil {
//		fmt.Println("create account fail,err:", err)
//		return
//	}
//	fmt.Println("Return  =", common.Bytes2Hex(ret))
//
//}
