package handlers

import (
	"go-finance/internal/client"

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
	c.JSON(200, res)
}
