package aggregates

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Domain interface {
	Run()
	QueryRoutes() map[string]gin.HandlerFunc
	MutationRoutes() map[string]gin.HandlerFunc
}

type Aggregate interface {
	GetID() string
}

type Event struct {
	bun.BaseModel `bun:"table:events,alias:e"`
	ID            string          `bun:"id" json:"id"`
	TypeID        int             `bun:"type_id" json:"type_id"`
	AggregateID   string          `bun:"aggregate_id" json:"aggregate_id"`
	Payload       json.RawMessage `bun:"payload" json:"payload"`
	CreatedAt     time.Time       `bun:"created_at" json:"created_at"`
}

type Locker struct {
	bun.BaseModel `bun:"table:locker,alias:loc"`
	EventID       uuid.UUID `bun:"event_id,type:uuid" json:"event_id"`
	AggregateID   uuid.UUID `bun:"aggregate_id,type:uuid" json:"aggregate_id"`
	LockDomain    string    `bun:"lock_domain" json:"lock_domain"`
}

type AggregateEvent int

type EventDelete struct {
	ID string `json:"id"`
}
