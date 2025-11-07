package auth

import (
	"fmt"
	"net/http"

	"github.com/openshift-online/rh-trex/pkg/errors"
)

type JWTMiddleware interface {
	AuthenticateAccountJWT(next http.Handler) http.Handler
}

type Middleware struct{}

var _ JWTMiddleware = &Middleware{}

func NewAuthMiddleware() (*Middleware, error) {
	middleware := Middleware{}
	return &middleware, nil
}

// AuthenticateAccountJWT Middleware handler to validate JWT tokens and authenticate users
func (a *Middleware) AuthenticateAccountJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		payload, err := GetAuthPayload(r)
		if err != nil {
			handleError(ctx, w, errors.ErrorUnauthorized, fmt.Sprintf("Unable to get payload details from JWT token: %s", err))
			return
		}

		// Append the username to the request context
		ctx = SetUsernameContext(ctx, payload.Username)
		*r = *r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
