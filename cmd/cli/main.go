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
	"fmt"
	"github.com/peterh/liner"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
)

const cmdHistoryPath = "/tmp/mis-cli"

var commandList = []string{
	"contact", "add", "run", "--abi", "--bin", "--addr", "--func", "--input",
}

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

type Request struct {
	Type       string // 类型
	Command    string // 命令
	Parameters []byte // 参数
}

type ContactInfo struct {
	Addr kernel.Address
	Abi  []byte
	Bin  []byte
}

type RunContact struct {
	Input       string // 输入
	AccountAddr string // 账户地址
	ContactAddr string // 合约地址
	Sign        string // 函数签名
}

type Response struct {
	ErrMsg      string // 需要返回的错误信息
	ContactAddr string // 合约地址
	Result      string // 运行结果
}

// new func of request
func makeNewRequest(Type, command string, param []byte) *Request {
	return &Request{
		Type:       Type,
		Command:    command,
		Parameters: param,
	}
}

func main() {
	// define a new liner
	line := liner.NewLiner()
	defer line.Close()
	// ctrl + c exit
	line.SetCtrlCAborts(true)
	// 自动补全功能
	line.SetCompleter(func(li string) (res []string) {
		for _, c := range commandList {
			if strings.HasPrefix(c, li) {
				res = append(res, strings.ToLower(c))
			}
		}
		return
	})

	// open and save cmd history.
	if f, err := os.Open(cmdHistoryPath); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	defer func() {
		// save history
		if f, err := os.Create(cmdHistoryPath); err != nil {
			fmt.Printf("writing cmd history err: %v\n", err)
		} else {
			line.WriteHistory(f)
			f.Close()
		}
	}()

	prompt := "mis-cli>>"
	for {
		cmd, err := line.Prompt(prompt)
		if err != nil {
			fmt.Println(err)
			break
		}
		// trim space
		cmd = strings.TrimSpace(cmd)
		if len(cmd) == 0 {
			continue
		}
		// transfer to low
		lowerCmd := strings.ToLower(cmd)

		c := strings.Split(cmd, " ")
		// print help or quit.
		if lowerCmd == "quit" {
			fmt.Println("bye")
			break
		} else {
			// execute the command and print the reply.
			line.AppendHistory(cmd)
			app := &cli.App{
				Name:  "mis-cli",
				Usage: "mic cmd tool for contact",
				Commands: []*cli.Command{
					addContactCommand(),
					runContactCommand(),
				},
			}

			err := app.Run(c)
			if err != nil {
				fmt.Printf("(error) %v \n", err)
			}
		}
	}
}

func addContactCommand() *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "create a new contact",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "address",
				Usage:    "the address of account",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "abi",
				Usage:    "the path of abi file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "bin",
				Usage:    "the path of bin file",
				Required: true,
			},
		},
		Action: addContact,
	}
}

func addContact(c *cli.Context) error {
	abiData, err := getDataFromFile(c.String("abi"))
	if err != nil {
		return err
	}
	binData, err := getDataFromFile(c.String("bin"))
	if err != nil {
		return err
	}
	addr := c.String("address")
	contact := ContactInfo{
		Addr: kernel.HexToAddress(addr),
		Abi:  abiData,
		Bin:  binData,
	}

	paramData, _ := json.Marshal(&contact)

	req := makeNewRequest(CONTACT, ADD, paramData)

	// send request
	data, _ := json.Marshal(req)
	conn, err := net.DialTimeout("tcp", "127.0.0.1:8080", 2*time.Second)
	if err != nil {
		return err
	}
	_, err = conn.Write(data)
	if err != nil {
		return err
	}
	// get response
	var response Response
	decoder := json.NewDecoder(conn)
	err = decoder.Decode(&response)
	if err != nil {
		return err
	}

	if response.ContactAddr != "" {
		fmt.Println("contact address: ", response.ContactAddr)
	} else {
		return errors.New(response.ErrMsg)
	}
	conn.Close()
	return nil
}

func runContactCommand() *cli.Command {
	return &cli.Command{
		Name:   "run",
		Usage:  "run contact",
		Hidden: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "accountAddr",
				Usage:    "the addr of contact",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "contactAddr",
				Usage:    "the addr of contact",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "func",
				Aliases:  []string{"f"},
				Usage:    "the function sign",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "input",
				Aliases: []string{"i"},
				Usage:   "the function input",
			},
		},
		Action: runContact,
	}
}

func runContact(c *cli.Context) error {
	runContact := RunContact{
		AccountAddr: c.String("accountAddr"),
		ContactAddr: c.String("contactAddr"),
		Input:       c.String("input"),
		Sign:        c.String("func"),
	}

	paramData, _ := json.Marshal(&runContact)

	req := makeNewRequest(CONTACT, RUN, paramData)

	// send request
	data, _ := json.Marshal(req)
	conn, err := net.DialTimeout("tcp", "127.0.0.1:8080", 2*time.Second)
	if err != nil {
		return err
	}
	_, err = conn.Write(data)
	if err != nil {
		return err
	}
	// get response
	var response Response
	decoder := json.NewDecoder(conn)
	err = decoder.Decode(&response)
	if err != nil {
		return err
	}

	if response.ErrMsg == "" {
		fmt.Println("result:", response.Result)
	} else {
		return errors.New(response.ErrMsg)
	}
	conn.Close()
	return nil
}

func getDataFromFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}
