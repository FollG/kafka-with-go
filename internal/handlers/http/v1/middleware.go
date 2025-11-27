package v1

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/FollG/kafka-with-go/internal/pkg/logger"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/time/rate"
)

// RateLimitMiddleware ограничивает количество запросов, число берется из docker
func RateLimitMiddleware(requestsPerSecond int) func(http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				render.Status(r, http.StatusTooManyRequests)
				render.JSON(w, r, ErrorResponse{
					Error:   "rate_limit_exceeded",
					Message: "Too many requests",
				})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// LoggingMiddleware логирует запросы
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Создаем request ID
		requestID := middleware.GetReqID(r.Context())
		if requestID == "" {
			requestID = fmt.Sprintf("%d", time.Now().UnixNano())
		}

		// Добавляем request ID в контекст
		ctx := context.WithValue(r.Context(), "request_id", requestID)
		r = r.WithContext(ctx)

		// Используем ResponseWriter для захвата статуса
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		// Логируем информацию о запросе
		duration := time.Since(start)

		logger.Info(r.Context(), "http_request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.Status(),
			"duration", duration.String(),
			"user_agent", r.UserAgent(),
			"request_id", requestID,
		)
	})
}

// RecoverMiddleware обрабатывает паники
func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(r.Context(), "panic_recovered",
					"error", err,
					"path", r.URL.Path,
				)

				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, ErrorResponse{
					Error:   "internal_error",
					Message: "Internal server error",
				})
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// ContentTypeJSONMiddleware устанавливает Content-Type для ответов
func ContentTypeJSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
