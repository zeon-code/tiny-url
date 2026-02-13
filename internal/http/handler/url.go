package handler

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	json "github.com/json-iterator/go"
	"github.com/zeon-code/tiny-url/internal/db"
	"github.com/zeon-code/tiny-url/internal/model"
	"github.com/zeon-code/tiny-url/internal/pkg/cache"
	"github.com/zeon-code/tiny-url/internal/pkg/observability"
	"github.com/zeon-code/tiny-url/internal/pkg/pagination"
	"github.com/zeon-code/tiny-url/internal/service"
)

type UrlHandler struct {
	UrlSvc service.URLService
	logger observability.Logger
}

func NewUrlHandler(services service.Services, observer observability.Observer) UrlHandler {
	return UrlHandler{
		UrlSvc: services.Url,
		logger: observer.Logger().With("handler", "url"),
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
	ctx := r.Context()

	defer r.Body.Close()

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	request := UrlCreateRequest{}
	body, err := io.ReadAll(r.Body)

	if err != nil {
		observability.TraceError(ctx, http.StatusText(http.StatusInternalServerError), err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(body, &request); err != nil {
		observability.TraceError(ctx, http.StatusText(http.StatusBadRequest), err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	url, err := h.UrlSvc.Create(ctx, request.Target)

	if err != nil {
		h.logger.Error(ctx, "error creating url", slog.Any("error", err))
		observability.TraceError(ctx, http.StatusText(http.StatusInternalServerError), err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(UrlCreateResponse{
		ID:     url.ID,
		Code:   url.Code,
		Target: url.Target,
	})

	if err != nil {
		observability.TraceError(ctx, http.StatusText(http.StatusInternalServerError), err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func (h UrlHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := 50
	ctx := r.Context()

	if r.Header.Get("Accept") != "application/json" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	direction, cursor := pagination.GetCursor(r)
	urls, err := h.UrlSvc.List(cache.WithCache(ctx), limit, direction, cursor)

	if err != nil {
		observability.TraceError(ctx, http.StatusText(http.StatusInternalServerError), err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	cursorKey := func(u model.URL) int64 { return u.ID }
	data, err := pagination.NewPagination(urls, limit, cursor).Encode(cursorKey)

	if err != nil {
		observability.TraceError(ctx, http.StatusText(http.StatusInternalServerError), err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h UrlHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Header.Get("Accept") != "application/json" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		observability.TraceError(ctx, http.StatusText(http.StatusBadRequest), err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	url, err := h.UrlSvc.GetByID(cache.WithCache(ctx), id)

	if errors.Is(err, db.ErrDBResourceNotFound) {
		observability.TraceError(ctx, http.StatusText(http.StatusNotFound), err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	} else if err != nil {
		observability.TraceError(ctx, http.StatusText(http.StatusInternalServerError), err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(url)

	if err != nil {
		observability.TraceError(ctx, http.StatusText(http.StatusInternalServerError), err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h UrlHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	code := r.PathValue("code")

	if code == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	url, err := h.UrlSvc.GetByCode(cache.WithCache(ctx), code)

	if errors.Is(err, db.ErrDBResourceNotFound) {
		observability.TraceError(ctx, http.StatusText(http.StatusNotFound), err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	} else if err != nil {
		observability.TraceError(ctx, http.StatusText(http.StatusInternalServerError), err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", url.Target)
	w.WriteHeader(http.StatusFound)
}
