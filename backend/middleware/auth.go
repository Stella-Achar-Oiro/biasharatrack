package middleware

import (
	"os"
	"strings"

	"github.com/OAthooh/BiasharaTrack.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorLogger("Missing authorization header")
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			utils.ErrorLogger("Invalid token format")
			c.JSON(401, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenString := bearerToken[1]

		// Get JWT secret from environment variable
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			utils.ErrorLogger("JWT_SECRET not set in environment")
			c.JSON(500, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		utils.InfoLogger("Token: %v", token)

		if err != nil || !token.Valid {
			utils.ErrorLogger("Invalid token: %v", err)
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userID, ok := claims["user_id"].(float64); ok {
				c.Set("userID", uint(userID))
			}
		}

		c.Next()
	}
}
