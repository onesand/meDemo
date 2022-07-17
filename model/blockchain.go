package model

import "time"

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
	Nonce    uint64 `json:"nonce"`
	Data     string `json:"data"`
	From     string `json:"from"`
	To       string `json:"to"`
	Receipt  uint64 `json:"receipt"`
}
