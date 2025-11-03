package main

import (
	"log"
	"os"

	"go-finance/internal/client"
	"go-finance/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	port := os.Getenv("GO_FINANCE_PORT")
	if port == "" {
		log.Fatal("Environment variable GO_FINANCE_PORT not set!")
	}

	// create clients once and reuse
	fc := client.NewFinmindClient("https://api.finmindtrade.com")
	bc := client.NewBinanceClient("https://api.binance.com")

	// create handlers
	fh := handlers.NewFinmindHandler(fc)
	bh := handlers.NewBinanceHandler(bc)

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "ok"}) })

	// each handler registers its own routes
	fh.RegisterRoutes(router)
	bh.RegisterRoutes(router)

	router.Run("0.0.0.0:" + port)
}
