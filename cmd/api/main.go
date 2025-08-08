package main

import (
	"log"

	"github.com/elokanugrah/go-financing-btpns/internal/config"
	"github.com/elokanugrah/go-financing-btpns/internal/database"
	"github.com/elokanugrah/go-financing-btpns/internal/repository/postgres"
	"github.com/elokanugrah/go-financing-btpns/internal/usecase"

	httpDelivery "github.com/elokanugrah/go-financing-btpns/internal/delivery/http"

	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg := config.Load()

	db := database.NewConnection(cfg)
	defer db.Close()

	// Initialize Repository Layer
	tenorRepo := postgres.NewTenorRepository(db)

	// Initialize Usecase Layer
	calculateUseCase := usecase.NewFinancingUsecase(tenorRepo)

	// Initialize Delivery Layer (Handler)
	apiHandler := httpDelivery.NewHandler(calculateUseCase)

	// Setup Router and Start Server
	router := httpDelivery.SetupRouter(apiHandler)

	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
