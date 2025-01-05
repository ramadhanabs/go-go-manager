package main

import (
	"fmt"
	"go-go-manager/config"
	"go-go-manager/db"
	"log"

	"go-go-manager/routes"
)

func main() {
	cfg := config.LoadConfig()
	database := db.InitDB(cfg)
	r := routes.SetupRouter()

	defer func() {
		if err := database.Close(); err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
		log.Println("Database connection closed.")
	}()

	fmt.Printf("Starting server on port %s...\n", cfg.AppPort)
	r.Run(":" + cfg.AppPort)
}
