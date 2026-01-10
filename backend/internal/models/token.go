package models

import "time"

type Token struct {
	ID            int64     `json:"id" db:"id"`
	ChainID       int64     `json:"chain_id" db:"chain_id"`
	Address       string    `json:"address" db:"address"`
	Type          string    `json:"type" db:"type"`
	Name          *string   `json:"name,omitempty" db:"name"`
	Symbol        *string   `json:"symbol,omitempty" db:"symbol"`
	Decimals      *int      `json:"decimals,omitempty" db:"decimals"`
	TotalSupply   *string   `json:"total_supply,omitempty" db:"total_supply"`
	HolderCount   int64     `json:"holder_count" db:"holder_count"`
	TransferCount int64     `json:"transfer_count" db:"transfer_count"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type TokenTransfer struct {
	ID              int64     `json:"id" db:"id"`
	ChainID         int64     `json:"chain_id" db:"chain_id"`
	TransactionHash string    `json:"transaction_hash" db:"transaction_hash"`
	LogIndex        int       `json:"log_index" db:"log_index"`
	TokenAddress    string    `json:"token_address" db:"token_address"`
	FromAddress     string    `json:"from_address" db:"from_address"`
	ToAddress       string    `json:"to_address" db:"to_address"`
	Value           *string   `json:"value,omitempty" db:"value"`
	TokenID         *string   `json:"token_id,omitempty" db:"token_id"`
	BlockNumber     int64     `json:"block_number" db:"block_number"`
	Timestamp       time.Time `json:"timestamp" db:"timestamp"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

type TokenBalance struct {
	ID            int64     `json:"id" db:"id"`
	ChainID       int64     `json:"chain_id" db:"chain_id"`
	TokenAddress  string    `json:"token_address" db:"token_address"`
	HolderAddress string    `json:"holder_address" db:"holder_address"`
	Balance       string    `json:"balance" db:"balance"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}