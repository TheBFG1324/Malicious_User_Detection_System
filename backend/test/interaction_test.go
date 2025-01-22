package test

import (
	"backend/models"
	"testing"
	"time"
)

func TestNewInteraction(t *testing.T) {
	t.Run("Create New Interaction", func(t *testing.T) {
		userID := "user123"
		endpoint := "/api/test"
		statusCode := 200
		honeytoken := true
		ipAddress := "192.168.1.1"

		interaction := models.NewInteraction(userID, endpoint, statusCode, honeytoken, ipAddress)

		// Verify fields
		if interaction.UserID != userID {
			t.Errorf("Expected UserID to be '%s', got '%s'", userID, interaction.UserID)
		}
		if interaction.Endpoint != endpoint {
			t.Errorf("Expected Endpoint to be '%s', got '%s'", endpoint, interaction.Endpoint)
		}
		if interaction.ResponseStatusCode != statusCode {
			t.Errorf("Expected ResponseStatusCode to be '%d', got '%d'", statusCode, interaction.ResponseStatusCode)
		}
		if interaction.HoneytokenTriggered != honeytoken {
			t.Errorf("Expected HoneytokenTriggered to be '%v', got '%v'", honeytoken, interaction.HoneytokenTriggered)
		}
		if interaction.IPAddress != ipAddress {
			t.Errorf("Expected IPAddress to be '%s', got '%s'", ipAddress, interaction.IPAddress)
		}

		// Verify timestamp is recent
		now := time.Now()
		if interaction.Timestamp.Before(now.Add(-time.Minute)) || interaction.Timestamp.After(now.Add(time.Minute)) {
			t.Errorf("Expected Timestamp to be within a minute of now, got '%s'", interaction.Timestamp)
		}
	})
}

func TestToMap(t *testing.T) {
	t.Run("Convert Interaction to Map", func(t *testing.T) {
		interaction := models.Interaction{
			UserID:              "user123",
			Endpoint:            "/api/test",
			Timestamp:           time.Date(2025, 1, 20, 23, 0, 0, 0, time.UTC),
			ResponseStatusCode:  200,
			HoneytokenTriggered: true,
			IPAddress:           "192.168.1.1",
		}

		result := interaction.ToMap()

		// Verify map values
		if result["user_id"] != interaction.UserID {
			t.Errorf("Expected 'user_id' to be '%s', got '%v'", interaction.UserID, result["user_id"])
		}
		if result["endpoint"] != interaction.Endpoint {
			t.Errorf("Expected 'endpoint' to be '%s', got '%v'", interaction.Endpoint, result["endpoint"])
		}
		if result["response_status_code"] != interaction.ResponseStatusCode {
			t.Errorf("Expected 'response_status_code' to be '%d', got '%v'", interaction.ResponseStatusCode, result["response_status_code"])
		}
		if result["honeytoken_triggered"] != interaction.HoneytokenTriggered {
			t.Errorf("Expected 'honeytoken_triggered' to be '%v', got '%v'", interaction.HoneytokenTriggered, result["honeytoken_triggered"])
		}
		if result["ip_address"] != interaction.IPAddress {
			t.Errorf("Expected 'ip_address' to be '%s', got '%v'", interaction.IPAddress, result["ip_address"])
		}

		// Verify timestamp formatting
		expectedTime := interaction.Timestamp.Format(time.RFC3339)
		if result["timestamp"] != expectedTime {
			t.Errorf("Expected 'timestamp' to be '%s', got '%v'", expectedTime, result["timestamp"])
		}
	})
}
