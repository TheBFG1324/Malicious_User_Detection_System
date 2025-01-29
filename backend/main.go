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
	AIIntegrationService := services.NewAIIntegrationService("http://localhost:5000", logger)
	userAnalysisService := services.NewUserAnalysisService(neo4jService, AIIntegrationService, logger)

	// Initialize handlers
	interactionHandler := handlers.NewInteractionHandler(neo4jService, logger)
	honeytokenHandler := handlers.NewHoneytokenHandler(neo4jService, logger)
	userAnalysisHandler := handlers.NewUserAnalysisHandler(userAnalysisService, logger)

	// Define routes
	http.HandleFunc("/api/log-interaction", interactionHandler.LogInteraction)
	http.HandleFunc("/api/detect-honeytoken", honeytokenHandler.DetectHoneytoken)
	http.HandleFunc("/api/analyze-user", userAnalysisHandler.AnalyzeUser)

	// Start the server
	logger.Info("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
