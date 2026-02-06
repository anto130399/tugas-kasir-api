package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"test1/models"
	"test1/services"

	"github.com/gorilla/mux"
)

type ProdukHandler struct {
	service *services.ProdukService
}

func NewProdukHandler(service *services.ProdukService) *ProdukHandler {
	return &ProdukHandler{service: service}
}

// GetAllProduk - GET /api/produk
func (h *ProdukHandler) GetAllProduk(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("nama")
	products, err := h.service.GetAllProduk(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// CreateProduk - POST /api/produk
func (h *ProdukHandler) CreateProduk(w http.ResponseWriter, r *http.Request) {
	var produk models.Produk
	err := json.NewDecoder(r.Body).Decode(&produk)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.CreateProduk(&produk)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(produk)
}

// GetProdukByID - GET /api/produk/{id}
func (h *ProdukHandler) GetProdukByID(w http.ResponseWriter, r *http.Request) {
	// Menggunakan mux vars untuk mengambil param ID jika tersedia, atau fallback parsing URL
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		// Fallback manual parsing jika tidak pakai mux router (meski main.go pakai mux)
		idStr = strings.TrimPrefix(r.URL.Path, "/api/produk/")
	}
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid produk ID", http.StatusBadRequest)
		return
	}

	product, err := h.service.GetProdukByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// UpdateProduk - PUT /api/produk/{id}
func (h *ProdukHandler) UpdateProduk(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		idStr = strings.TrimPrefix(r.URL.Path, "/api/produk/")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid produk ID", http.StatusBadRequest)
		return
	}

	var produk models.Produk
	err = json.NewDecoder(r.Body).Decode(&produk)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	produk.ID = id
	err = h.service.UpdateProduk(&produk)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(produk)
}

// DeleteProduk - DELETE /api/produk/{id}
func (h *ProdukHandler) DeleteProduk(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		idStr = strings.TrimPrefix(r.URL.Path, "/api/produk/")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid produk ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteProduk(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Produk deleted successfully",
	})
}