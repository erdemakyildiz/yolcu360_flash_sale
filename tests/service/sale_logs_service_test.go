package service

import (
	"flash_sale_management/config"
	"flash_sale_management/entity"
	"flash_sale_management/repository"
	"flash_sale_management/service"
	"flash_sale_management/tests/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var saleLog = &entity.SaleLog{
	ID:                    1,
	ProductID:             2,
	RemainingSaleStock:    10,
	RemainingProductStock: 10,
	Price:                 10,
	CreatedAt:             time.Now(),
}

func init() {
	config.LoadConfig()
}

func Test_CreateSaleLog_when_expect_success(t *testing.T) {
	repo := new(mocks.SaleLogRepository)

	repo.On("Save", saleLog).Return(repository.Result{Result: saleLog})

	saleLogService := service.NewSaleLogService(repo)

	err := saleLogService.SaveSaleLog(saleLog)

	assert.Nil(t, err)
	repo.AssertExpectations(t)
}
