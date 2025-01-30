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

func (s *MrtService) GetScheduleById(ctx context.Context, id int64, isHoliday bool, directionStationId int64) *GetStationById {
	schedules, err := s.repo.GetScheduleById(ctx, id, isHoliday, directionStationId)
	if err != nil {
		log.Printf("Error getting schedule: %s", err)
	}

	scheduleFirstIndex := schedules[0]
	result := GetStationById{
		Station: Station{
			StationID:   scheduleFirstIndex.ID,
			StationName: scheduleFirstIndex.Name,
		},
		Line: Line{
			LineID: scheduleFirstIndex.LinesID,
			StationStart: Station{
				StationID:   scheduleFirstIndex.StationsIDStart,
				StationName: scheduleFirstIndex.StationsStartName,
			},
			StationEnd: Station{
				StationID:   scheduleFirstIndex.StationsIDEnd,
				StationName: scheduleFirstIndex.StationsEndName,
			},
			ScheduleNormal:  []Schedule{},
			ScheduleHoliday: []Schedule{},
		},
	}

	for _, v := range schedules {
		now := time.Now()
		midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		elapsed := now.Sub(midnight)
		currentMicroseconds := elapsed.Microseconds()

		if v.Time.Microseconds > currentMicroseconds {
			schedule := Schedule{
				Time:      MicrosecondsToTimeString(v.Time.Microseconds),
				IsHoliday: v.IsHoliday,
			}

			if v.IsHoliday {
				result.Line.ScheduleHoliday = append(result.Line.ScheduleHoliday, schedule)
			} else {
				result.Line.ScheduleNormal = append(result.Line.ScheduleNormal, schedule)
			}
		}

	}

	return &result
}

func (s *MrtService) GetAllStation(ctx context.Context) []GetAllStation {
	schedules, err := s.repo.GetSchedule(ctx) // Ensure this uses your SQL query
	if err != nil {
		log.Printf("Error getting schedule: %s", err)
	}

	var result []GetAllStation        //To store the result
	var currentStation *GetAllStation //To track currentStation
	var currentLine *Line             //To track currentLine

	for _, v := range schedules {
		// Create a new station if it doesn't exist or new station is encountered
		if currentStation == nil || currentStation.Station.StationID != v.ID {
			result = append(result, GetAllStation{
				Station: Station{
					StationID:   v.ID,
					StationName: v.Name,
				},
				Line: []Line{},
			})
			currentStation = &result[len(result)-1]
			currentLine = nil
		}
		// Create a new line if it doesn't exist or new line is encountered
		if currentLine == nil || currentLine.LineID != v.LinesID {
			currentStation.Line = append(currentStation.Line, Line{
				LineID: v.LinesID,
				StationStart: Station{
					StationID:   v.StationsIDStart,
					StationName: v.StationsStartName,
				},
				StationEnd: Station{
					StationID:   v.StationsIDEnd,
					StationName: v.StationsEndName,
				},
				ScheduleHoliday: []Schedule{},
				ScheduleNormal:  []Schedule{},
			})
			currentLine = &currentStation.Line[len(currentStation.Line)-1]
		}

		// Append schedule to the current line
		schedule := Schedule{
			Time:      MicrosecondsToTimeString(v.Time.Microseconds),
			IsHoliday: v.IsHoliday,
		}
		if v.IsHoliday {
			currentLine.ScheduleHoliday = append(currentLine.ScheduleHoliday, schedule)
		} else {
			currentLine.ScheduleNormal = append(currentLine.ScheduleNormal, schedule)
		}
	}

	return result
}

func MicrosecondsToTimeString(microseconds int64) string {
	midnight := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)

	t := midnight.Add(time.Duration(microseconds) * time.Microsecond)

	return t.Format("15:04")
}
