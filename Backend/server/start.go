package server

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/colmedev/IA-KuroJam/Backend/api"
	"github.com/colmedev/IA-KuroJam/Backend/careers"
	"github.com/colmedev/IA-KuroJam/Backend/careertest"
	"github.com/colmedev/IA-KuroJam/Backend/llm"
	"github.com/colmedev/IA-KuroJam/Backend/users"
)

var (
	version = "0.0.1"
)

func StartServer() error {
	var config api.Config

	flag.IntVar(&config.Port, "port", 8000, "API server port")
	flag.StringVar(&config.Env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&config.DB.DSN, "db-dsn", "", "PostgreSQL DSN")
	flag.IntVar(&config.DB.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&config.DB.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&config.DB.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Float64Var(&config.Limiter.RPS, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&config.Limiter.Burst, "limiter-burst", 4, "Rate limiter maxumim burst")
	flag.BoolVar(&config.Limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		config.Cors.TrustedOrigins = strings.Fields(val)
		return nil
	})

	flag.StringVar(&config.Auth.SigningKey, "signing-key", "abc123", "JWT Tokens Signing Key")
	flag.IntVar(&config.Auth.TokenExpirationInMinutes, "token-expiration", 15, "Token Expiration in Minutes")

	flag.StringVar(&config.LlmApiKey, "llm-api-key", "", "LLM API Key")

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version\t%s\n", version)
	}

	// Dependencies

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))

	logger.Info("Testing")

	db, err := api.OpenDb(config)
	if err != nil {
		return err
	}

	defer db.Close()

	logger.Info("database connection pool established")

	// Services
	usersService, err := users.NewService(db)
	if err != nil {
		return fmt.Errorf("error initializing user service: %w", err)
	}

	// TODO: Add services
	llmService := llm.NewOpenAIService(config.LlmApiKey)
	careerTestService := careertest.NewService(db, llmService)
	careerService := careers.NewCareerService(db)

	// Application
	app := api.NewApplication(
		config,
		db,
		logger,
		api.WithUserService(usersService),
		api.WithLlmService(llmService),
		api.WithCareerTestService(careerTestService),
		api.WithCareerService(careerService),
	)

	// Handlers

	h := &Handlers{
		app: app,
	}

	err = h.Serve()
	if err != nil {
		return err
	}

	return nil

}
