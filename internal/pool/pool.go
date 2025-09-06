package pool

import "fmt"

type TransactionPool struct {
	Transactions []*Transaction `json:"transactions"`
}

func NewTransactionPool() *TransactionPool {
	return &TransactionPool{
		Transactions: make([]*Transaction, 0),
	}
}

func (tp *TransactionPool) AddTransaction(tx *Transaction) error {
	if !tx.IsValid() {
		return fmt.Errorf("invalid transaction")
	}

	for _, existingTx := range tp.Transactions {
		if existingTx.Hash == tx.Hash {
			return fmt.Errorf("transaction already exists in pool")
		}
	}

	tp.Transactions = append(tp.Transactions, tx)
	return nil
}

func (tp *TransactionPool) GetTransactions(limit int) []*Transaction {
	if limit <= 0 || limit > len(tp.Transactions) {
		return tp.Transactions
	}
	return tp.Transactions[:limit]
}

func (tp *TransactionPool) RemoveTransactions(txHashes []string) {
	for _, hash := range txHashes {
		for i, tx := range tp.Transactions {
			if tx.Hash == hash {
				tp.Transactions = append(tp.Transactions[:i], tp.Transactions[i+1:]...)
				break
			}
		}
	}
}
