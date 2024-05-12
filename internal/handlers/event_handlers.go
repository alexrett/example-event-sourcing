package handlers

import (
	"example-event-sourcing/internal/aggregates"
	"example-event-sourcing/internal/repositories"
	"github.com/gin-gonic/gin"
)

type EventHandlers struct {
	eventRepository *repositories.EventRepository
}

func NewEventHandlers(eventRepository *repositories.EventRepository) *EventHandlers {
	return &EventHandlers{
		eventRepository: eventRepository,
	}
}

func (e *EventHandlers) SaveEvent(c *gin.Context) {
	event := &aggregates.Event{}
	err := c.BindJSON(event)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	_, err = e.eventRepository.SaveEvent(c, event)
	c.JSON(200, gin.H{
		"result": "Event saved",
	})
}
