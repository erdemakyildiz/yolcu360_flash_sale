package repository

import (
	"flash_sale_management/entity"
	"flash_sale_management/repository"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var sale = entity.Sale{
	ID:        1,
	ProductID: 1,
	SaleStock: 30.0,
	Discount:  10,
	CreatedAt: time.Time{},
	UpdatedAt: time.Time{},
	StartTime: time.Time{},
	EndTime:   time.Time{},
	Active:    false,
}

func Test_when_requestFindSales_expect_returnAllSales(t *testing.T) {
	db, mock, err := CreateMock()
	if err != nil {
		t.Fatalf("failed to create mocks: %s", err)
	}

	salesRepository := repository.NewSaleRepository(db)

	rows := sqlmock.NewRows([]string{"id", "product_id", "sale_stock", "created_at", "start_time", "end_time", "active"}).
		AddRow(sale.ID, sale.ProductID, sale.SaleStock, sale.CreatedAt, sale.StartTime, sale.EndTime, sale.Active)

	mock.ExpectQuery(`^SELECT \* FROM "sales"`).
		WillReturnRows(rows)

	sales := salesRepository.FindAll()
	var data = sales.Result.(*entity.Sale)

	assert.NotEmpty(t, data)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_when_requestFindSale_expect_returnOneSale(t *testing.T) {
	db, mock, err := CreateMock()
	if err != nil {
		t.Fatalf("failed to create mocks: %s", err)
	}

	salesRepository := repository.NewSaleRepository(db)

	rows := sqlmock.NewRows([]string{"id", "product_id", "sale_stock", "created_at", "start_time", "end_time", "active"}).
		AddRow(sale.ID, sale.ProductID, sale.SaleStock, sale.CreatedAt, sale.StartTime, sale.EndTime, sale.Active)

	mock.ExpectQuery(`^SELECT \* FROM "sales"`).
		WithArgs(sale.ID, 1).
		WillReturnRows(rows)

	sales := salesRepository.FindOneById(sale.ID)
	var data = sales.Result.(*entity.Sale)

	assert.NotEmpty(t, data)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_when_requestFindSaleByProductID_expect_returnOneSale(t *testing.T) {
	db, mock, err := CreateMock()
	if err != nil {
		t.Fatalf("failed to create mocks: %s", err)
	}

	salesRepository := repository.NewSaleRepository(db)

	rows := sqlmock.NewRows([]string{"id", "product_id", "sale_stock", "created_at", "start_time", "end_time", "active"}).
		AddRow(sale.ID, sale.ProductID, sale.SaleStock, sale.CreatedAt, sale.StartTime, sale.EndTime, sale.Active)

	mock.ExpectQuery(`^SELECT \* FROM "sales"`).
		WithArgs(sale.ProductID, 1).
		WillReturnRows(rows)

	sales := salesRepository.FindOneByProduct(sale.ProductID)
	var data = sales.Result.(*entity.Sale)

	assert.NotEmpty(t, data)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_Save_when_validSale_expect_success(t *testing.T) {
	db, mock, err := CreateMock()
	if err != nil {
		t.Fatalf("failed to create mocks: %s", err)
	}

	repo := repository.NewSaleRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(`^INSERT INTO "sales" (.+) VALUES (.+) RETURNING "id"`).
		WithArgs(sale.ProductID, sale.SaleStock, sale.Discount, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sale.Active, sale.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(sale.ID))
	mock.ExpectCommit()

	result := repo.Save(&sale)
	data := result.Result.(*entity.Sale)

	assert.NoError(t, result.Error)
	assert.Equal(t, data, result.Result)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_update_when_validSale_expect_success(t *testing.T) {
	db, mock, err := CreateMock()
	if err != nil {
		t.Fatalf("failed to create mocks: %s", err)
	}

	repo := repository.NewSaleRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`^UPDATE "sales" SET (.+) WHERE "id" = ?`).
		WithArgs(sale.ProductID, sale.SaleStock, sale.Discount, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sale.Active, sale.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	result := repo.Update(&sale)
	data := result.Result.(*entity.Sale)

	assert.NoError(t, result.Error)
	assert.Equal(t, data, result.Result)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_when_deleteSaleById_expect_success(t *testing.T) {
	db, mock, err := CreateMock()
	if err != nil {
		t.Fatalf("failed to create mocks: %s", err)
	}

	saleRepository := repository.NewSaleRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`^DELETE FROM "sales" WHERE "sales"."id" = ?`).
		WithArgs(sale.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	deleteResult := saleRepository.DeleteOneById(sale.ID)

	assert.Nil(t, deleteResult.Result)
	assert.NoError(t, deleteResult.Error)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
