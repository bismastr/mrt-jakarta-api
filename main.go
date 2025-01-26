package main

import (
	"context"
	"log"
	"time"

	dbClient "github.com/bismastr/scrapper-example/internal/db"
	"github.com/bismastr/scrapper-example/internal/repository"
	colly "github.com/gocolly/colly/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	c := colly.NewCollector()

	db, err := dbClient.NewDbClient()
	if err != nil {
		log.Fatalf("Error creating DB client: %v", err)
	}

	repo := repository.New(db.Pool)
	ctx := context.Background()

	// Start transaction
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		log.Fatalf("Unable to start transaction: %v", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			log.Fatalf("Recovered from panic: %v", p)
		} else if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	// Using the transaction in the repository
	repoWithTx := repo.WithTx(tx)
	result, err := repo.GetLanes(ctx)
	if err != nil {
		log.Fatalf("Error getting lanes: %v", err)
	}
	lanesMap := lanesMapByStationName(result)

	// Log the lanes map for debugging
	log.Printf("Lanes map: %+v", lanesMap)

	c.OnHTML("[data-stasiun]", func(e *colly.HTMLElement) {
		station := e.Attr("data-stasiun")
		log.Printf("Processing station: %s", station)

		e.ForEach(".col-12.col-xl-6", func(i int, h *colly.HTMLElement) {
			direction := h.ChildText("h3")
			schedule := h.ChildTexts("span")

			var laneEndId int64
			switch direction {
			case "Arah Bundaran HI":
				laneEndId = int64(1)
			case "Arah Lebak Bulus":
				laneEndId = int64(13)
			// Add more cases if there are more directions
			default:
				log.Printf("Unknown direction: %s", direction)
				return
			}

			log.Printf("Station: %s, Direction: %s, LaneEndId: %d", station, direction, laneEndId)

			if lanes, ok := lanesMap[station]; ok {

				if lane, ok := lanes[laneEndId]; ok {
					log.Printf("LaneId: %s", lane.ID)
					deleteReq := repository.DeleteSchedule{
						LineID: lane.ID,
					}

					err = repoWithTx.DeleteSchedule(ctx, deleteReq)
					if err != nil {
						log.Printf("Unable to delete existing schedules: %v", err)
						return
					}

					for _, t := range schedule {
						parsedTime, err := time.Parse("15:04", t)
						if err != nil {
							log.Printf("Unable to parse time: %s, error: %v", t, err)
							return
						}

						scheduleTime := pgtype.Time{
							Microseconds: int64(parsedTime.Hour())*3600*1e6 + int64(parsedTime.Minute())*60*1e6,
							Valid:        true,
						}

						reqInsert := repository.InsertSchedule{
							LineID: lane.ID,
							Time:   scheduleTime,
						}

						err = repoWithTx.InsertSchedule(ctx, reqInsert)
						if err != nil {
							log.Printf("Unable to insert new schedule: %v", err)
							return
						}
					}
				} else {
					log.Printf("Lane not found for station: %s, laneEndId: %d", station, laneEndId)
				}
			} else {
				log.Printf("Station not found in lanes map: %s", station)
			}
		})
	})

	err = c.Visit("https://jakartamrt.co.id/id/jadwal-keberangkatan-mrt?dari=null")
	if err != nil {
		log.Fatalf("Error visiting website: %v", err)
	}
}

func lanesMapByStationName(lanes []repository.Station) map[string]map[int64]repository.Station {
	result := make(map[string]map[int64]repository.Station)

	for _, lane := range lanes {
		if _, exists := result[lane.StationName]; !exists {
			result[lane.StationName] = make(map[int64]repository.Station)
		}
		result[lane.StationName][lane.StationEndID] = lane
	}

	return result
}
