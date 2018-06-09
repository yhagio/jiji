package middlewares

import (
	"context"
	"jiji/models"
)

type privateKey string

const (
	userKey privateKey = "user"
)

func AssignUserToContext(ctx context.Context, user *models.User) context.Context {
	// we dont have to worry about invalid data type
	return context.WithValue(ctx, userKey, user)
}

func LookUpUserFromContext(ctx context.Context) *models.User {
	temp := ctx.Value(userKey)
	if temp != nil {
		user, ok := temp.(*models.User)
		if ok {
			return user
		}
	}
	return nil
}
