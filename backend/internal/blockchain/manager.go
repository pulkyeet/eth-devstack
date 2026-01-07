package blockchain

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
)

type ChainConfig struct {
	ChainID            int64    `json:"chain_id"`
	Name               string   `json:"name"`
	ShortName          string   `json:"short_name"`
	NativeSymbol       string   `json:"native_symbol"`
	RPCEndpoint        string   `json:"rpc_endpoint"`
	WSEndpoint         string   `json:"ws_endpoint"`
	BlockTimeSeconds   int      `json:"block_time_seconds"`
	IsTestnet          bool     `json:"is_testnet"`
	IsActive           bool     `json:"is_active"`
	SupportsEIP1559    bool     `json:"supports_eip1559"`
	GasPriceOracle     string   `json:"gas_price_oracle"`
	BackupRPCEndpoints []string `json:"backup_rpc_endpoints"`
}

type ChainsFile struct {
	Chains         []ChainConfig `json:"chains"`
	DefaultChainID int64         `json:"default_chain_id"`
}

type ChainManager struct {
	chains  map[int64]*ChainClient
	configs map[int64]*ChainConfig
	logger  *zap.SugaredLogger
	mu      sync.RWMutex
}

func NewChainManager(configPath string, logger *zap.Logger) (*ChainManager, error) {
	sugar := logger.Sugar()
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read chains config: %w", err)
	}
	var chainsFile ChainsFile
	if err := json.Unmarshal(data, &chainsFile); err != nil {
		return nil, fmt.Errorf("failed to parse chains config: %w", err)
	}
	manager := &ChainManager{
		chains:  make(map[int64]*ChainClient),
		configs: make(map[int64]*ChainConfig),
		logger:  sugar,
	}

	for _, config := range chainsFile.Chains {
		manager.configs[config.ChainID] = &config
		if config.IsActive {
			client, err := NewChainClient(&config, logger)
			if err != nil {
				sugar.Warnw("Failed to initialise chain client", "chain_id", config.ChainID, "name", config.Name, "error", err)
				continue
			}
			manager.chains[config.ChainID] = client
			sugar.Infow("Initialised chain client", "chain_id", config.ChainID, "name", config.Name)
		}
	}
	return manager, nil
}

func (m *ChainManager) GetClient(chainID int64) (*ChainClient, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	client, ok := m.chains[chainID]
	if !ok {
		return nil, fmt.Errorf("chain %d not found or not active", chainID)
	}
	return client, nil
}

func (m *ChainManager) GetConfig(chainID int64) (*ChainConfig, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	config, ok := m.configs[chainID]
	if !ok {
		return nil, fmt.Errorf("chain %d not configured", chainID)
	}
	return config, nil
}

func (m *ChainManager) GetActiveChains() []*ChainConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var active []*ChainConfig
	for _, config := range m.configs {
		if config.IsActive {
			active = append(active, config)
		}
	}
	return active
}

func (m *ChainManager) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for chainID, client := range m.chains {
		client.Close()
		m.logger.Infow("Closed chain client", "chain_id", chainID)
	}
}