package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Account struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
}

type LoadedAccount struct {
	Address    common.Address
	PrivateKey *ecdsa.PrivateKey
}

var (
	clients  []*ethclient.Client
	accounts []LoadedAccount
	chainID  = big.NewInt(1337)
)

func main() {
	fmt.Println("🚀 Starting transaction generator...")

	// Connect to all signers
	rpcEndpoints := []string{
		"http://signer1:8545",
		"http://signer2:8545",
		"http://signer3:8545",
		"http://signer4:8545",
		"http://signer5:8545",
		"http://signer6:8545",
		"http://signer7:8545",
	}

	clients = make([]*ethclient.Client, len(rpcEndpoints))
	for i, endpoint := range rpcEndpoints {
		client, err := ethclient.Dial(endpoint)
		if err != nil {
			panic(fmt.Sprintf("Failed to connect to %s: %v", endpoint, err))
		}
		clients[i] = client
		fmt.Printf("✅ Connected to %s\n", endpoint)
	}

	// Load accounts
	accountsData, err := os.ReadFile("/data/accounts.json")
	if err != nil {
		panic(err)
	}

	var accountsJSON []Account
	err = json.Unmarshal(accountsData, &accountsJSON)
	if err != nil {
		panic(err)
	}

	accounts = make([]LoadedAccount, len(accountsJSON))
	for i, acc := range accountsJSON {
		privateKey, err := crypto.HexToECDSA(acc.PrivateKey[2:])
		if err != nil {
			panic(err)
		}
		accounts[i] = LoadedAccount{
			Address:    common.HexToAddress(acc.Address),
			PrivateKey: privateKey,
		}
	}

	fmt.Printf("✅ Loaded %d accounts\n", len(accounts))

	// Wait for network
	fmt.Println("⏳ Waiting for network...")
	for {
		blockNumber, err := clients[0].BlockNumber(context.Background())
		if err == nil && blockNumber > 0 {
			fmt.Printf("✅ Network ready at block %d\n\n", blockNumber)
			break
		}
		time.Sleep(2 * time.Second)
	}

	fmt.Println("💸 Generating 40-120 txs every 2 seconds")
	fmt.Println("   Distributed round-robin across all signers\n")

	generateTransactions()
}

func generateTransactions() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		txCount := rand.Intn(81) + 40 // 40-120 txs
		fmt.Printf("📦 Batch: %d transactions\n", txCount)

		sent := 0
		for i := 0; i < txCount; i++ {
			// Random accounts
			fromIdx := rand.Intn(len(accounts))
			toIdx := rand.Intn(len(accounts))
			for toIdx == fromIdx {
				toIdx = rand.Intn(len(accounts))
			}

			from := accounts[fromIdx]
			to := accounts[toIdx]

			// Random amount
			amountETH := 0.001 + rand.Float64()*999.999
			amountWei := ethToWei(amountETH)

			// Get fresh nonce from chain
			nonce, err := clients[0].PendingNonceAt(context.Background(), from.Address)
			if err != nil {
				continue
			}

			// Gas price
			gasPrice, err := clients[0].SuggestGasPrice(context.Background())
			if err != nil {
				continue
			}
			
			priorityFeeGwei := int64(1 + rand.Intn(10000))
			priorityFeeWei := new(big.Int).Mul(big.NewInt(priorityFeeGwei), big.NewInt(1e9))
			totalGasPrice := new(big.Int).Add(gasPrice, priorityFeeWei)

			// Create and sign
			tx := types.NewTransaction(nonce, to.Address, amountWei, 21000, totalGasPrice, nil)
			signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), from.PrivateKey)
			if err != nil {
				continue
			}

			// Send to random signer
			targetClient := clients[rand.Intn(len(clients))]
			err = targetClient.SendTransaction(context.Background(), signedTx)
			if err == nil {
				sent++
			}
		}

		fmt.Printf("   ✅ Sent %d/%d\n", sent, txCount)
	}
}

func ethToWei(eth float64) *big.Int {
	weiFloat := eth * 1e18
	wei := new(big.Int)
	wei.SetString(fmt.Sprintf("%.0f", weiFloat), 10)
	return wei
}