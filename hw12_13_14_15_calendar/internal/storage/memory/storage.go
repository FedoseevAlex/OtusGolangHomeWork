package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/FedoseevAlex/OtusGolangHomeWork/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Storage struct {
	mu     sync.RWMutex
	events map[uuid.UUID]storage.Event
}

func New() *Storage {
	return &Storage{
		mu:     sync.RWMutex{},
		events: make(map[uuid.UUID]storage.Event),
	}
}

func (s *Storage) Create(ctx context.Context, e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[e.ID]; ok {
		return errors.Wrapf(storage.ErrEventAlreadyExists, "%s: %s", e.ID.String(), e.Title)
	}

	for _, event := range s.events {
		// Check that events doesn't overlap with others
		if e.EndTime.Before(event.StartTime) || e.StartTime.After(event.EndTime) {
			continue
		}
		return errors.Wrapf(storage.ErrTimeBusy, "\"%s\" overlaps with \"%s\"", e.Title, event.Title)
	}

	s.events[e.ID] = e

	return nil
}

func (s *Storage) Update(ctx context.Context, eventID uuid.UUID, e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.events[eventID]
	if !ok {
		return errors.Wrapf(storage.ErrEventNotFound, "%#v", e)
	}

	e.ID = eventID
	s.events[eventID] = e

	return nil
}

func (s *Storage) Remove(ctx context.Context, eventID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.events, eventID)

	if e, ok := s.events[eventID]; ok {
		return errors.Wrapf(storage.ErrEventDeleteFailed, "%s: %s", e.ID.String(), e.Title)
	}

	return nil
}

func (s *Storage) List(ctx context.Context, startDate, endDate time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]storage.Event, 0, 10)
	for _, event := range s.events {
		if event.StartTime.After(startDate) && event.EndTime.Before(endDate) {
			result = append(result, event)
		}
	}
	return result, nil
}
