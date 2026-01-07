package models

import "time"

type Address struct {
	ID              int64      `json:"id" db:"id"`
	ChainID         int64      `json:"chain_id" db:"chain_id"`
	Address         string     `json:"address" db:"address"`
	Balance         int64      `json:"balance" db:"balance"`
	Nonce           int64      `json:"nonce" db:"nonce"`
	IsContract      bool       `json:"is_contract" db:"is_contract"`
	ContractCreator *string    `json:"contract_creator,omitempty" db:"contract_creator"`
	CreationTxHash  *string    `json:"creation_tx_hash,omitempty" db:"creation_tx_hash"`
	CodeHash        *string    `json:"code_hash,omitempty" db:"code_hash"`
	TxCount         int64      `json:"tx_count" db:"tx_count"`
	FirstSeenBlock  *int64     `json:"first_seen_block,omitempty" db:"first_seen_block"`
	LastSeenBlock   *int64     `json:"last_seen_block,omitempty" db:"last_seen_block"`
	FirstSeenAt     *time.Time `json:"first_seen_at,omitempty" db:"first_seen_at"`
	LastSeenAt      *time.Time `json:"last_seen_at,omitempty" db:"last_seen_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}
