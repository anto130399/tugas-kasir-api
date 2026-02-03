package repositories

import (
	"context"
	"fmt"
	"test1/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
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

func (r *CategoryRepository) GetByID(id int) (*models.Category, error) {
	ctx := context.Background()
	var c models.Category
	// Query expects 'name' column in DB based on previous code
	err := r.db.QueryRow(ctx, `SELECT id, name AS nama FROM categories WHERE id=$1`, id).Scan(&c.ID, &c.Nama)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) Create(c *models.Category) error {
	ctx := context.Background()
	// Column is likely 'name' based on previous Selects
	return r.db.QueryRow(ctx, `
		INSERT INTO categories (name)
		VALUES ($1)
		RETURNING id
	`, c.Nama).Scan(&c.ID)
}

func (r *CategoryRepository) Update(c *models.Category) error {
	ctx := context.Background()
	res, err := r.db.Exec(ctx, `UPDATE categories SET name=$1 WHERE id=$2`, c.Nama, c.ID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("category dengan id %d tidak ditemukan", c.ID)
	}
	return nil
}

func (r *CategoryRepository) Delete(id int) error {
	ctx := context.Background()
	res, err := r.db.Exec(ctx, `DELETE FROM categories WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("category dengan id %d tidak ditemukan", id)
	}
	return nil
}
