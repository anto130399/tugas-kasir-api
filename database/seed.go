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
	// Load Config
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

	// 1. Insert Category
	var catID int
	fmt.Println("Creating Category 'Umum'...")
	// Try to find existing first
	err = db.QueryRow(ctx, "SELECT id FROM categories WHERE name = 'Umum'").Scan(&catID)
	if err != nil {
		// Insert if not found
		err = db.QueryRow(ctx, "INSERT INTO categories (name) VALUES ('Umum') RETURNING id").Scan(&catID)
		if err != nil {
			log.Fatalf("❌ Failed to insert category: %v", err)
		}
		fmt.Printf("✅ Created Category ID: %d\n", catID)
	} else {
		fmt.Printf("ℹ️  Category 'Umum' already exists (ID: %d)\n", catID)
	}

	// 2. Insert Product
	var prodID int
	fmt.Println("Creating Product 'Barang Contoh'...")
	err = db.QueryRow(ctx, "SELECT id FROM produk WHERE nama = 'Barang Contoh'").Scan(&prodID)
	if err != nil {
		err = db.QueryRow(ctx, "INSERT INTO produk (nama, harga, stok, category_id) VALUES ('Barang Contoh', 25000, 100, $1) RETURNING id", catID).Scan(&prodID)
		if err != nil {
			log.Fatalf("❌ Failed to insert product: %v", err)
		}
		fmt.Printf("✅ Created Product ID: %d\n", prodID)
	} else {
		fmt.Printf("ℹ️  Product 'Barang Contoh' already exists (ID: %d)\n", prodID)
	}
	
	fmt.Println("\n🎉 Database seeded successfully!")
	fmt.Printf("Use \"produk_id\": %d in your checkout request.\n", prodID)
}
