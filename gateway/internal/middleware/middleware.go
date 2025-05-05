package middleware

import (
	"context"
	"net/http"
	"strings"

	customContext "github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/utils/context"
	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/utils/response"
	"github.com/RecruitEase-Capstone/recruitEase-BE/pkg/jwt"
)

type IMiddleware interface {
	JwtAuthMiddleware(next http.Handler) http.Handler
}

type Middleware struct {
	jwt jwt.JWTItf
}

func NewMiddleware(jwt jwt.JWTItf) IMiddleware {
	return &Middleware{jwt: jwt}
}

func (m *Middleware) JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		bearerToken := r.Header.Get("Authorization")

		if bearerToken == "" {
			response.FailedResponse(w, http.StatusUnauthorized, "Authorization token is required", nil)
			return
		}

		token := strings.Split(bearerToken, " ")[1]

		userID, err := m.jwt.VerifyToken(token)
		if err != nil {
			response.FailedResponse(w, http.StatusUnauthorized, err.Error(), nil)
			return
		}

		// Set userID in context
		ctx := context.WithValue(r.Context(), customContext.UserIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
