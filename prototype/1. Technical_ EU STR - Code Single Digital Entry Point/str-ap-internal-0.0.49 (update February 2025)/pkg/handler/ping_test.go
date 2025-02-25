package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Setup router
	r := gin.Default()
	r.GET("/ping", Ping)

	// Create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	// Serve request
	r.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %v", w.Code)
	}

	// Check response body
	var response Status
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status != "ok" {
		t.Errorf("Expected status 'ok', got %v", response.Status)
	}
}
