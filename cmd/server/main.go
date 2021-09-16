/**
 * @Author: yzy
 * @Description:
 * @Version: 1.0.0
 * @Date: 2021/9/15 17:52
 * @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */
package main

import (
	"encoding/json"
	"errors"
	"evm/kernel"
	"evm/runtime"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"log"
	"math/big"
	"net"
	"strconv"
	"strings"
)

// type set
const (
	CONTACT = "contact"
	ACCOUNT = "account"
)

// command set
const (
	RUN = "run"
	ADD = "add"
)

// error when exec occurs
var (
	createContactErr  = errors.New("contact: create contact fail")
	unMarshalParamErr = errors.New("unmarshal param fail")
	parseAbiErr       = errors.New("parse abi file fail")
	abiNotExitErr     = errors.New("abi not exist")
	funcNotExitErr    = errors.New("func not exist")
)

var db = kernel.MakeNewStateDB(new(kernel.MockDB))

type Request struct {
	Type       string // 类型
	Command    string // 命令
	Parameters []byte // 参数
}

type Response struct {
	ErrMsg      string // 需要返回的错误信息
	ContactAddr string // 合约地址
	Result      string // 执行合约结果
}

type ContactInfo struct {
	Addr string
	Abi  []byte
	Bin  []byte
}

type RunContact struct {
	Input       string // 输入
	AccountAddr string // 账户地址
	ContactAddr string // 合约地址
	Sign        string // 函数签名
}

func makeNewResponse(ErrMsg string, addr string, result string) *Response {
	return &Response{
		ErrMsg:      ErrMsg,
		ContactAddr: addr,
		Result:      result,
	}
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("get conn fail,err is:", err)
			continue
		}
		fmt.Println("conn:", conn.RemoteAddr())
		var req Request
		jsonReader := json.NewDecoder(conn)

		err = jsonReader.Decode(&req)
		if err != nil {
			log.Println("parse json fail,err is:", err)
			continue
		}
		var res *Response
		switch req.Type {
		case CONTACT:
			switch req.Command {
			case ADD:
				addr, err := dealContactAddCmd(req.Parameters)

				if err != nil {
					res = makeNewResponse(err.Error(), "", "")
				} else {
					res = makeNewResponse("", addr.Hex(), "")
				}
			case RUN:
				ret, err := dealContactRunCmd(req.Parameters)
				if err != nil {
					res = makeNewResponse(err.Error(), "", "")
				} else {
					res = makeNewResponse("", "", ret)
				}
			default:
				fmt.Println("the command is null or not correct")
			}
		case ACCOUNT:
			switch req.Command {
			default:
				fmt.Println("the command is null or not correct")
			}

		default:
			fmt.Println("the type is null or not correct")
		}
		resData, _ := json.Marshal(res)

		conn.Write(resData)

		conn.Close()
	}

}

func dealContactAddCmd(param []byte) (kernel.Address, error) {
	var contactInfo ContactInfo
	err := json.Unmarshal(param, &contactInfo)
	if err != nil {
		fmt.Println("unmarshal param fail,the err is ", err)
		return kernel.Address{}, err
	}

	var abi abi.ABI
	err = json.Unmarshal(contactInfo.Abi, &abi)
	if err != nil {
		fmt.Println("parse abi file fail,the err is ", err)
		return kernel.Address{}, parseAbiErr
	}

	calleraddress := contactInfo.Addr
	evm := runtime.CreateExecuteRuntime(kernel.HexToAddress(calleraddress), db)
	caller := kernel.AccountRef(evm.Origin)
	_, contractAddr, _, err := evm.Create(caller, kernel.FromHex(string(contactInfo.Bin)), evm.GasLimit, new(big.Int))
	if err != nil {
		fmt.Println("create contact fail,the err is ", err)
		return kernel.Address{}, err
	}

	db.SetABI(contractAddr, &abi)

	fmt.Println("create contact success,the addr is ", contractAddr.Hex())
	return contractAddr, nil
}

func dealContactRunCmd(param []byte) (string, error) {
	var runContact RunContact
	err := json.Unmarshal(param, &runContact)
	if err != nil {
		fmt.Println("unmarshal param fail,the err is ", err)
		return "", err
	}

	var contactAddr = kernel.HexToAddress(runContact.ContactAddr)
	_abi := db.GetABI(contactAddr)

	if _abi == nil {
		fmt.Println("abi not exist")
		return "", abiNotExitErr
	}

	if _, ok := _abi.Methods[runContact.Sign]; !ok {
		fmt.Printf("func %s not exist\n", runContact.Sign)
		return "", funcNotExitErr
	}

	calleraddress := kernel.HexToAddress(runContact.AccountAddr)
	evm := runtime.CreateExecuteRuntime(calleraddress, db)
	caller := kernel.AccountRef(evm.Origin)
	input, _ := getInput(_abi, runContact.Sign, runContact.Input)

	ret, _, err := evm.Call(
		caller,
		contactAddr,
		input,
		evm.GasLimit,
		new(big.Int))
	if err != nil {
		fmt.Println("run contact fail,the err is ", err)
		return "", err
	}

	fmt.Println("run contact success,the result is ", kernel.Bytes2Hex(ret))
	return kernel.Bytes2Hex(ret), nil
}

func dealAccountCmd(param []byte) {

}

func getInput(abi *abi.ABI, sign string, inputRaw string) ([]byte, error) {
	if inputRaw == "" {
		return abi.Pack(sign)
	}
	var input []interface{}
	args := strings.Split(inputRaw, ",")
	for _, arg := range args {
		n, _ := strconv.ParseInt(arg, 10, 64)
		input = append(input, big.NewInt(n))
	}
	return abi.Pack(sign, input...)
}
