package service

import (
	"flash_sale_management/entity"
	"flash_sale_management/repository"
)

type SaleLogService struct {
	saleLogRepository repository.SaleLogRepositoryInterface
}

func NewSaleLogService(repo repository.SaleLogRepositoryInterface) SaleLogService {
	return SaleLogService{saleLogRepository: repo}
}

func (sl *SaleLogService) SaveSaleLog(saleLog *entity.SaleLog) (*entity.SaleLog, error) {
	result := sl.saleLogRepository.Save(saleLog)
	if result.Error != nil {
		return nil, result.Error
	}

	return saleLog, nil
}
