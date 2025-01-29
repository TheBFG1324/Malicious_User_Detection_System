package handlers

import (
	"encoding/json"
	"net/http"

	"backend/services"
	"backend/utils"
)

// UserAnalysisHandler handles user analysis-related requests
type UserAnalysisHandler struct {
	UserAnalysisService *services.UserAnalysisService
	Logger              *utils.Logger
}

// NewUserAnalysisHandler creates a new UserAnalysisHandler
func NewUserAnalysisHandler(userAnalysisService *services.UserAnalysisService, logger *utils.Logger) *UserAnalysisHandler {
	return &UserAnalysisHandler{
		UserAnalysisService: userAnalysisService,
		Logger:              logger,
	}
}

// AnalyzeUser handles requests to analyze a user
func (h *UserAnalysisHandler) AnalyzeUser(w http.ResponseWriter, r *http.Request) {

	type RequestBody struct {
		UserID string `json:"user_id"`
	}

	var requestBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		h.Logger.Error("Failed to decode request body: " + err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	analysisResult, err := h.UserAnalysisService.AnalyzeUser(requestBody.UserID)
	if err != nil {
		h.Logger.Error("Failed to analyze user: " + err.Error())
		http.Error(w, "Failed to analyze user", http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(analysisResult)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
