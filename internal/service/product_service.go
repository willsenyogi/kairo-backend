package service

import (
	"kairo-backend/internal/model"

	"gorm.io/gorm"
)

type ProductService interface {
	GetAllProducts() ([]model.Product, error)
}

type productService struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) ProductService {
	return &productService{db: db}
}

func (s *productService) GetAllProducts() ([]model.Product, error) {
	var products []model.Product
	if err := s.db.Preload("Seller").Preload("Category").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
