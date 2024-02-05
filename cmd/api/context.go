package main

import (
	"context"
	"net/http"

	"github.com/ipramudya/go-greenlight/internal/data"
)

type contextKey string

// use this as the key for getting and setting user information inside request context
const userContextKey = contextKey("user")

// return a new copy f the request with the provided user sturct added to the context
func (*application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (*application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
