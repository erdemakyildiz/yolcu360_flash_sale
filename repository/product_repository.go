package repository

import (
	"flash_sale_management/entity"
	"gorm.io/gorm"
)

type ProductRepositoryInterface interface {
	FindOneById(id int) Result
	Save(product *entity.Product) Result
	Update(product *entity.Product) Result
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
