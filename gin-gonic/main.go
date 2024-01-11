package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type player struct {
	id    string
	email string
}

func main() {
	r := gin.Default()
	r.POST("/players", CreatePlayer)
	r.Run()
}

func CreatePlayer(c *gin.Context) {
	var p player
	err := c.BindJSON(&p)

	if err != nil {
		fmt.Println("Error:", err)
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, p)
}
