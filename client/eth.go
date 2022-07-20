package client

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"log"
	"math/big"
	"meDemo/constant"
	token "meDemo/contract"
	"meDemo/model"
	"strings"
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

	SetupConnections()

	followAddress := []string{
		"0x8a42f0ab1dcbb65ca290d2b11bd3d88563569070",
		"0xA6f4fa9840Aa6825446c820aF6d5eCcf9f78BA05",
		"0x9c8F92bddF72b5B36Eaa4EA7F3d581CEc0802c13",
		"0x709bF4aC7ED6Bb2F9d60b1215d983496AB68efbc",
		"0xd640C898B0902bD02f69dE0FE8d0bd560956DB76",
		"0x84BDbEaB9Dd28C17C6E11702934E5E9cFe566462",
		"0x4cffe1FEa2B6918F6d9596B8274d0D859Ab1699e",
		"0x6868B90BA68E48b3571928A7727201B9efE1D374",
		"0x0fe60E55a8C0700b47d4a2663079c445Fc4A5893",
		"0xba69593F1F51D4b2ff0a73298c8bE0C8586be931",
		"0x6Eb5f7C3Aa91e974bE11f23CaBD3532458070CB9",
		"0x18cCC241CcE98a67564E367eBc1F1f0e692E3188",
		"0xA6C88571d0028f47ADba983A7240Bf12Af94633e",
		"0xd6F6E99c4905c6e8A751Bb0aFeEFA8Dcc56a30dC",
		"0x4c8F62f1498FA55D4158CdBFEA7783f84556a68e",
		"0x0BeDa5116cD204c428379b5D852DADc04F3Bc384",

		"0x6238872A0Bd9F0E19073695532a7Ed77ce93C69E",
	}
	var freeModels []model.FreeMintMode
	var contractIds []string
	DB().Table("free_mint_modes").Select("contract_address").Find(&contractIds)
	//marshal, err := json.Marshal(contractIds)
	//if err != nil {
	//	return
	//}
	//println(string(marshal))

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
				tokenId := vLog.Topics[3].Big().String()

				//已经记录，跳过
				if contains(contractIds, contractsAddress) {
					//println("已经处理，略过")
					continue
				}

				println("检测到结果,txHash:" + txHash)
				println("合约地址:" + contractsAddress)
				println("正在mint的tokenId：" + tokenId)

				transaction, isPending, err := ethWsClient.TransactionByHash(context.Background(), vLog.TxHash)
				if err != nil {
					println(err.Error())
					println("获取交易详情出错")
					println("\n\n")
					continue
				}

				if isPending {
					println("isPending:true")
				} else {
					println("isPending:false")
				}

				mint(contractsAddress, transaction.Data())

				// 判断是否free mint
				if transaction.Value().Int64() != 0 {
					println("此笔交易需要value：" + transaction.Value().String() + " 略过")
					println("\n\n")
					continue
				}

				//rinkeby 4
				msg, err := transaction.AsMessage(types.LatestSignerForChainID(big.NewInt(int64(4))), transaction.GasPrice())
				if err != nil {
					println(err.Error())
					println("AsMessage获取出错")
					println("\n\n")
					continue
				}

				// 收发不能是自己
				if strings.ToLower(transaction.To().String()) == strings.ToLower(msg.From().String()) {
					println("此笔交易是自己转账到自己，略过")
					println("\n\n")
					continue
				}
				// 判断是否跟单列表的地址
				if !contains(followAddress, msg.From().String()) {
					println("在跟单列表里不匹配，略过")
					println("\n\n")
					continue
				}

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
				println("name：" + tokenName)

				supply, _ := erc721.TotalSupply(&opts)
				if err != nil {
					println("此合约未开源或者是纯土狗，pass")
					supply = big.NewInt(0)
				}
				println("supply：" + supply.String())

				contractIds = append(contractIds, contractsAddress)

				//达到一定数量后，存库
				if len(freeModels) == 2 {
					println("入库")
					DB().Create(&freeModels)
					freeModels = nil
				}

				freeModel := model.FreeMintMode{
					TxHash:          txHash,
					ContractAddress: contractsAddress,
					TokenId:         tokenId,
					TokenName:       tokenName,
					TotalSupply:     supply.String(),
					FollowBy:        msg.From().String(),
				}
				freeModels = append(freeModels, freeModel)

				// 发送消息到 dc
				dcMessage := "监测到有新的FreeMint\n"
				dcMessage = dcMessage + "合约地址：https://etherscan.io/token/" + contractsAddress + "\n"
				send, err := BOT().ChannelMessageSend("999193783797293148", dcMessage)
				if err != nil {
					println("bot机器人发消息出错：" + err.Error() + "  " + send.Content)
					return
				}
				println("\n\n")
			}
		}
	}
}

func mint(contractAddress string, data []byte) {
	//fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19
	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := EthWsClient().PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(0)
	gasLimit := uint64(152818)
	gasPrice, err := EthWsClient().SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress(contractAddress)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    value,
		Gas:      gasLimit + uint64(68*len(data)),
		GasPrice: gasPrice,
		Data:     data,
	})

	//tx := types.NewTransaction(nonce, toAddress, value, nil, gasPrice, data)

	chainID, err := EthWsClient().NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	ts := types.Transactions{signedTx}
	bytedata := new(bytes.Buffer)
	ts.EncodeIndex(0, bytedata)
	rawTxHex := hex.EncodeToString(bytedata.Bytes())

	println("需要gas:" + ToDecimal(CalcGasCost(gasLimit, gasPrice), 18).String())
	println("0x" + rawTxHex)

	//rawTxBytes, err := hex.DecodeString(rawTxHex)
	//tx = new(types.Transaction)
	//err = rlp.DecodeBytes(rawTxBytes, &tx)
	//if err != nil {
	//	return
	//}
	//
	//err = EthClient().SendTransaction(context.Background(), tx)
	//if err != nil {
	//	println("SendTransaction 出错：" + err.Error())
	//}
	//
	//println("tx sent: ", tx.Hash().Hex())
}

func ToDecimal(ivalue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}

// CalcGasCost calculate gas cost given gas limit (units) and gas price (wei)
func CalcGasCost(gasLimit uint64, gasPrice *big.Int) *big.Int {
	gasLimitBig := big.NewInt(int64(gasLimit))
	return gasLimitBig.Mul(gasLimitBig, gasPrice)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.ToLower(a) == strings.ToLower(e) {
			return true
		}
	}
	return false
}
