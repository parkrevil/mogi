package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Suction Client Starting...")
	fmt.Println("Hello World from Suction Client!")

	// Simple HTTP client to test server
	client := &http.Client{Timeout: 10 * time.Second}

	// Test server health
	resp, err := client.Get("http://localhost:8081/health")
	if err != nil {
		log.Printf("Server not available: %v", err)
		fmt.Println("Client is running. Press Ctrl+C to exit.")
		select {}
	}
	defer resp.Body.Close()

	fmt.Printf("Server Status: %s\n", resp.Status)
	fmt.Println("Client is running. Press Ctrl+C to exit.")
	select {}
}
