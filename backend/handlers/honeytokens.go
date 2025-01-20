package handlers

import (
	"net/http"

	"backend/models"
	"backend/services"
	"backend/utils"
)

// HoneytokenHandler handles honeytoken-related requests
type HoneytokenHandler struct {
	Neo4jService *services.Neo4jService
	Logger       *utils.Logger
}

// NewHoneytokenHandler creates a new HoneytokenHandler
func NewHoneytokenHandler(neo4jService *services.Neo4jService, logger *utils.Logger) *HoneytokenHandler {
	return &HoneytokenHandler{
		Neo4jService: neo4jService,
		Logger:       logger,
	}
}

// DetectHoneytoken detects and logs honeytoken access
func (h *HoneytokenHandler) DetectHoneytoken(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	ipAddress := r.RemoteAddr

	if userID == "" {
		h.Logger.Error("Missing user_id in request")
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	honeytokenInteraction := models.NewInteraction(
		userID,
		"/honeytoken/trap",
		http.StatusForbidden,
		true,
		ipAddress,
	)

	if err := h.Neo4jService.SaveInteraction(honeytokenInteraction); err != nil {
		h.Logger.Error("Failed to log honeytoken access: " + err.Error())
		http.Error(w, "Failed to log honeytoken access", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Honeytoken access logged for user_id: " + userID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Honeytoken access detected and logged"))
}
