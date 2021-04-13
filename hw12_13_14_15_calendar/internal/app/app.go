package app

import (
	"context"
	"time"

	"github.com/FedoseevAlex/OtusGolangHomeWork/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type App struct { // TODO
	storage Storage
	log     Logger
}

type Logger interface {
	Debug(msg string, args ...map[string]interface{})
	Info(msg string, args ...map[string]interface{})
	Warn(msg string, args ...map[string]interface{})
	Error(msg string, args ...map[string]interface{})
	Trace(msg string, args ...map[string]interface{})
}

type Storage interface {
	Create(ctx context.Context, e storage.Event) error
	Update(ctx context.Context, eventID uuid.UUID, e storage.Event) error
	Remove(ctx context.Context, eventID uuid.UUID) error
	List(ctx context.Context, startDate, endDate time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{log: logger, storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, title, description string, userID uuid.UUID) error {
	event := storage.Event{
		ID:    uuid.New(),
		Title: title,
		//	StartTime,
		//	EndTime,
		Description: description,
		OwnerID:     userID,
		//	NotifyAt,
	}
	return a.storage.Create(ctx, event)
}

func (a *App) Logger() Logger {
	return a.log
}
