package user

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

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserPayload

	err := request.DecodeJSON(w, r, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}
	userResp, err := h.service.Create(r.Context(), req)
	if errors.Is(err, ErrUserPhoneNumberAlreadyExists) || errors.Is(err, ErrUserEmailAlreadyExists) {
		response.JSON(w, http.StatusConflict, response.ResponseBody{
			Message: "User already exists",
			Error:   err.Error(),
		})
		return
	}
	if errors.Is(err, ErrValidationFailed) {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{
			Message: "Internal server error",
			Error:   err.Error(),
		})
		return
	}
	response.JSON(w, http.StatusCreated, response.ResponseBody{
		Message: "User registered successfully",
		Data:    userResp,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginPayload

	err := request.DecodeJSON(w, r, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}
	userResp, err := h.service.Login(r.Context(), req)
	if errors.Is(err, ErrUserNotFound) {
		response.JSON(w, http.StatusNotFound, response.ResponseBody{
			Message: "Not found",
			Error:   err.Error(),
		})
		return
	}
	if errors.Is(err, ErrWrongPassword) {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}
	if errors.Is(err, ErrValidationFailed) {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{
			Message: "Internal server error",
			Error:   err.Error(),
		})
		return
	}
	response.JSON(w, http.StatusOK, response.ResponseBody{
		Message: "User logged successfully",
		Data:    userResp,
	})
}

func (h *Handler) LinkEmail(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		response.JSON(w, http.StatusUnauthorized, response.ResponseBody{
			Message: "Unauthorized",
			Error:   err.Error(),
		})
		return
	}
	var req LinkEmailPayload

	err = request.DecodeJSON(w, r, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}
	err = h.service.LinkEmail(r.Context(), req, userID)
	if errors.Is(err, ErrValidationFailed) {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}
	if errors.Is(err, ErrUserHasEmail) {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}
	if errors.Is(err, ErrUserEmailAlreadyExists) {
		response.JSON(w, http.StatusConflict, response.ResponseBody{
			Message: "Conflict",
			Error:   err.Error(),
		})
		return
	}
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{
			Message: "Internal server error",
			Error:   err.Error(),
		})
		return
	}
	response.JSON(w, http.StatusOK, response.ResponseBody{
		Message: "User email linked successfully",
	})
}

func (h *Handler) LinkPhoneNumber(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		response.JSON(w, http.StatusUnauthorized, response.ResponseBody{
			Message: "Unauthorized",
			Error:   err.Error(),
		})
		return
	}
	var req LinkPhoneNumberPayload

	err = request.DecodeJSON(w, r, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}
	err = h.service.LinkPhoneNumber(r.Context(), req, userID)
	if errors.Is(err, ErrValidationFailed) {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}
	if errors.Is(err, ErrUserHasPhoneNumber) {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}
	if errors.Is(err, ErrUserPhoneNumberAlreadyExists) {
		response.JSON(w, http.StatusConflict, response.ResponseBody{
			Message: "Conflict",
			Error:   err.Error(),
		})
		return
	}
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{
			Message: "Internal server error",
			Error:   err.Error(),
		})
		return
	}
	response.JSON(w, http.StatusOK, response.ResponseBody{
		Message: "Successfully linked phone number",
	})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		response.JSON(w, http.StatusUnauthorized, response.ResponseBody{
			Message: "Unauthorized",
			Error:   err.Error(),
		})
		return
	}
	var req UpdateUserPayload

	err = request.DecodeJSON(w, r, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}
	err = h.service.Update(r.Context(), req, userID)
	if errors.Is(err, ErrValidationFailed) {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{
			Message: "Internal server error",
			Error:   err.Error(),
		})
		return
	}
	response.JSON(w, http.StatusOK, response.ResponseBody{
		Message: "Successfully update user",
	})
}

func (h *Handler) ListUser(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		response.JSON(w, http.StatusUnauthorized, response.ResponseBody{
			Message: "Unauthorized",
			Error:   err.Error(),
		})
		return
	}

	var req ListUserPayload

	newSchema := schema.NewDecoder()
	newSchema.IgnoreUnknownKeys(true)
	if err := newSchema.Decode(&req, r.URL.Query()); err != nil {
		slog.Error(err.Error())
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{})
		return
	}

	req.UserID = userID

	usersResp, pagination, err := h.service.List(r.Context(), req)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{
			Message: "Internal server error",
			Error:   err.Error(),
		})
		return
	}
	response.JSON(w, http.StatusOK, response.ResponseBody{
		Message: "Users fetched successfully",
		Data:    usersResp,
		Meta:    pagination,
	})
}

func getUserID(r *http.Request) (string, error) {
	if authValue, ok := r.Context().Value(middleware.ContextAuthKey{}).(string); ok {
		return authValue, nil
	}
	slog.Error("cannot parse auth value from context")
	return "", errors.New("cannot parse auth value from context")
}
