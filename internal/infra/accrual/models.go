package accrual

import (
	"goffermart/internal/core/model"
)

type Accrual struct {
	OrderID string            `json:"order"`
	Status  model.OrderStatus `json:"status"`
	Accrual float64           `json:"accrual"`
}
