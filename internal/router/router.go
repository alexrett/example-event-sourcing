package router

import (
	"example-event-sourcing/internal/config"
	"github.com/gin-gonic/gin"
)

type Router struct {
	router *gin.Engine
}

func New() *Router {
	return &Router{router: gin.Default()}
}

func (r *Router) AddMutationRoutes(routes map[string]gin.HandlerFunc) {
	for path, handler := range routes {
		r.router.POST(path, handler)
	}
}

func (r *Router) AddQueryRoutes(routes map[string]gin.HandlerFunc) {
	for path, handler := range routes {
		r.router.GET(path, handler)
	}
}

func (r *Router) Run() {
	r.router.Run(config.EnvConfigs.AppServerPort)
}
