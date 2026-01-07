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
