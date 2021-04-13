package storage

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTimeBusy           = errors.New("given event overlaps with another event")
	ErrEventNotFound      = errors.New("requested event not found in storage")
	ErrEventAlreadyExists = errors.New("duplicate event is added or UUID's match (very likely the former)")
	ErrEventDeleteFailed  = errors.New("event delete failed")
)

type Event struct {
	ID          uuid.UUID
	Title       string
	StartTime   time.Time
	EndTime     time.Time
	Description string
	// User ID which own this event
	OwnerID uuid.UUID
	// When to notify about event
	NotifyAt time.Time
}
