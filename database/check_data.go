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
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.ReadInConfig()

	dbConn := viper.GetString("DB_CONN")
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

	// Check Columns
	fmt.Println("--- COLUMNS in transaction_details ---")
	rows, err := db.Query(ctx, "SELECT column_name FROM information_schema.columns WHERE table_name = 'transaction_details'")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var name string
		rows.Scan(&name)
		fmt.Println(name)
	}
	rows.Close()

	// Check Produk Data
	fmt.Println("\n--- DATA in produk ---")
	rows2, err := db.Query(ctx, "SELECT id, nama, stok FROM produk")
	if err != nil {
		log.Fatal(err)
	}
	count := 0
	for rows2.Next() {
		var id, stock int
		var nama string
		rows2.Scan(&id, &nama, &stock)
		fmt.Printf("ID: %d, Nama: %s, Stok: %d\n", id, nama, stock)
		count++
	}
	if count == 0 {
		fmt.Println("(Table is empty)")
	}
}
