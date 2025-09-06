package main

import (
	"MyCoinApp/config"
	"fmt"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Println("Config loaded:", cfg)

	if err := cfg.Validate(); err != nil {
		fmt.Println("Config validation error:", err)
		return
	}

}
