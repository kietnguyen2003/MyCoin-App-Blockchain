package models

import "MyCoinApp/internal/pool"

type CreateWalletResponse struct {
	Address    string `json:"address"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

type BalanceResponse struct {
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
}

type TransactionWithBlock struct {
	*pool.Transaction
	BlockIndex int64  `json:"block_index"`
	BlockHash  string `json:"block_hash"`
}

type TransactionHistoryResponse struct {
	Address      string                  `json:"address"`
	Transactions []*TransactionWithBlock `json:"transactions"`
}

type SendTransactionRequest struct {
	From       string  `json:"from"`
	To         string  `json:"to"`
	Amount     float64 `json:"amount"`
	Fee        float64 `json:"fee"`
	PrivateKey string  `json:"private_key"`
}
