package service

import (
	"fmt"
	"goffermart/internal/core/config"
	"goffermart/internal/core/model"
	"goffermart/internal/infra/repo"

	"github.com/gin-gonic/gin"
)

type OrderService struct {
	cfg  *config.ServerConfig
	repo *repo.OrderRepo
}

func NewOrdersService(repo *repo.OrderRepo, cfg *config.ServerConfig) *OrderService {
	return &OrderService{repo: repo, cfg: cfg}
}

// create new order in DB with default status
func (o *OrderService) CreateOrder(ctx *gin.Context, orderId int64, user *model.User) error {
	return o.repo.CreateOrder(ctx, orderId, user)
}

// load orders from DB
func (o *OrderService) ListOrders(ctx *gin.Context, user *model.User) (*[]model.Order, error) {
	orders, err := o.repo.ListOrders(ctx, user)
	if err != nil {
		return orders, fmt.Errorf("reading orders error: %w", err)
	}
	return orders, nil
}
