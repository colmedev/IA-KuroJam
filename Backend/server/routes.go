package server

import (
	"net/http"

	"github.com/colmedev/IA-KuroJam/Backend/api"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	app *api.Api
}

func (h *Handlers) routes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(h.app.Authenticate)
	// TODO: Add Panic Recover, RateLimiter, CORS

	router.Get("/v1/healthcheck", h.healthCheckHandler)
	router.Post("/v1/auth/login", h.authenticate)
	router.Post("/v1/auth/register", h.register)

	router.Group(func(protected chi.Router) {

		protected.Use(h.app.RequireAuthenticatedUser)

		protected.Get("/v1/protected_info", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Welcome to your info!"))
		})
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
