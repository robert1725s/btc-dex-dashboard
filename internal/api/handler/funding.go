package handler

import (
	"net/http"

	"btc-dex-dashboard/internal/service"

	"github.com/gin-gonic/gin"
)

type FundingHandler struct {
	fundingService *service.FundingService
}

func NewFundingHandler(fundingService *service.FundingService) *FundingHandler {
	return &FundingHandler{fundingService: fundingService}
}

func (h *FundingHandler) GetRates(c *gin.Context) {
	result, err := h.fundingService.GetLatestRates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}