package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/players", PostPlayer)
	router.GET("/players/:id", GetPlayer)
	router.GET("/players", GetPlayers)
	router.DELETE("/players/:id", DeletePlayer)
	router.DELETE("/players", DeletePlayers)

	return router
}
