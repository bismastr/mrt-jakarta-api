package repository

import "github.com/jackc/pgx/v5/pgtype"

type Schedule struct {
	ID                int64
	Name              string
	LinesID           int64
	StationsIDStart   int64
	StationsIDEnd     int64
	StationsEndName   string
	StationsStartName string
	Time              pgtype.Time
	IsHoliday         bool
}
