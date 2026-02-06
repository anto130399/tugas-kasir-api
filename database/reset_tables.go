package main

import (
	"context"
	"log"
	
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func main() {
	// 1. Load Config
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Println("⚠️  Warning: .env file not found")
	}

	dbConn := viper.GetString("DB_CONN")
	
	// Config for Supabase
	config, err := pgxpool.ParseConfig(dbConn)
	if err != nil {
		log.Fatal(err)
	}
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	ctx := context.Background()
	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("✅ Connected to Database")

	// 2. Drop Tables
	log.Println("🔥 Dropping old tables...")
	_, err = db.Exec(ctx, "DROP TABLE IF EXISTS transaction_details CASCADE")
	if err != nil {
		log.Fatal("Failed to drop transaction_details:", err)
	}
	_, err = db.Exec(ctx, "DROP TABLE IF EXISTS transactions CASCADE")
	if err != nil {
		log.Fatal("Failed to drop transactions:", err)
	}

	// 3. Recreate Tables
	log.Println("🏗️  Recreating tables...")
	
	// Transactions Header
	_, err = db.Exec(ctx, `
		CREATE TABLE transactions (
			id SERIAL PRIMARY KEY,
			total_amount INT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("Failed to create transactions:", err)
	}
	log.Println("✅ Table 'transactions' created")

	// Transaction Details (Ensure produk_id is used)
	_, err = db.Exec(ctx, `
		CREATE TABLE transaction_details (
			id SERIAL PRIMARY KEY,
			transaction_id INT REFERENCES transactions(id) ON DELETE CASCADE,
			produk_id INT REFERENCES produk(id),
			quantity INT NOT NULL,
			subtotal INT NOT NULL
		)
	`)
	if err != nil {
		log.Fatal("Failed to create transaction_details:", err)
	}
	log.Println("✅ Table 'transaction_details' created (referencing produk.id)")
	
	log.Println("🎉 Database schema reset successfully!")
}
