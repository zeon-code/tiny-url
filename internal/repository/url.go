package repository

import (
	"context"
	"fmt"

	"github.com/zeon-code/tiny-url/internal/db"
	"github.com/zeon-code/tiny-url/internal/model"
	"github.com/zeon-code/tiny-url/internal/pkg/base62"
	"github.com/zeon-code/tiny-url/internal/pkg/observability"
)

type URLRepository interface {
	Create(context.Context, string) (*model.URL, error)
	List(context.Context, int, string, *int64) ([]model.URL, error)
	GetByID(context.Context, int64) (*model.URL, error)
}

type URLStore struct {
	db     db.SQLClient
	memory db.SQLReader
	logger observability.Logger
}

func NewURLRepository(database db.SQLClient, memory db.SQLReader, observer observability.Observer) URLRepository {
	return URLStore{
		db:     database,
		memory: memory,
		logger: observer.Logger().WithGroup("url-repository"),
	}
}

func (s URLStore) Create(ctx context.Context, target string) (*model.URL, error) {
	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, err
	}

	var url model.URL
	query := "INSERT INTO urls (target, code) VALUES ($1, '') RETURNING id, target, code, created_at, updated_at"

	if err := tx.Get(ctx, &url, query, target); err != nil {
		tx.Rollback()
		return nil, err
	}

	query = "UPDATE urls SET code = $1 WHERE id = $2"
	url.Code = base62.Encode(url.ID)

	if err := tx.Exec(ctx, query, url.Code, url.ID); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	return &url, nil
}

func (s URLStore) List(ctx context.Context, limit int, direction string, cursor *int64) ([]model.URL, error) {
	var err error
	urls := []model.URL{}
	query := "SELECT id, code, target FROM urls"

	if cursor != nil {
		query = fmt.Sprintf("%s WHERE id %s $1 ORDER BY id DESC LIMIT $2", query, direction)
		err = s.memory.Select(ctx, &urls, query, cursor, limit)
	} else {
		query = fmt.Sprintf("%s ORDER BY id DESC LIMIT $1 ", query)
		err = s.memory.Select(ctx, &urls, query, limit)
	}

	if err != nil {
		return urls, err
	}

	return urls, nil
}

func (s URLStore) GetByID(ctx context.Context, id int64) (*model.URL, error) {
	var url model.URL
	query := "SELECT * FROM urls WHERE id = $1"

	if err := s.memory.Get(ctx, &url, query, id); err != nil {
		return nil, err
	}

	return &url, nil
}
