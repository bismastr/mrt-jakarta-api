package main

import (
	"context"
	"fmt"
	"log"
	"time"

	dbClient "github.com/bismastr/scrapper-example/internal/db"
	"github.com/bismastr/scrapper-example/internal/repository"
	colly "github.com/gocolly/colly/v2"
	"github.com/jackc/pgx/v5/pgtype"
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
		station := e.Attr("data-stasiun")
		e.ForEach(".col-12.col-xl-6", func(i int, h *colly.HTMLElement) {
			schedule := h.ChildTexts("span") //i got this

			if station == "Stasiun Bundaran HI Bank DKI" {
				//Delete existing schedules
				deleteReq := repository.DeleteSchedule{
					LineID: int64(1),
				}

				err = repoWithTx.DeleteSchedule(ctx, deleteReq)
				if err != nil {
					log.Fatalf("Unable to delete existing schedules: %v", err)
				}

				// Insert new schedule times
				for _, t := range schedule {
					parsedTime, err := time.Parse("15:04", t)
					if err != nil {
						log.Fatalf("Unable to parse time: %s, error: %v", t, err)
					}

					scheduleTime := pgtype.Time{
						Microseconds: int64(parsedTime.Hour())*3600*1e6 + int64(parsedTime.Minute())*60*1e6,
						Valid:        true,
					}

					reqInsert := repository.InsertSchedule{
						LineID: int64(1),
						Time:   scheduleTime,
					}

					err = repoWithTx.InsertSchedule(ctx, reqInsert)
					if err != nil {
						log.Fatalf("Unable to insert new schedule: %v", err)
					}
				}

				err = tx.Commit(ctx)
				if err != nil {
					log.Fatalf("Unable to commit transaction: %v", err)
				}
			}
		})
	})

	c.Visit("https://jakartamrt.co.id/id/jadwal-keberangkatan-mrt?dari=null")
}
