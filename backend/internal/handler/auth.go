package handler

import (
	"encoding/json"
	"net/http"

	"github.com/entaku0818/aso-compass/backend/internal/middleware"
	"github.com/entaku0818/aso-compass/backend/internal/model"
	"github.com/entaku0818/aso-compass/backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
		User  struct {
			*model.User
			IsPro bool `json:"is_pro"`
		} `json:"user"`
	}{
		Token: resp.Token,
		User: struct {
			*model.User
			IsPro bool `json:"is_pro"`
		}{User: resp.User, IsPro: resp.User.IsPro()},
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.authService.GetUserByID(r.Context(), userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, struct {
		*model.User
		IsPro bool `json:"is_pro"`
	}{User: user, IsPro: user.IsPro()})
}

func (h *AuthHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.authService.CreateUser(r.Context(), &req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if _, err := h.authService.CreateUser(r.Context(), &req); err != nil {
		handleServiceError(w, err)
		return
	}

	resp, err := h.authService.Login(r.Context(), &model.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, struct {
		Token string `json:"token"`
		User  struct {
			*model.User
			IsPro bool `json:"is_pro"`
		} `json:"user"`
	}{
		Token: resp.Token,
		User: struct {
			*model.User
			IsPro bool `json:"is_pro"`
		}{User: resp.User, IsPro: resp.User.IsPro()},
	})
}
