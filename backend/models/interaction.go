package models

import "time"

// Interaction represents a user's interaction with the system
type Interaction struct {
	UserID              string    `json:"user_id"`              // Unique ID of the user
	Endpoint            string    `json:"endpoint"`             // The endpoint accessed
	Timestamp           time.Time `json:"timestamp"`            // Time of the interaction
	ResponseStatusCode  int       `json:"response_status_code"` // HTTP response status code
	HoneytokenTriggered bool      `json:"honeytoken_triggered"` // Whether this interaction involved a honeytoken
	IPAddress           string    `json:"ip_address"`           // User's IP address
}

// NewInteraction creates a new Interaction instance
func NewInteraction(userID, endpoint string, statusCode int, honeytoken bool, ipAddress string) Interaction {
	return Interaction{
		UserID:              userID,
		Endpoint:            endpoint,
		Timestamp:           time.Now(),
		ResponseStatusCode:  statusCode,
		HoneytokenTriggered: honeytoken,
		IPAddress:           ipAddress,
	}
}

// ToMap converts the Interaction struct to a map for easier handling (e.g., for Neo4j insertion)
func (i Interaction) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"user_id":              i.UserID,
		"endpoint":             i.Endpoint,
		"timestamp":            i.Timestamp.Format(time.RFC3339),
		"response_status_code": i.ResponseStatusCode,
		"honeytoken_triggered": i.HoneytokenTriggered,
		"ip_address":           i.IPAddress,
	}
}
