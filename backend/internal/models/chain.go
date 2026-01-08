package models

import "time"

type Chain struct {
	ChainID          int64     `json:"chain_id" db:"chain_id"`
	Name             string    `json:"name" db:"name"`
	ShortName        string    `json:"short_name" db:"short_name"`
	NativeSymbol     string    `json:"native_symbol" db:"native_symbol"`
	RPCEndpoint      string    `json:"rpc_endpoint" db:"rpc_endpoint"`
	WSEndpoint       *string   `json:"ws_endpoint,omitempty" db:"ws_endpoint"`
	BlockTimeSeconds int       `json:"block_time_seconds" db:"block_time_seconds"`
	IsActive         bool      `json:"is_active" db:"is_active"`
	IsTestnet        bool      `json:"is_testnet" db:"is_testnet"`
	LastIndexedBlock int64     `json:"last_indexed_block" db:"last_indexed_block"`
	IconURL          *string   `json:"icon_url" db:"icon_url"`
	ExplorerURL      *string   `json:"explorer_url" db:"explorer_url"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}
