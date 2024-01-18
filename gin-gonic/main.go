package main

import (
	"example.com/gin-gonic/service"
	"github.com/gin-gonic/gin"
)

func main() {
	r := setupRouter("players.txt")
	r.Run(":8080")
}

func setupRouter(filename string) *gin.Engine {
	r := gin.Default()

	r.POST("/players", service.PostPlayer)
	r.GET("/players/:id", service.GetPlayer)
	r.GET("/players", service.GetPlayers)
	r.DELETE("/players/:id", service.DeletePlayer)
	r.DELETE("/players", service.DeletePlayers)

	service.Filename = filename

	return r
}
