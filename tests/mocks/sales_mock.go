package mocks

import (
	"flash_sale_management/entity"
	"flash_sale_management/repository"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type SaleRepository struct {
	mock.Mock
}

func (m *SaleRepository) Save(sale *entity.Sale) repository.Result {
	args := m.Called(sale)
	return args.Get(0).(repository.Result)
}

func (m *SaleRepository) Update(sale *entity.Sale) repository.Result {
	args := m.Called(sale)
	return args.Get(0).(repository.Result)
}

func (m *SaleRepository) FindAll() repository.Result {
	args := m.Called()
	return args.Get(0).(repository.Result)
}

func (m *SaleRepository) FindOneById(id int) repository.Result {
	args := m.Called(id)
	return args.Get(0).(repository.Result)
}

func (m *SaleRepository) FindOneByProduct(id int) repository.Result {
	args := m.Called(id)
	return args.Get(0).(repository.Result)
}

func (m *SaleRepository) DeleteOneById(id int) repository.Result {
	args := m.Called(id)
	return args.Get(0).(repository.Result)
}

func (m *SaleRepository) LockAndUpdateSale(tx *gorm.DB, sale *entity.Sale) repository.Result {
	args := m.Called(tx, sale)
	return args.Get(0).(repository.Result)
}

func (m *SaleRepository) BeginTransaction() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}
