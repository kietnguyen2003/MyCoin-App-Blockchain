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
			// walletApi.POST("/import", s.ImportWallet)
			walletApi.GET("/balance/:address", s.getBalance)
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
