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
	Service *services.ProdukService
}

func NewProdukHandler(s *services.ProdukService) *ProdukHandler {
	return &ProdukHandler{Service: s}
}

// ------------------ METHOD HANDLER ------------------

func (h *ProdukHandler) GetAllProduk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	produks, err := h.Service.GetAllProduk()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(produks)
}

func (h *ProdukHandler) GetProdukByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	produk, err := h.Service.GetProdukByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(produk)
}

func (h *ProdukHandler) CreateProduk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var p models.Produk
	json.NewDecoder(r.Body).Decode(&p)
	if err := h.Service.CreateProduk(&p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(p)
}

func (h *ProdukHandler) UpdateProduk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var p models.Produk
	json.NewDecoder(r.Body).Decode(&p)
	p.ID = id
	if err := h.Service.UpdateProduk(&p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(p)
}

func (h *ProdukHandler) DeleteProduk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := h.Service.DeleteProduk(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Produk berhasil dihapus"})
}

func (h *ProdukHandler) GetAllCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	categories, err := h.Service.GetAllCategory()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(categories)
}
