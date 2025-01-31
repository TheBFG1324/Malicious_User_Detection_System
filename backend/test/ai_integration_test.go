package test

import (
	"backend/services"
	"backend/utils"
	"fmt"
	"testing"
)

// TestPredictMaliciousness_RealServer calls the actual Flask AI model
// and prints the predicted maliciousness score.
func TestPredictMaliciousness_RealServer(t *testing.T) {
	// Initialize the logger
	logger := utils.NewLogger()

	// Use the actual Flask server running on localhost:5000
	flaskServerURL := "http://127.0.0.1:5000"

	// Instantiate the AIIntegrationService with the real Flask server URL
	service := services.NewAIIntegrationService(flaskServerURL, logger)

	// Define sample input data
	inputData := map[string]interface{}{
		"total_access_count":             2,
		"honeytoken_access_count":        2,
		"shared_ip_count":                4,
		"avg_associated_malicious_score": 0.858166666666667,
	}

	// Call PredictMaliciousness
	result, err := service.PredictMaliciousness(inputData)
	if err != nil {
		t.Fatalf("Failed to get prediction from Flask server: %v", err)
	}

	// Print the raw response
	fmt.Println("AI Model Response:", result)

	// Print maliciousness score if available
	if score, exists := result["maliciousness_score"]; exists {
		fmt.Println("Predicted Maliciousness Score:", score)
	} else {
		fmt.Println("No maliciousness_score key found in response")
	}
}
