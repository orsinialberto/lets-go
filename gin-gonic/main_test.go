package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/gin-gonic/model"
	"github.com/go-playground/assert/v2"
)

func TestPostPlayer(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/players", bytes.NewBufferString("{\"email\":\"unit-test@example.com\"}"))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetPlayers(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/players", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetPlayer(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/players", bytes.NewBufferString("{\"email\":\"unit-test-delete@example.com\"}"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var p model.Player
	err := json.Unmarshal(w.Body.Bytes(), &p)
	assert.Equal(t, nil, err)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/players/"+p.Id, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeletePlayer(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/players", bytes.NewBufferString("{\"email\":\"unit-test-delete@example.com\"}"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var p model.Player
	err := json.Unmarshal(w.Body.Bytes(), &p)
	assert.Equal(t, nil, err)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/players/"+p.Id, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
