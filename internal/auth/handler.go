package auth

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	// Check method
	// If the method does not match the expected method POST
	// return 405 Status Not Allowed
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode JSON request body into RegisterRequest struct
	// If the request body is invalid JSON or does not match the expected format
	// return 400 Bad Request
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if email and password are nil
	// If one of them is nil
	// return 400 Bad Request
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password required", http.StatusBadRequest)
		return
	}

	// Call service layer to register a new user
	// return 400 Bad Request if registration fails (e.g. email already exists)
	if err := h.Service.Register(req.Email, req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Response with 201 Created to indicate successful user registration
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	// Check method
	// If the method does not match the expected method POST
	// return 405 Status Not Allowed
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode JSON request body into LoginRequest struct
	// If the request body is invalid or does not match the expected format
	// return 400 Bad Request
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Authentication user credentials via service layer
	// return 401 Unauthorized if email or password is invalid
	token, err := h.Service.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Prepare login response containing the generated authentication token
	resp := LoginResponse{Token: token}

	// Set response content type to JSON and send the login response to client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get the Bearer token via the function in token.go
	tokenString := GetTokenFromHeader(r)

	// Return 401 Status Unauthorized if the token is missing
	if tokenString == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	// Parse and validate the authentication token
	// Return 401 Status Unauthorized if the token is invalid or expired
	claims, err := ParseToken(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Call service layer to revoke the current authentication token to invalidate it
	// Return 500 Internal Server Error if token revocation fails
	if err := h.Service.RevokeToken(tokenString, claims.ExpiresAt.Time); err != nil {
		http.Error(w, "Logout failed", http.StatusInternalServerError)
		return
	}

	// Response with 200 OK to indicate successful log out
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out"))
}
