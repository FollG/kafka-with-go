package v1

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	db    *sql.DB
	redis *redis.Client
}

func NewHealthHandler(db *sql.DB, redis *redis.Client) *HealthHandler {
	return &HealthHandler{
		db:    db,
		redis: redis,
	}
}

func (h *HealthHandler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.HealthCheck)
	r.Get("/ready", h.ReadyCheck)
}

// HealthCheck возвращает статус сервиса
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	healthStatus := "healthy"
	statusCode := http.StatusOK

	// Проверяем подключение к БД
	if err := h.db.PingContext(ctx); err != nil {
		healthStatus = "unhealthy"
		statusCode = http.StatusServiceUnavailable
	}

	// Проверяем подключение к Redis
	if err := h.redis.Ping(ctx).Err(); err != nil {
		healthStatus = "unhealthy"
		statusCode = http.StatusServiceUnavailable
	}

	response := HealthResponse{
		Status:  healthStatus,
		Message: "Service health check",
	}

	render.Status(r, statusCode)
	render.JSON(w, r, response)
}

// ReadyCheck проверяет готовность сервиса
func (h *HealthHandler) ReadyCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	readyStatus := "ready"
	statusCode := http.StatusOK

	// Проверяем все зависимости
	if err := h.db.PingContext(ctx); err != nil {
		readyStatus = "not_ready"
		statusCode = http.StatusServiceUnavailable
	}

	if err := h.redis.Ping(ctx).Err(); err != nil {
		readyStatus = "not_ready"
		statusCode = http.StatusServiceUnavailable
	}

	response := HealthResponse{
		Status:  readyStatus,
		Message: "Service readiness check",
	}

	render.Status(r, statusCode)
	render.JSON(w, r, response)
}

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
