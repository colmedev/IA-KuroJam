package server

import (
	"net/http"

	"github.com/colmedev/IA-KuroJam/Backend/api"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Handlers struct {
	app *api.Api
}

func (h *Handlers) routes() http.Handler {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"}, // Use this to allow specific origin hosts
		// AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(h.app.Authenticate)
	// TODO: Add Panic Recover, RateLimiter, CORS

	router.Get("/v1/healthcheck", h.healthCheckHandler)
	// router.Post("/v1/auth/login", h.authenticate)
	// router.Post("/v1/auth/register", h.register)

	router.Group(func(protected chi.Router) {

		protected.Use(h.app.RequireAuthenticatedUser)

		protected.Post("/v1/start-test", h.startTest)
		protected.Get("/v1/questions/{id}", h.getQuestion)
		protected.Post("/v1/answer/{id}", h.postAnswer)
		protected.Get("/v1/results", h.getResults)
		protected.Get("/v1/get-test", h.getActiveTest)
	})

	return router
}

func (h *Handlers) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	env := api.Envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": "Testing",
			"version":     "1.0.0",
		},
	}

	err := h.app.WriteJSON(w, http.StatusOK, env, nil)
	if err != nil {
		h.app.ServerErrorResponse(w, r, err)
	}
}
