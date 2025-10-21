package handlers

import (
	"encoding/json"
	"net/http"

	"main-services/internal/app/words"
)

// Inject the service
type WordHandler struct{ service words.Service }

func NewWordHandler(service words.Service) *WordHandler { return &WordHandler{service: service} }

// Helper to get current user ID from context (set by auth middleware)
func currentUserID(request *http.Request) string {
	// TODO: Get from context when auth middleware is implemented
	// For now, return empty string or a test user ID
	return "" // This will require auth middleware to work properly
}

func (handler *WordHandler) Create(writer http.ResponseWriter, request *http.Request) {
	ownerID := currentUserID(request)
	if ownerID == "" {
		// TODO: Remove this workaround when auth is implemented
		ownerID = "test-user" // Temporary for testing
	}

	var in words.CreateWordReq
	if err := json.NewDecoder(request.Body).Decode(&in); err != nil {
		http.Error(writer, "invalid json", http.StatusBadRequest)
		return
	}

	id, err := handler.service.Create(request.Context(), ownerID, in)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(writer, http.StatusCreated, map[string]any{"id": id})
}

// Minimal JSON writer (use your shared util if you have one)
func writeJSON(writer http.ResponseWriter, code int, v any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	_ = json.NewEncoder(writer).Encode(v)
}
