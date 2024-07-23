package accrual

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type Accrual struct {
	OrderID string      `json:"number"`
	Status  OrderStatus `json:"status"`
	Accrual float64     `json:"accrual"`
}
