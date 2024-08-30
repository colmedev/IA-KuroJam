package api

import (
	"log/slog"
	"sync"

	"github.com/colmedev/IA-KuroJam/Backend/careers"
	"github.com/colmedev/IA-KuroJam/Backend/careertest"
	"github.com/colmedev/IA-KuroJam/Backend/llm"
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
	UserService       users.Service
	LlmService        llm.Service
	CareerTestService careertest.Service
	CareerService     careers.Service
}

// Passing services

type Option func(*Services)

func WithUserService(service users.Service) Option {
	return func(s *Services) {
		s.UserService = service
	}
}

func WithLlmService(service llm.Service) Option {
	return func(s *Services) {
		s.LlmService = service
	}
}

func WithCareerTestService(service careertest.Service) Option {
	return func(s *Services) {
		s.CareerTestService = service
	}
}

func WithCareerService(service careers.Service) Option {
	return func(s *Services) {
		s.CareerService = service
	}
}
