package main

import (
	"log"
	"net/http"
	"time"

	"btc-dex-dashboard/internal/api/handler"
	"btc-dex-dashboard/internal/infrastructure/database"
	"btc-dex-dashboard/internal/repository"
	"btc-dex-dashboard/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.NewDB("dev.db")
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	if err := database.Seed(db); err != nil {
		log.Fatal("failed to seed database:", err)
	}
	log.Println("Database initialized successfully")

	// Repository
	exchangeRepo := repository.NewGormExchangeRepository(db)
	marketRepo := repository.NewGormMarketRepository(db)
	priceRepo := repository.NewGormPriceRepository(db)
	fundingRepo := repository.NewGormFundingRateRepository(db)

	// Service
	spreadService := service.NewSpreadService(marketRepo, priceRepo)
	fundingService := service.NewFundingService(marketRepo, fundingRepo)

	// Handler
	spreadHandler := handler.NewSpreadHandler(spreadService)
	fundingHandler := handler.NewFundingHandler(fundingService)

	r := gin.Default()

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	r.GET("/api/exchanges", func(c *gin.Context) {
		exchanges, err := exchangeRepo.FindAll(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"exchanges": exchanges})
	})

	r.GET("/api/spread", spreadHandler.GetSpread)
	r.GET("/api/funding-rates", fundingHandler.GetRates)

	log.Println("Server starting on :8080")
	r.Run(":8080")
}