package service

import (
	"errors"
	"flash_sale_management/config"
	"flash_sale_management/entity"
	"flash_sale_management/repository"
	"flash_sale_management/service"
	"flash_sale_management/tests/mocks"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var product = &entity.Product{
	ID:        2,
	Name:      "Test product",
	Price:     10,
	Stock:     20,
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

func init() {
	config.LoadConfig()
}

func Test_CreateProduct_when_expect_success(t *testing.T) {
	repo := new(mocks.ProductRepository)
	redisService := new(mocks.RedisService)

	repo.On("Save", product).Return(repository.Result{Result: product})

	productService := service.NewProductService(repo, redisService)

	result := productService.CreateProduct(*product)

	assert.NotNil(t, result)
	assert.Equal(t, result.ID, result.ID)
	repo.AssertExpectations(t)
}

func Test_GetProduct_when_return_product(t *testing.T) {
	repo := new(mocks.ProductRepository)
	redisService := new(mocks.RedisService)

	repo.On("FindOneById", product.ID).Return(repository.Result{Result: product})
	redisService.On("Get", fmt.Sprintf(service.ProductKey, product.ID)).Return(nil, errors.New("error"))

	productService := service.NewProductService(repo, redisService)

	pr, err := productService.GetProduct(product.ID)
	if err != nil {
		return
	}

	assert.NotNil(t, pr)
	assert.Equal(t, pr.ID, product.ID)
	repo.AssertExpectations(t)
}
