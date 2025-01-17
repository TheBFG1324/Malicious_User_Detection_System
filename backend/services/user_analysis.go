package services

import (
	"fmt"

	"backend/utils"
)

// UserAnalysisService handles business logic for user analysis
type UserAnalysisService struct {
	Neo4jService         *Neo4jService
	AIIntegrationService *AIIntegrationService
	Logger               *utils.Logger
}

// NewUserAnalysisService creates a new UserAnalysisService
func NewUserAnalysisService(neo4jService *Neo4jService, aiService *AIIntegrationService, logger *utils.Logger) *UserAnalysisService {
	return &UserAnalysisService{
		Neo4jService:         neo4jService,
		AIIntegrationService: aiService,
		Logger:               logger,
	}
}

// AnalyzeUser identifies and processes malicious users
func (s *UserAnalysisService) AnalyzeUser(userID string) (map[string]interface{}, error) {
	s.Logger.Info("Starting user analysis for user_id: " + userID)

	// Updated Cypher query to analyze interactions via relationships
	query := `
	MATCH (u:User {user_id: $user_id})-[:HAS_INTERACTION]->(i:Interaction)
	RETURN
		COUNT(i) AS total_access_count,
		SUM(CASE WHEN i.honeytoken_triggered THEN 1 ELSE 0 END) AS honeytoken_access_count,
		COUNT(DISTINCT i.ip_address) AS shared_ip_count
	`
	params := map[string]interface{}{"user_id": userID}
	s.Logger.Debug("Executing Cypher query for user analysis", true)

	// Run the query
	records, err := s.Neo4jService.RunQuery(query, params)
	if err != nil {
		s.Logger.Error("Failed to run Cypher query: " + err.Error())
		return nil, fmt.Errorf("failed to analyze user: %v", err)
	}

	// Extract features from the query results
	if len(records) == 0 {
		s.Logger.Info("No interactions found for user_id: " + userID)
		return nil, fmt.Errorf("no interactions found for user_id: %s", userID)
	}

	record := records[0]
	totalAccessCount, _ := record.Get("total_access_count")
	honeytokenAccessCount, _ := record.Get("honeytoken_access_count")
	sharedIPCount, _ := record.Get("shared_ip_count")

	features := map[string]interface{}{
		"total_access_count":      totalAccessCount.(int64),
		"honeytoken_access_count": honeytokenAccessCount.(int64),
		"shared_ip_count":         sharedIPCount.(int64),
	}
	s.Logger.Info(fmt.Sprintf("Extracted features for user_id %s: %+v", userID, features))

	// Predict maliciousness using the AI model
	s.Logger.Info("Sending features to AI model for prediction")
	prediction, err := s.AIIntegrationService.PredictMaliciousness(features)
	if err != nil {
		s.Logger.Error("Failed to get prediction from AI model: " + err.Error())
		return nil, fmt.Errorf("failed to predict user maliciousness: %v", err)
	}

	s.Logger.Info(fmt.Sprintf("Prediction result for user_id %s: %+v", userID, prediction))
	return prediction, nil
}
