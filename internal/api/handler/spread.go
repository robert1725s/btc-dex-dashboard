package handler

import (
	"net/http"

	"btc-dex-dashboard/internal/service"

	"github.com/gin-gonic/gin"
)

type SpreadHandler struct {
	spreadService *service.SpreadService
}

func NewSpreadHandler(spreadService *service.SpreadService) *SpreadHandler {
	return &SpreadHandler{spreadService: spreadService}
}

func (h *SpreadHandler) GetSpread(c *gin.Context) {
	result, err := h.spreadService.CalculateSpread(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}