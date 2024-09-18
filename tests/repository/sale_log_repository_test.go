package repository

import (
	"flash_sale_management/entity"
	"flash_sale_management/repository"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var saleLog = entity.SaleLog{
	ID:                    1,
	ProductID:             2,
	RemainingSaleStock:    10,
	RemainingProductStock: 10,
	Price:                 100,
	CreatedAt:             time.Now(),
}

func Test_Save_when_expect_success(t *testing.T) {
	db, mock, err := CreateMock()
	if err != nil {
		t.Fatalf("failed to create mocks: %s", err)
	}

	repo := repository.NewSaleLogRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(`^INSERT INTO "sale_logs" (.+) VALUES (.+) RETURNING "id"`).
		WithArgs(saleLog.ProductID, saleLog.RemainingSaleStock, saleLog.RemainingProductStock, saleLog.Price, sqlmock.AnyArg(), saleLog.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(sale.ID))
	mock.ExpectCommit()

	result := repo.Save(&saleLog)
	data := result.Result.(*entity.SaleLog)

	assert.NoError(t, result.Error)
	assert.Equal(t, data, result.Result)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
