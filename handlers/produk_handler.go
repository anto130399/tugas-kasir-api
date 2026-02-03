package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"test1/models"
	"test1/services"

	"github.com/gorilla/mux"
)

type ProdukHandler struct {
	service *services.ProdukService
}

func NewProdukHandler(s *services.ProdukService) *ProdukHandler {
	return &ProdukHandler{service: s}
}

// ====================== PRODUK ======================

func (h *ProdukHandler) GetAllProduk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	produks, err := h.service.GetAllProduk()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(produks)
}

func (h *ProdukHandler) GetProdukByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	produk, err := h.service.GetProdukByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if produk == nil {
		http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(produk)
}

func (h *ProdukHandler) CreateProduk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var p models.Produk
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Request body tidak valid: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateProduk(&p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(p)
}

func (h *ProdukHandler) UpdateProduk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	var p models.Produk
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Request body tidak valid: "+err.Error(), http.StatusBadRequest)
		return
	}

	p.ID = id
	if err := h.service.UpdateProduk(&p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(p)
}

func (h *ProdukHandler) DeleteProduk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteProduk(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
