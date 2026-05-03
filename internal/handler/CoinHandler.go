// Пакет handler обрабатывает HTTP-запросы и передаёт их в service слой.
package handler

import (
	"AICryptoAdvisor/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CoinHandler — обработчик HTTP-запросов, связанных с монетами.
type CoinHandler struct {
	service *service.Service
}

// NewCoinHandler — создаёт новый handler с переданным сервисом.
func NewCoinHandler(s *service.Service) *CoinHandler {
	return &CoinHandler{service: s}
}

// Analyze — обрабатывает запрос анализа монеты по имени.
func (h *CoinHandler) Analyze(c *gin.Context) {
	name := c.Query("name")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	coin, err := h.service.AnalysisCoin(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, coin)
}