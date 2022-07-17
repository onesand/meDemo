package test

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
	"math/big"
	"meDemo/client"
	"meDemo/model"
)

func TransitionsInBlock(num string) string {
	client.SetUpEthClient()

	blockNumber, _ := new(big.Int).SetString(num, 10)
	block, err := client.EthClient().BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	var result []model.TransactionsInBlockModel
	resultModel := model.TransactionsInBlockModel{}
	for _, tx := range block.Transactions() {

		receipt, err := client.EthClient().TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}

		msg, err := tx.AsMessage(types.NewEIP155Signer(big.NewInt(1)), big.NewInt(0))
		if err != nil {
			log.Fatal(err)
		}

		resultModel = model.TransactionsInBlockModel{
			TxHash:   tx.Hash().Hex(),
			Value:    tx.Value().String(),
			Gas:      tx.Gas(),
			GasPrice: tx.GasPrice().Uint64(),
			Nonce:    tx.Nonce(),
			Data:     string(tx.Data()),
			From:     msg.From().Hex(),
			To:       tx.To().Hex(),
			Receipt:  receipt.Status,
		}

		result = append(result, resultModel)

		if len(result) == 1 {
			break
		}
	}

	resultJson, err := json.Marshal(result)
	if err != nil {
		println(err.Error())
		return ""
	}
	println(string(resultJson))
	return string(resultJson)
}
