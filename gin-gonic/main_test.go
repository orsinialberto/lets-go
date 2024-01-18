package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

var (
	filename string = "players_test.txt"
	playerId string = "a2730fd5-f905-48ee-ad2b-04d30e5c5596"
	router   *gin.Engine
)

func setup(t *testing.T) func() {
	router = setupRouter(filename)

	file, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	file.WriteString("{\"id\":\"" + playerId + "\",\"email\":\"unit-test@example.com\"}\n")
	file.Close()

	return func() {
		os.Remove(filename)
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
