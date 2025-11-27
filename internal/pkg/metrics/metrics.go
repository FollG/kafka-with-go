package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// HTTP метрики
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "path", "status"})

	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})

	// Бизнес метрики
	productsCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "products_created_total",
		Help: "Total number of products created",
	})

	productsUpdated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "products_updated_total",
		Help: "Total number of products updated",
	})

	productsDeleted = promauto.NewCounter(prometheus.CounterOpts{
		Name: "products_deleted_total",
		Help: "Total number of products deleted",
	})

	// Kafka метрики
	kafkaMessagesProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "kafka_messages_processed_total",
		Help: "Total number of Kafka messages processed",
	}, []string{"topic", "status"})
)

func Init(port int) {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		fmt.Printf("Starting metrics server on port %d\n", port)
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if err != nil {
			panic(err)
		}
	}()
}

func RecordHTTPRequest(method, path string, status int, duration float64) {
	httpRequestsTotal.WithLabelValues(method, path, fmt.Sprintf("%d", status)).Inc()
	httpRequestDuration.WithLabelValues(method, path).Observe(duration)
}

func RecordProductCreated() {
	productsCreated.Inc()
}

func RecordProductUpdated() {
	productsUpdated.Inc()
}

func RecordProductDeleted() {
	productsDeleted.Inc()
}

func RecordKafkaMessageProcessed(topic, status string) {
	kafkaMessagesProcessed.WithLabelValues(topic, status).Inc()
}
