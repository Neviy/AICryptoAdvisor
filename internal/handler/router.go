// Package handler содержит HTTP-роутер и обработчики для API.
package handler

import "github.com/gin-gonic/gin"

func NewRouter(
	coinHandler *CoinHandler,
	userHandler *UserHandler,
) *gin.Engine {
	router := gin.Default()
	// Роуты для монет
	router.GET("/coin/analyze", coinHandler.Analyze)
	router.POST("/coin", coinHandler.Create)
	router.GET("/coin", coinHandler.GetAll)
	// Роуты для пользователей
	router.GET("/users/:email", userHandler.GetUserByEmail)
	router.PUT("/users/:email", userHandler.UpdateUser)
	router.DELETE("/users/:email", userHandler.DeleteUser)
	router.PUT("/users/:email/password", userHandler.ChangePassword)
	return router
}