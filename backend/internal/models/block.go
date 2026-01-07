package models

import "time"

type Block struct {
	ID               int64     `json:"id" db:"id"`
	ChainID          int64     `json:"chain_id" db:"chain_id"`
	BlockNumber      int64     `json:"block_number" db:"block_number"`
	Hash             string    `json:"hash" db:"hash"`
	ParentHash       string    `json:"parent_hash" db:"parent_hash"`
	Nonce            *string   `json:"nonce,omitempty" db:"nonce"`
	Sha3Uncles       *string   `json:"sha3_uncles,omitempty" db:"sha3_uncles"`
	Miner            string    `json:"miner" db:"miner"`
	StateRoot        *string   `json:"state_root,omitempty" db:"state_root"`
	TransactionsRoot *string   `json:"transactions_root,omitempty" db:"transactions_root"`
	ReceiptsRoot     *string   `json:"receipts_root,omitempty" db:"receipts_root"`
	Difficulty       *string   `json:"difficulty,omitempty" db:"difficulty"`
	TotalDifficulty  *string   `json:"total_difficulty,omitempty" db:"total_difficulty"`
	Size             *int64    `json:"size,omitempty" db:"size"`
	GasLimit         int64     `json:"gas_limit" db:"gas_limit"`
	GasUsed          int64     `json:"gas_used" db:"gas_used"`
	Timestamp        time.Time `json:"timestamp" db:"timestamp"`
	ExtraData        *string   `json:"extra_data,omitempty" db:"extra-data"`
	MixHash          *string   `json:"mix_hash,omitempty" db:"mix_hash"`
	BaseFeePerGas    *string   `json:"base_fee_per_gas,omitempty" db:"base_fee_per_gas"`
	TxCount          int       `json:"tx_count" db:"tx_count"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}
