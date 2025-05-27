package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// for gin recover if panic occurs
func Recover() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("GinPanic: %ve", err)
				ctx.Header("content-type", "application/json")
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"errors": "there are some internal error"})
			}
		}()
		ctx.Next()
	}
}
