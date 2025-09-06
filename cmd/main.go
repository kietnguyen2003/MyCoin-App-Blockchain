package main

import (
	"MyCoinApp/config"
	"MyCoinApp/internal/api"
	"MyCoinApp/internal/blockchain"
	"fmt"
	"log"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Println("Config loaded:", cfg)

	if err := cfg.Validate(); err != nil {
		fmt.Println("Config validation error:", err)
		return
	}

	bc := blockchain.NewBlockchain()
	log.Printf("Blockchain initialized with %d blocks", len(bc.Chain))

	srv := api.NewServer(bc, cfg)

	log.Printf("Server starting on %s", cfg.Port)
	log.Fatal(srv.Start())

}
