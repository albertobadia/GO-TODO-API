package users

import "context"

type UserCtxKeyType string

const UserCtxKey UserCtxKeyType = "user"

func GetUserFromContext(ctx context.Context) User {
	return ctx.Value(UserCtxKey).(User)
}
