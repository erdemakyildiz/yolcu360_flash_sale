package service

import (
	"errors"
	"flash_sale_management/config"
	"flash_sale_management/dto/request"
	"flash_sale_management/entity"
	"flash_sale_management/repository"
	"flash_sale_management/service"
	"flash_sale_management/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func init() {
	config.LoadConfig()
}

var saleEntity = entity.Sale{
	ID:        1,
	ProductID: saleProduct.ID,
	SaleStock: 20,
	Discount:  10,
	CreatedAt: time.Now(),
	StartTime: time.Now(),
	EndTime:   time.Now(),
	UpdatedAt: time.Now(),
	Active:    true,
}

var saleProduct = &entity.Product{
	ID:        2,
	Name:      "Test product sale",
	Price:     10,
	Stock:     20,
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

func Test_when_getFlashSales_expect_returnSales(t *testing.T) {
	saleRepo := new(mocks.SaleRepository)
	productRepo := new(mocks.ProductRepository)
	saleLogRepo := new(mocks.SaleLogRepository)
	redisService := new(mocks.RedisService)

	saleRepo.On("FindAll").Return(repository.Result{Result: &saleEntity})
	redisService.On("Get", service.SalesKey).Return(nil, errors.New("error"))

	productService := service.NewProductService(productRepo, redisService)
	logService := service.NewSaleLogService(saleLogRepo)
	saleService := service.NewSalesService(saleRepo, productService, logService, redisService)

	sale, err := saleService.FindSales()

	assert.Nil(t, err)
	assert.NotNil(t, sale)
	saleRepo.AssertExpectations(t)
}

func Test_when_getFlashSalesArray_expect_returnSales(t *testing.T) {
	saleRepo := new(mocks.SaleRepository)
	productRepo := new(mocks.ProductRepository)
	saleLogRepo := new(mocks.SaleLogRepository)
	redisService := new(mocks.RedisService)

	saleRepo.On("FindAll").Return(repository.Result{Result: &[]entity.Sale{saleEntity}})
	redisService.On("Get", service.SalesKey).Return(nil, errors.New("error"))

	productService := service.NewProductService(productRepo, redisService)
	logService := service.NewSaleLogService(saleLogRepo)
	saleService := service.NewSalesService(saleRepo, productService, logService, redisService)

	sale, err := saleService.FindSales()

	assert.Nil(t, err)
	assert.NotNil(t, sale)
	saleRepo.AssertExpectations(t)
}

func Test_when_getFlashSale_expect_returnSale(t *testing.T) {
	saleRepo := new(mocks.SaleRepository)
	productRepo := new(mocks.ProductRepository)
	saleLogRepo := new(mocks.SaleLogRepository)
	redisService := new(mocks.RedisService)

	saleRepo.On("FindOneById", saleEntity.ID).Return(repository.Result{Result: &saleEntity})
	redisService.On("Get", mock.Anything).Return(nil, errors.New("error"))

	productService := service.NewProductService(productRepo, redisService)
	logService := service.NewSaleLogService(saleLogRepo)
	saleService := service.NewSalesService(saleRepo, productService, logService, redisService)

	sale, err := saleService.FindSale(saleEntity.ID)

	assert.Nil(t, err)
	assert.NotNil(t, sale)
	assert.Equal(t, sale.ID, saleEntity.ID)
	saleRepo.AssertExpectations(t)
}

func Test_when_updateFlashSale_expect_returnSale(t *testing.T) {
	saleRepo := new(mocks.SaleRepository)
	productRepo := new(mocks.ProductRepository)
	saleLogRepo := new(mocks.SaleLogRepository)
	redisService := new(mocks.RedisService)

	saleRepo.On("FindOneById", saleEntity.ID).Return(repository.Result{Result: &saleEntity})
	saleRepo.On("Update", &saleEntity).Return(repository.Result{Result: saleEntity})
	redisService.On("Get", mock.Anything).Return(nil, errors.New("error"))

	productService := service.NewProductService(productRepo, redisService)
	logService := service.NewSaleLogService(saleLogRepo)
	saleService := service.NewSalesService(saleRepo, productService, logService, redisService)

	saleEntity.Active = false

	request := request.UpdateSaleRequest{
		ID:        1,
		Discount:  40,
		SaleStock: 50,
		StartTime: "2024-09-16T11:04",
		EndTime:   "2024-09-16T10:04",
		Active:    false,
	}

	sale, err := saleService.UpdateSale(request)

	assert.Nil(t, err)
	assert.NotNil(t, sale)
	saleRepo.AssertExpectations(t)
}

func Test_when_createFlashSale_expect_returnErrorNoStock(t *testing.T) {
	saleRepo := new(mocks.SaleRepository)
	productRepo := new(mocks.ProductRepository)
	saleLogRepo := new(mocks.SaleLogRepository)
	redisService := new(mocks.RedisService)

	saleProduct.Stock = 0
	productRepo.On("FindOneById", saleProduct.ID).Return(repository.Result{Result: saleProduct})
	redisService.On("Get", mock.Anything).Return(nil, errors.New("error"))

	productService := service.NewProductService(productRepo, redisService)
	logService := service.NewSaleLogService(saleLogRepo)
	saleService := service.NewSalesService(saleRepo, productService, logService, redisService)

	createSaleRequest := request.CreateSaleRequest{
		ProductID: saleProduct.ID,
		SaleStock: 20,
		Discount:  30,
		StartTime: "2024-09-16T11:04",
		EndTime:   "2024-09-16T12:04",
	}

	_, err := saleService.CreateSale(createSaleRequest)

	assert.NotNil(t, err)
	saleRepo.AssertExpectations(t)
}

func Test_when_createFlashSale_expect_returnAlreadyExist(t *testing.T) {
	saleRepo := new(mocks.SaleRepository)
	productRepo := new(mocks.ProductRepository)
	saleLogRepo := new(mocks.SaleLogRepository)
	redisService := new(mocks.RedisService)

	saleProduct.Stock = 10
	productRepo.On("FindOneById", saleProduct.ID).Return(repository.Result{Result: saleProduct})
	saleRepo.On("FindOneByProduct", saleProduct.ID).Return(repository.Result{Result: saleEntity})
	redisService.On("Get", mock.Anything).Return(nil, errors.New("error"))

	productService := service.NewProductService(productRepo, redisService)
	logService := service.NewSaleLogService(saleLogRepo)
	saleService := service.NewSalesService(saleRepo, productService, logService, redisService)

	createSaleRequest := request.CreateSaleRequest{
		ProductID: saleProduct.ID,
		SaleStock: 20,
		Discount:  30,
		StartTime: "2024-09-16T11:04",
		EndTime:   "2024-09-16T12:04",
	}

	_, err := saleService.CreateSale(createSaleRequest)

	assert.NotNil(t, err)
	saleRepo.AssertExpectations(t)
}

func Test_when_createFlashSale_expect_returnWrongTime(t *testing.T) {
	saleRepo := new(mocks.SaleRepository)
	productRepo := new(mocks.ProductRepository)
	saleLogRepo := new(mocks.SaleLogRepository)
	redisService := new(mocks.RedisService)

	saleProduct.Stock = 10
	productRepo.On("FindOneById", saleProduct.ID).Return(repository.Result{Result: saleProduct})
	saleRepo.On("FindOneByProduct", saleProduct.ID).Return(repository.Result{Result: nil})
	redisService.On("Get", mock.Anything).Return(nil, errors.New("error"))

	productService := service.NewProductService(productRepo, redisService)
	logService := service.NewSaleLogService(saleLogRepo)
	saleService := service.NewSalesService(saleRepo, productService, logService, redisService)

	createSaleRequest := request.CreateSaleRequest{
		ProductID: saleProduct.ID,
		SaleStock: 20,
		Discount:  30,
		StartTime: "2024-09-16T11:04",
		EndTime:   "2024-09-16T10:04",
	}

	_, err := saleService.CreateSale(createSaleRequest)

	assert.NotNil(t, err)
	saleRepo.AssertExpectations(t)
}

func Test_when_buyFlashSale_expect_returnProductNoStock(t *testing.T) {
	saleRepo := new(mocks.SaleRepository)
	productRepo := new(mocks.ProductRepository)
	saleLogRepo := new(mocks.SaleLogRepository)
	redisService := new(mocks.RedisService)

	saleProduct.Stock = 0

	productRepo.On("FindOneById", saleProduct.ID).Return(repository.Result{Result: saleProduct})
	saleRepo.On("FindOneById", saleEntity.ID).Return(repository.Result{Result: &saleEntity})
	redisService.On("Get", mock.Anything).Return(nil, errors.New("error"))

	productService := service.NewProductService(productRepo, redisService)
	logService := service.NewSaleLogService(saleLogRepo)
	saleService := service.NewSalesService(saleRepo, productService, logService, redisService)

	_, err := saleService.Buy(saleEntity.ID, 0)

	assert.NotNil(t, err)

	saleRepo.AssertExpectations(t)
}

func Test_when_buyFlashSale_expect_returnSaleNoStock(t *testing.T) {
	saleRepo := new(mocks.SaleRepository)
	productRepo := new(mocks.ProductRepository)
	saleLogRepo := new(mocks.SaleLogRepository)
	redisService := new(mocks.RedisService)

	saleProduct.Stock = 1
	saleEntity.SaleStock = 0

	productRepo.On("FindOneById", saleProduct.ID).Return(repository.Result{Result: saleProduct})
	saleRepo.On("FindOneById", saleEntity.ID).Return(repository.Result{Result: &saleEntity})
	redisService.On("Get", mock.Anything).Return(nil, errors.New("error"))

	productService := service.NewProductService(productRepo, redisService)
	logService := service.NewSaleLogService(saleLogRepo)
	saleService := service.NewSalesService(saleRepo, productService, logService, redisService)

	_, err := saleService.Buy(saleEntity.ID, 0)

	assert.NotNil(t, err)

	saleRepo.AssertExpectations(t)
}

func Test_when_buyFlashSale_expect_returnIncorrectTime(t *testing.T) {
	saleRepo := new(mocks.SaleRepository)
	productRepo := new(mocks.ProductRepository)
	saleLogRepo := new(mocks.SaleLogRepository)
	redisService := new(mocks.RedisService)

	saleProduct.Stock = 1
	saleEntity.SaleStock = 1
	saleEntity.EndTime = time.Now().Add(-10 * time.Minute)

	productRepo.On("FindOneById", saleProduct.ID).Return(repository.Result{Result: saleProduct})
	saleRepo.On("FindOneById", saleEntity.ID).Return(repository.Result{Result: &saleEntity})
	redisService.On("Get", mock.Anything).Return(nil, errors.New("error"))

	productService := service.NewProductService(productRepo, redisService)
	logService := service.NewSaleLogService(saleLogRepo)
	saleService := service.NewSalesService(saleRepo, productService, logService, redisService)

	_, err := saleService.Buy(saleEntity.ID, 0)

	assert.NotNil(t, err)

	saleRepo.AssertExpectations(t)
}
