package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const insertSchedule = `
INSERT INTO schedules (line_id, time)
VALUES ($1, $2)
`

type InsertSchedule struct {
	LineID int64
	Time   pgtype.Time
}

func (q *Queries) InsertSchedule(ctx context.Context, schedule InsertSchedule) error {
	_, err := q.db.Exec(ctx, insertSchedule, schedule.LineID, schedule.Time)
	return err
}

const deleteSchedule = `
DELETE FROM schedules
WHERE line_id = $1
`

type DeleteSchedule struct {
	LineID int64
}

func (q *Queries) DeleteSchedule(ctx context.Context, schedule DeleteSchedule) error {
	_, err := q.db.Exec(ctx, deleteSchedule, schedule.LineID)
	return err
}

type Station struct {
	ID             int64
	StationName    pgtype.Text
	LaneID         int64
	StationStartID int64
	StationEndID   int64
}

const getLanes = `
SELECT 
	stations.id,
	stations.name,
	lines.id as lane_id,
	lines.stations_id_start,
	lines.stations_id_end
FROM stations
LEFT JOIN lines ON stations.id = lines.stations_id_start
ORDER BY stations.id
`

func (q *Queries) GetLanes(ctx context.Context) ([]Station, error) {
	rows, err := q.db.Query(ctx, getLanes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Station
	for rows.Next() {
		var i Station
		if err := rows.Scan(
			&i.ID,
			&i.StationName,
			&i.LaneID,
			&i.StationStartID,
			&i.StationEndID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
