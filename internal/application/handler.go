package application

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/quynhanh/internship-tracker/internal/auth"
)

type Handler struct {
	Service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{Service: s}
}

// Create
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	// Get userID from context via middleware layer
	userID := auth.GetUserID(r.Context())

	// Decode JSON request body into CreateRequest struct
	// If the request body is invalid JSON or does not match the expected format
	// return 400 Bad Request
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	// Call service layer to create a new application
	// return 400 Bad Request if creation fails
	if err := h.Service.Create(userID, req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Response with 201 Created to indicate successful application creation
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Application created successfully"))
}

// List
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	// Get userID from context via middleware layer
	userID := auth.GetUserID(r.Context())

	// Parse pagination parameters from query string
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	// Optional filter by application status
	status := r.URL.Query().Get("status")

	// Call service layer by application status
	// Return 500 Status Internal Server Error if failed to fetch application
	apps, total, err := h.Service.List(userID, page, limit, status)
	if err != nil {
		http.Error(w, "Failed to fetch application", http.StatusInternalServerError)
		return
	}

	// Construct response payload with data and pagination metadata
	resp := map[string]interface{}{
		"data": apps,
		"meta": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}

	// Set response content type to JSON and send the response to client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Update
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	// Get userID from context via middleware layer
	userID := auth.GetUserID(r.Context())

	// Get application ID and convert it into int
	vars := mux.Vars(r)
	idStr := vars["applicationID"]
	appID, _ := strconv.Atoi(idStr)

	// Decode JSON request body into UpdateStatusRequest struct
	var req UpdateStatusRequest
	json.NewDecoder(r.Body).Decode(&req)

	// Call service layer to update the application
	// return 404 Status Not Found if the resource is not found
	if err := h.Service.UpdateStatus(userID, uint(appID), req.Status); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Response with 200 OK to indicate successful application update
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Application updated succsesfully"))
}

// Delete
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	// Get userID from context via middleware layer
	userID := auth.GetUserID(r.Context())

	// Get application ID and convert it into int
	vars := mux.Vars(r)
	idStr := vars["applicationID"]
	appID, _ := strconv.Atoi(idStr)

	// Call service layer to delete the application
	// return 404 Status Not Found if the resource is not found
	if err := h.Service.Delete(userID, uint(appID)); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Response with 200 OK to indicate successful application deletion
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Application deleted successfully"))
}
