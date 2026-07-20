package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"example.com/cli-first-go-app/internal/app"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	n, err := strconv.Atoi(port)
	if err != nil || n < 1 || n > 65535 {
		log.Fatal("PORT must be an integer between 1 and 65535")
	}

	addr := "127.0.0.1:" + port
	fmt.Printf("API listening on http://%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, app.NewHandler()))
}
