package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	var direction string
	flag.StringVar(&direction, "direction", "up", "Migration direction: up or down")
	flag.Parse()

	// Direct connection string
	connStr := "host=localhost port=5432 user=eth_user password=eth_pass_dev_only dbname=ethereum_explorer sslmode=disable"

	// Open DB connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Create driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("Failed to create migrate driver:", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/database/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal("Failed to create migrate instance:", err)
	}

	// Run migration
	if direction == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Failed to run down migration:", err)
		}
		log.Println("Migration down completed")
	} else {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Failed to run up migration:", err)
		}
		log.Println("Migration up completed")
	}
}
