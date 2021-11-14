package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingHandler(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestSnatchHandler(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	body := bytes.NewReader([]byte("uid=4444"))
	req, _ := http.NewRequest("POST", "/snatch", body)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestOpenHandler(t *testing.T) {

}

func TestWalletListHandler(t *testing.T) {

}