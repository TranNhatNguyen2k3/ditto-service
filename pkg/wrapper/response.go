package wrapper

import (
	"encoding/json"
	"net/http"

	"ditto/pkg/errors"

	"github.com/gin-gonic/gin"
)

// Response defines response format
type Response struct {
	StatusCode int         `json:"status_code"`
	Code       int         `json:"code"`
	Data       interface{} `json:"data,omitempty"`
	Message    string      `json:"message,omitempty"`
}

func (r *Response) String() string {
	data, _ := json.Marshal(r)
	return string(data)
}

// NewResponse creates a new response
func NewResponse(statusCode int, code int, data interface{}, message string) *Response {
	return &Response{
		StatusCode: statusCode,
		Code:       code,
		Data:       data,
		Message:    message,
	}
}

// NewSuccessResponse creates a success response
func NewSuccessResponse(data interface{}) *Response {
	return NewResponse(
		http.StatusOK,
		0,
		data,
		"Success",
	)
}

// NewErrorResponse creates an error response
func NewErrorResponse(err *errors.AppError) *Response {
	return NewResponse(
		err.Status,
		err.Code,
		nil,
		err.Message,
	)
}

// JSONOk sends a success response
func JSONOk(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, NewSuccessResponse(data))
}
