package config

import (
	"fmt"
	"strconv"
	"time"
)

type Config struct {
	// Server configuration
	Port        string
	Host        string
	Environment string

	// Blockchain settings
	NetWork              string
	MiningReward         float64
	BlockTime            int
	MaxBlockSize         int
	InitialWalletBalance float64

	// Database configuration
	Database     string
	DatabaseURL  string
	DatabaseName string

	// Security configuration
	TLSEnabled bool
	CertFile   string
	KeyFile    string
	APIKey     string
	JWTSecret  string

	// Rate limiting configuration
	RateLimit RateLimitConfig

	// Consensus configuration
	Consensus ConsensusConfig

	// Logging configuration
	LogLevel      string
	LogFile       string
	EnableLogging bool

	// CORS configuration
	CORS CORSConfig
}

type RateLimitConfig struct {
	GeneralRPM  int // Requests per minute for general endpoints
	FaucetRPH   int // Requests per hour for faucet
	MiningRPM   int // Requests per minute for mining
	WindowSize  time.Duration
	EnableLimit bool
}

type ConsensusConfig struct {
	DefaultType   string  // "pow" or "pos"
	MinStake      float64 // Minimum stake for PoS
	MaxValidators int     // Maximum number of validators
	SlashRate     float64 // Slashing rate for malicious behavior
	RewardRate    float64 // Reward rate for successful validation
}

type CORSConfig struct {
	Enabled        bool
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	MaxAge         int
}

func LoadConfig() *Config {
	return &Config{}
}

func (c *Config) Validate() error {
	if port, err := strconv.Atoi(c.Port); err != nil || port <= 0 || port > 65535 {
		return fmt.Errorf("invalid port number: %s", c.Port)
	}
	// Validate mining reward
	if c.MiningReward < 0 {
		return fmt.Errorf("mining reward cannot be negative: %f", c.MiningReward)
	}
	// Validate consensus configuration
	if c.Consensus.MinStake <= 0 {
		return fmt.Errorf("minimum stake must be positive")
	}

	if c.Consensus.MaxValidators <= 0 {
		return fmt.Errorf("maximum validators must be positive")
	}
	// Validate TLS configuration
	if c.TLSEnabled {
		if c.CertFile == "" || c.KeyFile == "" {
			return fmt.Errorf("TLS enabled but cert or key file not specified")
		}
	}

	return nil
}
