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
	FileName string `json:"fileName"`
}

var config *Config
var once sync.Once

func SetConfig(e string) *Config {
	once.Do(func() {

		fmt.Println("Setting config")
		filepath := "configs/config_" + e + ".json"
		if e == "test" {
			filepath = "../" + filepath
		}

		file, err := os.Open(filepath)
		if err != nil {
			fmt.Println("Error opening configuration file:", err)
			os.Exit(1)
		}
		defer file.Close()

		var c Config
		d := json.NewDecoder(file)
		if err := d.Decode(&c); err != nil {
			fmt.Println("Error decoding configuration file:", err)
			os.Exit(1)
		}

		config = &c
	})

	return config
}

func GetConfig() *Config {
	return config
}
