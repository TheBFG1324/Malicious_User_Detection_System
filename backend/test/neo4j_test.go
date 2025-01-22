package test

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"context"
	"testing"
)

// TestNeo4jService is an integration test for the Neo4jService.
func TestNeo4jService(t *testing.T) {
	// 1. Setup: create a test logger (you can use your own logger implementation)
	logger := utils.NewLogger()

	// 2. Provide your test Neo4j connection details (could be from environment variables).
	//    Adjust these to point to your running Neo4j instance.
	uri := "bolt://localhost:7687"
	username := "neo4j"
	password := "Password"

	// 3. Create the Neo4jService
	svc := services.NewNeo4jService(uri, username, password, logger)

	// 4. Defer closing the driver until after tests run
	defer func() {
		ctx := context.Background()
		if err := svc.Driver.Close(ctx); err != nil {
			t.Fatalf("Failed to close Neo4j driver: %v", err)
		}
	}()

	// 5. Define test data
	userID := "test_user_123"
	interactions := []models.Interaction{
		models.NewInteraction(userID, "/api/endpoint1", 200, false, "1.1.1.1"),
		models.NewInteraction(userID, "/api/endpoint2", 404, true, "2.2.2.2"),
		models.NewInteraction(userID, "/api/endpoint3", 500, true, "1.1.1.1"),
	}

	// 6. Save multiple interactions
	for _, interaction := range interactions {
		err := svc.SaveInteraction(interaction)
		if err != nil {
			t.Fatalf("Failed to save interaction for user %s: %v", userID, err)
		}
	}

	// 7. Query to validate the inserted data
	query := `
		MATCH (u:User {user_id: $user_id})-[:HAS_INTERACTION]->(i:Interaction)
		RETURN
			COUNT(i) AS total_access_count,
			SUM(CASE WHEN i.honeytoken_triggered THEN 1 ELSE 0 END) AS honeytoken_access_count,
			COUNT(DISTINCT i.ip_address) AS shared_ip_count
	`

	params := map[string]interface{}{
		"user_id": userID,
	}

	records, err := svc.RunQuery(query, params)
	if err != nil {
		t.Fatalf("Failed to run query to validate data: %v", err)
	}

	if len(records) == 0 {
		t.Fatal("No records returned from validation query")
	}

	// 8. Parse results
	record := records[0]

	totalAccessRaw, _ := record.Get("total_access_count")
	totalAccessCount, ok := totalAccessRaw.(int64)
	if !ok {
		t.Fatalf("Could not parse total_access_count as int64; got %T", totalAccessRaw)
	}

	honeytokenAccessRaw, _ := record.Get("honeytoken_access_count")
	honeytokenAccessCount, ok := honeytokenAccessRaw.(int64)
	if !ok {
		t.Fatalf("Could not parse honeytoken_access_count as int64; got %T", honeytokenAccessRaw)
	}

	sharedIPCountRaw, _ := record.Get("shared_ip_count")
	sharedIPCount, ok := sharedIPCountRaw.(int64)
	if !ok {
		t.Fatalf("Could not parse shared_ip_count as int64; got %T", sharedIPCountRaw)
	}

	// Now compare using the correct type
	if totalAccessCount != 3 {
		t.Errorf("Expected total_access_count to be 3, got %d", totalAccessCount)
	}
	if honeytokenAccessCount != 2 {
		t.Errorf("Expected honeytoken_access_count to be 2, got %d", honeytokenAccessCount)
	}
	if sharedIPCount != 2 {
		t.Errorf("Expected shared_ip_count to be 2, got %d", sharedIPCount)
	}

	svc.CloseDriver()
}
