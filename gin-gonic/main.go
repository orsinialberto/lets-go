package main

import (
	"os"

	"example.com/gin-gonic/api"
	"example.com/gin-gonic/model"
)

func main() {
	router := api.SetupRouter()
	router.Run("127.0.0.1:8080")
}

func init() {
	if _, err := os.Stat(model.PlayersFilePath); os.IsNotExist(err) {
		os.Create(model.PlayersFilePath)
	}

	if _, err := os.Stat(model.VersionFilePath); os.IsNotExist(err) {
		os.Create(model.VersionFilePath)
	}

	model.InitConfig("configs/config_dev.json")
}
