package mocks

import (
	"flash_sale_management/entity"
	"flash_sale_management/repository"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type ProductRepository struct {
	mock.Mock
}

func (m *ProductRepository) FindOneById(id int) repository.Result {
	args := m.Called(id)
	return args.Get(0).(repository.Result)
}

func (m *ProductRepository) Save(product *entity.Product) repository.Result {
	args := m.Called(product)
	return args.Get(0).(repository.Result)
}

func (m *ProductRepository) Update(product *entity.Product) repository.Result {
	args := m.Called(product)
	return args.Get(0).(repository.Result)
}

func (m *ProductRepository) LockAndUpdateProduct(tx *gorm.DB, product *entity.Product) repository.Result {
	args := m.Called(tx, product)
	return args.Get(0).(repository.Result)
}

func (m *ProductRepository) BeginTransaction() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}
