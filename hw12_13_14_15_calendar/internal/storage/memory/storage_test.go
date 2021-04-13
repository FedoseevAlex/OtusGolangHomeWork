package memorystorage

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/FedoseevAlex/OtusGolangHomeWork/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// Smoke tests to check basic functionality.
func TestSmoke(t *testing.T) {
	ctx := context.Background()
	memStore := New()
	id := uuid.New()

	t.Run("create event", func(t *testing.T) {
		e := storage.Event{
			ID:          id,
			Title:       "Test event",
			StartTime:   time.Now().Add(2 * time.Hour),
			EndTime:     time.Now().Add(3 * time.Hour),
			Description: "Really useless event. Just for testing purposes.",
			OwnerID:     uuid.New(),
			NotifyAt:    time.Now().Add(time.Hour),
		}

		err := memStore.Create(ctx, e)
		require.NoError(t, err)
		require.Len(t, memStore.events, 1)
	})

	t.Run("overlap event", func(t *testing.T) {
		e := storage.Event{
			ID:          uuid.New(),
			Title:       "Test event",
			StartTime:   time.Now().Add(2 * time.Hour),
			EndTime:     time.Now().Add(3 * time.Hour),
			Description: "Really useless event. Just for testing purposes.",
			OwnerID:     uuid.New(),
			NotifyAt:    time.Now().Add(time.Hour),
		}

		err := memStore.Create(ctx, e)
		require.ErrorIs(t, err, storage.ErrTimeBusy)
		require.Len(t, memStore.events, 1)
	})

	t.Run("duplicate event", func(t *testing.T) {
		e := storage.Event{
			ID:          id,
			Title:       "Test event",
			StartTime:   time.Now().Add(2 * time.Hour),
			EndTime:     time.Now().Add(3 * time.Hour),
			Description: "Really useless event. Just for testing purposes.",
			OwnerID:     uuid.New(),
			NotifyAt:    time.Now().Add(time.Hour),
		}

		err := memStore.Create(ctx, e)
		require.ErrorIs(t, err, storage.ErrEventAlreadyExists)
		require.Len(t, memStore.events, 1)
	})

	t.Run("update event", func(t *testing.T) {
		updatedTitle := "Updated test event"
		e := storage.Event{
			ID:          id,
			Title:       updatedTitle,
			StartTime:   time.Now().Add(2 * time.Hour),
			EndTime:     time.Now().Add(3 * time.Hour),
			Description: "Really useless event. Just for testing purposes.",
			OwnerID:     uuid.New(),
			NotifyAt:    time.Now().Add(time.Hour),
		}

		err := memStore.Update(ctx, id, e)
		require.NoError(t, err)

		updatedEvent := memStore.events[id]
		require.Equal(t, updatedTitle, updatedEvent.Title)
	})

	t.Run("list events", func(t *testing.T) {
		dayDuration := 24 * time.Hour
		today := time.Now().Truncate(dayDuration)
		tomorrow := today.Add(dayDuration)

		events, err := memStore.List(ctx, today, tomorrow)
		require.NoError(t, err)
		require.Len(t, events, 1)
	})

	t.Run("delete event", func(t *testing.T) {
		err := memStore.Remove(ctx, id)
		require.NoError(t, err)
		require.Len(t, memStore.events, 0)
	})

	t.Run("event not found", func(t *testing.T) {
		err := memStore.Update(ctx, id, storage.Event{})
		require.ErrorIs(t, err, storage.ErrEventNotFound)
	})
}

func TestParallelSafety(t *testing.T) {
	ctx := context.Background()
	memStore := New()
	eventsNum := 10
	eventTemplate := storage.Event{
		Title:       "Test event",
		Description: "Really useless event. Just for testing purposes.",
		OwnerID:     uuid.New(),
		NotifyAt:    time.Now().Add(time.Hour),
	}

	uuidPool := make([]uuid.UUID, 0, eventsNum)
	for i := 0; i < eventsNum; i++ {
		uuidPool = append(uuidPool, uuid.New())
	}
	require.Len(t, uuidPool, eventsNum, "uuid pool initialisation failed")

	t.Run("parallel creation of events", func(t *testing.T) {
		wg := sync.WaitGroup{}
		for i := 0; i < eventsNum; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				event := eventTemplate
				event.ID = uuidPool[i]
				event.StartTime = time.Now().Add(time.Duration(i) * time.Second)
				event.EndTime = event.StartTime.Add(time.Millisecond)
				memStore.Create(ctx, event)
			}(i)
		}
		wg.Wait()
		require.Len(t, memStore.events, eventsNum)
	})
}
