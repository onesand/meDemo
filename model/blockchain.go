package model

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"time"
)

type UserAddress struct {
	ID        uint64
	Address   string
	CreatedAt time.Time
}

type TransactionsInBlockModel struct {
	TxHash   string `json:"tx_hash"`
	Value    string `json:"value"`
	Gas      uint64 `json:"gas"`
	GasPrice uint64 `json:"gas_price"`
	GasUsed  uint64 `json:"gas_used"`
	Nonce    uint64 `json:"nonce"`
	Data     string `json:"data"`
	From     string `json:"from"`
	To       string `json:"to"`
	Receipt  uint64 `json:"receipt"`
}

//合约解析所用

// LogTransfer ..
type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}

// LogApproval ..
type LogApproval struct {
	TokenOwner common.Address
	Spender    common.Address
	Tokens     *big.Int
}
