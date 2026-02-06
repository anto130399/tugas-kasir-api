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
	
	// Force simple protocol for Supabase
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

	rows, err := db.Query(ctx, "SELECT column_name, data_type FROM information_schema.columns WHERE table_name = 'transaction_details'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Columns in transaction_details:")
	for rows.Next() {
		var name, dtype string
		rows.Scan(&name, &dtype)
		fmt.Printf("- %s (%s)\n", name, dtype)
	}
}
