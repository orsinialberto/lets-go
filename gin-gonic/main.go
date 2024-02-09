package main

import (
	"fmt"
	"os"
	"path/filepath"

	"example.com/gin-gonic/api"
	"example.com/gin-gonic/model"
)

func main() {
	router := api.SetupRouter()
	router.Run("127.0.0.1:8080")
}

func init() {

	model.InitConfig("configs/config_dev.json")

	if err := os.MkdirAll(filepath.Dir(model.PlayersFilePath), os.ModePerm); err != nil {
		fmt.Println("error creating player path", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(filepath.Dir(model.VersionFilePath), os.ModePerm); err != nil {
		fmt.Println("error creating version path", err)
		os.Exit(1)
	}
}
