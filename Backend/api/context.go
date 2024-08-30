package api

import (
	"context"
	"net/http"

	"github.com/colmedev/IA-KuroJam/Backend/users"
)

type contextKey string

const userContextKey = contextKey("user")

func (app *Api) ContextSetUser(r *http.Request, user *users.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *Api) ContextGetUser(r *http.Request) *users.User {
	user, ok := r.Context().Value(userContextKey).(*users.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
