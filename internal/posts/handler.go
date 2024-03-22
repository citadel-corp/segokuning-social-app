package posts

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/citadel-corp/segokuning-social-app/internal/common/middleware"
	"github.com/citadel-corp/segokuning-social-app/internal/common/request"
	"github.com/citadel-corp/segokuning-social-app/internal/common/response"
	"github.com/gorilla/schema"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req CreatePostPayload
	var resp Response
	var err error

	userID, err := getUserID(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	req.UserID = userID

	err = request.DecodeJSON(w, r, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}

	err = req.Validate()
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Error: err.Error(),
		})
		return
	}

	resp = h.service.Create(r.Context(), req)
	response.JSON(w, resp.Code, response.ResponseBody{
		Message: resp.Message,
		Error:   resp.Error,
	})
	return
}

func (h *Handler) CreatePostComment(w http.ResponseWriter, r *http.Request) {
	var req CreatePostCommentPayload
	var resp Response
	var err error

	userID, err := getUserID(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	req.UserID = userID

	err = request.DecodeJSON(w, r, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
		})
	}
	err = req.Validate()
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Error: err.Error(),
		})
		return
	}

	resp = h.service.CreatePostComment(r.Context(), req)
	response.JSON(w, resp.Code, response.ResponseBody{
		Message: resp.Message,
		Error:   resp.Error,
	})
}

func (h *Handler) ListPost(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var req ListPostPayload

	newSchema := schema.NewDecoder()
	newSchema.IgnoreUnknownKeys(true)
	if err := newSchema.Decode(&req, r.URL.Query()); err != nil {
		slog.Error(err.Error())
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{})
		return
	}

	req.UserID = userID

	resp := h.service.List(r.Context(), req)
	response.JSON(w, resp.Code, response.ResponseBody{
		Message: resp.Message,
		Data:    resp.Data,
		Meta:    resp.Meta,
		Error:   resp.Error,
	})
}

func getUserID(r *http.Request) (string, error) {
	if authValue, ok := r.Context().Value(middleware.ContextAuthKey{}).(string); ok {
		return authValue, nil
	}

	return "", errors.New("unauthorized")
}
