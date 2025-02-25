package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCheckJWT(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		token        string
		expectedCode int
	}{
		{
			name:         "Missing Token",
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Invalid Token",
			token:        "Bearer invalid-token",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.Default()
			r.Use(CheckJWT())
			r.GET("/protected", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "success"})
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/protected", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %v, got %v", tt.expectedCode, w.Code)
			}
		})
	}
}
