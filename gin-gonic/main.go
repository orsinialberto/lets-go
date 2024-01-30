package main

import (
	"example.com/gin-gonic/model"
	"example.com/gin-gonic/service"
	"github.com/gin-gonic/gin"
)

func main() {
	r := setupRouter("dev")
	r.Run(":8080")
}

func setupRouter(e string) *gin.Engine {
	model.SetConfig(e)
	r := gin.Default()

	r.POST("/players", service.PostPlayer)
	r.GET("/players/:id", service.GetPlayer)
	r.GET("/players", service.GetPlayers)
	r.DELETE("/players/:id", service.DeletePlayer)
	r.DELETE("/players", service.DeletePlayers)

	return r
}
