package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeon-code/tiny-url/internal/pkg/test"
	"github.com/zeon-code/tiny-url/internal/repository"
)

func TestRepositories(t *testing.T) {
	t.Run("Should define url repository", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		repositories := repository.NewRepositories(fake.DB(), fake.Memory(), fake.Logger())

		assert.NotNil(t, repositories.Url)
	})
}
