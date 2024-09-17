package entity

import (
	"errors"
	"flash_sale_management/dto"
	"fmt"
	"time"
)

type Sale struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	ProductID int       `gorm:"type:int;not null"`
	Quantity  int       `gorm:"type:int;not null"`
	Discount  float64   `gorm:"type:decimal(10,2);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoCreateTime"`
	StartTime time.Time `gorm:"type:timestamp;not null"`
	EndTime   time.Time `gorm:"type:timestamp;not null"`
	Active    bool      `gorm:"default:false"`
}

func (sale *Sale) FromDto(request dto.CreateSaleRequest) (*Sale, error) {
	sale.ProductID = request.ProductID
	sale.Quantity = request.Quantity
	sale.Discount = request.Discount

	sTime, err := formatTime(request.StartTime)
	if err != nil {
		return nil, err
	}

	sale.StartTime = *sTime

	eTime, err := formatTime(request.EndTime)
	if err != nil {
		return nil, err
	}

	sale.EndTime = *eTime

	sale.Active = false

	return sale, nil
}

func (sale *Sale) FromUpdateDto(request dto.UpdateSaleRequest) (*Sale, error) {
	if request.Quantity > 0 {
		sale.Quantity = request.Quantity
	}

	if request.Discount > 0 {
		sale.Discount = request.Discount
	}

	if request.StartTime != "" {
		t, err := formatTime(request.StartTime)
		if err != nil {
			return nil, err
		}

		sale.StartTime = *t
	}

	if request.EndTime != "" {
		t, err := formatTime(request.EndTime)
		if err != nil {
			return nil, err
		}

		sale.EndTime = *t
	}

	if request.Active != sale.Active {
		sale.Active = request.Active
	}

	return sale, nil
}

func formatTime(date string) (*time.Time, error) {
	layout := "2006-01-02T15:04"

	parsedTime, err := time.Parse(layout, date)
	if err != nil {
		msg := fmt.Sprintf("Error parsing date: %v", err)
		return nil, errors.New(msg)
	}

	return &parsedTime, nil
}
