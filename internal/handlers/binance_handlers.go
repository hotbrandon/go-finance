package handlers

import (
	"net/http"

	"go-finance/internal/client"

	"github.com/gin-gonic/gin"
)

type BinanceHandler struct {
	Client *client.BinanceClient
}

func NewBinanceHandler(bc *client.BinanceClient) *BinanceHandler {
	return &BinanceHandler{Client: bc}
}

func (h *BinanceHandler) GetTicker(c *gin.Context) {
	symbol := c.Param("symbol")
	// call h.Client.GetTicker(...) and return JSON
	c.JSON(http.StatusOK, gin.H{"symbol": symbol})
}

func (h *BinanceHandler) RegisterRoutes(router *gin.Engine) {
	bg := router.Group("/binance")
	bg.GET("/ticker/:symbol", h.GetTicker)
}
