package consensus

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sort"
	"time"
)

type Validator struct {
	Address       string  `json:"address"`
	StakedAmount  float64 `json:"staked_amount"`
	LastBlockTime int64   `json:"last_block_time"`
	SlashCount    int     `json:"slash_count"`
	IsActive      bool    `json:"is_active"`
	JoinTime      int64   `json:"join_time"`
	TotalRewards  float64 `json:"total_rewards"`
}

type StakingPool struct {
	Validators      map[string]*Validator `json:"validators"`
	MinStakeAmount  float64               `json:"min_stake_amount"`
	MaxValidators   int                   `json:"max_validators"`
	SlashingPenalty float64               `json:"slashing_penalty"`
	BlockReward     float64               `json:"block_reward"`
	StakingReward   float64               `json:"staking_reward"`
}

func NewStakingPool() *StakingPool {
	return &StakingPool{
		Validators:      make(map[string]*Validator),
		MinStakeAmount:  10.0, // Minimum 10 MYC to become validator
		MaxValidators:   100,  // Maximum 100 validators
		SlashingPenalty: 10.0, // 10% penalty for malicious behavior
		BlockReward:     5.0,  // 50 MYC reward for block creator
		StakingReward:   5.0,  // 5% annual staking reward
	}
}

func (sp *StakingPool) AddValidator(address string, stakeAmount float64) error {
	if stakeAmount < sp.MinStakeAmount {
		return fmt.Errorf("minimum stake amount is %.2f MYC", sp.MinStakeAmount)
	}

	if len(sp.Validators) >= sp.MaxValidators {
		return fmt.Errorf("maximum number of validators (%d) reached", sp.MaxValidators)
	}

	if _, exists := sp.Validators[address]; exists {
		return fmt.Errorf("validator already exists")
	}

	validator := &Validator{
		Address:       address,
		StakedAmount:  stakeAmount,
		LastBlockTime: 0,
		SlashCount:    0,
		IsActive:      true,
		JoinTime:      time.Now().Unix(),
		TotalRewards:  0,
	}

	sp.Validators[address] = validator
	return nil
}

func (sp *StakingPool) RemoveValidator(address string) error {
	validator, exists := sp.Validators[address]
	if !exists {
		return fmt.Errorf("validator not found")
	}

	// Return staked amount (minus any penalties)
	penalty := validator.StakedAmount * (float64(validator.SlashCount) * sp.SlashingPenalty / 100.0)
	_ = validator.StakedAmount - penalty // refund amount (handled by blockchain.go)

	delete(sp.Validators, address)
	return nil
}

func (sp *StakingPool) SelectValidator() (string, error) {
	activeValidators := sp.getActiveValidators()
	if len(activeValidators) == 0 {
		return "", fmt.Errorf("no active validators available")
	}

	// Weighted random selection based on stake amount
	totalStake := 0.0
	for _, validator := range activeValidators {
		totalStake += validator.StakedAmount
	}

	// Generate random number between 0 and totalStake
	randomBig, err := rand.Int(rand.Reader, big.NewInt(int64(totalStake*1000000)))
	if err != nil {
		return "", err
	}
	random := float64(randomBig.Int64()) / 1000000.0

	// Select validator based on weighted probability
	currentWeight := 0.0
	for address, validator := range activeValidators {
		currentWeight += validator.StakedAmount
		if random <= currentWeight {
			return address, nil
		}
	}

	// Fallback: return first validator
	for address := range activeValidators {
		return address, nil
	}

	return "", fmt.Errorf("failed to select validator")
}

func (sp *StakingPool) getActiveValidators() map[string]*Validator {
	active := make(map[string]*Validator)
	currentTime := time.Now().Unix()

	for address, validator := range sp.Validators {
		// Validator is active if:
		// 1. IsActive flag is true
		// 2. Not slashed too many times (less than 3)
		// 3. Haven't created block in last 60 seconds (to prevent monopoly)
		if validator.IsActive &&
			validator.SlashCount < 3 &&
			(currentTime-validator.LastBlockTime) > 60 {
			active[address] = validator
		}
	}

	return active
}

func (sp *StakingPool) RewardValidator(address string, blockReward float64) error {
	validator, exists := sp.Validators[address]
	if !exists {
		return fmt.Errorf("validator not found")
	}

	validator.TotalRewards += blockReward
	validator.LastBlockTime = time.Now().Unix()
	return nil
}

func (sp *StakingPool) SlashValidator(address string) error {
	validator, exists := sp.Validators[address]
	if !exists {
		return fmt.Errorf("validator not found")
	}

	validator.SlashCount++
	penalty := validator.StakedAmount * (sp.SlashingPenalty / 100.0)
	validator.StakedAmount -= penalty

	// Deactivate if slashed too many times
	if validator.SlashCount >= 3 {
		validator.IsActive = false
	}

	return nil
}

func (sp *StakingPool) GetValidatorInfo(address string) (*Validator, error) {
	validator, exists := sp.Validators[address]
	if !exists {
		return nil, fmt.Errorf("validator not found")
	}

	return validator, nil
}

func (sp *StakingPool) GetAllValidators() []*Validator {
	validators := make([]*Validator, 0, len(sp.Validators))
	for _, validator := range sp.Validators {
		validators = append(validators, validator)
	}

	// Sort by stake amount (highest first)
	sort.Slice(validators, func(i, j int) bool {
		return validators[i].StakedAmount > validators[j].StakedAmount
	})

	return validators
}

func (sp *StakingPool) GetTotalStaked() float64 {
	total := 0.0
	for _, validator := range sp.Validators {
		if validator.IsActive {
			total += validator.StakedAmount
		}
	}
	return total
}

func (sp *StakingPool) CalculateStakingRewards(address string) (float64, error) {
	validator, exists := sp.Validators[address]
	if !exists {
		return 0, fmt.Errorf("validator not found")
	}

	if !validator.IsActive {
		return 0, nil
	}

	// Calculate annual staking reward (5% per year)
	// Simplified: assume 1 year = 365 * 24 * 60 * 60 seconds
	stakingDuration := time.Now().Unix() - validator.JoinTime
	annualReward := validator.StakedAmount * (sp.StakingReward / 100.0)
	reward := annualReward * (float64(stakingDuration) / (365.0 * 24.0 * 60.0 * 60.0))

	return reward, nil
}
