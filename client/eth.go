package client

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"meDemo/constant"
)

var ethClient *ethclient.Client

func EthClient() *ethclient.Client {
	return ethClient
}

func ConnectEthNode() error {
	cli, err := ethclient.Dial(constant.NodeRpcUrl())
	if err != nil {
		return err
	}

	ethClient = cli
	return nil
}
