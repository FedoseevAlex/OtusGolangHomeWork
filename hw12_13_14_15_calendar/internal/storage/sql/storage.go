package sqlstorage

import (
	"context"
	"database/sql"
	"time"

	"github.com/FedoseevAlex/OtusGolangHomeWork/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"

	// Postgres sql driver import.
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Storage struct { // TODO
	DB *sqlx.DB
}

func New() *Storage {
	return &Storage{}
}

type dbEvent struct {
	ID          uuid.UUID `db:"id"`
	Title       string    `db:"title"`
	StartTime   time.Time `db:"start_time"`
	EndTime     time.Time `db:"end_time"`
	Description string    `db:"description"`
	OwnerID     uuid.UUID `db:"owner_id"`
	NotifyAt    time.Time `db:"notify_at"`
}

func (s *Storage) Connect(ctx context.Context, connStr string) error {
	db, err := sqlx.ConnectContext(ctx, "pgx", connStr)
	if err != nil {
		return err
	}
	s.DB = db
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.DB.Close()
}

func (s *Storage) Create(ctx context.Context, e storage.Event) error {
	// Check for date overlap
	selectQuery := `
	SELECT title FROM events
	WHERE start_time <= $1 AND end_time >= $2
	LIMIT 1;
	`
	row := s.DB.QueryRowContext(ctx, selectQuery, e.EndTime, e.StartTime)
	var title string
	err := row.Scan(&title)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		// That's a valid scenario. We don't found any overlapping events.
	case err != nil:
		// In case of any other error just return it
		return err
	default:
		// If no error occurred that means we found at least one event
		// with overlapping time
		return errors.Wrapf(storage.ErrTimeBusy, "\"%s\" overlaps with \"%s\"", e.Title, title)
	}

	event := dbEvent(e)
	insertQuery := `
	INSERT INTO events (id, title, start_time, end_time, description, owner_id, notify_at)
	VALUES (:id, :title, :start_time, :end_time, :description, :owner_id, :notify_at);`
	_, err = s.DB.NamedExecContext(ctx, insertQuery, event)

	var pgErr pgx.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return errors.Wrapf(storage.ErrEventAlreadyExists, pgErr.Error())
	}

	return err
}

func (s *Storage) Update(ctx context.Context, eventID uuid.UUID, e storage.Event) error {
	event := dbEvent(e)
	updateQuery := `
	UPDATE events
	SET (title,
		start_time,
		end_time,
		description,
		owner_id,
		notify_at) = (:title, :start_time, :end_time, :description, :owner_id, :notify_at)
	WHERE id = :id;`
	res, err := s.DB.NamedExecContext(ctx, updateQuery, event)
	if err != nil {
		return err
	}

	rowsUpdated, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsUpdated == 0 {
		return errors.Wrap(storage.ErrEventNotFound, "no row was updated")
	}
	return nil
}

func (s *Storage) Remove(ctx context.Context, eventID uuid.UUID) error {
	deleteQuery := `DELETE FROM events WHERE id = $1;`
	res, err := s.DB.ExecContext(ctx, deleteQuery, eventID)
	if err != nil {
		return err
	}

	rowsUpdated, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsUpdated == 0 {
		return errors.Wrap(storage.ErrEventDeleteFailed, "no row was deleted")
	}
	return nil
}

func (s *Storage) List(ctx context.Context, startDate, endDate time.Time) ([]storage.Event, error) {
	result := make([]storage.Event, 0, 10)

	selectQuery := `
	SELECT * FROM events
	WHERE start_time >= $1 AND end_time <= $2;
	`

	rows, err := s.DB.QueryxContext(ctx, selectQuery, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event dbEvent
		err := rows.StructScan(&event)
		if err != nil {
			return nil, err
		}

		result = append(result, storage.Event(event))
	}

	return result, nil
}
