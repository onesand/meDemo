package client

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"meDemo/constant"
	"meDemo/model"
)

var ethClient *ethclient.Client
var ethWsClient *ethclient.Client

func EthClient() *ethclient.Client {
	return ethClient
}

func EthWsClient() *ethclient.Client {
	return ethWsClient
}

func ConnectEthNode() error {
	cli, err := ethclient.Dial(constant.NodeRpcUrl())
	if err != nil {
		return err
	}

	ethClient = cli
	return nil
}

func ConnectEthWsNode() error {
	cli, err := ethclient.Dial(constant.NodeWsUrl())
	if err != nil {
		return err
	}

	ethWsClient = cli
	return nil
}

func TransitionsInBlock(num string) model.BaseResponse {
	blockNumber, _ := new(big.Int).SetString(num, 10)
	block, err := EthClient().BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	var result []model.TransactionsInBlockModel
	transactionModel := model.TransactionsInBlockModel{}
	response := model.BaseResponse{}
	for _, tx := range block.Transactions() {

		receipt, err := EthClient().TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}

		//(types.LatestSignerForChainID(big.NewInt(int64(1))), tx.GasPrice())
		msg, err := tx.AsMessage(types.LatestSignerForChainID(big.NewInt(int64(1))), tx.GasPrice())
		if err != nil {
			log.Fatal(err)
		}

		transactionModel = model.TransactionsInBlockModel{
			TxHash:   tx.Hash().Hex(),
			Value:    tx.Value().String(),
			Gas:      tx.Gas(),
			GasPrice: tx.GasPrice().Uint64(),
			GasUsed:  receipt.GasUsed,
			Nonce:    tx.Nonce(),
			//Data:   	string(tx.Data()),
			From:    msg.From().Hex(),
			To:      tx.To().Hex(),
			Receipt: receipt.Status,
		}

		result = append(result, transactionModel)
		//if len(result) == 1 {
		//	break
		//}
	}

	response.Code = 200
	response.Data = result
	return response
}

func SubscribeContractEvent() {
	err := ConnectEthWsNode()
	if err != nil {
		return
	}

	//0x3b484b82567a09e2588A13D54D032153f0c0aEe0
	contractAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := EthWsClient().SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Println(vLog) // pointer to event log
		}
	}
}
