package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"time"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

type ChainClient struct {
	config *ChainConfig
	rpcClient *ethclient.Client
	logger *zap.SugaredLogger
	currentRPC string
}

func NewChainClient(config *ChainConfig, logger *zap.Logger) (*ChainClient, error) {
	sugar := logger.Sugar()

	client, err := ethclient.Dial(config.RPCEndpoint)
	if err!=nil {
		sugar.Warnw("Failed to connect to primary RPC", "chain_id", config.ChainID, "endpoint", config.RPCEndpoint, "error", "err")
		for _, backup := range config.BackupRPCEndpoints {
			client, err = ethclient.Dial(backup)
			if err == nil {
				sugar.Infow("Connected to backup RPC",
					"chain_id", config.ChainID,
					"endpoint", backup,
				)
				return &ChainClient{
					config: config,
					rpcClient: client,
					logger: sugar,
					currentRPC: backup,
				}, nil
			}
		}
		return nil, fmt.Errorf("failed to connect to any RPC endpoint")
	}
	return &ChainClient{config:     config,
		rpcClient:  client,
		logger:     sugar,
		currentRPC: config.RPCEndpoint,
	}, nil
}

func (c *ChainClient) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	return c.rpcClient.BlockNumber(ctx)
}

func (c *ChainClient) GetBlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return c.rpcClient.BlockByNumber(ctx, number)
}

func (c *ChainClient) GetBlockByHash(ctx context.Context, hash string) (*types.Block, error) {
	return c.rpcClient.BlockByHash(ctx, common.HexToHash(hash))
}

func (c *ChainClient) GetTransactionReceipt(ctx context.Context, txHash string) (*types.Receipt, error) {
	return c.rpcClient.TransactionReceipt(ctx, common.HexToHash(txHash))
}

func (c *ChainClient) GetBalance(ctx context.Context, address string, blockNumber *big.Int) (*big.Int, error) {
	return c.rpcClient.BalanceAt(ctx, common.HexToAddress(address), blockNumber)
}

func (c *ChainClient) GetCode(ctx context.Context, address string, blockNumber *big.Int) ([]byte, error) {
	return c.rpcClient.CodeAt(ctx, common.HexToAddress(address), blockNumber)
}

func (c *ChainClient) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	return c.rpcClient.EstimateGas(ctx, msg)
}

func (c *ChainClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return c.rpcClient.SuggestGasPrice(ctx)
}

func (c *ChainClient) ChainID() int64 {
	return c.config.ChainID
}

func (c *ChainClient) Config() *ChainConfig {
	return c.config
}

func (c *ChainClient) Close() {
	c.rpcClient.Close()
}

func (c *ChainClient) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := c.rpcClient.BlockNumber(ctx)
	return err
}