package db_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/zeon-code/tiny-url/internal/pkg/test"
)

func TestPostgresClient(t *testing.T) {
	type Row struct {
		Name string `db:"name"`
	}

	ctx := context.Background()

	t.Run("proxy select query", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		query := "SELECT * FROM anything WHERE id > $1"
		rows := sqlmock.NewRows([]string{"name"}).AddRow("diego").AddRow("maria")

		fake.DBMock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
		err := fake.DB().Select(ctx, &[]Row{}, query, 1)

		assert.NoError(t, err)
		assert.Equal(t, fake.DBMetric.LastQuery, query)
		assert.NotNil(t, fake.DBMetric.LastDuration)
	})

	t.Run("proxy select query with error", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		query := "SELECT * FROM anything WHERE id > $1"

		fake.DBMock.ExpectQuery(query).WithArgs(1).WillReturnError(context.Canceled)
		err := fake.DB().Select(ctx, &[]Row{}, query, 1)

		assert.Equal(t, context.Canceled, err)
		assert.Equal(t, fake.DBMetric.LastDBQueryErr, query)
		assert.Equal(t, fake.DBMetric.LastDBErr, err.Error())
	})

	t.Run("proxy get query", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		query := "SELECT * FROM anything WHERE id = $1"
		rows := sqlmock.NewRows([]string{"name"}).AddRow("diego")

		fake.DBMock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
		err := fake.DB().Get(ctx, &Row{}, query, 1)

		assert.NoError(t, err)
		assert.Equal(t, fake.DBMetric.LastQuery, query)
		assert.NotNil(t, fake.DBMetric.LastDuration)
	})

	t.Run("proxy get query with error", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		query := "SELECT * FROM anything WHERE id = $1"

		fake.DBMock.ExpectQuery(query).WithArgs(1).WillReturnError(context.Canceled)
		err := fake.DB().Get(ctx, &Row{}, query, 1)

		assert.Equal(t, context.Canceled, err)
		assert.Equal(t, fake.DBMetric.LastDBQueryErr, query)
		assert.Equal(t, fake.DBMetric.LastDBErr, err.Error())
	})

	t.Run("proxy exec query", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		query := "UPDATE anything SET code = $1 WHERE id = $2"
		result := sqlmock.NewResult(1, 1)

		fake.DBMock.ExpectExec(query).WithArgs("code", 1).WillReturnResult(result)
		err := fake.DB().Exec(ctx, query, "code", 1)

		assert.NoError(t, err)
		assert.Equal(t, fake.DBMetric.LastQuery, query)
		assert.NotNil(t, fake.DBMetric.LastDuration)
	})

	t.Run("proxy exec query with error", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		query := "UPDATE anything SET code = $1 WHERE id = $2"

		fake.DBMock.ExpectExec(query).WithArgs("code", 1).WillReturnError(context.Canceled)
		err := fake.DB().Exec(ctx, query, "code", 1)

		assert.Equal(t, context.Canceled, err)
		assert.Equal(t, fake.DBMetric.LastDBQueryErr, query)
		assert.Equal(t, fake.DBMetric.LastDBErr, err.Error())
	})

	t.Run("proxy start transaction with error", func(t *testing.T) {
		fake := test.NewFakeDependencies()

		fake.DBMock.ExpectBegin()
		_, err := fake.DB().BeginTx(ctx, nil)

		assert.NoError(t, err)
		assert.Equal(t, fake.DBMetric.LastQuery, "START TRANSACTION;")
		assert.NotNil(t, fake.DBMetric.LastDuration)
	})

	t.Run("proxy start transaction with error", func(t *testing.T) {
		fake := test.NewFakeDependencies()

		fake.DBMock.ExpectBegin().WillReturnError(context.Canceled)
		_, err := fake.DB().BeginTx(ctx, nil)

		assert.Equal(t, context.Canceled, err)
		assert.Equal(t, fake.DBMetric.LastDBQueryErr, "START TRANSACTION;")
		assert.NotNil(t, fake.DBMetric.LastDBErr, err.Error())
	})
}

func TestPostgresTx(t *testing.T) {
	type Row struct {
		Name string `db:"name"`
	}

	ctx := context.Background()

	t.Run("proxy select query", func(t *testing.T) {
		fake := test.NewFakeDependencies()

		fake.DBMock.ExpectBegin()
		tx, err := fake.DB().BeginTx(ctx, nil)
		assert.NoError(t, err)

		query := "SELECT * FROM anything WHERE id > $1"
		rows := sqlmock.NewRows([]string{"name"}).AddRow("diego").AddRow("maria")

		fake.DBMock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
		err = tx.Select(ctx, &[]Row{}, query, 1)

		assert.NoError(t, err)
		assert.Equal(t, fake.DBMetric.LastQuery, query)
		assert.NotNil(t, fake.DBMetric.LastDuration)
	})

	t.Run("proxy select query with error", func(t *testing.T) {
		fake := test.NewFakeDependencies()

		fake.DBMock.ExpectBegin()
		tx, err := fake.DB().BeginTx(ctx, nil)
		assert.NoError(t, err)

		query := "SELECT * FROM anything WHERE id > $1"

		fake.DBMock.ExpectQuery(query).WithArgs(1).WillReturnError(context.Canceled)
		err = tx.Select(ctx, &[]Row{}, query, 1)

		assert.Equal(t, context.Canceled, err)
		assert.Equal(t, fake.DBMetric.LastDBQueryErr, query)
		assert.Equal(t, fake.DBMetric.LastDBErr, err.Error())
	})

	t.Run("proxy get query", func(t *testing.T) {
		fake := test.NewFakeDependencies()

		fake.DBMock.ExpectBegin()
		tx, err := fake.DB().BeginTx(ctx, nil)
		assert.NoError(t, err)

		query := "SELECT * FROM anything WHERE id = $1"
		rows := sqlmock.NewRows([]string{"name"}).AddRow("diego")

		fake.DBMock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
		err = tx.Get(ctx, &Row{}, query, 1)

		assert.NoError(t, err)
		assert.Equal(t, fake.DBMetric.LastQuery, query)
		assert.NotNil(t, fake.DBMetric.LastDuration)
	})

	t.Run("proxy get query with error", func(t *testing.T) {
		fake := test.NewFakeDependencies()

		fake.DBMock.ExpectBegin()
		tx, err := fake.DB().BeginTx(ctx, nil)
		assert.NoError(t, err)

		query := "SELECT * FROM anything WHERE id = $1"

		fake.DBMock.ExpectQuery(query).WithArgs(1).WillReturnError(context.Canceled)
		err = tx.Get(ctx, &Row{}, query, 1)

		assert.Equal(t, context.Canceled, err)
		assert.Equal(t, fake.DBMetric.LastDBQueryErr, query)
		assert.Equal(t, fake.DBMetric.LastDBErr, err.Error())
	})

	t.Run("proxy exec query", func(t *testing.T) {
		fake := test.NewFakeDependencies()

		fake.DBMock.ExpectBegin()
		tx, err := fake.DB().BeginTx(ctx, nil)
		assert.NoError(t, err)

		query := "UPDATE anything SET code = $1 WHERE id = $2"
		result := sqlmock.NewResult(1, 1)

		fake.DBMock.ExpectExec(query).WithArgs("code", 1).WillReturnResult(result)
		err = tx.Exec(ctx, query, "code", 1)

		assert.NoError(t, err)
		assert.Equal(t, fake.DBMetric.LastQuery, query)
		assert.NotNil(t, fake.DBMetric.LastDuration)
	})

	t.Run("proxy exec query with error", func(t *testing.T) {
		fake := test.NewFakeDependencies()

		fake.DBMock.ExpectBegin()
		tx, err := fake.DB().BeginTx(ctx, nil)
		assert.NoError(t, err)

		query := "UPDATE anything SET code = $1 WHERE id = $2"

		fake.DBMock.ExpectExec(query).WithArgs("code", 1).WillReturnError(context.Canceled)
		err = tx.Exec(ctx, query, "code", 1)

		assert.Equal(t, context.Canceled, err)
		assert.Equal(t, fake.DBMetric.LastDBQueryErr, query)
		assert.Equal(t, fake.DBMetric.LastDBErr, err.Error())
	})

	t.Run("proxy commit query", func(t *testing.T) {
		fake := test.NewFakeDependencies()

		fake.DBMock.ExpectBegin()
		tx, err := fake.DB().BeginTx(ctx, nil)
		assert.NoError(t, err)

		fake.DBMock.ExpectCommit()
		assert.NoError(t, tx.Commit())

		assert.NotNil(t, fake.DBMetric.LastDuration)
		assert.Equal(t, fake.DBMetric.LastQuery, "COMMIT TRANSACTION;")
	})

	t.Run("proxy commit query with error", func(t *testing.T) {
		fake := test.NewFakeDependencies()

		fake.DBMock.ExpectBegin()
		tx, err := fake.DB().BeginTx(ctx, nil)
		assert.NoError(t, err)

		fake.DBMock.ExpectCommit().WillReturnError(context.Canceled)
		assert.Equal(t, context.Canceled, tx.Commit())

		assert.Equal(t, fake.DBMetric.LastDBQueryErr, "COMMIT TRANSACTION;")
		assert.Equal(t, context.Canceled.Error(), fake.DBMetric.LastDBErr)
	})

	t.Run("proxy rollback query", func(t *testing.T) {
		fake := test.NewFakeDependencies()

		fake.DBMock.ExpectBegin()
		tx, err := fake.DB().BeginTx(ctx, nil)
		assert.NoError(t, err)

		fake.DBMock.ExpectRollback()
		assert.NoError(t, tx.Rollback())

		assert.NotNil(t, fake.DBMetric.LastDuration)
		assert.Equal(t, fake.DBMetric.LastQuery, "ROLLBACK;")
	})

	t.Run("proxy rollback query with error", func(t *testing.T) {
		fake := test.NewFakeDependencies()

		fake.DBMock.ExpectBegin()
		tx, err := fake.DB().BeginTx(ctx, nil)
		assert.NoError(t, err)

		fake.DBMock.ExpectRollback().WillReturnError(context.Canceled)
		assert.Equal(t, context.Canceled, tx.Rollback())

		assert.Equal(t, fake.DBMetric.LastDBQueryErr, "ROLLBACK;")
		assert.Equal(t, context.Canceled.Error(), fake.DBMetric.LastDBErr)
	})
}
