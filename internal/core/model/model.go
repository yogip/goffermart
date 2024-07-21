package model

import "time"

type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type User struct {
	ID           int64   `json:"id"`
	Login        string  `json:"login"`
	PasswordHash *[]byte `json:"password,omitempty"`
}

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type Order struct {
	ID        int64       `json:"number"`
	Status    OrderStatus `json:"status"`
	Accrual   int32       `json:"accrual"`
	CreatedAt *time.Time  `json:"uploaded_at"`
}

type Balance struct {
	Current   int64 `json:"current"`
	Withdrawn int64 `json:"withdrawn"`
}
