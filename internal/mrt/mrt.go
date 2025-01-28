package mrt

import (
	"context"
	"log"
	"time"

	"github.com/bismastr/scrapper-example/internal/repository"
)

type MrtService struct {
	repo *repository.Queries
}

func NewMrtService(repo *repository.Queries) *MrtService {
	return &MrtService{
		repo: repo,
	}
}

func (s *MrtService) GetAllStation(ctx context.Context) []GetAllStation {
	schedules, err := s.repo.GetSchedule(ctx)
	if err != nil {
		log.Printf("Error getting schedule: %s", err)
	}

	stationMap := make(map[int64]*GetAllStation)
	for _, v := range schedules {
		//Check if statinomaps already have the init stations
		stationValue, exists := stationMap[v.ID]
		if !exists {
			stationMap[v.ID] = &GetAllStation{
				Station: Station{
					StationID:   v.ID,
					StationName: v.Name,
				},
				Line: []Line{},
			}
			stationValue = stationMap[v.ID]
		}

		//make a map for line
		lineMap := make(map[int64]*Line)
		for i := range stationValue.Line {
			lineMap[stationValue.Line[i].LineID] = &stationValue.Line[i]
		}

		line, exist := lineMap[v.LinesID]
		if !exist {
			stationValue.Line = append(stationValue.Line, Line{
				LineID: v.LinesID,
				StationStart: Station{
					StationID:   v.StationsIDStart,
					StationName: v.StationsStartName,
				},
				StationEnd: Station{
					StationID:   v.StationsIDEnd,
					StationName: v.StationsEndName,
				},
				Schedule: []Schedule{},
			})
			line = &stationValue.Line[len(stationValue.Line)-1]
			lineMap[v.LinesID] = line
		}

		line.Schedule = append(line.Schedule, Schedule{
			Time:      MicrosecondsToTimeString(v.Time.Microseconds),
			IsHoliday: v.IsHoliday,
		})
	}

	result := make([]GetAllStation, 0, len(stationMap))
	for _, v := range stationMap {
		result = append(result, *v)
	}

	return result
}

func MicrosecondsToTimeString(microseconds int64) string {
	midnight := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)

	t := midnight.Add(time.Duration(microseconds) * time.Microsecond)

	return t.Format("15:04")
}
