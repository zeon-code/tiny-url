package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/zeon-code/tiny-url/internal/model"
	"github.com/zeon-code/tiny-url/internal/pkg/cache"
	"github.com/zeon-code/tiny-url/internal/pkg/pagination"
	"github.com/zeon-code/tiny-url/internal/service"
)

type UrlHandler struct {
	UrlSvc service.URLService
	logger *slog.Logger
}

func NewUrlHandler(services service.Services, logger *slog.Logger) UrlHandler {
	return UrlHandler{
		UrlSvc: services.Url,
		logger: logger,
	}
}

type UrlListResponse struct {
	Urls []model.URL     `json:"items"`
	Page pagination.Page `json:"page"`
}

type UrlCreateRequest struct {
	Target string `json:"target"`
}

type UrlCreateResponse struct {
	ID     int64  `json:"id"`
	Code   string `json:"code"`
	Target string `json:"target"`
}

func (h UrlHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	request := UrlCreateRequest{}
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(body, &request); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	url, err := h.UrlSvc.Create(r.Context(), request.Target)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(UrlCreateResponse{
		ID:     url.ID,
		Code:   url.Code,
		Target: url.Target,
	})

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func (h UrlHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := 50

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	direction, cursor := pagination.GetCursor(r)
	urls, err := h.UrlSvc.List(cache.WithCache(ctx), limit, direction, cursor)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	cursorKey := func(u model.URL) int64 { return u.ID }
	data, err := pagination.NewPagination(urls, limit, cursor).Encode(cursorKey)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h UrlHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	url, err := h.UrlSvc.GetByID(cache.WithCache(ctx), id)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(url)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
