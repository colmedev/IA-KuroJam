package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/colmedev/IA-KuroJam/Backend/api"
	"github.com/colmedev/IA-KuroJam/Backend/careertest"
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) startTest(w http.ResponseWriter, r *http.Request) {
	user := h.app.ContextGetUser(r)

	ct := &careertest.CareerTest{
		UserId: user.Id,
	}

	err := h.app.Services.CareerTestService.StartTest(r.Context(), ct)
	if err != nil {
		h.app.ServerErrorResponse(w, r, err)
		return
	}

	data := api.Envelope{
		"careerTest": ct,
	}

	err = h.app.WriteJSON(w, 201, data, nil)
	if err != nil {
		h.app.ServerErrorResponse(w, r, err)
	}
}

func (h *Handlers) getQuestion(w http.ResponseWriter, r *http.Request) {
	user := h.app.ContextGetUser(r)
	careerTestIdStr := chi.URLParam(r, "id")

	careerTestId, err := strconv.ParseInt(careerTestIdStr, 10, 64)
	if err != nil {
		h.app.BadRequestResponse(w, r, fmt.Errorf("invalid path value"))
		return
	}

	msg, err := h.app.Services.CareerTestService.GetQuestion(r.Context(), careerTestId, user.Id)
	if err != nil {
		switch {
		case errors.Is(err, careertest.ErrNotPermission):
			h.app.NotPermittedResponse(w, r)
			return
		default:
			h.app.ServerErrorResponse(w, r, err)
			return
		}
	}

	data := api.Envelope{
		"message": msg,
	}

	err = h.app.WriteJSON(w, 200, data, nil)
	if err != nil {
		h.app.ServerErrorResponse(w, r, err)
		return
	}
}

func (h *Handlers) postAnswer(w http.ResponseWriter, r *http.Request) {
	user := h.app.ContextGetUser(r)
	careerTestIdStr := chi.URLParam(r, "id")

	careerTestId, err := strconv.ParseInt(careerTestIdStr, 10, 64)
	if err != nil {
		h.app.BadRequestResponse(w, r, fmt.Errorf("invalid path value"))
		return
	}

	var input struct {
		Answer string `json:"answer"`
	}

	err = h.app.ReadJSON(w, r, &input)
	if err != nil {
		h.app.BadRequestResponse(w, r, err)
		return
	}

	if input.Answer == "" {
		h.app.BadRequestResponse(w, r, fmt.Errorf("answer is empty"))
		return
	}

	msg, err := h.app.Services.CareerTestService.PostAnswer(r.Context(), input.Answer, careerTestId, user.Id)
	if err != nil {
		switch {
		case errors.Is(err, careertest.ErrRecordNotFound):
			h.app.NotFoundResponse(w, r)
			return
		default:
			h.app.ServerErrorResponse(w, r, err)
			return
		}
	}

	data := api.Envelope{
		"message": msg,
	}

	err = h.app.WriteJSON(w, 200, data, nil)
	if err != nil {
		h.app.ServerErrorResponse(w, r, err)
	}
}

func (h *Handlers) getResults(w http.ResponseWriter, r *http.Request) {
	user := h.app.ContextGetUser(r)

	emb, err := h.app.Services.CareerTestService.GetResultsEmbedding(r.Context(), user.Id)
	if err != nil {
		h.app.ServerErrorResponse(w, r, err)
		return
	}

	fmt.Println(h.app.Services.CareerService)
	careers, err := h.app.Services.CareerService.GetSimilarity(r.Context(), emb)
	if err != nil {
		h.app.ServerErrorResponse(w, r, err)
		return
	}

	data := api.Envelope{
		"careers": careers,
	}

	err = h.app.WriteJSON(w, 200, data, nil)
	if err != nil {
		h.app.ServerErrorResponse(w, r, err)
	}
}

func (h *Handlers) getActiveTest(w http.ResponseWriter, r *http.Request) {
	user := h.app.ContextGetUser(r)

	careerTest, err := h.app.Services.CareerTestService.GetActiveTest(r.Context(), user.Id)
	if err != nil {
		switch {
		case errors.Is(err, careertest.ErrRecordNotFound):
			h.app.NotFoundResponse(w, r)
			return
		default:
			h.app.ServerErrorResponse(w, r, err)
			return
		}
	}

	data := api.Envelope{
		"careerTest": careerTest,
	}

	err = h.app.WriteJSON(w, 200, data, nil)
	if err != nil {
		h.app.ServerErrorResponse(w, r, err)
	}

}
