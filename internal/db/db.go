package db

import (
	"database/sql"
	"example-event-sourcing/internal/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func Connect() *bun.DB {
	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(config.EnvConfigs.DbDsn)))
	db := bun.NewDB(sqlDB, pgdialect.New())
	return db
}
