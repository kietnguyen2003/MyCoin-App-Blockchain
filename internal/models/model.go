package models

type CreateWalletResponse struct {
	Address    string `json:"address"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

type BalanceResponse struct {
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
}
