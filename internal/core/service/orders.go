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
func (o *OrderService) CreateOrder(ctx *gin.Context, orderID int64, user *model.User) error {
	return o.repo.CreateOrder(ctx, orderID, user)
}

// load orders from DB
func (o *OrderService) ListOrders(ctx *gin.Context, user *model.User) (*[]model.OrderResp, error) {
	orders, err := o.repo.ListOrders(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("reading orders error: %w", err)
	}
	orderResp := []model.OrderResp{}
	for _, o := range *orders {
		orderResp = append(orderResp, model.OrderResp{
			ID:        fmt.Sprint(o.ID),
			Status:    o.Status,
			Accrual:   o.Accrual,
			CreatedAt: *o.CreatedAt,
		})
	}
	return &orderResp, nil
}
