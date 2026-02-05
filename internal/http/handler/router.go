package handler

import (
	"log/slog"
	"net/http"

	"github.com/zeon-code/tiny-url/internal/http/middleware"
	"github.com/zeon-code/tiny-url/internal/pkg/metric"
	"github.com/zeon-code/tiny-url/internal/service"
)

func NewRouter(svc service.Services, metrics metric.MetricClient, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	url := NewUrlHandler(svc, logger.With("handler", "url"))

	mux.HandleFunc("GET /api/v1/url/", url.List)
	mux.HandleFunc("POST /api/v1/url/", url.Create)
	mux.HandleFunc("GET /api/v1/url/{id}", url.GetByID)

	return middleware.HTTPMetrics(metrics)(mux)
}
