package services

import (
	"test1/models"
	"test1/repositories"
)

type ProdukService struct {
	repo *repositories.ProdukRepository
}

func NewProdukService(repo *repositories.ProdukRepository) *ProdukService {
	return &ProdukService{repo: repo}
}

func (s *ProdukService) GetAllProduk() ([]models.Produk, error) {
	return s.repo.GetAll()
}

func (s *ProdukService) GetProdukByID(id int) (*models.Produk, error) {
	return s.repo.GetByID(id)
}

func (s *ProdukService) CreateProduk(p *models.Produk) error {
	return s.repo.Create(p)
}

func (s *ProdukService) UpdateProduk(p *models.Produk) error {
	return s.repo.Update(p)
}

func (s *ProdukService) DeleteProduk(id int) error {
	return s.repo.Delete(id)
}


