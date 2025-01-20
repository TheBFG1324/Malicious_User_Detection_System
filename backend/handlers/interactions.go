package handlers

import (
	"encoding/json"
	"net/http"

	"backend/models"
	"backend/services"
	"backend/utils"
)

// InteractionHandler handles interaction-related requests
type InteractionHandler struct {
	Neo4jService *services.Neo4jService
	Logger       *utils.Logger
}

// NewInteractionHandler creates a new InteractionHandler
func NewInteractionHandler(neo4jService *services.Neo4jService, logger *utils.Logger) *InteractionHandler {
	return &InteractionHandler{
		Neo4jService: neo4jService,
		Logger:       logger,
	}
}

// LogInteraction handles requests to log user interactions
func (h *InteractionHandler) LogInteraction(w http.ResponseWriter, r *http.Request) {
	var interaction models.Interaction

	if err := json.NewDecoder(r.Body).Decode(&interaction); err != nil {
		h.Logger.Error("Invalid request body: " + err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Neo4jService.SaveInteraction(interaction); err != nil {
		h.Logger.Error("Failed to save interaction: " + err.Error())
		http.Error(w, "Failed to save interaction", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Interaction logged successfully for user_id: " + interaction.UserID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Interaction logged successfully"))
}
