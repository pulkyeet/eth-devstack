package models

import "time"

type TransactionLog struct {
	ID               int64     `json:"id" db:"id"`
	ChainID          int64     `json:"chain_id" db:"chain_id"`
	TransactionHash  string    `json:"transaction_hash" db:"transaction_hash"`
	LogIndex         int       `json:"log_index" db:"log_index"`
	Address          string    `json:"address" db:"address"`
	Data             *string   `json:"data,omitempty" db:"data"`
	Topic0           *string   `json:"topic0,omitempty" db:"topic0"`
	Topic1           *string   `json:"topic1,omitempty" db:"topic1"`
	Topic2           *string   `json:"topic2,omitempty" db:"topic2"`
	Topic3           *string   `json:"topic3,omitempty" db:"topic3"`
	BlockNumber      int64     `json:"block_number" db:"block_number"`
	BlockHash        string    `json:"block_hash" db:"block_hash"`
	TransactionIndex int       `json:"transaction_index" db:"transaction_index"`
	Removed          bool      `json:"removed" db:"removed"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}