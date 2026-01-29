package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Produk struct {
	ID       int    `json:"id"`
	Nama     string `json:"daftar_menu"`
	Harga    int    `json:"harga"`
	Stok     int    `json:"stok"`
	Category string `json:"category"`
}

var produk = []Produk{
	{ID: 1, Nama: "Nasi Goreng", Harga: 15000, Stok: 10, Category: "Makanan"},
	{ID: 2, Nama: "Mie Tektek", Harga: 12000, Stok: 15, Category: "Makanan"},
	{ID: 3, Nama: "Es Teh Manis", Harga: 5000, Stok: 20, Category: "Minuman"},
	{ID: 4, Nama: "Jus Alpukat", Harga: 10000, Stok: 30, Category: "Minuman"},
}

func getProdukByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for _, p := range produk {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "Produk belum ada", http.StatusNotFound)

}

func getAllCategory(w http.ResponseWriter, r *http.Request) {
	categoryMap := make(map[string]bool)
	var categories []string

	for _, p := range produk {
		if p.Category != "" {
			categoryMap[p.Category] = true
		}
	}

	for c := range categoryMap {
		categories = append(categories, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func getProdukByCategory(w http.ResponseWriter, r *http.Request) {
	category := strings.TrimPrefix(r.URL.Path, "/api/category/")
	category = strings.ToLower(category)

	var result []Produk

	for _, p := range produk {
		if strings.ToLower(p.Category) == category {
			result = append(result, p)
		}
	}

	if len(result) == 0 {
		http.Error(w, "Category tidak ditemukan", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// PUT localhost:8090/api/produk/{id}
func updateProduk(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid produk ID", http.StatusBadRequest)
		return
	}

	var UpdateProduk Produk
	err = json.NewDecoder(r.Body).Decode(&UpdateProduk)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	for i := range produk {
		if produk[i].ID == id {
			UpdateProduk.ID = id
			produk[i] = UpdateProduk

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(UpdateProduk)
			return
		}
	}

	http.Error(w, "Produk belum ada", http.StatusNotFound)
}

func deleteProduk(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid produk ID", http.StatusBadRequest)
		return
	}

	for i, p := range produk {
		if p.ID == id {
			produk = append(produk[:i], produk[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "sukses delete",
			})
			return
		}
	}

	http.Error(w, "Produk belum ada", http.StatusNotFound)
}

func main() {

	// GET, PUT, DELETE by ID
	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getProdukByID(w, r)
		} else if r.Method == "PUT" {
			updateProduk(w, r)
		} else if r.Method == "DELETE" {
			deleteProduk(w, r)
		}
	})

	// GET all & POST
	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk)
		} else if r.Method == "POST" {
			var produkBaru Produk
			err := json.NewDecoder(r.Body).Decode(&produkBaru)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			produkBaru.ID = len(produk) + 1
			produk = append(produk, produkBaru)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(produkBaru)
		}
	})

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"massage": "API Running",
		})
	})
	// CATEGORY - GET ALL
	http.HandleFunc("/api/category", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getAllCategory(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// CATEGORY - BY NAME
	http.HandleFunc("/api/category/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getProdukByCategory(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	fmt.Println("server running di localhost:8090")

	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		fmt.Println("gagal running server:", err)
	}
}
