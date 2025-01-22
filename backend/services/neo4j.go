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

// SaveInteraction saves a user interaction to the database
func (s *Neo4jService) SaveInteraction(interaction models.Interaction) error {
	ctx := context.Background()

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
