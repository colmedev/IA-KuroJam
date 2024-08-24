package api

import (
	"fmt"
	"net/http"
)

func (app *Api) LogError(r *http.Request, err error) {
	app.Logger.Error(
		err.Error(),
		"request_method", r.Method,
		"request_url", r.URL.String(),
	)
}

func (app *Api) ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := Envelope{"error": message}

	err := app.WriteJSON(w, status, env, nil)
	if err != nil {
		app.LogError(r, err)
		w.WriteHeader(500)
	}
}

func (app *Api) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.LogError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.ErrorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *Api) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.ErrorResponse(w, r, http.StatusNotFound, message)
}

func (app *Api) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.ErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *Api) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *Api) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.ErrorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *Api) EditConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.ErrorResponse(w, r, http.StatusConflict, message)
}

func (app *Api) RateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	app.ErrorResponse(w, r, http.StatusTooManyRequests, message)
}

func (app *Api) InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	app.ErrorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *Api) InvalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	app.ErrorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *Api) AuthenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	app.ErrorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *Api) InactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated to access this resource"
	app.ErrorResponse(w, r, http.StatusForbidden, message)
}

func (app *Api) NotPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account doesn't have the necessary permissions to access this resource"
	app.ErrorResponse(w, r, http.StatusForbidden, message)
}
