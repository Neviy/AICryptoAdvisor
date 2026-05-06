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

type createCoinRequest struct {
	Name string `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required,gt=0"`
}

// NewCoinHandler — создаёт новый handler с переданным сервисом.
func NewCoinHandler(s *service.Service) *CoinHandler {
	return &CoinHandler{service: s}
}

// Analyze — обрабатывает запрос анализа монеты по имени.
func (h *CoinHandler) Analyze(c *gin.Context) {
	name:= c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name query parameter is required"})
		return
	}
	result, err := h.service.AnalysisCoin(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// Create — обрабатывает запрос создания новой монеты.
func (h *CoinHandler) Create(c *gin.Context){
	var coin createCoinRequest
	if err:=c.ShouldBindJSON(&coin);err !=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}
	coinCreated,err:= h.service.CreateCoin(c.Request.Context(), coin.Name, coin.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, coinCreated)
}

// GetAll — обрабатывает запрос получения всех монет.
func (h *CoinHandler) GetAll(c *gin.Context){
	coins, err := h.service.AllCoin(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, coins)	
}