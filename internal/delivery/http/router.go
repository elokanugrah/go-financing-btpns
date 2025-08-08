package http

import "github.com/gin-gonic/gin"

func SetupRouter(h *Handler) *gin.Engine {
	router := gin.Default()

	router.POST("/calculate-installments", h.Calculate)

	return router
}
