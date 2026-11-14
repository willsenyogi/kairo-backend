package service

import (
	"kairo-backend/internal/model"
	"kairo-backend/internal/repository"
)

type ProductService interface {
	GetAllProducts() ([]model.Product, error)
}

type productService struct {
	productRepo repository.ProductRepository
}

func InitProductService(
	productRepo repository.ProductRepository,
) ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

func (s *productService) GetAllProducts() ([]model.Product, error) {
	products, err := s.productRepo.GetAllProducts()
	if err != nil {
		return nil, err
	}

	return products, nil
}
