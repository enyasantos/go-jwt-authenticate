package middlewares

import (
	"authentication-jwt/internal/auth"
	"authentication-jwt/internal/repositories"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(userRepository repositories.UserRepositoryInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		secret := os.Getenv("JWT_SECRET")
		_, clains, err := auth.ValidateAccessToken(accessToken, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		userID, ok := clains["sub"].(string)
		if !ok || userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		log.Default().Println("Authenticated user ID:", userID)

		user, err := userRepository.FindById(c.Request.Context(), userID)
		if err != nil || user == nil {
			c.JSON(404, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
