package models

type Produk struct {
	ID         int      `json:"id"`
	Nama       string   `json:"nama"`
	Harga      int      `json:"harga"`
	Stok       int      `json:"stok"`
	CategoryID int      `json:"category_id"`
	Category   Category `json:"category"`
}
