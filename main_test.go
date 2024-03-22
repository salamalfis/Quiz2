package main

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/salamalfis/Golang-DTS/helper/handler"
)

func TestGetUsers(t *testing.T) {
    // Initialize Gin engine
    engine := gin.New()
    engine.GET("/api/v1/users", handler.GetUsers)

    // Create a request to test the endpoint
    req, err := http.NewRequest("GET", "/api/v1/users", nil)
    if err != nil {
        t.Fatal(err)
    }

    // Create a response recorder to record the response
    rr := httptest.NewRecorder()

    // Serve the request to the recorder
    engine.ServeHTTP(rr, req)

    // Check the response status code
    if rr.Code != http.StatusOK {
        t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
    }

    // Check the response body (optional)
    // You can parse the response body and validate it if needed
}
