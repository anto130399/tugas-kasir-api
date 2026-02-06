
package models

type Transaction struct {
	ID          int                 `json:"id"`
	TotalAmount int                 `json:"total_amount"`
	Details     []TransactionDetail `json:"details"`
}

type TransactionDetail struct {
	ID            int    `json:"id"`
	TransactionID int    `json:"transaction_id"`
	ProdukID      int    `json:"produk_id"`
	NamaProduk    string `json:"nama_produk"`
	Quantity      int    `json:"quantity"`
	Subtotal      int    `json:"subtotal"`
}

type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}

type CheckoutItem struct {
	ProdukID int `json:"produk_id"`
	Quantity int `json:"quantity"`
}
