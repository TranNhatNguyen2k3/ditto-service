package router

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"ditto/config"
	"ditto/internal/http/handler"
	"ditto/internal/middleware"
)

type Router struct {
	engine *gin.Engine
	proxy  *handler.ProxyHandler
	config *config.Config
}

func NewRouter(engine *gin.Engine, proxy *handler.ProxyHandler, config *config.Config) *Router {
	return &Router{
		engine: engine,
		proxy:  proxy,
		config: config,
	}
}

func (r *Router) Setup() {
	// Create auth config for proxy API
	authConfig := &middleware.AuthConfig{
		Username: r.config.Proxy.AuthUsername,
		Password: r.config.Proxy.AuthPassword,
	}

	// Health check endpoint (no auth required)
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Apply auth middleware to all routes under /api
	r.engine.Use(func(c *gin.Context) {
		// Skip auth for health check
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		// Apply auth middleware for all other routes
		middleware.BasicAuth(authConfig)(c)
	})

	// API routes group
	api := r.engine.Group("/api")
	{
		// Setup device routes
		SetupDeviceRoutes(api, r.config)

		// Proxy /api/things/* to Ditto
		api.Any("/things/*path", r.proxy.ProxyRequest)
	}

	// Print all registered routes
	log.Printf("All registered routes:")
	for _, route := range r.engine.Routes() {
		log.Printf("%s %s", route.Method, route.Path)
	}
}

var Module = fx.Options(
	fx.Provide(
		NewRouter,
	),
)
