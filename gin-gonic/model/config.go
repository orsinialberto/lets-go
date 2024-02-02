package model

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Config struct {
	Database Database `json:"database"`
}

type Database struct {
	Directory string `json:"directory"`
}

var once sync.Once

var config *Config
var PlayersFilePath string
var VersionFilePath string

func InitConfig(confPath string) {
	once.Do(func() {

		fmt.Println("Reading config:", confPath)

		file, err := os.Open(confPath)
		if err != nil {
			fmt.Println("Error opening configuration file:", err)
			os.Exit(1)
		}
		defer file.Close()

		var c Config
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&c); err != nil {
			fmt.Println("Error decoding configuration file:", err)
			os.Exit(1)
		}

		config = &c
		PlayersFilePath = config.Database.Directory + "players.txt"
		VersionFilePath = config.Database.Directory + "version.txt"

		fmt.Println("Config loaded successfully")
	})
}

func GetConfig() *Config {
	return config
}
