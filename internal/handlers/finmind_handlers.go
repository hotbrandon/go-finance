package handlers

import (
	"encoding/json"
	"go-finance/internal/client"
	"go-finance/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FinmindHandler holds dependencies for Finmind routes.
type FinmindHandler struct {
	Client *client.FinmindClient
}

func NewFinmindHandler(fc *client.FinmindClient) *FinmindHandler {
	return &FinmindHandler{Client: fc}
}

// RegisterRoutes registers Finmind routes on the provided engine.
func (h *FinmindHandler) RegisterRoutes(router *gin.Engine) {
	fg := router.Group("/finmind")
	fg.GET("/TaiwanStockPrice/:data_id/:start_date", h.GetTaiwanStockPrice)
}

// GetTaiwanStockPrice godoc
// @Summary Get Taiwan stock price
// @Produce json
// @Param data_id path string true "stock id"
// @Param start_date path string true "start date"
// @Success 200 {object} models.TaiwanStockPriceResponse
// @Failure 500 {object} map[string]string
// @Router /finmind/TaiwanStockPrice/{data_id}/{start_date} [get]
func (h *FinmindHandler) GetTaiwanStockPrice(c *gin.Context) {
	dataID := c.Param("data_id")
	startDate := c.Param("start_date")

	res, err := h.Client.GetTaiwanStockPrice(c.Request.Context(), "TaiwanStockPrice", dataID, startDate)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	// decode the raw "data" into the concrete slice
	var prices []models.StockPrice
	if err := json.Unmarshal(res.Data, &prices); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid data format: " + err.Error()})
		return
	}

	// return a typed response (good for clients and docs)
	c.JSON(http.StatusOK, models.TaiwanStockPriceResponse{
		Msg:    res.Msg,
		Status: res.Status,
		Data:   prices,
	})
}
