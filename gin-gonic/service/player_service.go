package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"example.com/gin-gonic/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	filename = "players.txt"
)

func PostPlayer(c *gin.Context) {
	var p model.Player
	err := c.BindJSON(&p)
	if err != nil {
		fmt.Println("Error:", err)
		c.IndentedJSON(http.StatusBadRequest, "Invalid player")
		return
	}

	p.Id = uuid.New().String()
	fmt.Println("Creating player:", p)

	jsonP, err := p.ToJsonString()
	if err != nil {
		fmt.Println("Error:", err)
		c.IndentedJSON(http.StatusBadRequest, "Invalid JSON player")
		return
	}

	if err := writePlayer(jsonP, filename); err != nil {
		fmt.Println("Error:", err)
		c.IndentedJSON(http.StatusBadRequest, "Error writing player")
		return
	}

	c.IndentedJSON(http.StatusCreated, p)
}

func GetPlayer(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("Search player", id)

	var p model.Player

	player, err := readPlayer(id)
	if err != nil {
		fmt.Println("Player not found", id)
		c.IndentedJSON(http.StatusNotFound, "player "+id+" not found")
		return
	}

	if err := json.Unmarshal([]byte(player), &p); err != nil {
		fmt.Println("Error:", err)
		c.IndentedJSON(http.StatusBadRequest, "error unmarshalling player")
		return
	}

	c.IndentedJSON(http.StatusOK, p)
}

func GetPlayers(c *gin.Context) {
	fmt.Println("Searching players")

	players, err := readPlayers()
	if err != nil {
		fmt.Println("Error:", err)
		c.IndentedJSON(http.StatusBadRequest, "error reading players")
		return
	}

	c.IndentedJSON(http.StatusOK, players)
}

func DeletePlayer(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("Delete player", id)

	if err := deletePlayer(id); err != nil {
		fmt.Println("Error:", err)
		c.IndentedJSON(http.StatusBadRequest, "error deleting player")
		return
	}

	c.IndentedJSON(http.StatusNoContent, "")
}

func writePlayer(s string, f string) error {
	file, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(s + "\n"); err != nil {
		fmt.Println("Error:", err)
		return err
	}

	return nil
}

func readPlayer(pId string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, pId) {
			return line, nil
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	return "", nil
}

func readPlayers() ([]model.Player, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer file.Close()

	players := []model.Player{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var p model.Player
		if err := json.Unmarshal([]byte(line), &p); err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}
		players = append(players, p)
	}

	return players, nil
}

func deletePlayer(pId string) error {
	file, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	scanner := bufio.NewScanner(file)
	fileTmp := filename + "tmp"

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, pId) {
			writePlayer(line, fileTmp)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
		return err
	}

	file.Close()

	if err := os.Remove(filename); err != nil {
		fmt.Println("Error:", err)
		return err
	}

	if err := os.Rename(fileTmp, filename); err != nil {
		fmt.Println("Error:", err)
		return err
	}
	return nil
}