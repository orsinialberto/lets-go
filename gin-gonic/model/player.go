package model

import (
	"encoding/json"
	"fmt"
)

type Player struct {
	Id    string
	Email string `json:"email"`
}

func (p Player) ToJsonString() (string, error) {
	jsonData, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	return string(jsonData), nil
}
