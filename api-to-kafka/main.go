package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.api2kafka/model"
)

func main() {
	resp, err := http.Get("http://127.0.0.1:8080/players")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Invalid request. Status:", resp.Status)
		return
	}

	var players []model.Player
	if err := json.NewDecoder(resp.Body).Decode(&players); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(players)
}
