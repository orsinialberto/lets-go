package main

import (
	"example.com/gin-gonic/api"
	"example.com/gin-gonic/model"
)

func main() {
	router := api.SetupRouter()
	router.Run("127.0.0.1:8080")
}

func init() {
	model.InitConfig("configs/config_dev.json")
}
