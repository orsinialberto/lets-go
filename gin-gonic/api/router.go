package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/players", PostPlayer)
	r.GET("/players/:id", GetPlayer)
	r.GET("/players", GetPlayers)
	r.DELETE("/players/:id", DeletePlayer)
	r.DELETE("/players", DeletePlayers)

	return r
}
