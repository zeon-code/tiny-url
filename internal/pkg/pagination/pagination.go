package pagination

import (
	"bytes"
	"fmt"
	"net/http"

	json "github.com/json-iterator/go"
	"github.com/zeon-code/tiny-url/internal/pkg/base62"
)

type CursorKey[T any] func(T) int64

type Page struct {
	Next     string `json:"next,omitempty"`
	Previous string `json:"previous,omitempty"`
	Size     int    `json:"size"`
}

type Pagination[T any] struct {
	Limit  int    `json:"-"`
	Cursor *int64 `json:"-"`
	Items  []T    `json:"items"`
	Page   Page   `json:"page"`
}

func NewPagination[T any](items []T, limit int, cursor *int64) Pagination[T] {
	return Pagination[T]{
		Items:  items,
		Limit:  limit,
		Cursor: cursor,
	}
}

func (p Pagination[T]) Encode(cursorKey CursorKey[T]) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	p.Page = Page{Size: len(p.Items)}

	enc.SetEscapeHTML(false)

	if len(p.Items) > 0 {
		if p.Limit <= len(p.Items) {
			last := p.Items[len(p.Items)-1]
			p.Page.Next = fmt.Sprintf("<%s", base62.Encode(cursorKey(last)))
		}

		if p.Cursor != nil {
			first := p.Items[0]
			p.Page.Previous = fmt.Sprintf(">%s", base62.Encode(cursorKey(first)))
		}
	}

	err := enc.Encode(p)
	return buf.Bytes(), err
}

func GetCursor(r *http.Request) (string, *int64) {
	if value := r.URL.Query().Get("cursor"); value != "" {
		direction := value[0:1]
		c := base62.Decode(value[1:])

		if direction != "<" && direction != ">" {
			direction = "<"
		}

		return direction, &c
	}

	return "<", nil
}
