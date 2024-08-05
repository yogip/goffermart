package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"goffermart/internal/core/service"
)

type AuthMiddleware struct {
	iam *service.IAM
}

func NewAuthMiddleware(iamService *service.IAM) *AuthMiddleware {
	return &AuthMiddleware{iam: iamService}
}

func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Token required"})
			return
		}

		parts := strings.Split(token, "Bearer ")
		if len(parts) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Token is not valid"})
			return
		}
		token = parts[1]

		user, err := m.iam.ParseToken(token) // todo rename
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Token is not valid"})
			return
		}

		c.Set("User", user)
		c.Next()
	}
}
