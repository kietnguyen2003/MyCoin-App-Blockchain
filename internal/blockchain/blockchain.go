package blockchain

import (
	"MyCoinApp/internal/consensus"
	"MyCoinApp/internal/models"
	"MyCoinApp/internal/pool"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type Blockchain struct {
	Chain               []*Block            `json:"chain"`
	PendingTransactions []*pool.Transaction `json:"pending_transactions"`
	MiningReward        float64             `json:"mining_reward"`

	Balances map[string]float64 `json:"balances"`

	StakingPool *consensus.StakingPool `json:"staking_pool"`
	mutex       sync.RWMutex           `json:"-"`
}

func NewBlockchain() *Blockchain {
	bc := &Blockchain{
		Chain:               []*Block{},
		PendingTransactions: []*pool.Transaction{},
		MiningReward:        50.0, // Default mining reward

		Balances:    make(map[string]float64),
		StakingPool: consensus.NewStakingPool(),
	}

	bc.CreateGenesisBlock()
	bc.LoadFromFile()
	return bc
}

// CreateGenesisBlock tạo block đầu tiên (genesis block) của blockchain
// Genesis block là block không có previous hash và chứa coin ban đầu
func (bc *Blockchain) CreateGenesisBlock() {
	// Tạo genesis block với thông tin cơ bản
	genesisBlock := &Block{
		Index:        0,
		Timestamp:    1609459200,
		Transactions: []*pool.Transaction{},
		PreviousHash: "0",
		Hash:         "0",
	}
	bc.Chain = append(bc.Chain, genesisBlock)
	bc.Balances["genesis"] = 1000000.0
}

func (bc *Blockchain) SaveToFile() error {
	data, err := json.MarshalIndent(bc, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile("blockchain.json", data, 0644)
}

func (bc *Blockchain) LoadFromFile() error {
	if _, err := os.Stat("blockchain.json"); os.IsNotExist(err) {
		return nil
	}

	data, err := ioutil.ReadFile("blockchain.json")
	if err != nil {
		return err
	}

	return json.Unmarshal(data, bc)
}

func (bc *Blockchain) AddBalance(address string, amount float64) {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()
	bc.Balances[address] += amount
	log.Printf("Balance of %s updated to %.2f", address, bc.Balances[address])
	bc.SaveToFile()
}

func (bc *Blockchain) GetBalance(address string) float64 {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()
	// Tìm balance trong map
	balance, exists := bc.Balances[address]
	if !exists {
		return 0.0
	}
	return balance
}

func (bc *Blockchain) IsChainValid() bool {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	for i := 1; i < len(bc.Chain); i++ {
		currentBlock := bc.Chain[i]
		previousBlock := bc.Chain[i-1]

		if currentBlock.Hash != currentBlock.CalculateHash() {
			return false
		}

		if currentBlock.PreviousHash != previousBlock.Hash {
			return false
		}
	}

	return true
}

func (bc *Blockchain) GetLatestBlock() *Block {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	if len(bc.Chain) == 0 {
		return nil
	}
	return bc.Chain[len(bc.Chain)-1]
}

func (bc *Blockchain) GetValidators() []*consensus.Validator {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	return bc.StakingPool.GetAllValidators()
}

func (bc *Blockchain) GetTransactionHistoryWithBlocks(address string) []*models.TransactionWithBlock {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	var transactions []*models.TransactionWithBlock

	for _, block := range bc.Chain {
		for _, tx := range block.Transactions {
			if tx.From == address || tx.To == address {
				txWithBlock := &models.TransactionWithBlock{
					Transaction: tx,
					BlockIndex:  block.Index,
					BlockHash:   block.Hash,
				}
				transactions = append(transactions, txWithBlock)
			}
		}
	}

	return transactions
}

func (bc *Blockchain) AddTransaction(transaction *pool.Transaction) error {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	if transaction.From != "genesis" && transaction.From != "" {
		// Lấy balance hiện tại của người gửi
		balance, exists := bc.Balances[transaction.From]
		if !exists { // Nếu chưa có trong map
			balance = 0.0 // Đặt balance = 0
		}
		// Kiểm tra đủ tiền (gồm cả amount + fee)
		if balance < transaction.Amount+transaction.Fee {
			return fmt.Errorf("insufficient balance") // Lỗi: Không đủ tiền
		}
	}

	// Thêm transaction vào pending pool
	bc.PendingTransactions = append(bc.PendingTransactions, transaction)
	return nil // Thành công
}
