package userfriends

import (
	"errors"
	"net/http"

	"github.com/citadel-corp/segokuning-social-app/internal/common/middleware"
	"github.com/citadel-corp/segokuning-social-app/internal/common/request"
	"github.com/citadel-corp/segokuning-social-app/internal/common/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateUserFriends(w http.ResponseWriter, r *http.Request) {
	var req CreateUserFriendPayload
	var err error

	userID, err := getUserID(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	req.LoggedUserID = userID

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

	resp := h.service.Create(r.Context(), req)
	response.JSON(w, resp.Code, response.ResponseBody{
		Message: resp.Message,
		Error:   resp.Error,
	})
	return
}

func (h *Handler) DeleteUserFriends(w http.ResponseWriter, r *http.Request) {
	var req DeleteUserFriendPayload
	var err error

	userID, err := getUserID(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	req.LoggedUserID = userID

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

	resp := h.service.Delete(r.Context(), req)
	response.JSON(w, resp.Code, response.ResponseBody{
		Message: resp.Message,
		Error:   resp.Error,
	})
	return
}

func getUserID(r *http.Request) (string, error) {
	if authValue, ok := r.Context().Value(middleware.ContextAuthKey{}).(string); ok {
		return authValue, nil
	}

	return "", errors.New("unauthorized")
}
