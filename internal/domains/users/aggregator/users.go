package aggregator

import (
	"context"
	"encoding/json"
	"example-event-sourcing/internal/aggregates"
	"example-event-sourcing/internal/domains/users/models"
	"example-event-sourcing/internal/domains/users/repository"
	"example-event-sourcing/internal/event_listener"
	"example-event-sourcing/internal/repositories"
	"github.com/uptrace/bun"
)

const Factor = 1000

const (
	UsersCreated aggregates.AggregateEvent = iota + Factor
	UsersUpdatePassword
	UsersUpdateUsername
	UsersUpdateEmail
	UsersDeleted

	UsersUnusedEvent
)

type EventUpdatePassword struct {
	Password string `json:"password"`
}

type EventUpdateUsername struct {
	Username string `json:"username"`
}

type EventUpdateEmail struct {
	Email string `json:"email"`
}

func ConvertFromInt(in int) aggregates.AggregateEvent {
	return aggregates.AggregateEvent(in)
}

func IsEventTypeValid(eventType aggregates.AggregateEvent) bool {
	return eventType >= UsersCreated && eventType < UsersUnusedEvent
}

type Aggregator struct {
	db              *bun.DB
	eventListener   *event_listener.EventListener
	repository      *repository.Repository
	eventRepository *repositories.EventRepository
}

func New(
	db *bun.DB,
	eventListener *event_listener.EventListener,
	repository *repository.Repository,
	eventRepository *repositories.EventRepository,
) *Aggregator {
	return &Aggregator{
		db:              db,
		eventListener:   eventListener,
		repository:      repository,
		eventRepository: eventRepository,
	}
}

func (a *Aggregator) Run() {
	ctx := context.Background()
	a.eventListener.Listen(ctx, event_listener.EventsCreated, a.buildUsers)
}

func (a *Aggregator) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	userAggregate, err := a.repository.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	user := userAggregate.(*models.User)
	return user, nil
}

func (a *Aggregator) ExtractUserAndUpdatedModel(ctx context.Context, eventID string, eventType int, aggregateID string, eventUpdate interface{}) (*models.User, interface{}, error) {
	event, err := a.eventRepository.GetEventByIdAndType(ctx, eventID, eventType)
	if err != nil {
		return nil, nil, err
	}

	user, err := a.GetUserByID(ctx, aggregateID)
	if err != nil {
		return nil, nil, err
	}

	payloadJson := event.Payload
	err = json.Unmarshal(payloadJson, eventUpdate)
	if err != nil {
		return nil, nil, err
	}

	return user, eventUpdate, nil
}

func (a *Aggregator) buildUsers(payload string) error {
	ctx := context.Background()
	id, eventType, aggregateID := event_listener.ExtractIdAndTypeFromPayload(payload)
	convertedType := ConvertFromInt(eventType)
	if !IsEventTypeValid(convertedType) {
		return nil
	}
	switch convertedType {
	case UsersCreated:
		event, err := a.eventRepository.GetEventByIdAndType(ctx, id, eventType)
		if err != nil {
			return err
		}
		user := &models.User{}
		payloadJson := event.Payload
		err = json.Unmarshal(payloadJson, user)
		if err != nil {
			return err
		}
		return a.repository.Save(ctx, user, id)

	case UsersUpdatePassword:
		user, userUpdate, err := a.ExtractUserAndUpdatedModel(ctx, id, eventType, aggregateID, &EventUpdatePassword{})
		if err != nil {
			return err
		}
		user.Password = userUpdate.(*EventUpdatePassword).Password
		return a.repository.Save(ctx, user, id)
	case UsersUpdateEmail:
		user, userUpdate, err := a.ExtractUserAndUpdatedModel(ctx, id, eventType, aggregateID, &EventUpdateEmail{})
		if err != nil {
			return err
		}
		user.Email = userUpdate.(*EventUpdateEmail).Email
		return a.repository.Save(ctx, user, id)
	case UsersUpdateUsername:
		user, userUpdate, err := a.ExtractUserAndUpdatedModel(ctx, id, eventType, aggregateID, &EventUpdateUsername{})
		if err != nil {
			return err
		}
		user.Username = userUpdate.(*EventUpdateUsername).Username
		return a.repository.Save(ctx, user, id)
	case UsersDeleted:
		user, _, err := a.ExtractUserAndUpdatedModel(ctx, id, eventType, aggregateID, &aggregates.EventDelete{})
		if err != nil {
			return err
		}
		return a.repository.Delete(ctx, user.ID.String(), id)
	default:
		panic("unhandled default case")
	}

	return nil
}
