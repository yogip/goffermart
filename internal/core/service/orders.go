package service

import (
	"errors"
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
func (o *OrderService) CreateOrder(ctx *gin.Context, orderId int64) error {
	u, ok := ctx.Get("User")
	if !ok {
		return errors.New("context must has a user")
	}
	user, ok := u.(*model.User)
	if !ok {
		return errors.New("context must contains valid user object")
	}

	err := o.repo.CreateOrder(ctx, orderId, user)
	return err
}
