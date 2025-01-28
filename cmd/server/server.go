package main

import (
	"log"

	dbClient "github.com/bismastr/scrapper-example/internal/db"
	"github.com/bismastr/scrapper-example/internal/mrt"
	"github.com/bismastr/scrapper-example/internal/repository"
	"github.com/bismastr/scrapper-example/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db, err := dbClient.NewDbClient()
	if err != nil {
		log.Fatalf("Error creating DB client: %v", err)
	}

	repo := repository.New(db.Pool)
	mrtService := mrt.NewMrtService(repo)
	handler := mrt.NewHandler(mrtService)

	server := server.NewServer()
	router := server.SetupRoutes(handler)
	server.Start(router)
}
