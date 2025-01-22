package handlers

import (
	"encoding/json"
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
	type RequestBody struct {
		UserID string `json:"user_id"`
	}

	var requestBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		h.Logger.Error("Failed to decode request body: " + err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	honeytokenInteraction := models.NewInteraction(
		requestBody.UserID,
		"/api/honeytoken-endpoint",
		http.StatusOK,
		false,
		r.RemoteAddr,
	)

	if err := h.Neo4jService.SaveInteraction(honeytokenInteraction); err != nil {
		h.Logger.Error("Failed to log honeytoken access: " + err.Error())
		http.Error(w, "Failed to log honeytoken access", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Honeytoken access logged for user_id: " + requestBody.UserID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Honeytoken access detected and logged"))
}
