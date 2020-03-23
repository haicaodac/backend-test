package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	routes "hanyny/app"
	jobs "hanyny/app/jobs"
	models "hanyny/app/models"

	"github.com/rs/cors"
)

func main() {
	// Init models
	models.Init()

	// Routes
	routes := routes.NewRouter()

	// crawl.Push()

	// Cors domain
	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"https://hanyny.com", "https://v3.hanyny.com", "http://localhost:3000", "http://localhost:8000", "https://b7511f0f.ngrok.io", "http://b7511f0f.ngrok.io"},
		AllowedMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders: []string{"*"},
		// AllowedHeaders:   []string{"Accept", "Accept-Language", "Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            false,
	}).Handler(routes)

	// Run cron
	jobs.CronRemoveCache()

	// Run server
	port := os.Getenv("PORT")

	fmt.Println("Server run port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
