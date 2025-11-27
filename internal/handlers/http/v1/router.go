package v1

import (
	"database/sql"
	"net/http"

	"github.com/FollG/kafka-with-go/internal/usecases"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
)

func NewRouter(
	productUC *usecases.ProductUseCase,
	db *sql.DB,
	redisClient *redis.Client,
	rateLimit int,
) http.Handler {
	r := chi.NewRouter()

	// Глобальные middleware
	r.Use(middleware.RequestID)
	r.Use(RecoverMiddleware)
	r.Use(LoggingMiddleware)
	r.Use(ContentTypeJSONMiddleware)
	r.Use(RateLimitMiddleware(rateLimit))

	// Health checks
	healthHandler := NewHealthHandler(db, redisClient)
	healthRouter := chi.NewRouter()
	healthHandler.RegisterRoutes(healthRouter)
	r.Mount("/health", healthRouter)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		productHandler := NewProductHandler(productUC)

		r.Route("/products", func(r chi.Router) {
			r.Post("/", productHandler.CreateProduct)
			r.Get("/", productHandler.ListProducts)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", productHandler.GetProduct)
				r.Put("/", productHandler.UpdateProduct)
				r.Delete("/", productHandler.DeleteProduct)
			})
		})
	})

	return r
}
