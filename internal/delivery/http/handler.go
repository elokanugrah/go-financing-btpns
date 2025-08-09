package http

import (
	"net/http"

	"github.com/elokanugrah/go-financing-btpns/internal/dto"
	"github.com/elokanugrah/go-financing-btpns/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	financingUsecase usecase.FinancingUsecase
}

func NewHandler(fuc usecase.FinancingUsecase) *Handler {
	return &Handler{
		financingUsecase: fuc,
	}
}

func (h *Handler) Calculate(c *gin.Context) {
	var req dto.CalculateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.financingUsecase.CalculateAllTenors(c.Request.Context(), req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) SubmitFinancing(c *gin.Context) {
	var req dto.SubmitFinancingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.financingUsecase.SubmitFinancing(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
