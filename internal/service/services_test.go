package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeon-code/tiny-url/internal/pkg/test"
	"github.com/zeon-code/tiny-url/internal/service"
)

func TestServices(t *testing.T) {
	t.Run("Should define url service", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		services := service.NewServices(fake.Repositories(), fake.Logger())

		assert.NotNil(t, services.Url)
	})
}
