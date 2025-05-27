package handler

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type ProxyHandler struct {
	targetURL string
}

func NewProxyHandler(targetURL string) *ProxyHandler {
	return &ProxyHandler{
		targetURL: targetURL,
	}
}

func (h *ProxyHandler) Proxy(c *gin.Context) {
	target, err := url.Parse(h.targetURL)
	if err != nil {
		log.Printf("Error parsing target URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid target URL"})
		return
	}

	log.Printf("Proxying request to: %s", target.String())

	proxy := httputil.NewSingleHostReverseProxy(target)

	// Modify the request
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host

		// Forward the original path
		req.URL.Path = singleJoiningSlash(target.Path, c.Param("proxyPath"))

		// Forward query parameters
		if target.RawQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = target.RawQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = target.RawQuery + "&" + req.URL.RawQuery
		}

		log.Printf("Proxying request to: %s %s", req.Method, req.URL.String())
	}

	// Modify the response
	proxy.ModifyResponse = func(resp *http.Response) error {
		log.Printf("Received response: %s", resp.Status)
		return nil
	}

	// Error handling
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Proxy error: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "Proxy error: " + err.Error()})
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
