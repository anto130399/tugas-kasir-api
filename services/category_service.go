package services

import (
	"test1/models"
	"test1/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAllCategory() ([]models.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoryService) GetCategoryByID(id int) (*models.Category, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) CreateCategory(c *models.Category) error {
	return s.repo.Create(c)
}

func (s *CategoryService) UpdateCategory(c *models.Category) error {
	return s.repo.Update(c)
}

func (s *CategoryService) DeleteCategory(id int) error {
	return s.repo.Delete(id)
}
