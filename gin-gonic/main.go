package main

import (
	"encoding/json"
	"fmt"
	"os"

	"example.com/gin-gonic/model"
	"example.com/gin-gonic/service"
	"github.com/gin-gonic/gin"
)

func main() {
	file, err := os.Open("resource/configuration/config.json")
	if err != nil {
		fmt.Println("Error opening configuration file:", err)
		return
	}
	defer file.Close()

	var c model.Config
	d := json.NewDecoder(file)
	if err := d.Decode(&c); err != nil {
		fmt.Println("Error decoding configuration file:", err)
		return
	}

	fileName := c.Database.FileName
	if fileName == "" {
		fmt.Println("Storage file not exists")
		return
	}

	r := setupRouter(fileName)
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
