package service

import (
	"encoding/json"
	"flash_sale_management/entity"
	"flash_sale_management/repository"
	"flash_sale_management/utils"
	"fmt"
	"gorm.io/gorm"
	"time"
)

const ProductKey = "KEY_PRODUCT:%d"

type ProductService struct {
	productRepository repository.ProductRepositoryInterface
	redisService      RedisServiceInterface
}

func NewProductService(repo repository.ProductRepositoryInterface, redis RedisServiceInterface) ProductService {
	return ProductService{productRepository: repo, redisService: redis}
}

func (ps *ProductService) CreateProduct(product entity.Product) *entity.Product {
	ps.productRepository.Save(&product)

	return &product
}

func (ps *ProductService) UpdateProduct(product entity.Product) error {
	product.UpdatedAt = time.Now()
	result := ps.productRepository.Update(&product)

	if result.Error != nil {
		utils.CreateLogMessage("error updating product from db", result.Error)
		return result.Error
	}

	if err := ps.redisService.Set(fmt.Sprintf(ProductKey, product.ID), product); err != nil {
		utils.CreateLogMessage("error setting product to redis", err)
		return err
	}

	return nil
}

func (ps *ProductService) GetProduct(id int) (*entity.Product, error) {
	productCache, err := ps.redisService.Get(fmt.Sprintf(ProductKey, id))
	if err == nil {
		var product entity.Product
		if json.Unmarshal([]byte(productCache), &product) == nil {
			return &product, nil
		}
	}

	result := ps.productRepository.FindOneById(id)
	if result.Error != nil {
		utils.CreateLogMessage("error getting product from db", result.Error)
		return nil, result.Error
	}

	data := result.Result.(*entity.Product)
	if err := ps.redisService.Set(fmt.Sprintf(ProductKey, id), data); err != nil {
		utils.CreateLogMessage("error updating product to redis", err)
		return nil, err
	}

	return data, nil
}

func (ps *ProductService) BeginTransaction() *gorm.DB {
	return ps.productRepository.BeginTransaction()
}

func (ps *ProductService) updateProductWithLock(tx *gorm.DB, product *entity.Product) error {
	if err := ps.productRepository.LockAndUpdateProduct(tx, product); err.Error != nil {
		return err.Error
	}

	err := ps.InvalidateProductCache(product.ID)
	if err != nil {
		return err
	}

	return nil
}

func (ps *ProductService) InvalidateProductCache(productID int) error {
	// invalidate product redis key
	if err := ps.redisService.Delete(fmt.Sprintf(ProductKey, productID)); err != nil {
		utils.CreateLogMessage("error deleting redis key", err)
		return err
	}

	return nil
}
