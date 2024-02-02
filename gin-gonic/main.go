package main

import (
	"example.com/gin-gonic/api"
	"example.com/gin-gonic/model"
)

func main() {
	model.InitConfig("configs/config_dev.json")

	router := api.SetupRouter()
	router.Run("127.0.0.1:8080")
}
