package test

import (
	"backend/services"
	"backend/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestPredictMaliciousness_ScoreResponse verifies that AIIntegrationService
// successfully parses and returns the AI model response containing
// {"maliciousness_score": 1.0}.
func TestPredictMaliciousness_ScoreResponse(t *testing.T) {
	// Create a test logger (adjust to your own logger if needed).
	logger := utils.NewLogger()

	// 1. Define a mock handler that returns the desired JSON: {"maliciousness_score": 1.0}.
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"maliciousness_score": 1.0}`))
	}

	// 2. Spin up a local mock HTTP server with the above handler.
	mockServer := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer mockServer.Close()

	// 3. Instantiate the AIIntegrationService with the mock server's URL.
	//    The service will send POST requests to: mockServer.URL + "/predict"
	service := services.NewAIIntegrationService(mockServer.URL, logger)

	// 4. Provide sample input data.
	inputData := map[string]interface{}{
		"user_id":  "123",
		"endpoint": "/api/test",
	}

	// 5. Call PredictMaliciousness.
	result, err := service.PredictMaliciousness(inputData)
	if err != nil {
		t.Fatalf("Unexpected error calling PredictMaliciousness: %v", err)
	}

	// 6. Check the response has the key "maliciousness_score".
	scoreRaw, exists := result["maliciousness_score"]
	if !exists {
		t.Fatalf("Expected 'maliciousness_score' in response, but key was missing")
	}

	// 7. Assert the type is float64 and the value is 1.0.
	score, ok := scoreRaw.(float64)
	if !ok {
		t.Fatalf("Expected maliciousness_score to be float64, got %T", scoreRaw)
	}

	if score != 1.0 {
		t.Errorf("Expected maliciousness_score to be 1.0, got %v", score)
	}
}

// Below is an *optional* additional test that shows you can handle errors, invalid JSON, etc.
func TestPredictMaliciousness_ServerError(t *testing.T) {
	logger := utils.NewLogger()

	// Mock handler that returns a 500 Internal Server Error
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}

	mockServer := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer mockServer.Close()

	service := services.NewAIIntegrationService(mockServer.URL, logger)

	inputData := map[string]interface{}{
		"user_id":  "456",
		"endpoint": "/api/with-error",
	}

	_, err := service.PredictMaliciousness(inputData)
	if err == nil || !strings.Contains(err.Error(), "failed to contact AI model") {
		t.Errorf("Expected an error with 'failed to contact AI model', got: %v", err)
	}
}
