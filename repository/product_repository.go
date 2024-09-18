package repository

import (
	"flash_sale_management/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepositoryInterface interface {
	FindOneById(id int) Result
	Save(product *entity.Product) Result
	Update(product *entity.Product) Result
	LockAndUpdateProduct(tx *gorm.DB, product *entity.Product) Result
	BeginTransaction() *gorm.DB
}

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) FindOneById(id int) Result {
	var product entity.Product

	err := r.db.Where(&entity.Product{ID: id}).Take(&product).Error

	if err != nil {
		return Result{Error: err}
	}

	return Result{Result: &product}
}

func (r *ProductRepository) Save(product *entity.Product) Result {
	err := r.db.Create(product).Error

	if err != nil {
		return Result{Error: err}
	}

	return Result{Result: product}
}

func (r *ProductRepository) Update(product *entity.Product) Result {
	err := r.db.Save(product).Error

	if err != nil {
		return Result{Error: err}
	}

	return Result{Result: product}
}

func (r *ProductRepository) LockAndUpdateProduct(tx *gorm.DB, updatedProduct *entity.Product) Result {
	var product entity.Product
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", updatedProduct.ID).Take(&product).Error
	if err != nil {
		tx.Rollback()
		return Result{Error: err}
	}

	err = tx.Save(&updatedProduct).Error
	if err != nil {
		tx.Rollback()
		return Result{Error: err}
	}

	return Result{Result: &product}
}

func (r *ProductRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *ProductRepository) EndTransaction() *gorm.DB {
	return r.db.Commit()
}

func (r *ProductRepository) RollbackTransaction() *gorm.DB {
	return r.db.Rollback()
}
