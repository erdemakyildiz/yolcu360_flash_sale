package request

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type CreateSaleRequest struct {
	ProductID int     `json:"product_id" validate:"required"`
	SaleStock int     `json:"saleStock" validate:"required,gt=1"`
	Discount  float64 `json:"discount" validate:"required,gt=1"`
	StartTime string  `json:"startTime" validate:"required"`
	EndTime   string  `json:"endTime" validate:"required"`
}

func (req *CreateSaleRequest) Validate() error {
	return validate.Struct(req)
}

type UpdateSaleRequest struct {
	ID        int     `json:"id" validate:"required"`
	Discount  float64 `json:"discount"`
	SaleStock int     `json:"saleStock"`
	StartTime string  `json:"startTime"`
	EndTime   string  `json:"endTime"`
	Active    bool    `json:"active"`
}

func (req *UpdateSaleRequest) Validate() error {
	return validate.Struct(req)
}
