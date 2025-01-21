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
