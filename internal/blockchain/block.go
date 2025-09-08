package blockchain

import (
	"MyCoinApp/internal/pool"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"time"
)

type Block struct {
	Index        int64               `json:"index"`
	Timestamp    int64               `json:"timestamp"`
	Transactions []*pool.Transaction `json:"transactions"`
	PreviousHash string              `json:"previous_hash"`
	Hash         string              `json:"hash"`
}

func NewBlock(transactions []*pool.Transaction, previousHash string, validator string, blockNumber int64) *Block {
	block := &Block{
		Index:        blockNumber,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
		PreviousHash: previousHash,
	}

	// Calculate hash directly (no mining)
	block.Hash = block.CalculateHashPOS(validator)
	return block
}

func (b *Block) CalculateHashPOS(validator string) string {
	data := strconv.FormatInt(b.Index, 10) +
		strconv.FormatInt(b.Timestamp, 10) +
		b.PreviousHash +
		validator
	for _, tx := range b.Transactions {
		txBytes, _ := json.Marshal(tx)
		data += string(txBytes)
	}

	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (b *Block) CalculateHash() string {
	data := strconv.FormatInt(b.Index, 10) +
		strconv.FormatInt(b.Timestamp, 10) +
		b.PreviousHash
	for _, tx := range b.Transactions {
		txBytes, _ := json.Marshal(tx)
		data += string(txBytes)
	}

	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
