package services

import (
	"test1/models"
	"test1/repositories"
)

type ProdukService struct {
	Repo *repositories.ProdukRepository
}

func NewProdukService(repo *repositories.ProdukRepository) *ProdukService {
	return &ProdukService{Repo: repo}
}

func (s *ProdukService) GetAllProduk() ([]models.Produk, error)       { return s.Repo.GetAll() }
func (s *ProdukService) GetProdukByID(id int) (*models.Produk, error) { return s.Repo.GetByID(id) }
func (s *ProdukService) CreateProduk(p *models.Produk) error          { return s.Repo.Create(p) }
func (s *ProdukService) UpdateProduk(p *models.Produk) error          { return s.Repo.Update(p) }
func (s *ProdukService) DeleteProduk(id int) error                    { return s.Repo.Delete(id) }
func (s *ProdukService) GetAllCategory() ([]models.Category, error)   { return s.Repo.GetAllCategory() }
