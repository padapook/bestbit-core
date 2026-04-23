package middleware

import (
	// "net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/padapook/bestbit-core/internal/utils"
	"github.com/padapook/bestbit-core/internal/utils/auth"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.HandleError(c, utils.ErrUnauthorized)
			c.Abort()
			return
		}

		// ตัด "Bearer"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.HandleError(c, utils.ErrUnauthorized)
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			utils.HandleError(c, utils.ErrUnauthorized)
			c.Abort()
			return
		}

		c.Set("account_id", claims.AccountID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
