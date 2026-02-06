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
	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found")
	}

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

	// 1. Check Columns
	fmt.Println("=== COLUMNS in transaction_details ===")
	rows, err := db.Query(ctx, `
		SELECT column_name, data_type 
		FROM information_schema.columns 
		WHERE table_name = 'transaction_details'
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name, dtype string
		rows.Scan(&name, &dtype)
		fmt.Printf("- %s (%s)\n", name, dtype)
	}
	rows.Close()

	// 2. Check Constraints
	fmt.Println("\n=== CONSTRAINTS in transaction_details ===")
	// This query lists constraints and the table they reference
	q := `
		SELECT
			tc.constraint_name, 
			kcu.column_name, 
			ccu.table_name AS foreign_table_name,
			ccu.column_name AS foreign_column_name 
		FROM 
			information_schema.table_constraints AS tc 
			JOIN information_schema.key_column_usage AS kcu
			  ON tc.constraint_name = kcu.constraint_name
			  AND tc.table_schema = kcu.table_schema
			JOIN information_schema.constraint_column_usage AS ccu
			  ON ccu.constraint_name = tc.constraint_name
			  AND ccu.table_schema = tc.table_schema
		WHERE tc.table_name = 'transaction_details' AND tc.constraint_type = 'FOREIGN KEY';
	`
	rows2, err := db.Query(ctx, q)
	if err != nil {
		log.Printf("Error checking constraints: %v", err)
	} else {
		defer rows2.Close()
		for rows2.Next() {
			var cname, col, ftable, fcol string
			rows2.Scan(&cname, &col, &ftable, &fcol)
			fmt.Printf("- Constraint: %s\n  Column: %s -> References: %s(%s)\n", cname, col, ftable, fcol)
		}
	}
}
