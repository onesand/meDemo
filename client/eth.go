package client

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"meDemo/constant"
	token "meDemo/contract"
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

	var freeModels []model.FreeMintMode
	var contractIds []string
	filter := ethereum.FilterQuery{}
	logs := make(chan types.Log)
	sub, err := ethWsClient.SubscribeFilterLogs(context.Background(), filter, logs)
	if err != nil {
		return
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			// erc721
			// topics是否Transfer，并且只包含三个参数   是否从黑洞地址发出
			if len(vLog.Topics) == 4 && vLog.Topics[0] == common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef") && vLog.Topics[1] == common.HexToHash("0x0000000000000000000000000000000000000000") {

				txHash := vLog.TxHash.String()
				contractsAddress := vLog.Address.String()
				//contractsAddress := "0x57f1887a8BF19b14fC0dF6Fd9B2acc9Af147eA85"
				tokenId := vLog.Topics[3].Big().String()

				//已经记录，跳过
				if contains(contractIds, contractsAddress) {
					println("已经处理，略过")
					continue
				}

				println("检测到结果,txHash:" + txHash)
				println("合约地址:" + contractsAddress)
				println("正在mint的tokenId：" + tokenId)

				erc721, err := token.NewERC721(common.HexToAddress(contractsAddress), ethWsClient)
				if err != nil {
					return
				}

				opts := bind.CallOpts{
					Pending:     false,
					BlockNumber: nil,
					Context:     nil,
				}
				tokenName, err := erc721.Name(&opts)
				if err != nil {
					println("无name，pass")
					tokenName = ""
				}
				println("name==>>>" + tokenName)

				supply, _ := erc721.TotalSupply(&opts)
				if err != nil {
					println("此合约未开源或者是纯土狗，pass")
					supply = big.NewInt(0)
				}
				println("supply==>>>" + supply.String())

				contractIds = append(contractIds, contractsAddress)

				//达到一定数量后，存库
				if len(contractIds) == 10 {
					println("入库")
				}

				freeModel := model.FreeMintMode{
					TxHash:          txHash,
					ContractAddress: contractsAddress,
					TokenId:         tokenId,
					TokenName:       tokenName,
				}
				freeModels = append(freeModels, freeModel)
				println("\n\n")
			}
		}
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
