package mocks

import (
	"flash_sale_management/entity"
	"flash_sale_management/repository"
	"github.com/stretchr/testify/mock"
)

type SaleLogRepository struct {
	mock.Mock
}

func (m *SaleLogRepository) Save(saleLog *entity.SaleLog) repository.Result {
	args := m.Called(saleLog)
	return args.Get(0).(repository.Result)
}
