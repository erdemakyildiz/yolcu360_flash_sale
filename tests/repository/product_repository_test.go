package repository

import (
	"flash_sale_management/entity"
	"flash_sale_management/repository"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var product = entity.Product{
	ID:        1,
	Name:      "Test product",
	Price:     10,
	Stock:     20,
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

func Test_when_requestFindProduct_expect_returnOneProduct(t *testing.T) {
	db, mockProduct, err := CreateMock()
	if err != nil {
		t.Fatalf("failed to create mocks: %s", err)
	}

	productRepository := repository.NewProductRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "price", "stock", "created_at", "updated_at"}).
		AddRow(product.ID, product.Name, product.Price, product.Stock, product.CreatedAt, product.UpdatedAt)

	mockProduct.ExpectQuery(`^SELECT \* FROM "products"`).
		WithArgs(product.ID, 1).
		WillReturnRows(rows)

	productResult := productRepository.FindOneById(product.ID)
	var data = productResult.Result.(*entity.Product)

	assert.NotEmpty(t, data)
	assert.NoError(t, err)

	if err := mockProduct.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_when_saveProduct_expect_returnSuccess(t *testing.T) {
	db, mockProduct, err := CreateMock()
	if err != nil {
		t.Fatalf("failed to create mocks: %s", err)
	}

	productRepository := repository.NewProductRepository(db)

	mockProduct.ExpectBegin()
	mockProduct.ExpectQuery(`^INSERT INTO "products" (.+) VALUES (.+) RETURNING "id"`).
		WithArgs(product.Name, product.Price, product.Stock, product.CreatedAt, product.UpdatedAt, product.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(product.ID))
	mockProduct.ExpectCommit()

	productResult := productRepository.Save(&product)
	data := productResult.Result.(*entity.Product)

	assert.NotEmpty(t, data)
	assert.NoError(t, productResult.Error)

	if err := mockProduct.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_when_updateProduct_expect_returnUpdatedProduct(t *testing.T) {
	db, mockProduct, err := CreateMock()
	if err != nil {
		t.Fatalf("failed to create mocks: %s", err)
	}

	productRepository := repository.NewProductRepository(db)

	mockProduct.ExpectBegin()
	mockProduct.ExpectExec(`^UPDATE "products" SET (.+) WHERE "id" = ?`).
		WithArgs(product.Name, product.Price, product.Stock, sqlmock.AnyArg(), sqlmock.AnyArg(), product.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockProduct.ExpectCommit()

	productResult := productRepository.Update(&product)
	data := productResult.Result.(*entity.Product)

	assert.NotEmpty(t, data)
	assert.Equal(t, product.Name, data.Name)
	assert.NoError(t, productResult.Error)

	if err := mockProduct.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
