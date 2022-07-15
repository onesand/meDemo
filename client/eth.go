package client

import (
	"github.com/ethereum/go-ethereum/ethclient"
)

var ethClient *ethclient.Client

func EthClient() *ethclient.Client {
	return ethClient
}

func ConnectEthNode() error {
	cli, err := ethclient.Dial("https://mainnet.infura.io/v3/05dd7c9be52e4582a455da059ffd106a")
	if err != nil {
		return err
	}

	ethClient = cli
	return nil
}
