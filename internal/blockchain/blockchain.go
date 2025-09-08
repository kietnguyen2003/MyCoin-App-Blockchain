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

func (bc *Blockchain) MinePendingTransactions(miningRewardAddress string) *Block {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	// Log thông tin debug
	log.Printf("MinePendingTransactions called with address: %s", miningRewardAddress)
	log.Printf("Consensus type: POS")

	// Tạo PoS block với validator được đề xuất
	block, err := bc.createPoS(miningRewardAddress)
	if err != nil {
		log.Printf("Error creating PoS block: %v", err)
		return nil
	}
	return block // Trả về block đã tạo thành công

}

func (bc *Blockchain) createPoS(proposedValidator string) (*Block, error) {
	// Bootstrap mode: Nếu chưa có validator nào, cho phép bất kỳ ai tạo block
	if len(bc.StakingPool.Validators) == 0 {
		log.Println("Bootstrap mode: No validators found, allowing anyone to create block")
		selectedValidator := proposedValidator

		// Reduced reward for PoS (50 MYC instead of 100)
		rewardAmount := bc.StakingPool.BlockReward
		log.Printf("Creating reward transaction: %s -> %.2f MYC", selectedValidator, rewardAmount)
		rewardTransaction := pool.NewTransaction("", selectedValidator, rewardAmount, 0)
		bc.PendingTransactions = append(bc.PendingTransactions, rewardTransaction)

		log.Printf("Total transactions for block: %d", len(bc.PendingTransactions))

		// Create block without validator requirements (bootstrap mode)
		previousHash := "0"
		if len(bc.Chain) > 0 {
			// Direct access since we already have the lock
			latestBlock := bc.Chain[len(bc.Chain)-1]
			if latestBlock != nil {
				previousHash = latestBlock.Hash
			}
		}

		blockNumber := int64(len(bc.Chain))
		log.Printf("Creating PoS block #%d with previous hash: %s", blockNumber, previousHash)
		block := NewBlock(bc.PendingTransactions, previousHash, selectedValidator, blockNumber)
		log.Printf("PoS block #%d created with hash: %s", blockNumber, block.Hash)

		log.Printf("Block created, adding to chain...")
		bc.Chain = append(bc.Chain, block)

		log.Printf("Updating balances...")
		bc.UpdateBalances(block)

		bc.PendingTransactions = []*pool.Transaction{}

		log.Printf("Saving blockchain to file...")
		bc.SaveToFile()

		log.Printf("Bootstrap block creation completed!")
		return block, nil
	}

	// Kiểm tra proposedValidator có phải là validator hợp lệ không
	_, err := bc.StakingPool.GetValidatorInfo(proposedValidator)
	if err != nil {
		log.Printf("Mining denied: address %s is not an active validator", proposedValidator)
		return nil, fmt.Errorf("mining denied: address %s is not an active validator", proposedValidator)
	}

	selectedValidator := proposedValidator
	rewardAmount := bc.StakingPool.BlockReward
	rewardTransaction := pool.NewTransaction("", selectedValidator, rewardAmount, 0)
	bc.PendingTransactions = append(bc.PendingTransactions, rewardTransaction)

	previousHash := "0"
	if len(bc.Chain) > 0 {
		latestBlock := bc.Chain[len(bc.Chain)-1]
		if latestBlock != nil {
			previousHash = latestBlock.Hash
		}
	}
	blockNumber := int64(len(bc.Chain))
	block := NewBlock(bc.PendingTransactions, previousHash, selectedValidator, blockNumber)
	bc.Chain = append(bc.Chain, block)

	bc.UpdateBalances(block)

	bc.StakingPool.RewardValidator(selectedValidator, rewardAmount)

	bc.PendingTransactions = []*pool.Transaction{}
	bc.SaveToFile()

	return block, nil
}

func (bc *Blockchain) UpdateBalances(block *Block) {
	for _, tx := range block.Transactions {
		if tx.From != "" && tx.From != "genesis" {
			bc.Balances[tx.From] -= (tx.Amount + tx.Fee)
		}
		if tx.To != "" {
			bc.Balances[tx.To] += tx.Amount
		}
	}
}

func (bc *Blockchain) StakeCoins(address string, amount float64) error {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	// Check if user has enough balance (direct access to avoid mutex deadlock)
	balance, exists := bc.Balances[address]
	if !exists {
		balance = 0.0
	}
	if balance < amount {
		return fmt.Errorf("insufficient balance for staking")
	}

	// Deduct staked amount from balance
	bc.Balances[address] -= amount

	// Add validator to staking pool
	err := bc.StakingPool.AddValidator(address, amount)
	if err != nil {
		// Refund if failed
		bc.Balances[address] += amount
		return err
	}

	bc.SaveToFile()
	return nil
}

func (bc *Blockchain) GetStakingInfo() map[string]interface{} {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	return map[string]interface{}{
		"total_staked":      bc.StakingPool.GetTotalStaked(),
		"min_stake_amount":  bc.StakingPool.MinStakeAmount,
		"max_validators":    bc.StakingPool.MaxValidators,
		"active_validators": len(bc.StakingPool.Validators),
		"block_reward":      bc.StakingPool.BlockReward,
		"staking_reward":    bc.StakingPool.StakingReward,
		"slashing_penalty":  bc.StakingPool.SlashingPenalty,
	}
}

func (bc *Blockchain) GetValidatorInfo(address string) (*consensus.Validator, error) {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	return bc.StakingPool.GetValidatorInfo(address)
}

func (bc *Blockchain) UnstakeCoins(address string) error {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	validator, err := bc.StakingPool.GetValidatorInfo(address)
	if err != nil {
		return err
	}

	// Return staked coins to balance
	bc.Balances[address] += validator.StakedAmount

	// Remove validator
	err = bc.StakingPool.RemoveValidator(address)
	if err != nil {
		bc.Balances[address] -= validator.StakedAmount // Revert
		return err
	}

	bc.SaveToFile()
	return nil
}
