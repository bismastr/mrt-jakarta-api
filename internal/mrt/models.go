package mrt

type GetAllStation struct {
	Station Station
	Line    []Line
}

type Line struct {
	LineID       int64
	StationStart Station
	StationEnd   Station
	Schedule     []Schedule
}

type Station struct {
	StationID   int64
	StationName string
}

type Schedule struct {
	Time      string
	IsHoliday bool
}
