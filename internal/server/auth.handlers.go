package server

import (
	"authentication-jwt/internal/auth"
	"authentication-jwt/internal/models"
	"authentication-jwt/internal/repositories"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userRepository         repositories.UserRepositoryInterface
	refreshTokenRepository repositories.RefreshTokenRepositoryInterface
}

func newAuthHandler(userRepository repositories.UserRepositoryInterface, refreshTokenRepository repositories.RefreshTokenRepositoryInterface) *AuthHandler {
	return &AuthHandler{
		userRepository:         userRepository,
		refreshTokenRepository: refreshTokenRepository,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Username string `json:"username" binding:"required,min=6"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepository.FindByEmail(c, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing user"})
		return
	}

	if user != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	hashPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user, err = models.NewUser(req.Username, req.Email, hashPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.userRepository.Create(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *AuthHandler) Logon(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepository.FindByEmail(c, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	if user == nil || !auth.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	accessToken, err := auth.GenerateAccessToken(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.SetCookie(
		"access_token",
		accessToken,
		15*60, // 15 minutes
		"/",
		"localhost", // domain
		false,       // secure
		true,        // httpOnly
	)

	c.SetCookie(
		"refresh_token",
		refreshToken,
		7*24*60*60, // 7 days
		"/api/auth/refresh",
		"localhost", // domain
		false,       // secure
		true,        // httpOnly
	)

	refreshTokenModel := models.NewRefreshToken(refreshToken, user.ID, time.Now().Add(7*24*time.Hour))

	err = h.refreshTokenRepository.Create(c.Request.Context(), refreshTokenModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not found"})
		return
	}

	refreshTokenModel, err := h.refreshTokenRepository.FindByToken(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve refresh token"})
		return
	}

	if refreshTokenModel == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	newAccessToken, err := auth.GenerateAccessToken(refreshTokenModel.UserID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new access token"})
		return
	}

	newRefreshToken, err := auth.GenerateRefreshToken(refreshTokenModel.UserID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new refresh token"})
		return
	}

	c.SetCookie(
		"access_token",
		newAccessToken,
		15*60, // 15 minutes
		"/",
		"localhost", // domain
		false,       // secure
		true,        // httpOnly
	)

	c.SetCookie(
		"refresh_token",
		newRefreshToken,
		7*24*60*60, // 7 days
		"/api/auth/refresh",
		"localhost", // domain
		false,       // secure
		true,        // httpOnly
	)

	newRefreshTokenModel := models.NewRefreshToken(newRefreshToken, refreshTokenModel.UserID, time.Now().Add(7*24*time.Hour))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store refresh token"})
		return
	}

	err = h.refreshTokenRepository.Create(c.Request.Context(), newRefreshTokenModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tokens refreshed successfully",
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie(
		"access_token",
		"",
		-1,
		"/",
		"localhost", // domain
		false,       // secure
		true,        // httpOnly
	)

	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/api/auth/refresh",
		"localhost", // domain
		false,       // secure
		true,        // httpOnly
	)
}
