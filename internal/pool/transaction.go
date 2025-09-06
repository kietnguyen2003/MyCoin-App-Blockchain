package pool

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"strconv"
	"time"
)

type Transaction struct {
	From      string  `json:"from"`
	To        string  `json:"to"`
	Amount    float64 `json:"amount"`
	Fee       float64 `json:"fee"`
	Timestamp int64   `json:"timestamp"`
	Hash      string  `json:"hash"`
	Signature string  `json:"signature"`
}

func NewTransaction(from, to string, amount, fee float64) *Transaction {
	tx := &Transaction{
		From:      from,
		To:        to,
		Amount:    amount,
		Fee:       fee,
		Timestamp: time.Now().Unix(),
	}

	tx.Hash = tx.CalculateHash()
	return tx
}

func (tx *Transaction) CalculateHash() string {
	data := tx.From + tx.To +
		strconv.FormatFloat(tx.Amount, 'f', -1, 64) +
		strconv.FormatFloat(tx.Fee, 'f', -1, 64) +
		strconv.FormatInt(tx.Timestamp, 10)

	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (tx *Transaction) SignTransaction(privateKey *ecdsa.PrivateKey) error {
	if tx.From == "" || tx.From == "genesis" {
		return nil
	}

	hashBytes, err := hex.DecodeString(tx.Hash)
	if err != nil {
		return err
	}

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashBytes)
	if err != nil {
		return err
	}

	signature := append(r.Bytes(), s.Bytes()...)
	tx.Signature = hex.EncodeToString(signature)

	return nil
}

func (tx *Transaction) VerifySignature(publicKey []byte) bool {
	if tx.From == "" || tx.From == "genesis" {
		return true
	}

	if tx.Signature == "" {
		return false
	}

	signatureBytes, err := hex.DecodeString(tx.Signature)
	if err != nil {
		return false
	}

	if len(signatureBytes) != 64 {
		return false
	}

	r := big.NewInt(0).SetBytes(signatureBytes[:32])
	s := big.NewInt(0).SetBytes(signatureBytes[32:])

	hashBytes, err := hex.DecodeString(tx.Hash)
	if err != nil {
		return false
	}

	x := big.NewInt(0).SetBytes(publicKey[:32])
	y := big.NewInt(0).SetBytes(publicKey[32:])

	pubKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	return ecdsa.Verify(&pubKey, hashBytes, r, s)
}

func (tx *Transaction) IsValid() bool {
	if tx.From == tx.To {
		return false
	}

	if tx.Amount <= 0 {
		return false
	}

	if tx.Fee < 0 {
		return false
	}

	if tx.Hash != tx.CalculateHash() {
		return false
	}

	return true
}

func (tx *Transaction) ToJSON() (string, error) {
	data, err := json.MarshalIndent(tx, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func FromJSON(jsonStr string) (*Transaction, error) {
	var tx Transaction
	err := json.Unmarshal([]byte(jsonStr), &tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}
