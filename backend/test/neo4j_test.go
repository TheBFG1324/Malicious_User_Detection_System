package test

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"context"
	"fmt"
	"testing"
)

func TestNeo4jService(t *testing.T) {
	// 1. Setup: create a test logger (you can use your own logger implementation)
	logger := utils.NewLogger()

	// 2. Provide your test Neo4j connection details (could be from environment variables).
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

	// 5. Define test data (main user)
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
	params := map[string]interface{}{"user_id": userID}

	records, err := svc.RunQuery(query, params)
	if err != nil {
		t.Fatalf("Failed to run query to validate data: %v", err)
	}
	if len(records) == 0 {
		t.Fatal("No records returned from validation query")
	}

	// 8. Parse and check results
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

	if totalAccessCount != 3 {
		t.Errorf("Expected total_access_count to be 3, got %d", totalAccessCount)
	}
	if honeytokenAccessCount != 2 {
		t.Errorf("Expected honeytoken_access_count to be 2, got %d", honeytokenAccessCount)
	}
	if sharedIPCount != 2 {
		t.Errorf("Expected shared_ip_count to be 2, got %d", sharedIPCount)
	}

	// 9. Test GetMaliciousScore (should be 0.0 by default)
	initialScore, err := svc.GetMaliciousScore(userID)
	if err != nil {
		t.Fatalf("Failed to get malicious score for user %s: %v", userID, err)
	}
	if initialScore != 0.0 {
		t.Errorf("Expected malicious score to be 0.0 initially, got %f", initialScore)
	}

	// 10. Test UpdateMaliciousScore
	newScore := 3.5
	err = svc.UpdateMaliciousScore(userID, newScore)
	if err != nil {
		t.Fatalf("Failed to update malicious score for user %s: %v", userID, err)
	}

	updatedScore, err := svc.GetMaliciousScore(userID)
	if err != nil {
		t.Fatalf("Failed to re-fetch malicious score for user %s: %v", userID, err)
	}
	if updatedScore != newScore {
		t.Errorf("Expected malicious score to be %f, got %f", newScore, updatedScore)
	}

	// 11. Create associated users and set malicious scores
	associatedUserID1 := "test_user_associated_1"
	associatedUserID2 := "test_user_associated_2"

	// Associate them
	if err := svc.AssociatedWith(userID, associatedUserID1); err != nil {
		t.Fatalf("Failed to associate user %s with %s: %v", userID, associatedUserID1, err)
	}
	if err := svc.AssociatedWith(userID, associatedUserID2); err != nil {
		t.Fatalf("Failed to associate user %s with %s: %v", userID, associatedUserID2, err)
	}

	// Update malicious scores for associated users
	if err := svc.UpdateMaliciousScore(associatedUserID1, 6.0); err != nil {
		t.Fatalf("Failed to update malicious score for %s: %v", associatedUserID1, err)
	}
	if err := svc.UpdateMaliciousScore(associatedUserID2, 4.0); err != nil {
		t.Fatalf("Failed to update malicious score for %s: %v", associatedUserID2, err)
	}

	// 12. Check average malicious score of the associated users
	avgQuery := `
		MATCH (u:User {user_id: $user_id})-[:ASSOCIATED_WITH]->(p:User)
		RETURN AVG(p.malicious_score) AS avg_associated_malicious_score
	`
	avgRecords, err := svc.RunQuery(avgQuery, map[string]interface{}{"user_id": userID})
	if err != nil {
		t.Fatalf("Failed to run avg malicious score query: %v", err)
	}
	if len(avgRecords) == 0 {
		t.Fatalf("No records returned when checking avg associated malicious score")
	}
	avgRec := avgRecords[0]
	avgScoreRaw, _ := avgRec.Get("avg_associated_malicious_score")
	avgScore, ok := avgScoreRaw.(float64)
	if !ok {
		t.Fatalf("Could not parse avg_associated_malicious_score as float64; got %T", avgScoreRaw)
	}

	// We expect average of 6.0 and 4.0 => 5.0
	expectedAvg := 5.0
	if avgScore != expectedAvg {
		t.Errorf("Expected avg_associated_malicious_score to be %f, got %f", expectedAvg, avgScore)
	}

	// 13. Clean up the database (delete all nodes and relationships)
	cleanupQuery := `MATCH (n) DETACH DELETE n`
	_, err = svc.RunQuery(cleanupQuery, nil)
	if err != nil {
		t.Fatalf("Failed to clean up database: %v", err)
	}

	// 14. Close the driver (if not already done by defer)
	svc.CloseDriver()
	fmt.Println("TestNeo4jService completed successfully.")
}
