package service

import (
	"encoding/json"
	"errors"
	"flash_sale_management/dto"
	"flash_sale_management/entity"
	"flash_sale_management/repository"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"reflect"
	"time"
)

type SalesService struct {
	saleRepository repository.SaleRepositoryInterface
	productService ProductService
	saleLogService SaleLogService
	redisService   RedisServiceInterface
}

var SalesKey = "SALES"
var SaleKey = "SALE:%d"
var ProductKey = "PRODUCT:%d"

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
		return nil, err
	}

	return &salesFromDB, nil
}

func (ss *SalesService) CreateSale(request dto.CreateSaleRequest) (*entity.Sale, error) {
	if err := request.Validate(); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	product, err := ss.productService.GetProduct(request.ProductID)
	if err != nil {
		return nil, err
	}

	if product.Stock <= 0 {
		msg := fmt.Sprintf("product doesn't have stock. id: %d", product.ID)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}

	if ss.saleRepository.FindOneByProduct(request.ProductID).Result != nil {
		msg := fmt.Sprintf("flash sale already exists for this product: %d", request.ProductID)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	sale, err := (&entity.Sale{}).FromDto(request)
	if err != nil {
		return nil, err
	}

	if sale.StartTime.After(sale.EndTime) || sale.EndTime.Before(time.Now()) {
		msg := "incorrect time information"
		log.Error(msg)
		return nil, errors.New(msg)
	}

	return sale, nil
}

func (ss *SalesService) SaveSale(sale *entity.Sale) (*entity.Sale, error) {
	result := ss.saleRepository.Save(sale)
	if result.Error != nil {
		return nil, result.Error
	}

	if err := ss.redisService.Delete(SalesKey); err != nil {
		return nil, err
	}

	if err := ss.redisService.Set(fmt.Sprintf(SaleKey, sale.ID), sale); err != nil {
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
		return nil, result.Error
	}

	data := result.Result.(*entity.Sale)
	if err := ss.redisService.Set(fmt.Sprintf(SaleKey, id), data); err != nil {
		return nil, err
	}

	return data, nil
}

func (ss *SalesService) UpdateSale(request dto.UpdateSaleRequest) (*entity.Sale, error) {
	if err := request.Validate(); err != nil {
		return nil, nil
	}

	sale, err := ss.FindSale(request.ID)
	if err != nil {
		return nil, err
	}

	sale, err = sale.FromUpdateDto(request)
	if err != nil {
		return nil, err
	}

	sale, err = ss.Update(sale)
	if err != nil {
		return nil, err
	}

	if err := ss.redisService.Set(fmt.Sprintf(SaleKey, sale.ID), sale); err != nil {
		return nil, err
	}

	return sale, nil
}

func (ss *SalesService) Update(sale *entity.Sale) (*entity.Sale, error) {
	sale.UpdatedAt = time.Now()
	result := ss.saleRepository.Update(sale)
	if result.Error != nil {
		return nil, result.Error
	}

	if err := ss.redisService.Set(fmt.Sprintf(SaleKey, sale.ID), sale); err != nil {
		return nil, err
	}
	if err := ss.redisService.Delete(SalesKey); err != nil {
		return nil, err
	}

	return sale, nil
}

func (ss *SalesService) DeleteSale(id int) (*entity.Sale, error) {
	sale, err := ss.FindSale(id)
	if err != nil {
		return nil, err
	}

	result := ss.saleRepository.DeleteOneById(id)
	if result.Error != nil {
		return nil, result.Error
	}

	if err := ss.redisService.Delete(fmt.Sprintf(SaleKey, id)); err != nil {
		return nil, err
	}

	if err := ss.redisService.Delete(SalesKey); err != nil {
		return nil, err
	}

	return sale, nil
}

func (ss *SalesService) Buy(id int, wait int) (*entity.SaleLog, error) {
	sale, product, err := ss.getSalesAndProduct(id)
	if err != nil {
		return nil, err
	}

	if !sale.Active || product.Stock <= 0 || sale.Quantity <= 0 || time.Now().After(sale.EndTime) {
		return nil, errors.New("purchase failed: insufficient product stock, sale quantity, or the sale period has ended")
	}

	if err := ss.redisService.Set(fmt.Sprintf(ProductKey, product.ID), product); err != nil {
		return nil, err
	}
	if err := ss.redisService.Set(fmt.Sprintf(SaleKey, sale.ID), sale); err != nil {
		return nil, err
	}

	price := product.Price * (1 - sale.Discount/100)

	time.Sleep(1 * time.Second)
	time.Sleep(time.Duration(wait) * time.Second)

	if err := ss.buyProduct(product, sale); err != nil {
		return nil, err
	}

	saleLog := entity.SaleLog{
		ProductID: sale.ProductID,
		Quantity:  1,
		Price:     price,
	}
	ss.saleLogService.SaveSaleLog(&saleLog)

	return &saleLog, nil
}

func (ss *SalesService) buyProduct(product *entity.Product, sale *entity.Sale) error {
	cachedProduct, err := ss.redisService.Get(fmt.Sprintf(ProductKey, product.ID))
	if err != nil {
		return err
	}
	cachedSale, err := ss.redisService.Get(fmt.Sprintf(SaleKey, sale.ID))
	if err != nil {
		return err
	}

	var productTime entity.Product
	var saleTime entity.Sale
	if err := json.Unmarshal([]byte(cachedProduct), &productTime); err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(cachedSale), &saleTime); err != nil {
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

	if sale.Active && product.Stock > 0 && sale.Quantity > 0 && time.Now().Before(sale.EndTime) {
		product.Stock--
		sale.Quantity--

		ss.productService.UpdateProduct(*product)
		if _, err := ss.Update(sale); err != nil {
			return err
		}
	} else {
		return errors.New("purchase failed: insufficient product stock, sale quantity, or the sale period has ended")
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
