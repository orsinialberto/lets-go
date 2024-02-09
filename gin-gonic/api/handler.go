package api

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

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
	p.CreatedAt = time.Now().UTC()
	p.Version, err = getVersion()
	if err != nil {
		fmt.Println("Error:", err)
		c.IndentedJSON(http.StatusBadRequest, "Invalid JSON player")
		return
	}

	fmt.Println("Creating player:", p)

	jsonP, err := p.ToJsonString()
	if err != nil {
		fmt.Println("Error:", err)
		c.IndentedJSON(http.StatusBadRequest, "Invalid JSON player")
		return
	}

	if err := SavePlayer(jsonP); err != nil {
		fmt.Println("Error:", err)
		c.IndentedJSON(http.StatusBadRequest, "Error writing player")
		return
	}

	c.IndentedJSON(http.StatusCreated, p)
}

func GetPlayer(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("Search player", id)

	player, err := FindPlayerById(id)
	if err != nil || player.Id == "" {
		fmt.Println("Player not found", id)
		c.IndentedJSON(http.StatusNotFound, "player "+id+" not found")
		return
	}

	c.IndentedJSON(http.StatusOK, player)
}

func GetPlayers(c *gin.Context) {
	fmt.Println("Searching players")

	size, _ := strconv.Atoi(c.Query("size"))
	if size == 0 {
		size = 500
	}

	from, _ := strconv.Atoi(c.Query("from"))

	players, err := ReadPlayersFromVersionLimit(from, size)
	if err != nil {
		fmt.Println("Error:", err)
		c.IndentedJSON(http.StatusBadRequest, "error reading players")
		return
	}

	if players == nil {
		fmt.Println("None player found")
		c.IndentedJSON(http.StatusOK, []model.Player{})
		return
	}

	c.IndentedJSON(http.StatusOK, players)
}

func DeletePlayer(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("Delete player", id)

	if err := DeletePlayerById(id); err != nil {
		fmt.Println("Error:", err)
		c.IndentedJSON(http.StatusBadRequest, "error deleting player")
		return
	}

	c.IndentedJSON(http.StatusNoContent, "")
}

func DeletePlayers(c *gin.Context) {
	fmt.Println("Delete players")

	if err := os.Remove(model.PlayersFilePath); err != nil {
		fmt.Println("Error:", err)
		c.IndentedJSON(http.StatusBadRequest, "internal server error")
		return
	}

	if err := os.Remove(model.VersionFilePath); err != nil {
		fmt.Println("Error:", err)
		c.IndentedJSON(http.StatusBadRequest, "internal server error")
		return
	}

	c.IndentedJSON(http.StatusNoContent, "")
}

func LastVersion(c *gin.Context) {
	fmt.Println("Get last version")
	version, err := GetLastVersion()
	if err != nil {
		fmt.Println("Error:", err)
		c.IndentedJSON(http.StatusBadRequest, "internal server error")
		return
	}

	c.IndentedJSON(http.StatusOK, version)
}

func getVersion() (int, error) {
	version, err := GetLastVersion()
	if err != nil {
		fmt.Println("Error:", err)
		return 0, err
	}

	if err := UpdateVersion(version + 1); err != nil {
		fmt.Println("Error:", err)
		return 0, err
	}

	return version + 1, nil
}
