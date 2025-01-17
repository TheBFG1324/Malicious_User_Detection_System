package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"backend/utils"
)

// AIIntegrationService handles AI model predictions
type AIIntegrationService struct {
	ModelEndpoint string
	Logger        *utils.Logger
}

// NewAIIntegrationService creates a new AIIntegrationService
func NewAIIntegrationService(endpoint string, logger *utils.Logger) *AIIntegrationService {
	return &AIIntegrationService{
		ModelEndpoint: endpoint,
		Logger:        logger,
	}
}

// PredictMaliciousness sends user interaction data to the AI model and gets a prediction
func (s *AIIntegrationService) PredictMaliciousness(data map[string]interface{}) (map[string]interface{}, error) {
	s.Logger.Info("Preparing to send prediction request to AI model")

	// Serialize data
	jsonData, err := json.Marshal(data)
	if err != nil {
		s.Logger.Error("Failed to serialize data: " + err.Error())
		return nil, fmt.Errorf("failed to serialize data: %v", err)
	}

	// Send HTTP POST request to AI model
	s.Logger.Debug("Sending data to AI model: "+s.ModelEndpoint+"/predict", true)
	resp, err := http.Post(s.ModelEndpoint+"/predict", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		s.Logger.Error("Failed to contact AI model: " + err.Error())
		return nil, fmt.Errorf("failed to contact AI model: %v", err)
	}
	defer resp.Body.Close()

	// Decode response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		s.Logger.Error("Failed to decode AI model response: " + err.Error())
		return nil, fmt.Errorf("failed to decode AI model response: %v", err)
	}

	s.Logger.Info("Received prediction response from AI model")
	return result, nil
}
