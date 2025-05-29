package middleware

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthConfig holds the authentication configuration
type AuthConfig struct {
	Username string
	Password string
	// Add more auth-related configs here if needed
}

// BasicAuth middleware for basic authentication
func BasicAuth(config *AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.Header("WWW-Authenticate", "Basic realm=Authorization Required")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Check if it's a Basic auth header
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Basic" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Decode the credentials
		payload, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			log.Printf("Failed to decode auth credentials: %v", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Split username and password
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Check credentials
		if pair[0] != config.Username || pair[1] != config.Password {
			log.Printf("Invalid credentials for user: %s", pair[0])
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Add user info to context
		c.Set("username", pair[0])
		c.Next()
	}
}

// APIKeyAuth middleware for API key authentication
func APIKeyAuth(config *AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get API key from header
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// TODO: Implement API key validation logic
		// For now, we'll just check if it matches the password
		if apiKey != config.Password {
			log.Printf("Invalid API key")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}

// CombinedAuth middleware that supports both Basic Auth and API Key
func CombinedAuth(config *AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for API key first
		apiKey := c.GetHeader("X-API-Key")
		if apiKey != "" {
			if apiKey == config.Password {
				c.Next()
				return
			}
		}

		// If no API key or invalid, try Basic Auth
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.Header("WWW-Authenticate", "Basic realm=Authorization Required")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Basic" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		payload, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			log.Printf("Failed to decode auth credentials: %v", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if pair[0] != config.Username || pair[1] != config.Password {
			log.Printf("Invalid credentials for user: %s", pair[0])
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("username", pair[0])
		c.Next()
	}
}
