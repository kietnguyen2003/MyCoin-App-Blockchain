package api

import (
	"MyCoinApp/config"
	"MyCoinApp/internal/blockchain"
	"MyCoinApp/internal/models"
	"MyCoinApp/internal/pool"
	"MyCoinApp/internal/wallet"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	blockchain *blockchain.Blockchain
	txPool     *pool.TransactionPool
	config     *config.Config
}

func NewServer(bc *blockchain.Blockchain, cfg *config.Config) *Server {
	return &Server{
		blockchain: bc,
		txPool:     pool.NewTransactionPool(),
		config:     cfg,
	}
}

func AllowedHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, Accept, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func (s *Server) Start() error {
	router := gin.Default()

	router.Use(AllowedHeaders())
	api := router.Group("/api")
	{
		// wallet endpoints
		walletApi := api.Group("/wallet")
		{
			walletApi.POST("/create", s.createWallet)
			walletApi.POST("/import", s.importWallet)
			walletApi.GET("/balance/:address", s.getBalance)
		}

		blockChainApi := api.Group("/blockchain")
		{
			blockChainApi.POST("/mine", s.mineBlock)
			blockChainApi.GET("/info", s.getBlockchainInfo)
			// blockChainApi.GET("/blocks", s.getAllBlocks)
			// blockChainApi.GET("/block/:index", s.getBlock)
		}

		stakingApi := api.Group("/staking")
		{
			stakingApi.POST("/stake", s.stakeCoins)
			stakingApi.POST("/unstake", s.unstakeCoins)
			stakingApi.GET("/validators", s.getValidators)
			stakingApi.GET("/validator/:address", s.getValidatorInfo)
			stakingApi.GET("/info", s.getStakingInfo)
		}

		transactionApi := api.Group("/transaction")
		{
			transactionApi.POST("/send", s.sendTransaction)
			transactionApi.GET("/history/:address", s.getTransactionHistory)
		}

	}

	router.Static("/static", "./web/static")
	router.StaticFile("/", "./web/static/index.html")
	router.StaticFile("/favicon.ico", "./web/static/favicon.ico")

	fmt.Printf("MyCoin server starting on %s\n", s.config.Port)

	return router.Run(s.config.Port)
}

// Wallet handlers
func (s *Server) createWallet(c *gin.Context) {
	newWallet := wallet.NewWallet()

	log.Printf("New wallet created: %s", newWallet.Address)
	s.blockchain.AddBalance(newWallet.Address, s.config.InitialWalletBalance)
	log.Printf("Initial balance of %.2f added to wallet %s", s.config.InitialWalletBalance, newWallet.Address)

	response := models.CreateWalletResponse{
		Address:    newWallet.Address,
		PublicKey:  newWallet.GetPublicKeyHex(),
		PrivateKey: newWallet.GetPrivateKeyHex(),
	}
	c.JSON(http.StatusOK, response)
}

func (s *Server) getBalance(c *gin.Context) {
	address := c.Param("address")
	balance := s.blockchain.GetBalance(address)

	response := models.BalanceResponse{
		Address: address,
		Balance: balance,
	}
	c.JSON(http.StatusOK, response)

}

func (s *Server) importWallet(c *gin.Context) {
	var request struct {
		PrivateKey string `json:"private_key,omitempty"`
		Passphrase string `json:"passphrase,omitempty"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var importedWallet *wallet.Wallet
	var err error

	if request.PrivateKey != "" {
		importedWallet, err = wallet.LoadWalletFromPrivateKey(request.PrivateKey)
	} else if request.Passphrase != "" {
		importedWallet = wallet.NewWalletFromPassphrase(request.Passphrase)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either private_key or passphrase is required"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to import wallet: " + err.Error()})
		return
	}

	// Check if wallet already has balance, if not give free coins directly
	currentBalance := s.blockchain.GetBalance(importedWallet.Address)
	if currentBalance == 0 {
		s.blockchain.AddBalance(importedWallet.Address, s.config.InitialWalletBalance)
		log.Printf("Added %.2f MYC directly to imported wallet balance", s.config.InitialWalletBalance)
	}

	response := models.CreateWalletResponse{
		Address:    importedWallet.Address,
		PrivateKey: importedWallet.GetPrivateKeyHex(),
		PublicKey:  importedWallet.GetPublicKeyHex(),
	}

	c.JSON(http.StatusOK, response)
}

// Blockchain handlers
func (s *Server) getBlockchainInfo(c *gin.Context) {
	info := gin.H{
		"chain_length":         len(s.blockchain.Chain),
		"pending_transactions": len(s.blockchain.PendingTransactions),
		"mining_reward":        s.blockchain.MiningReward,
		"is_valid":             s.blockchain.IsChainValid(),
		"latest_block":         s.blockchain.GetLatestBlock(),
	}

	c.JSON(http.StatusOK, info)
}

func (s *Server) mineBlock(c *gin.Context) {
	log.Println("=== CREATE BLOCK REQUEST ===")

	var request struct {
		MinerAddress string `json:"miner_address"` // ❌ User chooses who mines!
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// No user input - system selects validator
	selectedValidator, err := s.blockchain.StakingPool.SelectValidator()
	if selectedValidator == "" {
		log.Printf("Failed to select validator: %v", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Không có validator nào khả dụng để tạo block. Vui lòng stake coin để trở thành validator.",
		})
		return
	}
	if err != nil {
		log.Printf("Failed to select validator: %v", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Không có validator nào khả dụng để tạo block. Vui lòng stake coin để trở thành validator.",
		})
		return
	}
	if selectedValidator != request.MinerAddress {
		log.Printf("Mining denied for address: %s (selected: %s)", request.MinerAddress, selectedValidator)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Bạn không phải validator được chọn để tạo block!!",
		})
		return
	}
	if len(s.blockchain.PendingTransactions) == 0 {
		log.Println("No pending transactions - will create block with reward transaction only")
	}

	log.Println("Creating block...")

	// Add timeout protection
	done := make(chan *blockchain.Block, 1)
	go func() {
		block := s.blockchain.MinePendingTransactions(selectedValidator)
		done <- block
	}()

	select {
	case block := <-done:
		if block == nil {
			log.Printf("Mining denied for address: %s", selectedValidator)
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Bạn không phải validator",
			})
			return
		}

		log.Printf("Block created successfully!")
		log.Printf("- Hash: %s", block.Hash)
		log.Printf("- Index: %d", block.Index)
		log.Printf("- Transactions: %d", len(block.Transactions))
		log.Printf("- Reward recipient: %s", selectedValidator)

		c.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"block_hash": block.Hash,
			"block":      block,
		})
		return

	case <-time.After(10 * time.Second):
		log.Println("Block creation timeout!")
		c.JSON(http.StatusRequestTimeout, gin.H{"error": "Block creation timeout"})
		return
	}
}

// Staking handlers
func (s *Server) getValidators(c *gin.Context) {
	validators := s.blockchain.GetValidators()

	c.JSON(http.StatusOK, gin.H{
		"validators": validators,
		"count":      len(validators),
	})
}

func (s *Server) stakeCoins(c *gin.Context) {
	log.Println("=== STAKE COINS REQUEST ===")
	var request struct {
		Address    string  `json:"address"`
		Amount     float64 `json:"amount"`
		PrivateKey string  `json:"private_key"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	log.Printf("Stake request - Address: %s, Amount: %.2f", request.Address, request.Amount)

	// Verify private key matches address
	log.Printf("Loading wallet from private key...")
	userWallet, err := wallet.LoadWalletFromPrivateKey(request.PrivateKey)
	if err != nil {
		log.Printf("Invalid private key error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid private key"})
		return
	}
	log.Printf("Wallet loaded successfully. Address: %s", userWallet.Address)

	if userWallet.Address != request.Address {
		log.Printf("Address mismatch: wallet=%s, request=%s", userWallet.Address, request.Address)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Private key does not match address"})
		return
	}
	log.Printf("Address verification passed")

	// Stake coins
	log.Printf("Starting stake operation...")
	err = s.blockchain.StakeCoins(request.Address, request.Amount)
	if err != nil {
		log.Printf("Stake operation failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Stake operation completed successfully")

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": fmt.Sprintf("Successfully staked %.2f MYC", request.Amount),
		"address": request.Address,
		"amount":  request.Amount,
	})
}

func (s *Server) getStakingInfo(c *gin.Context) {
	info := s.blockchain.GetStakingInfo()

	c.JSON(http.StatusOK, info)
}

func (s *Server) getValidatorInfo(c *gin.Context) {
	address := c.Param("address")

	validator, err := s.blockchain.GetValidatorInfo(address)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Validator not found"})
		return
	}

	c.JSON(http.StatusOK, validator)
}

func (s *Server) unstakeCoins(c *gin.Context) {
	var request struct {
		Address    string `json:"address"`
		PrivateKey string `json:"private_key"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Verify private key matches address
	userWallet, err := wallet.LoadWalletFromPrivateKey(request.PrivateKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid private key"})
		return
	}

	if userWallet.Address != request.Address {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Private key does not match address"})
		return
	}

	// Get validator info first
	validator, err := s.blockchain.GetValidatorInfo(request.Address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validator not found"})
		return
	}

	// Unstake coins
	err = s.blockchain.UnstakeCoins(request.Address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": fmt.Sprintf("Successfully unstaked %.2f MYC", validator.StakedAmount),
		"address": request.Address,
		"amount":  validator.StakedAmount,
	})
}

// Transaction handlers
func (s *Server) getTransactionHistory(c *gin.Context) {
	address := c.Param("address")

	transactions := s.blockchain.GetTransactionHistoryWithBlocks(address)

	response := models.TransactionHistoryResponse{
		Address:      address,
		Transactions: transactions,
	}

	c.JSON(http.StatusOK, response)
}

func (s *Server) sendTransaction(c *gin.Context) {
	/* Steps to handle a transaction:
	1. Parse and validate request body
	2. Load wallet from private key and verify 'from' address
	3. Check balance of 'from' address
	4. Create transaction and add to pending pool
	5. Return success response with transaction hash
	*/

	var request models.SendTransactionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Error decoding request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	log.Printf("Transaction request: From=%s, To=%s, Amount=%.2f, Fee=%.2f",
		request.From, request.To, request.Amount, request.Fee)

	userWallet, err := wallet.LoadWalletFromPrivateKey(request.PrivateKey)
	if err != nil {
		log.Printf("Error loading wallet from private key: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid private key"})
		return
	}

	log.Printf("Wallet loaded: %s", userWallet.Address)

	if userWallet.Address != request.From {
		log.Printf("Address mismatch: wallet=%s, request=%s", userWallet.Address, request.From)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Private key does not match from address"})
		return
	}

	// Check current balance before transaction
	currentBalance := s.blockchain.GetBalance(request.From)
	log.Printf("Current balance for %s: %.2f MYC", request.From, currentBalance)
	log.Printf("Required amount: %.2f MYC (%.2f + %.2f fee)", request.Amount+request.Fee, request.Amount, request.Fee)

	tx := pool.NewTransaction(request.From, request.To, request.Amount, request.Fee)
	log.Printf("Transaction created with hash: %s", tx.Hash)

	if err := s.blockchain.AddTransaction(tx); err != nil {
		log.Printf("Error adding transaction: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Transaction successfully added to pending pool")
	log.Printf("Pending transactions count: %d", len(s.blockchain.PendingTransactions))

	c.JSON(http.StatusOK, gin.H{
		"status":           "success",
		"transaction_hash": tx.Hash,
		"message":          "Transaction added to pending pool",
	})
}
