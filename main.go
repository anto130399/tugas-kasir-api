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
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, dbConn)
	if err != nil {
		log.Fatal("❌ Gagal konek DB:", err)
	}
	defer dbpool.Close()

	log.Println("✅ Database terkoneksi")

	// ===== INIT REPO, SERVICE, HANDLER =====
	produkRepo := repositories.NewProdukRepository(dbpool)
	produkService := services.NewProdukService(produkRepo)
	produkHandler := handlers.NewProdukHandler(produkService)

	// ===== ROUTER =====
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.HandleFunc("/api/produk", produkHandler.GetAllProduk).Methods("GET")
	r.HandleFunc("/api/produk/{id}", produkHandler.GetProdukByID).Methods("GET")
	r.HandleFunc("/api/produk", produkHandler.CreateProduk).Methods("POST")
	r.HandleFunc("/api/produk/{id}", produkHandler.UpdateProduk).Methods("PUT")
	r.HandleFunc("/api/produk/{id}", produkHandler.DeleteProduk).Methods("DELETE")
	r.HandleFunc("/api/produk/category", produkHandler.GetAllCategory).Methods("GET")

	// ===== START SERVER =====
	log.Println("🚀 Server running on 0.0.0.0:" + port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, r))
}
