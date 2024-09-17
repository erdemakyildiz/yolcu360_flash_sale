package service

import (
	"flash_sale_management/entity"
	"flash_sale_management/repository"
	"time"
)

type ProductService struct {
	productRepository repository.ProductRepositoryInterface
}

func NewProductService(repo repository.ProductRepositoryInterface) ProductService {
	return ProductService{productRepository: repo}
}

func (ps *ProductService) CreateProduct(product entity.Product) *entity.Product {
	ps.productRepository.Save(&product)

	return &product
}

func (ps *ProductService) UpdateProduct(product entity.Product) *entity.Product {
	product.UpdatedAt = time.Now()
	ps.productRepository.Update(&product)

	return &product
}

func (ps *ProductService) GetProduct(id int) (*entity.Product, error) {
	result := ps.productRepository.FindOneById(id)
	if result.Error != nil {
		return nil, result.Error
	}

	data := result.Result.(*entity.Product)
	return data, nil
}
