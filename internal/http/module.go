package http

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"ditto/config"
	"ditto/internal/http/handler"
	"ditto/internal/http/router"
)

func NewGinEngine() *gin.Engine {
	engine := gin.Default()
	return engine
}

func NewProxyHandler(cfg *config.Config) *handler.ProxyHandler {
	return handler.NewProxyHandler(cfg.Proxy.TargetURL, cfg.Proxy.Username, cfg.Proxy.Password)
}

var Module = fx.Options(
	fx.Provide(
		NewGinEngine,
		NewProxyHandler,
		router.NewRouter,
	),
	fx.Invoke(func(r *router.Router) {
		r.Setup()
	}),
)
