package tools

import (
	"context"
	goecdsa "crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/FISCO-BCOS/go-sdk/abi"
	"github.com/FISCO-BCOS/go-sdk/abi/bind"
	"github.com/FISCO-BCOS/go-sdk/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kkvbAugust/fiscobcos-go/whole"
	"strings"
)

/*
公用发送交易组装器
*/

// SendTransaction 发送交易
func SendTransaction(name, method string, params ...interface{}) any {
	_, receipt, ok := whole.GoSdk.Contract[name].Transact(whole.GoSdk.Client.GetTransactOpts(), method, params...)

	if ok != nil {
		fmt.Println("txError=>", ok)
		panic("SendTransaction 发送交易失败")
	}
	json, wrong := abi.JSON(strings.NewReader(whole.Config.Contract[name].Abi))
	if wrong != nil {
		fmt.Println("wrong==>", wrong)
		panic("解析abi失败")
	}

	var (
		result = new(any)
	)
	//合约方法名
	task := json.Unpack(&result, method, common.FromHex(receipt.Output))
	if task != nil {
		fmt.Println("task==>", task)
		panic("解析函数返回值失败")
	}
	return *result
}

// SendTransactionAndGetReceipt 发送交易并获取交易回执
func SendTransactionAndGetReceipt(name, method string, params ...interface{}) (any, *types.Receipt) {
	_, receipt, ok := whole.GoSdk.Contract[name].Transact(whole.GoSdk.Client.GetTransactOpts(), method, params...)

	if ok != nil {
		fmt.Println("txError=>", ok)
		panic("SendTransaction 发送交易失败")
	}
	json, wrong := abi.JSON(strings.NewReader(whole.Config.Contract[name].Abi))
	if wrong != nil {
		fmt.Println("wrong==>", wrong)
		panic("解析abi失败")
	}

	var (
		result = new(any)
	)
	//合约方法名
	task := json.Unpack(&result, method, common.FromHex(receipt.Output))
	if task != nil {
		fmt.Println("task==>", task)
		panic("解析函数返回值失败")
	}
	return *result, receipt
}

// SendCall 获取链上信息
func SendCall(name, method string, out interface{}, params ...interface{}) any {
	err := whole.GoSdk.Contract[name].Call(whole.GoSdk.Client.GetCallOpts(), out, method, params...)
	if err != nil {
		fmt.Println("SendCall err==>", err)
		panic("SendCall 获取链上信息失败")
	}
	return out
}

func NewKeyedTransactor(key *goecdsa.PrivateKey) *bind.TransactOpts {
	//key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	keyAddr := crypto.PubkeyToAddress(key.PublicKey)
	return &bind.TransactOpts{
		From: keyAddr,
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != keyAddr {
				return nil, errors.New("not authorized to sign this account")
			}
			signature, err := crypto.Sign(signer.Hash(tx).Bytes(), key)
			if err != nil {
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
	}
}

// SendTransactionByKey 使用私钥发送交易
func SendTransactionByKey(name, method string, privateKey *goecdsa.PrivateKey, params ...interface{}) any {
	_, receipt, ok := whole.GoSdk.Contract[name].Transact(NewKeyedTransactor(privateKey), method, params...)

	if ok != nil {
		fmt.Println("txError=>", ok)
		panic("SendTransactionByKey 发送交易失败")
	}
	json, wrong := abi.JSON(strings.NewReader(whole.Config.Contract[name].Abi))
	if wrong != nil {
		fmt.Println("wrong==>", wrong)
		panic("解析abi失败")
	}

	var (
		result = new(any)
	)
	//合约方法名
	task := json.Unpack(&result, method, common.FromHex(receipt.Output))
	if task != nil {
		fmt.Println("task==>", task)
		panic("解析函数返回值失败")
	}
	return *result
}

// SendTransactionAndGetReceiptByKey 使用私钥发送交易,并获取交易回执
func SendTransactionAndGetReceiptByKey(name, method string, privateKey *goecdsa.PrivateKey, params ...interface{}) (any, *types.Receipt) {
	_, receipt, ok := whole.GoSdk.Contract[name].Transact(NewKeyedTransactor(privateKey), method, params...)

	if ok != nil {
		fmt.Println("txError=>", ok)
		panic("SendTransactionAndGetReceiptByKey 发送交易失败")
	}
	json, wrong := abi.JSON(strings.NewReader(whole.Config.Contract[name].Abi))
	if wrong != nil {
		fmt.Println("wrong==>", wrong)
		panic("解析abi失败")
	}

	var (
		result = new(any)
	)
	//合约方法名
	task := json.Unpack(&result, method, common.FromHex(receipt.Output))
	if task != nil {
		fmt.Println("task==>", task)
		panic("解析函数返回值失败")
	}
	return *result, receipt
}

// SendCallByKey SendCall 获取链上信息
func SendCallByKey(name, method string, privateKey *goecdsa.PrivateKey, out interface{}, params ...interface{}) any {
	clientCallOpts := &bind.CallOpts{From: NewKeyedTransactor(privateKey).From}
	err := whole.GoSdk.Contract[name].Call(clientCallOpts, out, method, params...)
	if err != nil {
		fmt.Println("SendCall err==>", err)
		panic("SendCallByKey 获取链上信息失败")
	}
	return out
}

// GetLatestBlocks 获取最新区块信息
func GetLatestBlocks(curBlockNumber int64) []types.Block {
	var blocks []types.Block
	client := whole.GoSdk.Client
	highestNumber, err := client.GetBlockNumber(context.Background())
	if err != nil {
		fmt.Println("txError=>", err.Error())
		return nil
	}

	for blockNumber := curBlockNumber + 1; blockNumber <= highestNumber; blockNumber++ {
		var block *types.Block
		block, err = client.GetBlockByNumber(context.Background(), blockNumber, true)
		if err != nil {
			fmt.Printf("Error getting block number %d: %v\n", blockNumber, err)
			break
		}
		blocks = append(blocks, *block)
	}
	return blocks
}
