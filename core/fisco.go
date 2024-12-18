package core

import (
	"fmt"
	"github.com/FISCO-BCOS/go-sdk/abi"
	"github.com/FISCO-BCOS/go-sdk/abi/bind"
	"github.com/FISCO-BCOS/go-sdk/client"
	"github.com/FISCO-BCOS/go-sdk/conf"
	"github.com/ethereum/go-ethereum/common"
	"github.com/kkvbAugust/fiscobcos-go/whole"
	"github.com/sirupsen/logrus"
	"strings"
)

/*
*
全局初始化client对象
*/

func InitClient(path string) *client.Client {
	configs, err := conf.ParseConfigFile(path)
	if err != nil {
		//log.Fatal("resources.ParseConfigFile ERR==>", err)
		logrus.Errorln("resources.ParseConfigFile ERR==>" + err.Error())
		panic("resources.ParseConfigFile ERR==>" + err.Error())
	}
	config := &configs[0]
	clientObj, ok := client.Dial(config)
	if ok != nil {
		//log.Fatal("client.Dial ERR===>", ok)
		logrus.Errorln("client.Dial ERR===>" + ok.Error())
		panic("client.Dial ERR===>" + ok.Error())
	}
	whole.GoSdk.Client = clientObj
	//开辟空间
	whole.GoSdk.Contract = make(map[string]*bind.BoundContract)
	fmt.Println("Client初始化完成")
	return clientObj
}

/**
全局初始化sdk对象
*/

func InitSession(name string) {

	contract, mask := bindContract(whole.Config.Contract[name].Abi, common.HexToAddress(whole.Config.Contract[name].Address), whole.GoSdk.Client, whole.GoSdk.Client, whole.GoSdk.Client)
	if mask != nil {
		//fmt.Println("err==>", mask)
		logrus.Errorln("InitSession Err==>" + mask.Error())
		panic("InitSession Err==>" + mask.Error())
	}
	whole.GoSdk.Contract[name] = contract
	fmt.Println("Session初始化完成")
}

/*
构造合约操作对象
*/
func bindContract(ABI string, address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ABI))
	if err != nil {
		//fmt.Println(err)
		logrus.Errorln("构造合约对象失败==>", err.Error())
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}
