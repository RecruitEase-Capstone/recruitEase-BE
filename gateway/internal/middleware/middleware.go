package middleware

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	customContext "github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/utils/context"
	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/utils/response"
	"github.com/RecruitEase-Capstone/recruitEase-BE/pkg/jwt"
	"github.com/rs/zerolog"
)

type IMiddleware interface {
	JwtAuthMiddleware(next http.Handler) http.Handler
	LoggingMiddleware(next http.Handler) http.Handler
	RateLimiter(next http.Handler) http.Handler
}

type Middleware struct {
	jwt jwt.JWTItf
	log zerolog.Logger
}

func NewMiddleware(jwt jwt.JWTItf, log zerolog.Logger) IMiddleware {
	return &Middleware{jwt: jwt, log: log}
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

var (
	rateLimit  = 15
	window     = 10 * time.Second
	requests   = make(map[string]int)
	timestamps = make(map[string]time.Time)
	mu         sync.Mutex
)

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

func (m *Middleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rr := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rr, r)

		duration := time.Since(start)

		m.log.Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Str("remote_addr", r.RemoteAddr).
			Int("status", rr.statusCode).
			Dur("duration", duration).
			Msg("Incoming HTTP request")
	})
}

func (m *Middleware) RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		userID, ok := r.Context().Value("userID").(string)
		if !ok || userID == "" {
			userID = getRealIP(r)
		}

		if _, exists := requests[userID]; !exists {
			requests[userID] = 0
			timestamps[userID] = time.Now()
		}

		elapsed := time.Since(timestamps[userID])

		if elapsed > window {
			requests[userID] = 0
			timestamps[userID] = time.Now()
		}

		if requests[userID] >= rateLimit {
			response.FailedResponse(w, http.StatusTooManyRequests, "rate limit exceeded", nil)
			return
		}

		requests[userID]++
		next.ServeHTTP(w, r)
	})
}

func getRealIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	return ""
}
