package accrual

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type Accrual struct {
	OrderID int64       `json:"number"`
	Status  OrderStatus `json:"status"`
	Accrual int32       `json:"accrual"`
}
