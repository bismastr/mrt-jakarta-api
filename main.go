package main

import (
	"context"
	"fmt"
	"log"

	dbClient "github.com/bismastr/scrapper-example/internal/db"
	"github.com/bismastr/scrapper-example/internal/repository"
	colly "github.com/gocolly/colly/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	c := colly.NewCollector()

	db, err := dbClient.NewDbClient()
	if err != nil {
		fmt.Println(err)
		return
	}

	repo := repository.New(db.Pool)
	ctx := context.Background()

	//TX
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		log.Fatalf("Unable to start transaction: %v", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// Using the transaction in the repository
	repoWithTx := repo.WithTx(tx)

	c.OnHTML("[data-stasiun]", func(e *colly.HTMLElement) {
		e.ForEach(".col-12.col-xl-6", func(i int, h *colly.HTMLElement) {
			result, _ := repoWithTx.GetLanes(ctx)
			fmt.Println(result[len(result)-1].StationName)
		})
	})

	c.Visit("https://jakartamrt.co.id/id/jadwal-keberangkatan-mrt?dari=null")
}
