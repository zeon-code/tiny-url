package handler

import (
	"context"
	"net/http"
	"time"

	json "github.com/json-iterator/go"
	"github.com/zeon-code/tiny-url/internal/model"
	"github.com/zeon-code/tiny-url/internal/pkg/observability"
	"github.com/zeon-code/tiny-url/internal/service"
)

type HealthHandler struct {
	HealthSvc service.HealthService
	logger    observability.Logger
}

func NewHealthHandler(services service.Services, observer observability.Observer) HealthHandler {
	return HealthHandler{
		HealthSvc: services.Health,
		logger:    observer.Logger().With("handler", "health"),
	}
}

func (h HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	var health model.Health = model.Health{Status: "alive"}

	data, err := json.Marshal(health)

	if err != nil {
		observability.TraceError(r.Context(), http.StatusText(http.StatusInternalServerError), err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	var statusCode int = http.StatusOK
	var health model.Health = model.Health{Status: "ready"}
	reason, err := h.HealthSvc.Ping(ctx)

	if err != nil {
		health.Reason = reason
		health.Status = "not_ready"
		statusCode = http.StatusServiceUnavailable
		observability.TraceError(ctx, http.StatusText(http.StatusServiceUnavailable), err)
	}

	data, err := json.Marshal(health)

	if err != nil {
		observability.TraceError(ctx, http.StatusText(http.StatusInternalServerError), err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}
