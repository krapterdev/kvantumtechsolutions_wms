package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

type Handler struct {
	service        *Service
	sessionManager *SessionManager
}

func NewHandler(service *Service, sessionManager *SessionManager) *Handler {
	return &Handler{
		service:        service,
		sessionManager: sessionManager,
	}
}

type loginRequest struct {
	OrganizationID string `json:"organization_id"`
	Email          string `json:"email"`
	Password       string `json:"password"`
}

type loginResponse struct {
	User struct {
		ID             string `json:"id"`
		OrganizationID string `json:"organization_id"`
		Email          string `json:"email"`
		FirstName      string `json:"first_name"`
	} `json:"user"`
}

type createUserRequest struct {
	OrganizationID string  `json:"organization_id"`
	Email          string  `json:"email"`
	Password       string  `json:"password"`
	FirstName      string  `json:"first_name"`
	LastName       *string `json:"last_name"`
}

type createUserResponse struct {
	User struct {
		ID             string `json:"id"`
		OrganizationID string `json:"organization_id"`
		Email          string `json:"email"`
		FirstName      string `json:"first_name"`
		LastName       string `json:"last_name,omitempty"`
		IsActive       bool   `json:"is_active"`
	} `json:"user"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if h.service == nil || h.sessionManager == nil {
		http.Error(w, "authentication service unavailable", http.StatusInternalServerError)
		return
	}

	var req loginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.OrganizationID == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "organization_id, email and password are required", http.StatusBadRequest)
		return
	}

	var organizationID pgtype.UUID

	if err := organizationID.Scan(req.OrganizationID); err != nil {
		http.Error(w, "invalid organization_id", http.StatusBadRequest)
		return
	}

	result, err := h.service.Login(r.Context(), h.sessionManager, LoginRequest{
		OrganizationID: organizationID,
		Email:          req.Email,
		Password:       req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
		case errors.Is(err, ErrUserInactive):
			http.Error(w, "user is inactive", http.StatusForbidden)
		default:
			http.Error(w, "authentication failed", http.StatusInternalServerError)
		}

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    result.SessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(sessionDuration.Seconds()),
	})

	loginResponse := loginResponse{}
	loginResponse.User.ID = result.User.ID.String()
	loginResponse.User.OrganizationID = result.User.OrganizationID.String()
	loginResponse.User.Email = result.User.Email
	loginResponse.User.FirstName = result.User.FirstName

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(loginResponse)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		http.Error(w, "authentication service unavailable", http.StatusInternalServerError)
		return
	}

	var req createUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.FirstName = strings.TrimSpace(req.FirstName)

	if req.LastName != nil {
		lastName := strings.TrimSpace(*req.LastName)
		req.LastName = &lastName
	}

	if req.OrganizationID == "" ||
		req.Email == "" ||
		req.Password == "" ||
		req.FirstName == "" {
		http.Error(
			w,
			"organization_id, email, password and first_name are required",
			http.StatusBadRequest,
		)
		return
	}

	var organizationID pgtype.UUID

	if err := organizationID.Scan(req.OrganizationID); err != nil {
		http.Error(w, "invalid organization_id", http.StatusBadRequest)
		return
	}

	user, err := h.service.CreateUser(r.Context(), CreateUserRequest{
		OrganizationID: organizationID,
		Email:          req.Email,
		Password:       req.Password,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
	})
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	response := createUserResponse{}
	response.User.ID = user.ID.String()
	response.User.OrganizationID = user.OrganizationID.String()
	response.User.Email = user.Email
	response.User.FirstName = user.FirstName
	response.User.IsActive = user.IsActive

	if user.LastName.Valid {
		response.User.LastName = user.LastName.String
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if h.sessionManager == nil {
		http.Error(w, "authentication service unavailable", http.StatusInternalServerError)
		return
	}

	cookie, err := r.Cookie("session")
	if err == nil && cookie.Value != "" {
		session, sessionErr := h.sessionManager.GetSessionByToken(r.Context(), cookie.Value)
		if sessionErr == nil {
			if revokeErr := h.sessionManager.RevokeSession(r.Context(), session.ID); revokeErr != nil {
				http.Error(w, "logout failed", http.StatusInternalServerError)
				return
			}
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "logged out successfully",
	})
}
