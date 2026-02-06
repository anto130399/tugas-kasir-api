package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"test1/handlers"
	"test1/repositories"
	"test1/services"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func main() {
	// ===== Viper CONFIG =====
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Println("⚠️ Tidak menemukan file .env, pakai ENV variable")
	}

	dbConn := viper.GetString("DB_CONN")
	if dbConn == "" {
		log.Fatal("❌ DB_CONN kosong")
	}

	port := viper.GetString("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("DB Host:", dbConn)
	log.Println("Port:", port)

	// ===== DB CONNECT =====
	// ===== DB CONNECT =====
	ctx := context.Background()
	
	// Parse URL ke Config object
	dbConfig, err := pgxpool.ParseConfig(dbConn)
	if err != nil {
		log.Fatal("❌ Gagal parsing config DB:", err)
	}

	// Force Simple Protocol (Solusi Ampuh untuk Supabase Transacion Pooler)
	dbConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	dbpool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		log.Fatal("❌ Gagal konek DB:", err)
	}
	defer dbpool.Close()

	log.Println("✅ Database terkoneksi")

	// ===== INIT REPO, SERVICE, HANDLER =====
	produkRepo := repositories.NewProdukRepository(dbpool)
	produkService := services.NewProdukService(produkRepo)
	produkHandler := handlers.NewProdukHandler(produkService)

	categoryRepo := repositories.NewCategoryRepository(dbpool)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	transactionRepo := repositories.NewTransactionRepository(dbpool)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// ===== ROUTER =====
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// PRODUK
	r.HandleFunc("/api/produk", produkHandler.GetAllProduk).Methods("GET")
	r.HandleFunc("/api/produk/{id}", produkHandler.GetProdukByID).Methods("GET")
	r.HandleFunc("/api/produk", produkHandler.CreateProduk).Methods("POST")
	r.HandleFunc("/api/produk/{id}", produkHandler.UpdateProduk).Methods("PUT")
	r.HandleFunc("/api/produk/{id}", produkHandler.DeleteProduk).Methods("DELETE")

	// CATEGORY
	r.HandleFunc("/api/categories", categoryHandler.GetAllCategory).Methods("GET")
	r.HandleFunc("/api/categories/{id}", categoryHandler.GetCategoryByID).Methods("GET")
	r.HandleFunc("/api/categories", categoryHandler.CreateCategory).Methods("POST")
	r.HandleFunc("/api/categories/{id}", categoryHandler.UpdateCategory).Methods("PUT")
	r.HandleFunc("/api/categories/{id}", categoryHandler.DeleteCategory).Methods("DELETE")

	// TRANSACTION
	r.HandleFunc("/api/checkout", transactionHandler.HandleCheckout)
	r.HandleFunc("/api/report/hari-ini", transactionHandler.GetDailyReport).Methods("GET")

	// ===== DEBUG ROUTES =====
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		log.Printf("Route: %s Methods: %v\n", pathTemplate, methods)
		return nil
	})

	// ===== START SERVER =====
	log.Println("🚀 Server running on 0.0.0.0:" + port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, r))
}
