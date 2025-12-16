package middleware

import (
	"go-article/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			response := utils.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil, "Missing or invalid token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := utils.ValidateToken(tokenString)

		if err != nil {
			response := utils.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil, "Invalid token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Set user ID dan role ke context
			c.Set("user_id", claims["user_id"])
		} else {
			response := utils.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil, "Invalid token claims")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Next()
	}
}
