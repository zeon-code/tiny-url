package test

import (
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func (d FakeDependencies) MockUrlCreate() {
	now := time.Now()

	updateQuery := "UPDATE urls SET code = $1 WHERE id = $2"
	insertQuery := "INSERT INTO urls (target, code) VALUES ($1, '') RETURNING id, target, code, created_at, updated_at"

	rows := sqlmock.NewRows([]string{"id", "target", "code", "created_at", "updated_at"}).
		AddRow(int64(1), "target", "", now, now)

	d.DBMock.ExpectBegin()
	d.DBMock.ExpectQuery(insertQuery).WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
	d.DBMock.ExpectExec(updateQuery).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	d.DBMock.ExpectCommit()
}

func (d FakeDependencies) MockUrlList() {
	query := "SELECT id, code, target FROM urls ORDER BY id DESC LIMIT $1"

	rows := sqlmock.NewRows([]string{"id", "code", "target"}).
		AddRow(int64(5), "5", "target5").
		AddRow(int64(4), "4", "target4").
		AddRow(int64(3), "3", "target3").
		AddRow(int64(2), "2", "target2").
		AddRow(int64(1), "1", "target1")

	d.DBMock.ExpectQuery(query).WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
}

func (d FakeDependencies) MockPaginatedUrlList() {
	query := "SELECT id, code, target FROM urls WHERE id > $1 ORDER BY id DESC LIMIT $2"

	rows := sqlmock.NewRows([]string{"id", "code", "target"}).
		AddRow(int64(6), "6", "target6").
		AddRow(int64(5), "5", "target5").
		AddRow(int64(4), "4", "target4").
		AddRow(int64(3), "3", "target3").
		AddRow(int64(2), "2", "target2")

	d.DBMock.ExpectQuery(query).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(rows)
}

func (d FakeDependencies) MockUrlGetById() {
	at, _ := time.Parse(time.RFC3339, "2026-01-29T15:23:24Z")
	query := "SELECT * FROM urls WHERE id = $1"

	rows := sqlmock.NewRows([]string{"id", "code", "target", "created_at", "updated_at"}).
		AddRow(int64(1), "1", "target1", at, at)

	d.DBMock.ExpectQuery(query).WillReturnRows(rows)
}
