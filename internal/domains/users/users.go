package users

import (
	"example-event-sourcing/internal/domains/users/aggregator"
	"example-event-sourcing/internal/domains/users/repository"
	"example-event-sourcing/internal/event_listener"
	"example-event-sourcing/internal/handlers"
	"example-event-sourcing/internal/repositories"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type UserDomain struct {
	db              *bun.DB
	eventListener   *event_listener.EventListener
	aggregator      *aggregator.Aggregator
	repository      *repository.Repository
	eventRepository *repositories.EventRepository
	eventHandlers   *handlers.EventHandlers
}

func New(
	db *bun.DB,
	eventListener *event_listener.EventListener,
	eventRepository *repositories.EventRepository,
	eventHandlers *handlers.EventHandlers,
) *UserDomain {
	repo := repository.New(db)
	agg := aggregator.New(db, eventListener, repo, eventRepository)

	return &UserDomain{
		db:              db,
		eventListener:   eventListener,
		repository:      repo,
		aggregator:      agg,
		eventRepository: eventRepository,
		eventHandlers:   eventHandlers,
	}
}

func (u *UserDomain) Run() {
	u.aggregator.Run()
}

func (u *UserDomain) QueryRoutes() map[string]gin.HandlerFunc {
	routes := make(map[string]gin.HandlerFunc)
	routes["/users"] = func(c *gin.Context) {
		data, err := u.repository.All(c)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"result": data,
		})
	}
	return routes
}

func (u *UserDomain) MutationRoutes() map[string]gin.HandlerFunc {
	routes := make(map[string]gin.HandlerFunc)
	routes["/users"] = u.eventHandlers.SaveEvent
	return routes
}
