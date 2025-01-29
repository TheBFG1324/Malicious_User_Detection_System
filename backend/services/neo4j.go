package services

import (
	"backend/models"
	"backend/utils"
	"context"
	"fmt"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Neo4jService handles Neo4j interactions
type Neo4jService struct {
	Driver neo4j.DriverWithContext
	Logger *utils.Logger
}

// NewNeo4jService creates a new Neo4jService
func NewNeo4jService(uri, username, password string, logger *utils.Logger) *Neo4jService {
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		logger.Error("Failed to create Neo4j driver: " + err.Error())
		log.Fatalf("Failed to create Neo4j driver: %v", err)
	}
	logger.Info("Successfully connected to Neo4j")
	return &Neo4jService{Driver: driver, Logger: logger}
}

func (s *Neo4jService) SaveInteraction(interaction models.Interaction) error {
	ctx := context.Background()
	session := s.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MERGE (u:User {user_id: $user_id})
			ON CREATE SET u.malicious_score = 0.0

			CREATE (i:Interaction {
				endpoint: $endpoint,
				timestamp: $timestamp,
				response_status_code: $response_status_code,
				honeytoken_triggered: $honeytoken_triggered,
				ip_address: $ip_address
			})
			CREATE (u)-[:HAS_INTERACTION]->(i)
		`

		s.Logger.Debug("Executing SaveInteraction query", true)

		return tx.Run(ctx, query, interaction.ToMap())
	})

	if err != nil {
		s.Logger.Error("Failed to save interaction: " + err.Error())
		return err
	}

	s.Logger.Info("Interaction saved successfully for user_id: " + interaction.UserID)
	return nil
}

// Associated with associates two user_ids in the Neo4j database
func (s *Neo4jService) AssociatedWith(user1, user2 string) error {
	ctx := context.Background()
	session := s.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MERGE (u1:User {user_id: $user1})
			ON CREATE SET u1.malicious_score = 0.0
			MERGE (u2:User {user_id: $user2})
			ON CREATE SET u2.malicious_score = 0.0
			CREATE (u1)-[:ASSOCIATED_WITH]->(u2)
			CREATE (u2)-[:ASSOCIATED_WITH]->(u1)
		`

		s.Logger.Debug("Executing AssociatedWith query", true)

		return tx.Run(ctx, query, map[string]interface{}{"user1": user1, "user2": user2})
	})

	if err != nil {
		s.Logger.Error("Failed to associate users: " + err.Error())
		return err
	}

	s.Logger.Info("Users associated successfully: " + user1 + ", " + user2)
	return nil
}

// GetMaliciousScore returns the malicious_score of a user
func (s *Neo4jService) GetMaliciousScore(userID string) (float64, error) {
	ctx := context.Background()
	session := s.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (u:User {user_id: $user_id})
			RETURN u.malicious_score AS malicious_score
		`

		s.Logger.Debug("Executing GetMaliciousScore query", true)
		res, err := tx.Run(ctx, query, map[string]interface{}{"user_id": userID})
		if err != nil {
			return nil, err
		}

		if res.Next(ctx) {
			record := res.Record()
			maliciousScore, _ := record.Get("malicious_score")
			return maliciousScore, nil
		}

		return nil, fmt.Errorf("user not found")
	})

	if err != nil {
		s.Logger.Error("Failed to get malicious score for user_id: " + userID + " - " + err.Error())
		return 0.0, err
	}

	score, ok := result.(float64)
	if !ok {
		return 0.0, fmt.Errorf("malicious score retrieval failed: expected float64, got %T", result)
	}

	s.Logger.Info(fmt.Sprintf("Retrieved malicious score for user_id %s: %f", userID, score))
	return score, nil
}

// UpdateMaliciousScore updates the malicious_score of a user
func (s *Neo4jService) UpdateMaliciousScore(userID string, newScore float64) error {
	ctx := context.Background()
	session := s.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (u:User {user_id: $user_id})
			SET u.malicious_score = $new_score
		`

		s.Logger.Debug("Executing UpdateMaliciousScore query", true)
		return tx.Run(ctx, query, map[string]interface{}{
			"user_id":   userID,
			"new_score": newScore,
		})
	})

	if err != nil {
		s.Logger.Error("Failed to update malicious score for user_id: " + userID + " - " + err.Error())
		return err
	}

	s.Logger.Info(fmt.Sprintf("Updated malicious score for user_id %s to %f", userID, newScore))
	return nil
}

// RunQuery executes a Cypher query on the Neo4j database
func (s *Neo4jService) RunQuery(query string, params map[string]interface{}) ([]neo4j.Record, error) {
	ctx := context.Background()
	session := s.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	s.Logger.Debug("Running query: "+query, true)
	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		var records []neo4j.Record
		for res.Next(ctx) {
			recPtr := res.Record()
			if recPtr != nil {
				records = append(records, *recPtr)
			}
		}

		if err := res.Err(); err != nil {
			return nil, err
		}

		return records, nil
	})

	if err != nil {
		s.Logger.Error("Failed to execute query: " + err.Error())
		return nil, err
	}

	records, ok := result.([]neo4j.Record)
	if !ok {
		return nil, fmt.Errorf("expected []neo4j.Record but got %T", result)
	}

	s.Logger.Info("Query executed successfully")
	return records, nil
}

// Close closes the Neo4j driver
func (s *Neo4jService) CloseDriver() {
	ctx := context.Background()
	if err := s.Driver.Close(ctx); err != nil {
		s.Logger.Error("Failed to close Neo4j driver: " + err.Error())
		log.Fatalf("Failed to close Neo4j driver: %v", err)
	}
	s.Logger.Info("Neo4j driver closed successfully")
}
