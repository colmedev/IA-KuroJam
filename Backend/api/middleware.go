package api

import (
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/colmedev/IA-KuroJam/Backend/users"
	"github.com/felixge/httpsnoop"
	"github.com/golang-jwt/jwt"
	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

func (app *Api) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				app.ServerErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *Api) RateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()

		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.Config.Limiter.Enabled {
			ip := realip.FromRequest(r)

			mu.Lock()

			if _, found := clients[ip]; !found {
				clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(app.Config.Limiter.RPS), app.Config.Limiter.Burst)}
			}

			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.RateLimitExceededResponse(w, r)
				return
			}

			mu.Unlock()

		}
		next.ServeHTTP(w, r)
	})
}

func (app *Api) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			r = app.ContextSetUser(r, &users.AnonymousUser)
			next.ServeHTTP(w, r)

			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.Logger.Debug("Empty token")
			app.InvalidCredentialsResponse(w, r)

			return
		}

		tokenString := headerParts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {

				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return app.Config.Auth.SigningKey, nil
		})

		if err != nil || !token.Valid {

			if err.Error() != "Token is not valid yet" {
				fmt.Println(err.Error())
				app.Logger.Debug("invalid token", err)
				app.InvalidAuthenticationTokenResponse(w, r)
				return
			}
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {

			if err.Error() != "Token is not valid yet" {
				app.Logger.Debug("invalid claims", err)
				app.InvalidAuthenticationTokenResponse(w, r)
				return
			}
		}

		userID := claims["sub"].(string)

		user, err := app.Services.UserService.FindByGoogleId(r.Context(), userID)

		if err != nil {
			if errors.Is(err, users.ErrRecordNotFound) {
				user := &users.User{
					ClerkId: userID,
				}

				err, _ = app.Services.UserService.Insert(r.Context(), user)
				if err != nil {
					fmt.Println("Here inside", err)

					app.Logger.Debug("error creating user", err)
					app.InvalidCredentialsResponse(w, r)
					return
				}
			} else {
				fmt.Println("Here inside 2", err)

				app.Logger.Debug("error getting user", err)
				app.InvalidCredentialsResponse(w, r)
				return
			}
		}

		r = app.ContextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

// Create a new requireAuthenticatedUser() middleware to check that a user is not anonymous.
func (app *Api) RequireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.ContextGetUser(r)

		if user.IsAnonymous() {
			app.AuthenticationRequiredResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// func (app *Api) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		user := app.contextGetUser(r)
//
// 		permissions, err := app.Services.Permissions.GetAllForUser(user.ID)
// 		if err != nil {
// 			app.ServerErrorResponse(w, r, err)
// 			return
// 		}
//
// 		if !permissions.Include(code) {
// 			app.NotPermittedResponse(w, r)
// 			return
// 		}
//
// 		next.ServeHTTP(w, r)
// 	}
//
// 	return app.requireActivatedUser(fn)
// }

func (app *Api) EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")

		w.Header().Add("Vary", "Origin")

		origin := r.Header.Get("Origin")

		if origin != "" {
			for i := range app.Config.Cors.TrustedOrigins {
				if origin == app.Config.Cors.TrustedOrigins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)

					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

						w.WriteHeader(http.StatusOK)
						return
					}
					break
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (app *Api) Metrics(next http.Handler) http.Handler {
	totalRequestsReceived := expvar.NewInt("total_requests_received")
	totalResponsesSent := expvar.NewInt("total_responses_sent")
	totalProcessingTimeMicroseconds := expvar.NewInt("total_processing_time_us")
	totalResponsesSentByStatus := expvar.NewMap("total_responses_sent_by_status")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		totalRequestsReceived.Add(1)

		metrics := httpsnoop.CaptureMetrics(next, w, r)

		totalResponsesSent.Add(1)

		totalProcessingTimeMicroseconds.Add(metrics.Duration.Microseconds())

		totalResponsesSentByStatus.Add(strconv.Itoa(metrics.Code), 1)
	})
}
