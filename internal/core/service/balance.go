package service

import (
	"goffermart/internal/core/config"
	"goffermart/internal/core/model"
	"goffermart/internal/infra/repo"

	"github.com/gin-gonic/gin"
)

type BalanceService struct {
	cfg  *config.ServerConfig
	repo *repo.BalanceRepo
}

func NewBalanceService(repo *repo.BalanceRepo, cfg *config.ServerConfig) *BalanceService {
	return &BalanceService{repo: repo, cfg: cfg}
}

// get user balance
func (o *BalanceService) GetBalance(ctx *gin.Context, user *model.User) (*model.Balance, error) {
	return o.repo.GetBalance(ctx, user)
}

func (o *BalanceService) ProcessWithdrawl(ctx *gin.Context, user *model.User, req *model.WithdrawRequest) error {
	return o.repo.ProcessWithdrawl(ctx, user, req)
}

func (o *BalanceService) ListWithdrawls(ctx *gin.Context, user *model.User) (*[]model.Withdrawal, error) {
	return o.repo.ListWithdrawls(ctx, user)
}
