package router

import (
	"ditto/internal/http/handler"

	"github.com/gin-gonic/gin"
)

func SetupProxyRouter(r *gin.Engine, targetURL string) {
	proxyHandler := handler.NewProxyHandler(targetURL)

	// Group all proxy routes under /proxy
	proxyGroup := r.Group("/proxy")
	{
		// Catch-all route for proxy requests
		proxyGroup.Any("/*proxyPath", proxyHandler.Proxy)
	}
}
