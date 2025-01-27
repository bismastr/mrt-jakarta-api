package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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

	repoWithTx := repo.WithTx(tx)
	result, err := repo.GetLanes(ctx)
	if err != nil {
		log.Fatalf("Error getting lanes: %v", err)
	}
	lanesMap := lanesMapByStationName(result)

	c.OnHTML(".tab-content.glasspaneltabscontent", func(e *colly.HTMLElement) {

		e.ForEach("[data-stasiun]", func(i int, stasiun *colly.HTMLElement) {
			station := stasiun.Attr("data-stasiun")
			class := stasiun.Attr("class")
			holidayClass := fmt.Sprintf("#holiday .%s", strings.ReplaceAll(class, " ", "."))

			stasiun.ForEach(".col-12.col-xl-6", func(i int, h *colly.HTMLElement) {
				laneEndId := getLaneEndId(h)
				lane := getLane(lanesMap, station, laneEndId)
				log.Printf("Runnning Lane %s", lane.StationName)
				log.Printf("With ID %v", lane.LaneID)
				deleteSchedule(ctx, repoWithTx, lane)
				insertScheduleForEach(ctx, h, repoWithTx, lane)
			})

			//holiday
			e.ForEach(holidayClass, func(i int, holiday *colly.HTMLElement) {
				holiday.ForEach(".col-12.col-xl-6", func(i int, h *colly.HTMLElement) {
					laneEndId := getLaneEndId(h)
					lane := getLane(lanesMap, station, laneEndId)
					log.Printf("Runnning Lane %s", lane.StationName)
					log.Printf("With ID %v", lane.LaneID)
					insertScheduleForEach(ctx, h, repoWithTx, lane)
				})
			})
		})

	})

	err = c.Visit("https://jakartamrt.co.id/id/jadwal-keberangkatan-mrt")
	if err != nil {
		log.Fatalf("Error visiting website: %v", err)
	}
}

func deleteSchedule(ctx context.Context, repoWithTx *repository.Queries, lane *repository.Station) {
	deleteReq := repository.DeleteSchedule{
		LineID: lane.LaneID,
	}
	log.Printf("Deleting laneId: %v", lane.LaneID)
	err := repoWithTx.DeleteSchedule(ctx, deleteReq)
	if err != nil {
		log.Printf("Unable to delete existing schedules: %v", err)
		return
	}
}

func insertScheduleForEach(ctx context.Context, h *colly.HTMLElement, repoWithTx *repository.Queries, lane *repository.Station) {
	schedule := h.ChildTexts("span")
	isHoliday := false
	fmt.Println("Running ")
	h.DOM.Parents().Each(func(i int, s *goquery.Selection) {
		if id, exists := s.Attr("id"); exists {
			if id == "holiday" {
				isHoliday = true
			}
		}
	})

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
			LineID: lane.LaneID,
			Time:   scheduleTime,
			IsHoliday: pgtype.Bool{
				Valid: true,
				Bool:  isHoliday,
			},
		}

		err = repoWithTx.InsertSchedule(ctx, reqInsert)
		if err != nil {
			log.Printf("Unable to insert new schedule: %v", err)
			return
		}
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

func getLane(lanesMap map[string]map[int64]repository.Station, station string, laneEndId int64) *repository.Station {
	if lanes, ok := lanesMap[station]; ok {
		if lane, ok := lanes[laneEndId]; ok {
			return &lane
		} else {
			log.Printf("Lane not found for station: %s, laneEndId: %d", station, laneEndId)
		}
	} else {
		log.Printf("Station not found in lanes map: %s", station)
	}

	return nil
}

func getLaneEndId(h *colly.HTMLElement) int64 {
	direction := h.ChildText("h3")

	var laneEndId int64
	switch direction {
	case "Arah Bundaran HI":
		laneEndId = int64(1)
	case "Arah Lebak Bulus":
		laneEndId = int64(13)

	default:
		log.Printf("Unknown direction: %s", direction)
		return -1
	}

	return laneEndId

}
