package main

import (
	"example.com/gin-gonic/service"
	"github.com/gin-gonic/gin"
)

var fileStorage = "resource/players.txt"

func main() {
	r := setupRouter(fileStorage)
	r.Run(":8080")
}

func setupRouter(fileStorage string) *gin.Engine {
	r := gin.Default()

	r.POST("/players", service.PostPlayer)
	r.GET("/players/:id", service.GetPlayer)
	r.GET("/players", service.GetPlayers)
	r.DELETE("/players/:id", service.DeletePlayer)
	r.DELETE("/players", service.DeletePlayers)

	service.Filename = fileStorage

	return r
}
