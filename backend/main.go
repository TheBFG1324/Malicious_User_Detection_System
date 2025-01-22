package main

import (
	"backend/handlers"
	"backend/services"
	"backend/utils"
	"log"
	"net/http"
)

func main() {
	// Initialize services
	logger := utils.NewLogger()
	neo4jService := services.NewNeo4jService("bolt://localhost:7687", "neo4j", "Password", logger)

	// Initialize handlers
	interactionHandler := handlers.NewInteractionHandler(neo4jService, logger)
	honeytokenHandler := handlers.NewHoneytokenHandler(neo4jService, logger)

	// Define routes
	http.HandleFunc("/api/log-interaction", interactionHandler.LogInteraction)
	http.HandleFunc("/api/detect-honeytoken", honeytokenHandler.DetectHoneytoken)

	// Start the server
	logger.Info("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
