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
	PrevHash     string              `json:"prev_hash"`
	Hash         string              `json:"hash"`
}

func NewBlock(transactions []*pool.Transaction, previousHash string, validator string, blockNumber int64) *Block {
	block := &Block{
		Index:        blockNumber,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
		PrevHash:     previousHash,
	}

	// Calculate hash directly (no mining)
	block.Hash = block.CaculateHash(validator)
	return block
}

func (b *Block) CaculateHash(validator string) string {
	data := strconv.FormatInt(b.Index, 10) +
		strconv.FormatInt(b.Timestamp, 10) +
		b.PrevHash +
		validator

	for _, tx := range b.Transactions {
		txBytes, _ := json.Marshal(tx)
		data += string(txBytes)
	}

	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
