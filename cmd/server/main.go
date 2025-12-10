package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"btc-dex-dashboard/internal/api/handler"
	"btc-dex-dashboard/internal/api/middleware"
	"btc-dex-dashboard/internal/config"
	"btc-dex-dashboard/internal/infrastructure/database"
	"btc-dex-dashboard/internal/infrastructure/dex"
	"btc-dex-dashboard/internal/job"
	"btc-dex-dashboard/internal/repository"
	"btc-dex-dashboard/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 設定読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}
	log.Printf("Config loaded: port=%s, db=%s, interval=%ds",
		cfg.Server.Port, cfg.Database.Path, cfg.Job.IntervalSeconds)

	db, err := database.NewDB(cfg.Database.Path)
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

	// MarketID を取得（DEX名 → MarketID のマッピング）
	ctx := context.Background()
	markets, err := marketRepo.FindAll(ctx)
	if err != nil {
		log.Fatal("failed to get markets:", err)
	}
	marketIDs := make(map[string]uint)
	for _, m := range markets {
		marketIDs[m.Exchange.Key] = m.ID
	}

	// DEX クライアント
	clients := []dex.DexClient{
		dex.NewHyperliquidClient(),
		dex.NewLighterClient(),
		dex.NewAsterClient(),
	}

	// 定期ジョブ
	interval := time.Duration(cfg.Job.IntervalSeconds) * time.Second
	fetcher := job.NewPriceFetcher(clients, priceRepo, marketIDs)
	scheduler := job.NewScheduler(fetcher, interval)
	go scheduler.Start(ctx)

	// Service
	spreadService := service.NewSpreadService(marketRepo, priceRepo)
	fundingService := service.NewFundingService(marketRepo, fundingRepo)

	// Handler
	spreadHandler := handler.NewSpreadHandler(spreadService)
	fundingHandler := handler.NewFundingHandler(fundingService)

	r := gin.Default()

	// CORS ミドルウェア
	r.Use(middleware.CORS(cfg.CORS.AllowedOrigins))

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

	log.Printf("Server starting on :%s", cfg.Server.Port)
	r.Run(":" + cfg.Server.Port)
}
