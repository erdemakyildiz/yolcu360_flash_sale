package entity

import "time"

type SaleLog struct {
	ID                    int       `gorm:"primaryKey;autoIncrement"`
	ProductID             int       `gorm:"type:int;not null"`
	RemainingSaleStock    int       `gorm:"type:int;not null"`
	RemainingProductStock int       `gorm:"type:int;not null"`
	Price                 float64   `gorm:"type:decimal(10,2);not null"`
	CreatedAt             time.Time `gorm:"autoCreateTime"`
}
