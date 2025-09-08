package blockchain

import (
	"MyCoinApp/internal/pool"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strconv"
)

type Block struct {
	Index        int64               `json:"index"`
	Timestamp    int64               `json:"timestamp"`
	Transactions []*pool.Transaction `json:"transactions"`
	PreviousHash string              `json:"previous_hash"`
	Hash         string              `json:"hash"`
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
