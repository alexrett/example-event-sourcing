package event_listener

import (
	"context"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
	"log"
	"strings"
)

const (
	EventsCreated = "events:created"
)

type EventListener struct {
	db           *bun.DB
	listener     *pgdriver.Listener
	retryCounter map[string]int
}

func New(db *bun.DB) *EventListener {
	return &EventListener{db: db, listener: pgdriver.NewListener(db), retryCounter: make(map[string]int)}
}

func (el *EventListener) Listen(ctx context.Context, channel string, callback func(string) error) {
	if err := el.listener.Listen(ctx, channel); err != nil {
		panic(err)
	}

	for notification := range el.listener.Channel() {
		err := callback(notification.Payload)
		if err != nil && !strings.Contains(err.Error(), `duplicate key value violates unique constraint "locker_pk"`) {
			log.Println(err)
			if _, ok := el.retryCounter[notification.Payload]; !ok {
				el.retryCounter[notification.Payload] = 0
			}
			if el.retryCounter[notification.Payload] < 3 {
				el.retryCounter[notification.Payload]++
				el.Notify(ctx, channel, notification.Payload)
			}
		} else {
			// fixme unsafe, but will see how it goes
			go func() {
				delete(el.retryCounter, notification.Payload)
			}()
		}
	}
}

func (el *EventListener) Notify(ctx context.Context, channel string, payload string) {
	if err := pgdriver.Notify(ctx, el.db, channel, payload); err != nil {
		panic(err)
	}
}
