package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/zeon-code/tiny-url/internal/db"
	"github.com/zeon-code/tiny-url/internal/model"
	"github.com/zeon-code/tiny-url/internal/pkg/test"
	"github.com/zeon-code/tiny-url/internal/repository"
)

func TestUrlRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("list urls", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		repo := repository.NewURLRepository(fake.DB(), fake.Memory(), nil)
		query := "SELECT id, code, target FROM urls ORDER BY id DESC LIMIT $1"

		rows := sqlmock.NewRows([]string{"id", "code", "target"}).
			AddRow(int64(5), "5", "target5").
			AddRow(int64(4), "4", "target4").
			AddRow(int64(3), "3", "target3").
			AddRow(int64(2), "2", "target2").
			AddRow(int64(1), "1", "target1")

		fake.DBMock.ExpectQuery(query).WithArgs(5).WillReturnRows(rows)
		urls, err := repo.List(ctx, 5, ">", nil)

		assert.NoError(t, err)
		assert.Len(t, urls, 5)
	})

	t.Run("list urls with no matches", func(t *testing.T) {
		limit := 5
		cursor := int64(8888)
		fake := test.NewFakeDependencies()

		repo := repository.NewURLRepository(fake.DB(), fake.Memory(), nil)
		query := "SELECT id, code, target FROM urls WHERE id > $1 ORDER BY id DESC LIMIT $2"

		rows := sqlmock.NewRows([]string{"id", "code", "target"})

		fake.DBMock.ExpectQuery(query).WithArgs(cursor, limit).WillReturnRows(rows)
		urls, err := repo.List(ctx, limit, ">", &cursor)

		assert.NoError(t, err)
		assert.Len(t, urls, 0)
	})

	t.Run("list urls with cursor", func(t *testing.T) {
		cursor := int64(1)
		fake := test.NewFakeDependencies()
		repo := repository.NewURLRepository(fake.DB(), fake.Memory(), nil)
		query := "SELECT id, code, target FROM urls WHERE id > $1 ORDER BY id DESC LIMIT $2"

		rows := sqlmock.NewRows([]string{"id", "code", "target"}).
			AddRow(int64(1), "6", "target6").
			AddRow(int64(5), "5", "target5").
			AddRow(int64(4), "4", "target4").
			AddRow(int64(3), "3", "target3").
			AddRow(int64(2), "2", "target2")

		fake.DBMock.ExpectQuery(query).WithArgs(&cursor, 5).WillReturnRows(rows)
		urls, err := repo.List(ctx, 5, ">", &cursor)

		assert.NoError(t, err)
		assert.Len(t, urls, 5)
	})

	t.Run("create url", func(t *testing.T) {
		now := time.Now()
		target := "target"
		fake := test.NewFakeDependencies()
		repo := repository.NewURLRepository(fake.DB(), fake.Memory(), nil)

		updateQuery := "UPDATE urls SET code = $1 WHERE id = $2"
		insertQuery := "INSERT INTO urls (target, code) VALUES ($1, '') RETURNING id, target, code, created_at, updated_at"

		rows := sqlmock.NewRows([]string{"id", "target", "code", "created_at", "updated_at"}).
			AddRow(int64(1), target, "", now, now)

		fake.DBMock.ExpectBegin()
		fake.DBMock.ExpectQuery(insertQuery).WithArgs(target).WillReturnRows(rows)
		fake.DBMock.ExpectExec(updateQuery).WithArgs("1", int64(1)).WillReturnResult(sqlmock.NewResult(1, 1))
		fake.DBMock.ExpectCommit()

		url, err := repo.Create(ctx, target)

		assert.NoError(t, err)
		assert.Equal(t, model.URL{ID: 1, Code: "1", Target: target, CreatedAt: url.CreatedAt, UpdatedAt: url.UpdatedAt}, *url)
	})

	t.Run("create url with insert error should rollback", func(t *testing.T) {
		target := "target"
		fake := test.NewFakeDependencies()

		repo := repository.NewURLRepository(fake.DB(), fake.Memory(), nil)
		insertQuery := "INSERT INTO urls (target, code) VALUES ($1, '') RETURNING id, target, code, created_at, updated_at"

		fake.DBMock.ExpectBegin()
		fake.DBMock.ExpectQuery(insertQuery).WithArgs(target).WillReturnError(db.ErrDBInvalidBackend)
		fake.DBMock.ExpectRollback()

		url, err := repo.Create(ctx, target)

		assert.Equal(t, db.ErrDBInvalidBackend, err)
		assert.Nil(t, url)
	})

	t.Run("create url with update error should rollback", func(t *testing.T) {
		now := time.Now()
		target := "target"
		fake := test.NewFakeDependencies()

		repo := repository.NewURLRepository(fake.DB(), fake.Memory(), nil)
		updateQuery := "UPDATE urls SET code = $1 WHERE id = $2"
		insertQuery := "INSERT INTO urls (target, code) VALUES ($1, '') RETURNING id, target, code, created_at, updated_at"

		rows := sqlmock.NewRows([]string{"id", "target", "code", "created_at", "updated_at"}).
			AddRow(int64(9999), target, "", now, now)

		fake.DBMock.ExpectBegin()
		fake.DBMock.ExpectQuery(insertQuery).WithArgs(target).WillReturnRows(rows)
		fake.DBMock.ExpectExec(updateQuery).WithArgs("2bH", int64(9999)).WillReturnError(db.ErrDBInvalidBackend)
		fake.DBMock.ExpectRollback()

		url, err := repo.Create(ctx, target)

		assert.Equal(t, db.ErrDBInvalidBackend, err)
		assert.Nil(t, url)
	})

	t.Run("create url with commit error should rollback", func(t *testing.T) {
		now := time.Now()
		target := "target"
		fake := test.NewFakeDependencies()

		repo := repository.NewURLRepository(fake.DB(), fake.Memory(), nil)
		updateQuery := "UPDATE urls SET code = $1 WHERE id = $2"
		insertQuery := "INSERT INTO urls (target, code) VALUES ($1, '') RETURNING id, target, code, created_at, updated_at"

		rows := sqlmock.NewRows([]string{"id", "target", "code", "created_at", "updated_at"}).
			AddRow(int64(9999), target, "", now, now)

		fake.DBMock.ExpectBegin()
		fake.DBMock.ExpectQuery(insertQuery).WithArgs(target).WillReturnRows(rows)
		fake.DBMock.ExpectExec(updateQuery).WithArgs("2bH", int64(9999)).WillReturnResult(sqlmock.NewResult(1, 1))
		fake.DBMock.ExpectCommit().WillReturnError(db.ErrDBInvalidBackend)
		fake.DBMock.ExpectRollback()

		url, err := repo.Create(ctx, target)

		assert.Equal(t, db.ErrDBInvalidBackend, err)
		assert.Nil(t, url)
	})

	t.Run("get by id", func(t *testing.T) {
		now := time.Now()
		fake := test.NewFakeDependencies()
		repo := repository.NewURLRepository(fake.DB(), fake.Memory(), nil)
		query := "SELECT * FROM urls WHERE id = $1"

		rows := sqlmock.NewRows([]string{"id", "code", "target", "created_at", "updated_at"}).
			AddRow(int64(1), "1", "target1", now, now)

		fake.DBMock.ExpectQuery(query).WillReturnRows(rows)
		url, err := repo.GetByID(ctx, int64(1))

		assert.NoError(t, err)
		assert.Equal(t, model.URL{ID: 1, Code: "1", Target: "target1", CreatedAt: url.CreatedAt, UpdatedAt: url.UpdatedAt}, *url)
	})

	t.Run("get by id when not found", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		repo := repository.NewURLRepository(fake.DB(), fake.Memory(), nil)
		query := "SELECT * FROM urls WHERE id = $1"

		rows := sqlmock.NewRows([]string{"id", "code", "target", "created_at", "updated_at"})

		fake.DBMock.ExpectQuery(query).WillReturnRows(rows)
		url, err := repo.GetByID(ctx, int64(4))

		assert.Nil(t, url)
		assert.Equal(t, db.ErrDBNotFound, err)
	})
}
