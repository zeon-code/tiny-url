package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	json "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeon-code/tiny-url/internal/http/handler"
	"github.com/zeon-code/tiny-url/internal/model"
	"github.com/zeon-code/tiny-url/internal/pkg/test"
)

func TestHealthHandler(t *testing.T) {
	t.Run("healthcheck ready", func(t *testing.T) {
		var payload model.Health
		fake := test.NewFakeDependencies()
		router := handler.NewRouter(fake.Services(), fake.Observer())

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)

		fake.MockUrlCreate()
		router.ServeHTTP(rec, req)

		err := json.NewDecoder(rec.Body).Decode(&payload)
		require.NoError(t, err)

		assert.Equal(t, model.Health{Status: "ready", Reason: ""}, payload)
	})

	t.Run("healthcheck live", func(t *testing.T) {
		var payload model.Health
		fake := test.NewFakeDependencies()
		router := handler.NewRouter(fake.Services(), fake.Observer())

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/health/live", nil)

		fake.MockUrlCreate()
		router.ServeHTTP(rec, req)

		err := json.NewDecoder(rec.Body).Decode(&payload)
		require.NoError(t, err)

		assert.Equal(t, model.Health{Status: "alive", Reason: ""}, payload)
	})
}
