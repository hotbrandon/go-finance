package main

import (
	"log"
	"os"

	"go-finance/docs"
	"go-finance/internal/client"
	"go-finance/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Go Finance API
// @version 1.0
// @BasePath /api/v1
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
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "ok"}) })

	// create API group once
	api := router.Group("/api/v1")

	// each handler registers its own routes
	fh.RegisterRoutes(api)
	bh.RegisterRoutes(api)

	router.Run("0.0.0.0:" + port)
}
