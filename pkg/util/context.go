package util

import (
	"context"

	"github.com/google/uuid"
)

const (
	UserIDCTX    = "user_id"
	UserEmailCTX = "email"
)

func UserIDFromCTX(ctx context.Context) (userID uuid.UUID) {
	user, ok := ctx.Value(UserIDCTX).(string)
	if !ok || user == "" {
		return uuid.Nil
	}

	userID, err := uuid.Parse(user)
	if err != nil {
		return uuid.Nil
	}

	return
}

func UserEmailFromCTX(ctx context.Context) (email string) {
	user := ctx.Value(UserEmailCTX)
	email, _ = user.(string)
	return
}
