package client

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"meDemo/constant"
	token "meDemo/contract"
	"meDemo/model"
	"strconv"
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
		println("Connect to eth ws error." + err.Error())
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

func Erc20LogsInBlock(start *big.Int, end *big.Int) {
	err := ConnectEthWsNode()
	if err != nil {
		return
	}

	//sos:0x3b484b82567a09e2588A13D54D032153f0c0aEe0
	//bayc:0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D
	contractAddress := common.HexToAddress("0x3b484b82567a09e2588A13D54D032153f0c0aEe0")

	if end == big.NewInt(0) {
		end = nil
	}
	filter := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		FromBlock: start,
		ToBlock:   end,
		Topics:    nil,
	}

	filterLogs, err := EthWsClient().FilterLogs(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	logTransferSig := []byte("Transfer(address,address,uint256)")
	LogApprovalSig := []byte("Approval(address,address,uint256)")
	swapSig := []byte("Swap(address,uint256,uint256,uint256,uint256,address)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	logApprovalSigHash := crypto.Keccak256Hash(LogApprovalSig)
	swapSigHash := crypto.Keccak256Hash(swapSig)

	println("swapSigHash:" + swapSigHash.Hex())
	for _, vLog := range filterLogs {

		fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
		fmt.Printf("Log Index: %d\n", vLog.Index)
		fmt.Printf("Log txHash: " + vLog.TxHash.String() + "\n")

		for _, topic := range vLog.Topics {
			println("topic:" + topic.Hex())
		}
		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			fmt.Printf("Log Name: Transfer\n")

			erc20Filterer, err := token.NewERC20Filterer(contractAddress, nil)
			if err != nil {
				return
			}
			parseApproval, err := erc20Filterer.ParseTransfer(vLog)
			if err != nil {
				return
			}
			fmt.Printf("From: %s\n", parseApproval.From.Hex())
			fmt.Printf("To: %s\n", parseApproval.To.Hex())
			fmt.Printf("Tokens: %s\n", parseApproval.Value.String())

		case logApprovalSigHash.Hex():
			fmt.Printf("Log Name: Approval\n")

			erc20Filterer, err := token.NewERC20Filterer(contractAddress, nil)
			if err != nil {
				log.Panic(err)
				return
			}
			parseApproval, err := erc20Filterer.ParseApproval(vLog)
			if err != nil {
				log.Panic(err)
				return
			}
			fmt.Printf("Token Owner: %s\n", parseApproval.Owner.Hex())
			fmt.Printf("Spender: %s\n", parseApproval.Spender.Hex())
			fmt.Printf("Tokens: %s\n", parseApproval.Value.String())

		case swapSigHash.Hex():
			fmt.Printf("Log Name: Swap\n")
		}

		fmt.Printf("\n\n")
	}

}

func FreeMintMonitor() {
	err := ConnectEthWsNode()
	if err != nil {
		return
	}

	number, err := ethWsClient.BlockNumber(context.Background())
	if err != nil {
		return
	}
	println("当前最新区块：" + strconv.FormatUint(number, 10))

	block, err2 := ethWsClient.BlockByNumber(context.Background(), big.NewInt(int64(number)))
	if err2 != nil {
		return
	}

	var txs = block.Transactions()
	print("当前区块包含交易数:")
	println(txs.Len())
	//contractMap := []string{}
	for _, tx := range txs {
		receipt, err := ethWsClient.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			return
		}
		if receipt.Status == 1 && tx.Value().String() == "0" {
			// 判断是否erc721
			// topics是否Transfer，并且只包含三个参数   是否从黑洞地址发出
			if len(receipt.Logs) > 0 && len(receipt.Logs[0].Topics) == 4 && receipt.Logs[0].Topics[0] == common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef") && receipt.Logs[0].Topics[1] == common.HexToHash("0x0000000000000000000000000000000000000000") {
				//amount := len(receipt.Logs)
				//collection := []string{}
				//if contractMap. {
				//
				//}

				println("检测到结果,txHash:" + receipt.TxHash.String())
				println("合约地址:" + receipt.Logs[0].Address.String())
				println("正在mint的tokenId：")
				print(receipt.Logs[0].Topics[3].Big().String())
			}
		}
	}

}
