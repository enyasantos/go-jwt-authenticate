package main

import (
	"authentication-jwt/internal/server"
	"fmt"
	"net/http"
)

func main() {
	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
