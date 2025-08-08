package http

import (
	"net/http"

	"github.com/elokanugrah/go-financing-btpns/internal/domain"
	"github.com/elokanugrah/go-financing-btpns/internal/dto"
	"github.com/elokanugrah/go-financing-btpns/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	calculateUsecase usecase.FinancingUsecase
}

func NewHandler(cuc usecase.FinancingUsecase) *Handler {
	return &Handler{
		calculateUsecase: cuc,
	}
}

func (h *Handler) Calculate(c *gin.Context) {
	var req dto.CalculateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Tenor hardcoded sementara, nanti bisa ambil dari DB
	tenors := []domain.Tenor{
		{TenorValue: 6},
		{TenorValue: 12},
		{TenorValue: 18},
		{TenorValue: 24},
		{TenorValue: 30},
		{TenorValue: 36},
	}

	resp, err := h.calculateUsecase.CalculateAllTenors(c.Request.Context(), req.Amount, tenors)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
