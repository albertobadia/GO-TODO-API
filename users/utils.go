package users

import "context"

// Getting a user from context needs a type for the key
type UserCtxKeyType string

// The key is a string, but golang shows a warning if you use a string directly as a key
const UserCtxKey UserCtxKeyType = "user"

func GetUserFromContext(ctx context.Context) User {
	return ctx.Value(UserCtxKey).(User)
}
