package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	json "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeon-code/tiny-url/internal/http/handler"
	"github.com/zeon-code/tiny-url/internal/model"
	"github.com/zeon-code/tiny-url/internal/pkg/pagination"
	"github.com/zeon-code/tiny-url/internal/pkg/test"
)

func TestUrlHandler(t *testing.T) {

	t.Run("create url", func(t *testing.T) {
		var payload handler.UrlCreateResponse
		fake := test.NewFakeDependencies()
		router := handler.NewRouter(fake.Services(), fake.Observer())

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/url/", bytes.NewBufferString(`{"target":"target"}`))
		req.Header.Set("Content-Type", "application/json")

		fake.MockUrlCreate()
		router.ServeHTTP(rec, req)

		err := json.NewDecoder(rec.Body).Decode(&payload)
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, handler.UrlCreateResponse{ID: 1, Code: "1", Target: "target"}, payload)
	})

	t.Run("list urls", func(t *testing.T) {
		var payload pagination.Pagination[model.URL]
		fake := test.NewFakeDependencies()
		router := handler.NewRouter(fake.Services(), fake.Observer())

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/url/", nil)
		req.Header.Set("Accept", "application/json")

		fake.MockUrlList()
		router.ServeHTTP(rec, req)

		err := json.NewDecoder(rec.Body).Decode(&payload)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, pagination.Pagination[model.URL]{
			Items: []model.URL{
				{ID: 5, Target: "target5", Code: "5"},
				{ID: 4, Target: "target4", Code: "4"},
				{ID: 3, Target: "target3", Code: "3"},
				{ID: 2, Target: "target2", Code: "2"},
				{ID: 1, Target: "target1", Code: "1"},
			},
			Page: pagination.Page{
				Size: 5,
			},
		}, payload)
	})

	t.Run("list urls with cursor", func(t *testing.T) {
		var payload pagination.Pagination[model.URL]
		fake := test.NewFakeDependencies()
		router := handler.NewRouter(fake.Services(), fake.Observer())

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/url/", nil)
		req.Header.Set("Accept", "application/json")

		query := req.URL.Query()
		query.Add("cursor", ">1")
		req.URL.RawQuery = query.Encode()

		fake.MockPaginatedUrlList()
		router.ServeHTTP(rec, req)

		err := json.NewDecoder(rec.Body).Decode(&payload)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, pagination.Pagination[model.URL]{
			Items: []model.URL{
				{ID: 6, Target: "target6", Code: "6"},
				{ID: 5, Target: "target5", Code: "5"},
				{ID: 4, Target: "target4", Code: "4"},
				{ID: 3, Target: "target3", Code: "3"},
				{ID: 2, Target: "target2", Code: "2"},
			},
			Page: pagination.Page{
				Previous: ">6",
				Size:     5,
			},
		}, payload)
	})

	t.Run("url get by id", func(t *testing.T) {
		var payload model.URL
		fake := test.NewFakeDependencies()
		router := handler.NewRouter(fake.Services(), fake.Observer())

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/url/1", nil)
		req.Header.Set("Accept", "application/json")

		fake.MockUrlGetById()
		router.ServeHTTP(rec, req)

		err := json.NewDecoder(rec.Body).Decode(&payload)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		at, _ := time.Parse(time.RFC3339, "2026-01-29T15:23:24Z")
		assert.Equal(t, model.URL{ID: 1, Code: "1", Target: "target1", CreatedAt: &at, UpdatedAt: &at}, payload)
	})

	t.Run("url get by code", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		router := handler.NewRouter(fake.Services(), fake.Observer())

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/r/1", nil)

		fake.MockUrlGetById()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusFound, rec.Code)
		assert.Equal(t, "target1", rec.Header().Get("Location"))
	})
}
