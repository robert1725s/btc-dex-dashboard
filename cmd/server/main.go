package main

import (
	"log"
	"net/http"
	"time"

	"btc-dex-dashboard/internal/domain/model"
	"btc-dex-dashboard/internal/infrastructure/database"

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

	r := gin.Default()

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	r.GET("/api/exchanges", func(c *gin.Context) {
		var exchanges []model.Exchange
		result := db.Find(&exchanges)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"exchanges": exchanges})
	})

	log.Println("Server starting on :8080")
	r.Run(":8080")
}
