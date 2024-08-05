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
func (o *OrderService) CreateOrder(ctx *gin.Context, orderID string, user *model.User) error {
	return o.repo.CreateOrder(ctx, orderID, user)
}

// load orders from DB
func (o *OrderService) ListOrders(ctx *gin.Context, user *model.User) (*[]model.Order, error) {
	orders, err := o.repo.ListOrders(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("reading orders error: %w", err)
	}
	return orders, nil
}
