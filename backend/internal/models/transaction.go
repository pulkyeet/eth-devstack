package models

import "time"

type Transaction struct {
	ID                   int64     `json:"id" db:"id"`
	ChainID              int64     `json:"chain_id" db:"chain_id"`
	Hash                 string    `json:"hash" db:"hash"`
	BlockNumber          int64     `json:"block_number" db:"block_number"`
	BlockHash            string    `json:"block_hash" db:"block_hash"`
	TransactionIndex     int       `json:"transaction_index" db:"transaction_index"`
	FromAddress          string    `json:"from_address" db:"from_address"`
	ToAddress            *string   `json:"to_address,omitempty" db:"to_address"`
	Value                string    `json:"value" db:"value"`
	Gas                  int64     `json:"gas" db:"gas"`
	GasPrice             *string   `json:"gas_price,omitempty" db:"gas_price"`
	MaxfeePerGas         *string   `json:"max_fee_per_gas,omitempty" db:"max_fee_per_gas"`
	MaxPriorityFeePerGas *string   `json:"max_prioority_fee_per_gas,omitempty" db:"max_priority_fee_per_gas"`
	Input                *string   `json:"input,omitempty" db:"input"`
	Nonce                int64     `json:"nonce" db:"nonce"`
	TransactionType      int       `json:"transaction_type" db:"transaction_type"`
	Status               *int      `json:"status,omitempty" db:"status"`
	GasUsed              *int64    `json:"gas_used,omitempty" db:"gas_used"`
	CumulativeGasUsed    *int64    `json:"cumulative_gas_used,omitempty" db:"cumulative_gas_sued"`
	EffectiveGasUsed     *string   `json:"effective_gas_used,omitempty" db:"effective_gas_used"`
	ContractAddress      *string   `json:"contract_address,omitempty" db:"contract_address"`
	LogsBloom            *string   `json:"logs_bloom,omitempty" db:"logs_bloom"`
	Timestamp            time.Time `json:"timestamp" db:"timestamp"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
}
