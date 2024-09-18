package service

import (
	"flash_sale_management/entity"
	"flash_sale_management/repository"
	"flash_sale_management/utils"
)

type SaleLogService struct {
	saleLogRepository repository.SaleLogRepositoryInterface
}

func NewSaleLogService(repo repository.SaleLogRepositoryInterface) SaleLogService {
	return SaleLogService{saleLogRepository: repo}
}

func (sl *SaleLogService) SaveSaleLog(saleLog *entity.SaleLog) error {
	result := sl.saleLogRepository.Save(saleLog)
	if result.Error != nil {
		utils.CreateLogMessage("error inserting log to db", result.Error)
		return result.Error
	}

	return nil
}
