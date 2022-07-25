package model

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type UserAddress struct {
	ID        uint64
	Address   string
	CreatedAt int
	UpdatedAt int
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

// FreeMintMode FreeMint 项目信息
type FreeMintMode struct {
	ID              uint64
	CreatedAt       int
	UpdatedAt       int
	TxHash          string `json:"tx_hash"`
	ContractAddress string `gorm:"primaryKey" json:"contract_address"`
	TokenId         string `json:"token_id"`
	TokenName       string `json:"token_name"`
	TotalSupply     string `json:"total_supply"`
	FollowBy        string `json:"follow_by"`
}

// Mints address mint 的信息
type Mints struct {
	ID              uint64
	CreatedAt       int
	TxHash          string `json:"tx_hash"` // mint txHash
	Address         string `json:"address"` // mint address
	ContractAddress string `json:"contract_address"`
	TokenName       string `json:"token_name"`
	TotalSupply     string `json:"total_supply"`
	FollowBy        string `json:"follow_by"`
}

type Nft struct {
	ID              string `json:"id" 'gorm:"primaryKey" json:"id"`
	CreatedAt       int
	UpdatedAt       int
	ContractAddress string `json:"contract_address"`
	TokenId         string `json:"token_id"`
	Owner           string `json:"owner"`
	TxHash          string `json:"tx_hash"`
}

type Asset struct {
	ID              string `gorm:"primaryKey" json:"id"`
	CreatedAt       int
	UpdatedAt       int
	ContractAddress string `json:"contract_address"`
	AssetType       string `json:"asset_type"`
	Balance         string `json:"balance"`
	Address         string `json:"address"`
	Nfts            []Nft  `gorm:"-:all"`
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
