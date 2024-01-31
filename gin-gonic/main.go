package main

import (
	"example.com/gin-gonic/api"
	"example.com/gin-gonic/model"
)

func main() {
	model.SetConfig("dev")
	r := api.SetupRouter()
	r.Run("127.0.0.1:8080")
}
