package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-finance/docs"
	"go-finance/internal/client"
	"go-finance/internal/handlers"
	"go-finance/internal/jobs"

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

	tz := os.Getenv("TZ")
	if tz == "" {
		tz = "Asia/Taipei"
	}
	loc, _ := time.LoadLocation(tz)
	sched := jobs.NewScheduler(fc, loc)

	// create a context that is cancelled on SIGINT/SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go sched.Start(ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down server...")
		cancel()
	}()

	// create http.Server so we can shut down gracefully
	srv := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: router,
	}

	// run server in background
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server listen error: %v", err)
		}
	}()

	// wait for cancellation (signal)
	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
	log.Println("Shutdown complete")
}
