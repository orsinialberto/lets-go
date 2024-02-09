package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"example.com/gin-gonic/api"
	"example.com/gin-gonic/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

var (
	playerId string = "a2730fd5-f905-48ee-ad2b-04d30e5c5596"
	router   *gin.Engine
)

func setup(t *testing.T) func() {
	model.InitConfig("../configs/config_test.json")
	if err := os.MkdirAll(filepath.Dir(model.PlayersFilePath), os.ModePerm); err != nil {
		fmt.Println("error creating player path", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(filepath.Dir(model.VersionFilePath), os.ModePerm); err != nil {
		fmt.Println("error creating version path", err)
		os.Exit(1)
	}

	router = api.SetupRouter()

	file, _ := os.OpenFile(model.PlayersFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	file.WriteString("{\"id\":\"" + playerId + "\",\"email\":\"unit-test@example.com\",\"version\":1}\n")
	file.Close()

	fileVersion, _ := os.OpenFile(model.VersionFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	fileVersion.WriteString("1")
	fileVersion.Close()

	return func() {
		os.Remove(model.PlayersFilePath)
		os.Remove(model.VersionFilePath)
	}
}

func TestPostPlayer(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/players", bytes.NewBufferString("{\"email\":\"post-unit-test@example.com\"}"))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetPlayers(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/players", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetPlayer(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/players/"+playerId, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeletePlayer(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/players/"+playerId, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeletePlayers(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/players", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestCheckPlayersVersion(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/players", bytes.NewBufferString("{\"email\":\"version-1-unit-test@example.com\"}"))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var player model.Player
	err := json.NewDecoder(w.Body).Decode(&player)

	assert.Equal(t, nil, err)
	assert.Equal(t, 2, player.Version)

	req, _ = http.NewRequest("POST", "/players", bytes.NewBufferString("{\"email\":\"version-2-unit-test@example.com\"}"))
	router.ServeHTTP(w, req)

	err = json.NewDecoder(w.Body).Decode(&player)

	assert.Equal(t, nil, err)
	assert.Equal(t, 3, player.Version)
}

func TestGetLastVersion(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/last-version", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
