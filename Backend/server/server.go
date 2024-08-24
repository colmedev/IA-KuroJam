package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (h *Handlers) Serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", h.app.Config.Port),
		Handler:      h.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		h.app.Logger.Info("shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		h.app.Logger.Info("completing background tasks", "addr", srv.Addr)

		h.app.Wg.Wait()

		shutdownError <- nil
	}()

	h.app.Logger.Info(
		"starting server",
		"addr", srv.Addr,
		"env", h.app.Config.Env,
	)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	h.app.Logger.Info("stopped server", "addr", srv.Addr)

	return nil
}
