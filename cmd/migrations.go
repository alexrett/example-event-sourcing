package cmd

import (
	"context"
	"example-event-sourcing/internal/db"
	"example-event-sourcing/migrations"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
	"log"
)

const MigrationsGroup = "migrations"

var migrationNameSql *string

func getMigrator(db *bun.DB) *migrate.Migrator {
	return migrate.NewMigrator(db, migrations.Migrations)
}

var migrationGroup = &cobra.Group{
	ID:    MigrationsGroup,
	Title: "Migrations",
}

var migrateInit = &cobra.Command{
	Use:     "migrate/init",
	Short:   "create migration tables",
	GroupID: MigrationsGroup,
	Run: func(cmd *cobra.Command, args []string) {
		database := db.Connect()
		defer func(database *bun.DB) {
			err := database.Close()
			if err != nil {
				panic(err)
			}
		}(database)
		migrator := getMigrator(database)
		c := context.Background()
		err := migrator.Init(c)
		if err != nil {
			log.Fatalln(err.Error())
		}
	},
}

var migrateUp = &cobra.Command{
	Use:     "migrate/up",
	Short:   "migrate database",
	GroupID: MigrationsGroup,
	Run: func(cmd *cobra.Command, args []string) {
		database := db.Connect()
		defer func(database *bun.DB) {
			err := database.Close()
			if err != nil {
				panic(err)
			}
		}(database)
		migrator := getMigrator(database)
		c := context.Background()

		if err := migrator.Lock(c); err != nil {
			log.Fatalln(err.Error())
		}
		defer func(migrator *migrate.Migrator, ctx context.Context) {
			err := migrator.Unlock(ctx)
			if err != nil {
				log.Fatalln(err.Error())
			}
		}(migrator, c)

		group, err := migrator.Migrate(c)
		if err != nil {
			log.Fatalln(err.Error())
		}
		if group.IsZero() {
			fmt.Printf("there are no new migrations to run (database is up to date)\n")
		}
		fmt.Printf("migrated to %s\n", group)
	},
}

var migrateDown = &cobra.Command{
	Use:     "migrate/down",
	Short:   "rollback the last migration group",
	GroupID: MigrationsGroup,
	Run: func(cmd *cobra.Command, args []string) {
		database := db.Connect()
		defer func(database *bun.DB) {
			err := database.Close()
			if err != nil {
				panic(err)
			}
		}(database)
		migrator := getMigrator(database)
		c := context.Background()

		if err := migrator.Lock(c); err != nil {
			log.Fatalln(err.Error())
		}
		defer func(migrator *migrate.Migrator, ctx context.Context) {
			err := migrator.Unlock(ctx)
			if err != nil {
				log.Fatalln(err.Error())
			}
		}(migrator, c)

		group, err := migrator.Rollback(c)
		if err != nil {
			log.Fatalln(err.Error())
		}
		if group.IsZero() {
			fmt.Printf("there are no groups to roll back\n")
		}
		fmt.Printf("rolled back %s\n", group)
	},
}

var migrateCreateSql = &cobra.Command{
	Use:     "migrate/create_sql",
	Short:   "create up and down SQL migrations",
	GroupID: MigrationsGroup,
	Run: func(cmd *cobra.Command, args []string) {
		database := db.Connect()
		defer func(database *bun.DB) {
			err := database.Close()
			if err != nil {
				panic(err)
			}
		}(database)
		migrator := getMigrator(database)
		c := context.Background()
		if migrationNameSql == nil {
			log.Fatalln("name not specified")
		}
		files, err := migrator.CreateSQLMigrations(c, *migrationNameSql)
		if err != nil {
			log.Fatalln(err.Error())
		}

		for _, mf := range files {
			fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)
		}

	},
}

func init() {
	rootCmd.AddGroup(migrationGroup)
	rootCmd.AddCommand(migrateInit)
	rootCmd.AddCommand(migrateUp)
	rootCmd.AddCommand(migrateDown)
	rootCmd.AddCommand(migrateCreateSql)

	migrationNameSql = migrateCreateSql.Flags().StringP("name", "n", "", "migration name")
}
