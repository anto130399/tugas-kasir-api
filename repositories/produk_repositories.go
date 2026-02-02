package repositories

import (
	"context"
	"fmt"
	"test1/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProdukRepository struct {
	db *pgxpool.Pool
}

func NewProdukRepository(db *pgxpool.Pool) *ProdukRepository {
	return &ProdukRepository{db: db}
}

func (r *ProdukRepository) GetAll() ([]models.Produk, error) {
	ctx := context.Background()
	rows, err := r.db.Query(ctx, `
		SELECT p.id, p.nama, p.harga, p.stok, p.category_id,
		       c.id, c.name AS nama
		FROM produk p
		JOIN categories c ON c.id = p.category_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var produks []models.Produk
	for rows.Next() {
		var p models.Produk
		if err := rows.Scan(&p.ID, &p.Nama, &p.Harga, &p.Stok, &p.CategoryID, &p.Category.ID, &p.Category.Nama); err != nil {
			return nil, err
		}
		produks = append(produks, p)
	}
	return produks, nil
}

func (r *ProdukRepository) GetByID(id int) (*models.Produk, error) {
	ctx := context.Background()
	var p models.Produk
	err := r.db.QueryRow(ctx, `
		SELECT p.id, p.nama, p.harga, p.stok, p.category_id,
		       c.id, c.name AS nama
		FROM produk p
		JOIN categories c ON c.id = p.category_id
		WHERE p.id=$1
	`, id).Scan(&p.ID, &p.Nama, &p.Harga, &p.Stok, &p.CategoryID, &p.Category.ID, &p.Category.Nama)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProdukRepository) Create(p *models.Produk) error {
	ctx := context.Background()
	return r.db.QueryRow(ctx, `
		INSERT INTO produk (nama, harga, stok, category_id)
		VALUES ($1,$2,$3,$4)
		RETURNING id
	`, p.Nama, p.Harga, p.Stok, p.CategoryID).Scan(&p.ID)
}

func (r *ProdukRepository) Update(p *models.Produk) error {
	ctx := context.Background()
	res, err := r.db.Exec(ctx, `
		UPDATE produk SET nama=$1, harga=$2, stok=$3, category_id=$4 WHERE id=$5
	`, p.Nama, p.Harga, p.Stok, p.CategoryID, p.ID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("produk dengan id %d tidak ditemukan", p.ID)
	}
	return nil
}

func (r *ProdukRepository) Delete(id int) error {
	ctx := context.Background()
	res, err := r.db.Exec(ctx, `DELETE FROM produk WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("produk dengan id %d tidak ditemukan", id)
	}
	return nil
}

func (r *ProdukRepository) GetAllCategory() ([]models.Category, error) {
	ctx := context.Background()
	rows, err := r.db.Query(ctx, `SELECT id, name AS nama FROM categories ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Nama); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}
