package services

import (
	"backend/models"
	"backend/utils"
	"context"
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

// SaveInteraction saves a user interaction to the database
func (s *Neo4jService) SaveInteraction(interaction models.Interaction) error {
	ctx := context.Background()
	defer s.Driver.Close(ctx)

	session := s.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
		MERGE (u:User {user_id: $user_id}) 
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

// RunQuery runs a custom query on the database
func (s *Neo4jService) RunQuery(query string, params map[string]interface{}) ([]neo4j.Record, error) {
	ctx := context.Background()
	session := s.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	s.Logger.Debug("Running query: "+query, true)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return tx.Run(ctx, query, params)
	})

	if err != nil {
		s.Logger.Error("Failed to execute query: " + err.Error())
		return nil, err
	}

	s.Logger.Info("Query executed successfully")
	return result.([]neo4j.Record), nil
}
