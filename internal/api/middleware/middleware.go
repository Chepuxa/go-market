package middleware

import (
	"fmt"
	"net/http"
	"training/proj/internal/customerrors"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				customerrors.ServerErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func Authenticator(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())

			if err != nil {
				customerrors.AuthenticationRequiredResponse(w, r)
				return
			}

			if token == nil || jwt.Validate(token, ja.ValidateOptions()...) != nil {
				customerrors.AuthenticationRequiredResponse(w, r)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}
