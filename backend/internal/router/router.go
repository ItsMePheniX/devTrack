package router

import (
	"devtrack/internal/config"
	"devtrack/internal/database"
	"devtrack/internal/handlers"

	"github.com/gin-gonic/gin"
)

func Setup(cfg *config.Config, queries *database.Queries) *gin.Engine {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	r.GET("/health", handlers.HealthCheck)

	authHandler := handlers.NewAuthHandler(queries, cfg)

	v1 := r.Group("/api/v1")
	{
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
		}
	}
	return r
}
