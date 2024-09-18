package service

import (
	"encoding/json"
	"errors"
	"flash_sale_management/dto/request"
	"flash_sale_management/entity"
	"flash_sale_management/repository"
	"flash_sale_management/utils"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"reflect"
	"time"
)

type SalesService struct {
	saleRepository repository.SaleRepositoryInterface
	productService ProductService
	saleLogService SaleLogService
	redisService   RedisServiceInterface
}

const SalesKey = "KEY_SALES"
const SaleKey = "KEY_SALE:%d"

func NewSalesService(repo repository.SaleRepositoryInterface, productService ProductService, saleLogService SaleLogService, service RedisServiceInterface) SalesService {
	return SalesService{
		saleRepository: repo,
		productService: productService,
		saleLogService: saleLogService,
		redisService:   service,
	}
}

func (ss *SalesService) FindSales() (*[]entity.Sale, error) {
	salesCache, err := ss.redisService.Get(SalesKey)
	if err == nil {
		var sales []entity.Sale
		if json.Unmarshal([]byte(salesCache), &sales) == nil {
			return &sales, nil
		}
	}

	result := ss.saleRepository.FindAll()
	if result.Error != nil {
		utils.CreateLogMessage("error getting all sales from db", result.Error)
		return nil, result.Error
	}

	var salesFromDB []entity.Sale
	if reflect.TypeOf(result.Result).Elem().Kind() == reflect.Slice {
		salesFromDB = *(result.Result.(*[]entity.Sale))
	} else {
		singleData := result.Result.(*entity.Sale)
		salesFromDB = []entity.Sale{*singleData}
	}

	if err := ss.redisService.Set(SalesKey, salesFromDB); err != nil {
		utils.CreateLogMessage("error setting all sales to redis", result.Error)
		return nil, err
	}

	return &salesFromDB, nil
}

func (ss *SalesService) CreateSale(request request.CreateSaleRequest) (*entity.Sale, error) {
	if err := request.Validate(); err != nil {
		utils.CreateLogMessage("body validation error", err)
		return nil, err
	}

	product, err := ss.productService.GetProduct(request.ProductID)
	if err != nil {
		return nil, err
	}

	if product.Stock <= 0 {
		err = errors.New(fmt.Sprintf("product doesn't have stock. id: %d", product.ID))
		utils.CreateLogMessage(err.Error(), err)
		return nil, err
	}

	if ss.saleRepository.FindOneByProduct(request.ProductID).Result != nil {
		err = errors.New(fmt.Sprintf("flash sale already exists for this product: %d", request.ProductID))
		utils.CreateLogMessage(err.Error(), err)
		return nil, err
	}

	sale, err := (&entity.Sale{}).FromDto(request)
	if err != nil {
		return nil, err
	}

	if sale.StartTime.After(sale.EndTime) || sale.EndTime.Before(time.Now()) {
		err = errors.New("incorrect time information")
		utils.CreateLogMessage(err.Error(), err)
		return nil, err
	}

	return sale, nil
}

func (ss *SalesService) SaveSale(sale *entity.Sale) (*entity.Sale, error) {
	result := ss.saleRepository.Save(sale)
	if result.Error != nil {
		utils.CreateLogMessage("create sale error", result.Error)
		return nil, result.Error
	}

	err := ss.InvalidateSalesCache(0)
	if err != nil {
		return nil, err
	}

	return sale, nil
}

func (ss *SalesService) UpdateSale(request request.UpdateSaleRequest) (*entity.Sale, error) {
	if err := request.Validate(); err != nil {
		utils.CreateLogMessage("body validation error", err)
		return nil, err
	}

	sale, err := ss.FindSale(request.ID)
	if err != nil {
		return nil, err
	}

	sale, err = sale.FromUpdateDto(request)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}

	sale, err = ss.Update(sale)
	if err != nil {
		return nil, err
	}

	return sale, nil
}

func (ss *SalesService) FindSale(id int) (*entity.Sale, error) {
	saleCache, err := ss.redisService.Get(fmt.Sprintf(SaleKey, id))
	if err == nil {
		var sale entity.Sale
		if json.Unmarshal([]byte(saleCache), &sale) == nil {
			return &sale, nil
		}
	}

	result := ss.saleRepository.FindOneById(id)
	if result.Error != nil {
		utils.CreateLogMessage("error finding sale", err)
		return nil, result.Error
	}

	data := result.Result.(*entity.Sale)
	if err := ss.redisService.Set(fmt.Sprintf(SaleKey, id), data); err != nil {
		utils.CreateLogMessage("error setting sale to redis", err)
		return nil, err
	}

	return data, nil
}

func (ss *SalesService) Update(sale *entity.Sale) (*entity.Sale, error) {
	sale.UpdatedAt = time.Now()
	result := ss.saleRepository.Update(sale)
	if result.Error != nil {
		utils.CreateLogMessage("error updating sale", result.Error)
		return nil, result.Error
	}

	err := ss.InvalidateSalesCache(sale.ID)
	if err != nil {
		return nil, err
	}

	if err := ss.redisService.Set(fmt.Sprintf(SaleKey, sale.ID), sale); err != nil {
		utils.CreateLogMessage("error updating product to redis", err)
		return nil, err
	}

	return sale, nil
}

func (ss *SalesService) InvalidateSalesCache(saleID int) error {
	// invalidate sales redis key
	if err := ss.redisService.Delete(SalesKey); err != nil {
		utils.CreateLogMessage("error deleting sales redis key", err)
		return err
	}

	if err := ss.redisService.Delete(fmt.Sprintf(SaleKey, saleID)); err != nil {
		utils.CreateLogMessage("error delete sale redis key", err)
		return err
	}

	return nil
}

func (ss *SalesService) DeleteSale(id int) error {
	_, err := ss.FindSale(id)
	if err != nil {
		return err
	}

	result := ss.saleRepository.DeleteOneById(id)
	if result.Error != nil {
		utils.CreateLogMessage("error deleting sales from db", err)
		return result.Error
	}

	err = ss.InvalidateSalesCache(id)
	if err != nil {
		return err
	}

	return nil
}

func (ss *SalesService) Buy(id int, wait int) (*entity.SaleLog, error) {
	sale, product, err := ss.getSalesAndProduct(id)
	if err != nil {
		return nil, err
	}

	// check eligible for sales
	if !sale.Active || product.Stock <= 0 || sale.SaleStock <= 0 || time.Now().After(sale.EndTime) {
		err = errors.New("start sale failed: insufficient product stock, sale stock, or the sale period has ended")
		utils.CreateLogMessage(err.Error(), err)
		return nil, err
	}

	// discounted price
	price := product.Price * (1 - sale.Discount/100)

	// wait for testing (checkout)
	time.Sleep(time.Duration(wait) * time.Second)

	saleTx := ss.saleRepository.BeginTransaction()
	productTx := ss.productService.BeginTransaction()

	if err := ss.buyProduct(product, sale, saleTx, productTx); err != nil {
		saleTx.Rollback()
		productTx.Rollback()

		return nil, err
	}

	// create sale log (order)
	saleLog := entity.SaleLog{
		ProductID:             sale.ProductID,
		RemainingSaleStock:    sale.SaleStock,
		RemainingProductStock: product.Stock,
		Price:                 price,
	}
	err = ss.saleLogService.SaveSaleLog(&saleLog)
	if err != nil {
		utils.CreateLogMessage("error creating order", err)
		saleTx.Rollback()
		productTx.Rollback()

		return nil, err
	}

	saleTx.Commit()
	productTx.Commit()

	if err := ss.redisService.Set(fmt.Sprintf(SaleKey, sale.ID), sale); err != nil {
		utils.CreateLogMessage("error updating product to redis", err)
		return nil, err
	}

	if err := ss.redisService.Set(fmt.Sprintf(ProductKey, product.ID), product); err != nil {
		utils.CreateLogMessage("error updating product to redis", err)
		return nil, err
	}

	return &saleLog, nil
}

func (ss *SalesService) buyProduct(product *entity.Product, sale *entity.Sale, saleTx *gorm.DB, productTx *gorm.DB) error {
	// get cached product and sale
	cachedProduct, err := ss.redisService.Get(fmt.Sprintf(ProductKey, product.ID))
	if err != nil {
		utils.CreateLogMessage("error getting product from redis ", err)
		return err
	}
	cachedSale, err := ss.redisService.Get(fmt.Sprintf(SaleKey, sale.ID))
	if err != nil {
		utils.CreateLogMessage("error getting sale from redis ", err)
		return err
	}

	var productTime entity.Product
	var saleTime entity.Sale
	if err := json.Unmarshal([]byte(cachedProduct), &productTime); err != nil {
		utils.CreateLogMessage("json unmarshall error", err)
		return err
	}
	if err := json.Unmarshal([]byte(cachedSale), &saleTime); err != nil {
		utils.CreateLogMessage("json unmarshall error", err)
		return err
	}

	if !product.UpdatedAt.Equal(productTime.UpdatedAt) {
		product, err = ss.productService.GetProduct(sale.ProductID)
		if err != nil {
			return err
		}
	}

	if !sale.UpdatedAt.Equal(saleTime.UpdatedAt) {
		sale, err = ss.FindSale(sale.ID)
		if err != nil {
			return err
		}
	}

	if sale.Active && product.Stock > 0 && sale.SaleStock > 0 && time.Now().Before(sale.EndTime) {
		product.Stock--
		sale.SaleStock--

		err = ss.productService.updateProductWithLock(productTx, product)
		if err != nil {
			return err
		}

		err = ss.updateSaleWithLock(saleTx, sale)
		if err != nil {
			return err
		}

	} else {
		err = errors.New("purchase failed: insufficient product stock, sale stock, or the sale period has ended")
		utils.CreateLogMessage(err.Error(), err)
		return err
	}

	return nil
}

func (ss *SalesService) getSalesAndProduct(id int) (*entity.Sale, *entity.Product, error) {
	sale, err := ss.FindSale(id)
	if err != nil {
		return nil, nil, err
	}
	product, err := ss.productService.GetProduct(sale.ProductID)
	if err != nil {
		return nil, nil, err
	}
	return sale, product, nil
}

func (ss *SalesService) updateSaleWithLock(tx *gorm.DB, sale *entity.Sale) error {
	if err := ss.saleRepository.LockAndUpdateSale(tx, sale); err.Error != nil {
		return err.Error
	}

	err := ss.InvalidateSalesCache(sale.ID)
	if err != nil {
		return err
	}

	return nil
}
