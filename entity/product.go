package entity

import (
	"time"
)

type Product struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Price     float64   `gorm:"type:decimal(10,2);not null"`
	Stock     int       `gorm:"not null;check:stock >= 0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
