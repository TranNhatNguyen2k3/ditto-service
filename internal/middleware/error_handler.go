package middleware

import (
	"net/http"
	"runtime/debug"

	"ditto/pkg/errors"
	"ditto/pkg/logger"
	"ditto/pkg/wrapper"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

// ErrorHandler handles application errors
type ErrorHandler struct {
	logger logger.Logger
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(logger logger.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

// Handle is the main error handling middleware
func (h *ErrorHandler) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				h.logger.Error("panic recovered",
					zap.Any("error", err),
					zap.String("stack", string(debug.Stack())),
				)
				c.JSON(http.StatusInternalServerError, wrapper.NewErrorResponse(
					errors.NewInternalServerError("Internal server error"),
				))
			}
		}()

		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			lastErr := c.Errors.Last()
			if lastErr == nil {
				return
			}
			err := lastErr.Err

			h.logger.Error("request error", zap.Any("error", err))

			// Handle different types of errors
			switch e := err.(type) {
			case *errors.AppError:
				c.JSON(e.Status, wrapper.NewErrorResponse(e))
			case validator.ValidationErrors:
				if len(e) > 0 {
					c.JSON(http.StatusBadRequest, wrapper.NewErrorResponse(
						errors.NewBadRequestError("validation error"),
					))
				} else {
					c.JSON(http.StatusBadRequest, wrapper.NewErrorResponse(
						errors.NewBadRequestError("invalid request"),
					))
				}
			case *pq.Error:
				h.handleDatabaseError(c, e)
			default:
				c.JSON(http.StatusInternalServerError, wrapper.NewErrorResponse(
					errors.NewInternalServerError(e.Error()),
				))
			}
		}
	}
}

// handleDatabaseError handles database specific errors
func (h *ErrorHandler) handleDatabaseError(c *gin.Context, err *pq.Error) {
	if err == nil {
		c.JSON(http.StatusInternalServerError, wrapper.NewErrorResponse(
			errors.NewInternalServerError("database error"),
		))
		return
	}

	switch err.Code {
	case "23505": // unique_violation
		c.JSON(http.StatusConflict, wrapper.NewErrorResponse(
			errors.NewConflictError("resource already exists"),
		))
	case "23503": // foreign_key_violation
		c.JSON(http.StatusBadRequest, wrapper.NewErrorResponse(
			errors.NewBadRequestError("invalid reference"),
		))
	default:
		h.logger.Error("database error", zap.Any("error", err))
		c.JSON(http.StatusInternalServerError, wrapper.NewErrorResponse(
			errors.NewInternalServerError("database error"),
		))
	}
}
