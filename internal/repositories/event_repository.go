package repositories

import (
	"context"
	"example-event-sourcing/internal/aggregates"
	"github.com/uptrace/bun"
)

type EventRepository struct {
	db *bun.DB
}

func NewEventRepository(db *bun.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) GetEventByIdAndType(ctx context.Context, id string, eventType int) (*aggregates.Event, error) {
	events := &aggregates.Event{}
	err := r.db.NewSelect().Model(events).Where("id = ? and type_id = ?", id, eventType).Scan(ctx)
	return events, err
}

func (r *EventRepository) SaveEvent(ctx context.Context, event *aggregates.Event) (*aggregates.Event, error) {
	_, err := r.db.NewInsert().Model(event).Exec(ctx)
	return event, err
}
