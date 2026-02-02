package main

import (
	"context"
	"log"
	"net/http"

	"test1/handlers"
	"test1/repositories"
	"test1/services"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func main() {
	// ===== LOAD ENV =====
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("⚠️ .env tidak terbaca:", err)
	}

	dbConn := viper.GetString("DB_CONN")
	if dbConn == "" {
		log.Fatal("❌ DB_CONN kosong, cek .env")
	}

	// DEBUG (sementara)
	log.Println("DB HOST =", dbConn)

	// ===== DB CONNECT =====
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, dbConn)
	if err != nil {
		log.Fatal("❌ Gagal konek DB:", err)
	}
	defer dbpool.Close()

	log.Println("✅ Database terkoneksi")

	// ===== DEPENDENCY =====
	produkRepo := repositories.NewProdukRepository(dbpool)
	produkService := services.NewProdukService(produkRepo)
	produkHandler := handlers.NewProdukHandler(produkService)

	// ===== ROUTER =====
	r := mux.NewRouter()
	r.HandleFunc("/api/produk", produkHandler.GetAllProduk).Methods("GET")
	r.HandleFunc("/api/produk/{id}", produkHandler.GetProdukByID).Methods("GET")
	r.HandleFunc("/api/produk", produkHandler.CreateProduk).Methods("POST")
	r.HandleFunc("/api/produk/{id}", produkHandler.UpdateProduk).Methods("PUT")
	r.HandleFunc("/api/produk/{id}", produkHandler.DeleteProduk).Methods("DELETE")
	r.HandleFunc("/api/produk/category", produkHandler.GetAllCategory).Methods("GET")

	port := viper.GetString("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("🚀 Server running at http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
