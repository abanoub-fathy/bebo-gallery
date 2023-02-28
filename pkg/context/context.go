package context

import (
	"context"

	"github.com/abanoub-fathy/bebo-gallery/model"
)

type privateKey string

const userKey privateKey = "user"

// WithUser is used to create a new context with user value
func WithUser(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// UserValue is used to get the user from ctx
// it will return pointer to user or nil if user
// is not set in the context
func UserValue(ctx context.Context) *model.User {
	if valFromKey := ctx.Value(userKey); valFromKey != nil {
		if user, ok := valFromKey.(*model.User); ok {
			return user
		}
	}

	return nil
}
