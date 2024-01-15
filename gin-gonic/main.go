package main

import (
	"example.com/gin-gonic/service"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/players", service.PostPlayer)
	r.GET("/players/:id", service.GetPlayer)

	r.Run("localhost:8080")
}
