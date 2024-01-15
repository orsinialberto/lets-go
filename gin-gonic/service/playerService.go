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

	writePlayer(jsonP)

	c.IndentedJSON(http.StatusCreated, p)
}

func GetPlayer(c *gin.Context) {
	pId := c.Param("id")

	fmt.Println("Searching player " + pId)

	var p model.Player
	err := json.Unmarshal([]byte(readPlayer(pId)), &p)

	if err != nil {
		fmt.Println("Player not found " + pId)
		c.IndentedJSON(http.StatusNotFound, "player "+pId+" not found")
		return
	}

	c.IndentedJSON(http.StatusCreated, p)
}

func writePlayer(p string) {

	file, err := os.OpenFile("players.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(p + "\n"); err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func readPlayer(pId string) string {
	file, err := os.Open("players.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, pId) {
			return line
		}
	}

	if err := scanner.Err(); err != nil {
		return ""
	}

	return ""
}
