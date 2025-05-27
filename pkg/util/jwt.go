package util

import (
	"fmt"
	"strings"

	"ditto/config"
	"ditto/pkg/constant"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cast"
)

// ExtractToken lấy token từ query params hoặc header Authorization
func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	tokenString := strings.Split(bearerToken, " ")
	if len(tokenString) == 2 {
		if tokenString[0] != "Bearer" {
			return ""
		}
		return tokenString[1]
	}
	return ""
}

// ExtractTokenID phân tích token và gán userID và email vào context
func ExtractTokenID(c *gin.Context, config config.JWTConfig) error {
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Secret), nil
	})
	if err != nil {
		return fmt.Errorf("can't parse token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userID := fmt.Sprint(claims[constant.UserID])
		if len(userID) == 0 {
			return nil
		}
		email := fmt.Sprint(claims[constant.Email])
		if len(email) == 0 {
			return nil
		}
		// Set userID và email vào context
		c.Set(constant.UserID, userID)
		c.Set(constant.Email, email)
		return nil
	}
	return nil
}

// TokenValid kiểm tra tính hợp lệ của token
func TokenValid(c *gin.Context, config config.JWTConfig) error {
	err := ExtractTokenID(c, config)
	if err != nil {
		return fmt.Errorf("can't extract token")
	}
	return nil
}

// ExtractUserIDFromContext lấy userID từ context
func ExtractUserIDFromContext(c *gin.Context) int {
	val := c.GetString(constant.UserID)
	return cast.ToInt(val)
}
