package constant

import "github.com/ethereum/go-ethereum/common"

var envConfig struct {
	Port        uint16 `env:"PORT" envDefault:"8080"`
	DatabaseURL string "postgres://igtquyeeikrosg:aeb181b46165c422f5c85037b635d6476a1b78edb846e5ab4c2bc8570ffe55a0@ec2-3-219-52-220.compute-1.amazonaws.com:5432/d93gh9lpdei3co"
	RedisURL    string `env:"REDIS_URL" envDefault:"redis://localhost:6379"`
	Environment string `env:"ENVIRONMENT" envDefault:"DEV"`
	NodeRpcURL  string "https://mainnet.infura.io/v3/05dd7c9be52e4582a455da059ffd106a"

	ExchangeAddress common.Address `env:"EXCHANGE_ADDRESS"`
	ChainID         uint           `env:"CHAIN_ID"`

	ReservoirApiEndpoint string `env:"RESERVOIR_API_ENDPOINT"`
	ReservoirApiKey      string `env:"RESERVOIR_API_KEY"`
}

func NodeRpcUrl() string {
	return "https://mainnet.infura.io/v3/05dd7c9be52e4582a455da059ffd106a"
}
func Port() uint16 {
	return envConfig.Port
}

func RedisURL() string {
	return envConfig.RedisURL
}

func DatabaseURL() string {
	return "postgres://igtquyeeikrosg:aeb181b46165c422f5c85037b635d6476a1b78edb846e5ab4c2bc8570ffe55a0@ec2-3-219-52-220.compute-1.amazonaws.com:5432/d93gh9lpdei3co"
}

func ReservoirApiEndpoint() string {
	return envConfig.ReservoirApiEndpoint
}

func ReservoirApiKey() string {
	return envConfig.ReservoirApiKey
}
