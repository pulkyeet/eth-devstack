package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type DB struct {
	conn   *sql.DB
	logger *zap.SugaredLogger
}

func NewDB(connStr string, maxConns, maxIdleConns int, logger *zap.Logger) (*DB, error) {
	sugar := logger.Sugar()

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("Failed to open database: %w", err)
	}

	conn.SetMaxOpenConns(maxConns)
	conn.SetMaxIdleConns(maxIdleConns)
	conn.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := conn.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("Failed to ping database: %w", err)
	}
	sugar.Infow("Database connected", "max_connections", maxConns, "max_idle_connections", maxIdleConns)

	return &DB{
		conn:   conn,
		logger: sugar,
	}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) Ping(ctx context.Context) error {
	return db.conn.PingContext(ctx)
}

func (db *DB) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return db.conn.BeginTx(ctx, nil)
}

func (db *DB) GetConn() *sql.DB {
	return db.conn
}

type Wallet struct {
	ID string `json:"id"`
	Address              string    `json:"address"`
	EncryptedPrivateKey  string    `json:"-"` 
	EncryptionNonce      string    `json:"-"`
	EncryptionTag        string    `json:"-"`
	EncryptionSalt       string    `json:"-"`
	Name                 string    `json:"name,omitempty"`
	DerivationPath       string    `json:"derivation_path,omitempty"`
	IsImported           bool      `json:"is_imported"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type WalletBalance struct {
	WalletID string `json:"wallet_id"`
	ChainID     int64     `json:"chain_id"`
	ChainName   string    `json:"chain_name"`
	Balance     string    `json:"balance"`
	Symbol      string    `json:"symbol"`
	LastUpdated time.Time `json:"last_updated"`
}

type WalletTransaction struct {
	ID          int64     `json:"id"`
	WalletID    string    `json:"wallet_id"`
	ChainID     int64     `json:"chain_id"`
	TxHash      string    `json:"tx_hash"`
	Direction   string    `json:"direction"`
	FromAddress string    `json:"from_address"`
	ToAddress   *string   `json:"to_address"`
	Value       string    `json:"value"`
	GasUsed     *int64    `json:"gas_used"`
	GasPrice    *string   `json:"gas_price"`
	Status      int       `json:"status"`
	BlockNumber *int64    `json:"block_number"`
	Timestamp   *time.Time `json:"timestamp"`
	CreatedAt   time.Time `json:"created_at"`
}

func (db *DB) CreateWallet(ctx context.Context, w *Wallet) error {
	query:= `INSERT INTO wallets (address, encrypted_private_key, encryption_nonce, 
			encryption_tag, encryption_salt, name, derivation_path, is_imported) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, updated_at`
	return db.conn.QueryRowContext(ctx, query,
		w.Address, w.EncryptedPrivateKey, w.EncryptionNonce,
		w.EncryptionTag, w.EncryptionSalt, w.Name, w.DerivationPath, w.IsImported,
	).Scan(&w.ID, &w.CreatedAt, &w.UpdatedAt)
}

func (db *DB) GetWalletByID(ctx context.Context, id string) (*Wallet, error) {
	query := `
		SELECT id, address, encrypted_private_key, encryption_nonce,
			encryption_tag, encryption_salt, name, derivation_path,
			is_imported, created_at, updated_at
		FROM wallets
		WHERE id = $1
	`
	var w Wallet
	err := db.conn.QueryRowContext(ctx, query, id).Scan(
		&w.ID, &w.Address, &w.EncryptedPrivateKey, &w.EncryptionNonce,
		&w.EncryptionTag, &w.EncryptionSalt, &w.Name, &w.DerivationPath,
		&w.IsImported, &w.CreatedAt, &w.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("wallet not found")
	}
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (db *DB) GetWalletByAddress(ctx context.Context, address string) (*Wallet, error) {
	query := `
		SELECT id, address, encrypted_private_key, encryption_nonce,
			encryption_tag, encryption_salt, name, derivation_path,
			is_imported, created_at, updated_at
		FROM wallets
		WHERE LOWER(address) = LOWER($1)
	`
	var w Wallet
	err := db.conn.QueryRowContext(ctx, query, address).Scan(
		&w.ID, &w.Address, &w.EncryptedPrivateKey, &w.EncryptionNonce,
		&w.EncryptionTag, &w.EncryptionSalt, &w.Name, &w.DerivationPath,
		&w.IsImported, &w.CreatedAt, &w.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("wallet not found")
	}
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (db *DB) GetWalletBalances(ctx context.Context, walletID string) ([]WalletBalance, error) {
	query := `
		SELECT wb.wallet_id, wb.chain_id, c.name, wb.balance, c.native_symbol, wb.last_updated
		FROM wallet_balances wb
		JOIN chains c ON wb.chain_id = c.chain_id
		WHERE wb.wallet_id = $1
		ORDER BY c.chain_id
	`
	rows, err := db.conn.QueryContext(ctx, query, walletID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []WalletBalance
	for rows.Next() {
		var b WalletBalance
		if err := rows.Scan(&b.WalletID, &b.ChainID, &b.ChainName, &b.Balance, &b.Symbol, &b.LastUpdated); err != nil {
			return nil, err
		}
		balances = append(balances, b)
	}
	return balances, rows.Err()
}

func (db *DB) UpsertWalletBalance(ctx context.Context, walletID string, chainID int64, balance string) error {
	query := `
		INSERT INTO wallet_balances (wallet_id, chain_id, balance, last_updated)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (wallet_id, chain_id)
		DO UPDATE SET balance = $3, last_updated = NOW()
	`
	_, err := db.conn.ExecContext(ctx, query, walletID, chainID, balance)
	return err
}

func (db *DB) CreateWalletTransaction(ctx context.Context, tx *WalletTransaction) error {
	query := `
		INSERT INTO wallet_transactions (wallet_id, chain_id, tx_hash, direction,
			from_address, to_address, value, gas_used, gas_price, status,
			block_number, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at
	`
	return db.conn.QueryRowContext(ctx, query,
		tx.WalletID, tx.ChainID, tx.TxHash, tx.Direction,
		tx.FromAddress, tx.ToAddress, tx.Value, tx.GasUsed,
		tx.GasPrice, tx.Status, tx.BlockNumber, tx.Timestamp,
	).Scan(&tx.ID, &tx.CreatedAt)
}

func (db *DB) GetWalletTransactions(ctx context.Context, walletID string, chainID *int64, limit, offset int) ([]WalletTransaction, error) {
	query := `
		SELECT id, wallet_id, chain_id, tx_hash, direction, from_address,
			to_address, value, gas_used, gas_price, status, block_number,
			timestamp, created_at
		FROM wallet_transactions
		WHERE wallet_id = $1
	`
	args := []interface{}{walletID}
	argCount := 1

	if chainID != nil {
		argCount++
		query += fmt.Sprintf(" AND chain_id = $%d", argCount)
		args = append(args, *chainID)
	}

	query += " ORDER BY created_at DESC"
	
	if limit > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, limit)
	}
	
	if offset > 0 {
		argCount++
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, offset)
	}

	rows, err := db.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []WalletTransaction
	for rows.Next() {
		var tx WalletTransaction
		if err := rows.Scan(
			&tx.ID, &tx.WalletID, &tx.ChainID, &tx.TxHash, &tx.Direction,
			&tx.FromAddress, &tx.ToAddress, &tx.Value, &tx.GasUsed,
			&tx.GasPrice, &tx.Status, &tx.BlockNumber, &tx.Timestamp, &tx.CreatedAt,
		); err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}
	return txs, rows.Err()
}