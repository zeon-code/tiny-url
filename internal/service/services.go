package service

import (
	"log/slog"

	"github.com/zeon-code/tiny-url/internal/repository"
)

type Services struct {
	Url URLService
}

func NewServices(r repository.Repositories, l *slog.Logger) Services {
	return Services{
		Url: NewUrlService(r, l.With("service", "url")),
	}
}
