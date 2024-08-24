package api

import (
	"log/slog"
	"sync"

	"github.com/colmedev/IA-KuroJam/Backend/users"
	"github.com/jmoiron/sqlx"
)

type Api struct {
	Config   Config
	Logger   *slog.Logger
	Wg       sync.WaitGroup
	Services *Services
}

func NewApplication(cfg Config, db *sqlx.DB, logger *slog.Logger, options ...Option) *Api {

	s := &Services{}
	for _, option := range options {
		option(s)
	}

	return &Api{
		Logger:   logger,
		Config:   cfg,
		Services: s,
	}
}

type Services struct {
	UserService users.Service
}

// Passing services

type Option func(*Services)

func WithUserService(service users.Service) Option {
	return func(s *Services) {
		s.UserService = service
	}
}
