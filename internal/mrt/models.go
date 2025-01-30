package mrt

type GetAllStation struct {
	Station Station `json:"station"`
	Line    []Line  `json:"line"`
}

type GetStationById struct {
	Station Station `json:"station"`
	Line    Line    `json:"line"`
}

type Line struct {
	LineID          int64      `json:"line_id"`
	StationStart    Station    `json:"start_station"`
	StationEnd      Station    `json:"end_station"`
	ScheduleNormal  []Schedule `json:"schedule_normal"`
	ScheduleHoliday []Schedule `json:"schedule_holiday"`
}

type Station struct {
	StationID   int64  `json:"station_id"`
	StationName string `json:"station_name"`
}

type Schedule struct {
	Time      string `json:"time"`
	IsHoliday bool   `json:"is_holiday"`
}
