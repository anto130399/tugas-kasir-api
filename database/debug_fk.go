package main

import (
	"context"
	"fmt"
	"log"
	
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func main() {
	// 1. Setup Connection
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.ReadInConfig()

	dbConn := viper.GetString("DB_CONN")
	// Ensure SimpleProtocol
	config, err := pgxpool.ParseConfig(dbConn)
	if err != nil {
		log.Fatal("ParseConfig error:", err)
	}
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	ctx := context.Background()
	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal("Connect error:", err)
	}
	defer db.Close()
	fmt.Println("✅ Connected to DB")

	// 2. Start Transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		log.Fatal("Begin Tx error:", err)
	}
	defer tx.Rollback(ctx)
	fmt.Println("🔄 Transaction Started")

	// 3. Insert Transaction Header
	var transactionID int
	totalAmount := 12345
	fmt.Println("➡️  Inserting into transactions...")
	err = tx.QueryRow(ctx, `INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id`, totalAmount).Scan(&transactionID)
	if err != nil {
		log.Fatal("❌ Insert transactions failed:", err)
	}
	fmt.Printf("✅ Header inserted. Got TransactionID: %d\n", transactionID)

	// 4. Insert Transaction Detail
	// We use hardcoded valid product_id 1 (assuming it exists from previous seed)
	// If product 1 doesn't exist, this might fail on product_id FK, but we are testing transaction_id FK here.
	produkID := 1 
	
	fmt.Printf("➡️  Inserting into transaction_details with TransactionID=%d, ProdukID=%d...\n", transactionID, produkID)
	_, err = tx.Exec(ctx, `
		INSERT INTO transaction_details (transaction_id, produk_id, quantity, subtotal) 
		VALUES ($1, $2, $3, $4)
	`, transactionID, produkID, 1, 10000)
	
	if err != nil {
		log.Printf("❌ Insert detail failed: %v\n", err)
		// Don't fatal yet, let's see why
		return
	}
	fmt.Println("✅ Detail inserted successfully!")

	// 5. Commit
	err = tx.Commit(ctx)
	if err != nil {
		log.Fatal("❌ Commit failed:", err)
	}
	fmt.Println("🎉 Transaction Committed!")
}
