package repositories

import (
	"context"
	"fmt"
	"test1/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	ctx := context.Background()

	// Mulai database transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// inisialisasi subtotal -> jumlah total transaksi keseluruhan
	totalAmount := 0
	// inisialisasi modeling transactionDetails -> nanti kita insert ke db
	details := make([]models.TransactionDetail, 0)
	
	// loop setiap item
	for _, item := range items {
		var namaProduk string
		var produkID, price, stock int
		
		// get product dapet pricing from table 'produk'
		// Columns: id, nama, harga, stok
		err := tx.QueryRow(ctx, `SELECT id, nama, harga, stok FROM produk WHERE id=$1`, item.ProdukID).Scan(&produkID, &namaProduk, &price, &stock)
		
		// Handle product not found properly
		if err != nil {
			return nil, fmt.Errorf("produk id %d not found or error: %v", item.ProdukID, err)
		}

		if stock < item.Quantity {
			return nil, fmt.Errorf("stok tidak cukup untuk produk %s (sisa: %d, diminta: %d)", namaProduk, stock, item.Quantity)
		}

		// hitung current total = quantity * pricing
		// ditambahin ke dalam subtotal
		subtotal := item.Quantity * price
		totalAmount += subtotal

		// kurangi jumlah stok di tabel 'produk'
		_, err = tx.Exec(ctx, `UPDATE produk SET stok = stok - $1 WHERE id = $2`, item.Quantity, produkID)
		if err != nil {
			return nil, err
		}

		// item nya dimasukkin ke transactionDetails struct sementara
		details = append(details, models.TransactionDetail{
			ProdukID:    produkID,
			NamaProduk:  namaProduk,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// insert transaction header
	var transactionID int
	// NOTE: Ensure 'transactions' table exists and has 'id' serial primary key
	err = tx.QueryRow(ctx, `INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id`, totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// insert transaction details to db
	for i, detail := range details {
		// Assign TransactionID yang baru didapat
		details[i].TransactionID = transactionID
		details[i].ID = 0 // akan diisi db SERIAL, tapi kita tidak scan balik id detailnya disini untuk efisiensi

		_, err := tx.Exec(ctx, `
			INSERT INTO transaction_details (transaction_id, produk_id, quantity, subtotal) 
			VALUES ($1, $2, $3, $4)
		`, transactionID, detail.ProdukID, detail.Quantity, detail.Subtotal)
		
		if err != nil {
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}