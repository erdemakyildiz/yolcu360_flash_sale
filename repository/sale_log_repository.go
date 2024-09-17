package repository

import (
	"flash_sale_management/entity"
	"gorm.io/gorm"
)

type SaleLogRepository struct {
	db *gorm.DB
}

type SaleLogRepositoryInterface interface {
	Save(sale *entity.SaleLog) Result
}

func NewSaleLogRepository(db *gorm.DB) *SaleLogRepository {
	return &SaleLogRepository{db: db}
}

func (r *SaleLogRepository) Save(sale *entity.SaleLog) Result {
	err := r.db.Create(sale).Error

	if err != nil {
		return Result{Error: err}
	}

	return Result{Result: sale}
}
