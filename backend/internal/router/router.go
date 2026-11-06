package router

import (
	"devtrack/internal/config"
	"devtrack/internal/handlers"

	"github.com/gin-gonic/gin"
)

func Setup(cfg *config.Config) *gin.Engine {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	r.GET("/health", handlers.HealthCheck)

	v1 := r.Group("api/v1")
	{
		_ = v1 //place holder until i add routes
	}
	return r
}
