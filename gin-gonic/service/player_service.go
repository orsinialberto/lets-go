package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"example.com/gin-gonic/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var Filename string

func PostPlayer(c *gin.Context) {
	var p model.Player
	err := c.BindJSON(&p)
	if err != nil {
		fmt.Println("Error:", err)
		c.IndentedJSON(http.StatusBadRequest, "Invalid player")
		return
	}

	p.Id = uuid.New().String()
	p.CreatedAt = time.Now().UTC()
	fmt.Println("Creating player:", p)

	jsonP, err := p.ToJsonString()
	if err != nil {
		fmt.Println("Error:", err)
		c.IndentedJSON(http.StatusBadRequest, "Invalid JSON player")
		return
	}

	if err := writePlayer(jsonP); err != nil {
		fmt.Println("Error:", err)
		c.IndentedJSON(http.StatusBadRequest, "Error writing player")
		return
	}

	c.IndentedJSON(http.StatusCreated, p)
}

func GetPlayer(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("Search player", id)

	p, err := readPlayer(id)
	if err != nil {
		fmt.Println("Player not found", id)
		c.IndentedJSON(http.StatusNotFound, "player "+id+" not found")
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

func DeletePlayers(c *gin.Context) {
	fmt.Println("Delete players")

	if err := os.Remove(Filename); err != nil {
		fmt.Println("Error:", err)
		c.IndentedJSON(http.StatusBadRequest, "internal server error")
		return
	}

	c.IndentedJSON(http.StatusNoContent, "")
}

func writePlayer(s string) error {
	file, err := os.OpenFile(Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
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

func readPlayer(pId string) (model.Player, error) {
	var p model.Player
	file, err := os.Open(Filename)
	if err != nil {
		fmt.Println("Error:", err)
		return p, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if err := json.Unmarshal([]byte(line), &p); err != nil {
			fmt.Println("Error:", err)
			return p, err
		}
		if p.Id == pId {
			return p, nil
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
		return p, err
	}

	return p, nil
}

func readPlayers() ([]model.Player, error) {
	file, err := os.Open(Filename)
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

	filenameTmp, err := copyFileWithoutId(pId)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	if err := os.Remove(Filename); err != nil {
		fmt.Println("Error:", err)
		return err
	}

	if err := os.Rename(filenameTmp, Filename); err != nil {
		fmt.Println("Error:", err)
		return err
	}
	return nil
}

func copyFileWithoutId(pId string) (string, error) {
	file, err := os.OpenFile(Filename, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	filenameTmp := Filename + ".tmp"

	fileTmp, err := os.Create(filenameTmp)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	defer fileTmp.Close()

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, pId) {
			if _, err := fileTmp.WriteString(line + "\n"); err != nil {
				fmt.Println("Error:", err)
				return "", err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	return filenameTmp, nil
}
