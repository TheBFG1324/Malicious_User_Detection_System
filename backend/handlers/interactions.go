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

	type RequestBody struct {
		UserID string `json:"user_id"`
	}

	var requestBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		h.Logger.Error("Failed to decode request body: " + err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	interaction := models.NewInteraction(
		requestBody.UserID,
		"/api/valid-endpoint",
		http.StatusOK,
		false,
		r.RemoteAddr,
	)

	// Save the interaction to Neo4j
	if err := h.Neo4jService.SaveInteraction(interaction); err != nil {
		h.Logger.Error("Failed to save interaction: " + err.Error())
		http.Error(w, "Failed to save interaction", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Interaction logged successfully for user_id: " + interaction.UserID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Interaction logged successfully"))
}

// LogAssociation handles requests to associate two users
func (h *InteractionHandler) LogAssociation(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		User1 string `json:"user1"`
		User2 string `json:"user2"`
	}

	var requestBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		h.Logger.Error("Failed to decode request body: " + err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the service to associate the two users in Neo4j
	if err := h.Neo4jService.AssociatedWith(requestBody.User1, requestBody.User2); err != nil {
		h.Logger.Error("Failed to associate users: " + err.Error())
		http.Error(w, "Failed to associate users", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Users associated successfully: " + requestBody.User1 + " <-> " + requestBody.User2)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Users associated successfully"))
}
