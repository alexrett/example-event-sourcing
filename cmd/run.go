package cmd

import (
	"example-event-sourcing/internal/aggregates"
	"example-event-sourcing/internal/db"
	"example-event-sourcing/internal/domains/users"
	"example-event-sourcing/internal/event_listener"
	"example-event-sourcing/internal/handlers"
	"example-event-sourcing/internal/repositories"
	"example-event-sourcing/internal/router"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
)

func getAllDomain(dbConnection *bun.DB) []aggregates.Domain {
	eventRepository := repositories.NewEventRepository(dbConnection)
	eventHandlers := handlers.NewEventHandlers(eventRepository)
	eventListener := event_listener.New(dbConnection)
	domains := []aggregates.Domain{
		users.New(dbConnection, eventListener, eventRepository, eventHandlers),
	}

	return domains
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run application",
	Run: func(cmd *cobra.Command, args []string) {
		dbConnection := db.Connect()
		defer dbConnection.Close()
		domains := getAllDomain(dbConnection)
		server := router.New()
		for _, domain := range domains {
			go domain.Run()
			server.AddQueryRoutes(domain.QueryRoutes())
			server.AddMutationRoutes(domain.MutationRoutes())
		}

		server.Run()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
