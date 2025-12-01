package handlers

import (
	"encoding/json"
	"errors"
	"go-finance/internal/client"
	"go-finance/internal/models"
	"net/http"
	"time"

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
func (h *FinmindHandler) RegisterRoutes(rg *gin.RouterGroup) {
	fg := rg.Group("/finmind")
	fg.GET("/TaiwanStockPrice/:data_id/:start_date", h.GetTaiwanStockPrice)
}

// GetTaiwanStockPrice godoc
// @Summary Get Taiwan stock price
// @Tags Finmind
// @Produce json
// @Param data_id path string true "stock id"
// @Param start_date path string true "start date (format YYYY-MM-DD)"
// @Success 200 {object} models.TaiwanStockPriceAPIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /finmind/TaiwanStockPrice/{data_id}/{start_date} [get]
func (h *FinmindHandler) GetTaiwanStockPrice(c *gin.Context) {
	dataID := c.Param("data_id")
	startDate := c.Param("start_date")

	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorAPIResponse(errors.New("invalid start_date format, expected YYYY-MM-DD")))
		return
	}

	res, err := h.Client.GetTaiwanStockPrice(c.Request.Context(), "TaiwanStockPrice", dataID, startDate)
	if err != nil {
		// Using the helper function for consistency
		c.JSON(http.StatusInternalServerError, models.NewErrorAPIResponse(err))
		return
	}
	// decode the raw "data" into the concrete slice
	var prices []models.StockPrice
	if err := json.Unmarshal(res.Data, &prices); err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorAPIResponse(errors.New("invalid data format from upstream API")))
		return
	}

	// return a typed response (good for clients and docs)
	// Now returning a unified APIResponse with the actual data
	c.JSON(http.StatusOK, models.NewSuccessResponse(prices))
}
