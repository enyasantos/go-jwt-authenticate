package server

import (
	"authentication-jwt/internal/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepository repositories.UserRepositoryInterface
}

func newUserHandler(userRepository repositories.UserRepositoryInterface) *UserHandler {
	return &UserHandler{
		userRepository: userRepository,
	}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userID, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	user, err := h.userRepository.FindById(c, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found on database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user.ToResponse(),
	})
}
