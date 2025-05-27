package middleware

import (
	"net/http"

	"ditto/config"
	"ditto/pkg/errors"
	"ditto/pkg/util"
	"ditto/pkg/wrapper"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := util.TokenValid(c, config.JWT)
		if err != nil {
			c.JSON(http.StatusUnauthorized, wrapper.NewErrorResponse(
				errors.NewUnauthorizedError("Unauthorized"),
			))
			c.Abort()
			return
		}
		c.Next()
	}
}
