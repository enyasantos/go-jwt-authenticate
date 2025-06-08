package server

import (
	"authentication-jwt/internal/database"
	"authentication-jwt/internal/middlewares"
	"authentication-jwt/internal/repositories"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	userRepository         repositories.UserRepositoryInterface
	refreshTokenRepository repositories.RefreshTokenRepositoryInterface
}

func NewServer() http.Server {
	port := os.Getenv("PORT")

	db := database.NewDatabase()
	userRepository := repositories.NewUserRepository(db)
	refreshTokenRepository := repositories.NewRefreshTokenRepository(db)

	server := &Server{
		userRepository:         userRepository,
		refreshTokenRepository: refreshTokenRepository,
	}

	return http.Server{
		Addr:         ":" + port,
		Handler:      server.RegisterRoutes(),
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Minute,
	}
}

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	authHandler := newAuthHandler(s.userRepository, s.refreshTokenRepository)
	userHandler := newUserHandler(s.userRepository)

	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("/register", authHandler.Register)

		authRoutes.POST("/logon", authHandler.Logon)

		authRoutes.POST("/refresh", authHandler.Refresh)

		authRoutes.POST("/logout", authHandler.Logout)
	}

	protectedRoutes := r.Group("/api")
	protectedRoutes.Use(middlewares.AuthMiddleware(s.userRepository))
	{
		protectedRoutes.GET("/user", userHandler.GetUser)
	}

	return r
}
