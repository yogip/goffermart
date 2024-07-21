package handlers

import (
	"errors"
	"goffermart/internal/core/model"

	"github.com/gin-gonic/gin"
)

type BaseHandler struct {
}

func (b *BaseHandler) getUser(ctx *gin.Context) (*model.User, error) {
	u, ok := ctx.Get("User")
	if !ok {
		return nil, errors.New("context must has a user")
	}
	user, ok := u.(*model.User)
	if !ok {
		return nil, errors.New("context must contains valid user object")
	}
	return user, nil
}
