package main

import (
	"fmt"
	"log"
	"net/http"

	"ai-knowledge-base/internal/database"
	"ai-knowledge-base/internal/handlers"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, loading from environment")
	}

	db := database.InitDB("./search.db")
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok"}`))
	})
	mux.HandleFunc("/api/search-query", handlers.SearchHandler(db))
	corsHandler := handlers.CORSMiddleware(mux)
	port := ":8080"
	fmt.Printf("Server is starting and listening on port %s...\n", port)
	if err := http.ListenAndServe(port, corsHandler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
