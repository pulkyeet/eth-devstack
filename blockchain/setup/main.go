package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type Account struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
}

type GenesisAlloc struct {
	Balance string `json:"balance"`
}

func main() {
	fmt.Println("🔑 Generating 128 funded accounts...")

	// Read existing genesis.json
	genesisBytes, err := os.ReadFile("../genesis.json")
	if err != nil {
		panic(err)
	}

	var genesis map[string]interface{}
	if err := json.Unmarshal(genesisBytes, &genesis); err != nil {
		panic(err)
	}

	// Generate 128 funded accounts
	accounts := make([]Account, 0, 128)
	alloc := make(map[string]GenesisAlloc)
	balance := new(big.Int)
	balance.SetString("1000000000000000000000000", 10) // 1M ETH

	for i := 1; i <= 128; i++ {
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			panic(err)
		}

		address := crypto.PubkeyToAddress(privateKey.PublicKey)
		privateKeyHex := hexutil.Encode(crypto.FromECDSA(privateKey))

		accounts = append(accounts, Account{
			Address:    address.Hex(),
			PrivateKey: privateKeyHex,
		})

		alloc[address.Hex()] = GenesisAlloc{
			Balance: balance.String(),
		}

		if i%32 == 0 {
			fmt.Printf("  Generated %d/128 accounts...\n", i)
		}
	}

	// Update genesis alloc
	genesis["alloc"] = alloc

	// Save accounts.json
	accountsJSON, err := json.MarshalIndent(accounts, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("accounts.json", accountsJSON, 0644); err != nil {
		panic(err)
	}

	// Save genesis.json
	genesisJSON, err := json.MarshalIndent(genesis, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("../genesis.json", genesisJSON, 0644); err != nil {
		panic(err)
	}

	fmt.Println("\n✅ Done! 128 funded accounts added to genesis.json")
}