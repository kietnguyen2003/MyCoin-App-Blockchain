package config

import "fmt"

type Config struct {
	Port                 string
	InitialWalletBalance float64
}

func LoadConfig() *Config {
	return &Config{
		Port:                 ":8080",
		InitialWalletBalance: 100.0,
	}
}

func (c *Config) Validate() error {
	if c.Port == "" {
		return fmt.Errorf("port cannot be empty")
	}
	if c.InitialWalletBalance < 0 {
		return fmt.Errorf("initial wallet balance cannot be negative")
	}
	return nil
}
