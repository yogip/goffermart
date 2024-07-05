package service

import (
	"context"
	"fmt"
	"goffermart/internal/core/config"
	"goffermart/internal/core/model"
	"goffermart/internal/infra/repo"
	"goffermart/internal/logger"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	UserID    int64
	UserLogin string
}

type IAM struct {
	cfg      *config.ServerConfig
	userRepo *repo.UserRepo // todo Implement repo !!!
}

func NewIAMService(userRepo *repo.UserRepo, cfg *config.ServerConfig) *IAM {
	return &IAM{userRepo: userRepo, cfg: cfg}
}

func (iam *IAM) Login(token string) (string, error) {
	// todo validate token and fetch user from DB
	// 200 — пользователь успешно зарегистрирован и аутентифицирован;
	// 400 — неверный формат запроса;
	// 409 — логин уже занят;
	// 500 — внутренняя ошибка сервера.
	return "token", nil
}

// Create new user and auth them (return token)
func (iam *IAM) Register(ctx context.Context, user *model.UserRequest) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hasing password error: %w", err)
	}

	u, err := iam.userRepo.CreateUser(ctx, user.Login, hashedPassword)
	if err != nil {
		return "", err
	}

	token, err := iam.buildToken(u)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (iam *IAM) buildToken(user *model.User) (string, error) {
	logger.Log.Info(fmt.Sprintf("buildToken u: %v", user))
	logger.Log.Info(fmt.Sprintf("buildToken ID: %v", user.ID))
	logger.Log.Info(fmt.Sprintf("buildToken LOGin: %v", user.Login))
	logger.Log.Info(fmt.Sprintf("buildToken TokenTTL: %v", iam.cfg.TokenTTL))
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		TokenClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(iam.cfg.TokenTTL)),
			},
			UserID:    user.ID,
			UserLogin: user.Login,
		},
	)

	logger.Log.Info(fmt.Sprintf("buildToken key: %s", iam.cfg.SecretKey))
	tokenString, err := token.SignedString([]byte(iam.cfg.SecretKey))
	if err != nil {
		return "", fmt.Errorf("build token error: %w", err)
	}

	logger.Log.Info(fmt.Sprintf("buildToken token: %s", tokenString))
	return tokenString, nil
}
