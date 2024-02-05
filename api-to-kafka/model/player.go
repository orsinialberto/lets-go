package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type Player struct {
	Id        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	Version   int       `json:"version"`
}

func (p Player) ToJsonString() (string, error) {
	jsonData, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	return string(jsonData), nil
}
