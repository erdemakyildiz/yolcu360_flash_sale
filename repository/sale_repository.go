package repository

import (
	"flash_sale_management/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SaleRepository struct {
	db *gorm.DB
}

type SaleRepositoryInterface interface {
	Save(sale *entity.Sale) Result
	Update(sale *entity.Sale) Result
	FindAll() Result
	FindOneById(id int) Result
	FindOneByProduct(id int) Result
	DeleteOneById(id int) Result
	LockAndUpdateSale(tx *gorm.DB, updatedSale *entity.Sale) Result
	BeginTransaction() *gorm.DB
}

func NewSaleRepository(db *gorm.DB) *SaleRepository {
	return &SaleRepository{db: db}
}

func (r *SaleRepository) Save(sale *entity.Sale) Result {
	err := r.db.Create(sale).Error

	if err != nil {
		return Result{Error: err}
	}

	return Result{Result: sale}
}

func (r *SaleRepository) Update(sale *entity.Sale) Result {
	err := r.db.Save(sale).Error

	if err != nil {
		return Result{Error: err}
	}

	return Result{Result: sale}
}

func (r *SaleRepository) LockAndUpdateSale(tx *gorm.DB, updatedSale *entity.Sale) Result {
	var sale entity.Sale
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", updatedSale.ID).Take(&sale).Error
	if err != nil {
		tx.Rollback()
		return Result{Error: err}
	}

	err = tx.Save(&updatedSale).Error
	if err != nil {
		tx.Rollback()
		return Result{Error: err}
	}

	return Result{Result: &sale}
}

func (r *SaleRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *SaleRepository) FindAll() Result {
	var sales entity.Sale

	err := r.db.Find(&sales).Error

	if err != nil {
		return Result{Error: err}
	}

	return Result{Result: &sales}
}

func (r *SaleRepository) FindOneById(id int) Result {
	var sale entity.Sale

	err := r.db.Where(&entity.Sale{ID: id}).Take(&sale).Error

	if err != nil {
		return Result{Error: err}
	}

	return Result{Result: &sale}
}

func (r *SaleRepository) FindOneByProduct(id int) Result {
	var sale entity.Sale

	err := r.db.Where(&entity.Sale{ProductID: id}).Take(&sale).Error

	if err != nil {
		return Result{Error: err}
	}

	return Result{Result: &sale}
}

func (r *SaleRepository) DeleteOneById(id int) Result {
	err := r.db.Delete(&entity.Sale{ID: id}).Error

	if err != nil {
		return Result{Error: err}
	}

	return Result{Result: nil}
}
