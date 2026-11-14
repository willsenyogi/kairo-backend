package repository

import (
	"kairo-backend/internal/model"

	"gorm.io/gorm"
)

type ProductRepository interface {
	GetAllProducts() ([]model.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func InitProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetAllProducts() ([]model.Product, error) {
	var products []model.Product

	if err := r.db.
		Preload("Seller").
		Preload("Category").
		Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}
