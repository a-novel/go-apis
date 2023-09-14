package apis

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

func GetCORS(allowedOrigins []string) cors.Config {
	return cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     cors.DefaultConfig().AllowMethods,
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Type", "Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
}

type HealthChecker func() error

type RouterConfig struct {
	Logger zerolog.Logger
	// ProjectID is the ID of the Google Cloud project the application is deployed to.
	// It is optional, and should only be used in production.
	ProjectID string
	CORS      cors.Config
	// Prod tells gin to use production settings.
	Prod bool
	// Health is a list of dependencies, to check for availability.
	Health map[string]HealthChecker
}

// GetRouter returns a new gin.Engine with custom defaults from the config object. It automatically allocates the
// "/healthcheck" and "/ping" routes.
func GetRouter(cfg RouterConfig) *gin.Engine {
	router := gin.New()
	router.Use(gin.RecoveryWithWriter(cfg.Logger), Logger(cfg.Logger, cfg.ProjectID), cors.New(cfg.CORS))

	if cfg.Prod {
		gin.SetMode(gin.ReleaseMode)
		router.TrustedPlatform = gin.PlatformGoogleAppEngine
	}

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	router.GET("/healthcheck", func(c *gin.Context) {
		dependencies := gin.H{}

		for name, checker := range cfg.Health {
			if err := checker(); err != nil {
				dependencies[name] = fmt.Sprintf("error: %q", err.Error())
			} else {
				dependencies[name] = "ok"
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"production":   cfg.Prod,
			"dependencies": dependencies,
		})
	})

	return router
}
