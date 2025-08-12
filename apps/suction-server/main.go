package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "suction-server",
			"version": "1.0.0",
			"message": "Hello World from Suction Server!",
		})
	})

	// Hello endpoint
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World from Suction Server!",
		})
	})

	port := ":8081"
	fmt.Printf("Suction server starting on port %s\n", port)
	if err := r.Run(port); err != nil {
		log.Fatal("Failed to start suction server:", err)
	}
}
