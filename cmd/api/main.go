package main

import (
	"net/http"

	"github.com/zeon-code/tiny-url/internal/http/handler"
	"github.com/zeon-code/tiny-url/internal/pkg/config"
	"github.com/zeon-code/tiny-url/internal/pkg/log"
	"github.com/zeon-code/tiny-url/internal/pkg/metric"
	"github.com/zeon-code/tiny-url/internal/repository"
	"github.com/zeon-code/tiny-url/internal/service"
)

func main() {
	conf := config.NewConfiguration()

	logger := log.NewLogger(conf)
	metrics := metric.NewMetricClient(conf, logger.With("client", "metric"))

	repo := repository.NewRepositoriesFromConfig(conf, metrics, logger.With("package", "repository"))
	svc := service.NewServices(repo, logger.With("package", "service"))

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler.NewRouter(svc, metrics, logger.With("package", "handler")),
	}

	logger.Info("Starting server")
	server.ListenAndServe()
}
