package db

import (
	"abstractions/config"
	"abstractions/logger"
	"context"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "gorm.io/driver/postgres"
)

// Migrate takes the number of steps to migrate the database and runs the migrations
// if the steps is 0, it will run all the migrations and bring our database up to date
// to migrate down we need to pass migrations steps as negative -1, -2, -3 etc
// to migrate up we need to pass migrations steps as 1, 2, 3 etc
func Migrate(steps int) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Info(ctx, "running migrations...")

	env, err := config.LoadEnvironment()
	if err != nil {
		logger.Error(ctx, err)
		return
	}

	runMigrations(ctx, env, steps)
}

// MigrateWithConfig runs migrations with a specific configuration
func MigrateWithConfig(env *config.ENV, steps int) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Info(ctx, "running migrations with custom configuration...")

	runMigrations(ctx, env, steps)
}

func runMigrations(ctx context.Context, env *config.ENV, steps int) {
	if steps != 0 {
		logger.Info(ctx, fmt.Sprintf("migration steps = %d", steps))
	}

	logger.Info(ctx, "database host: ", env.DatabaseHost)
	logger.Info(ctx, "database name: ", env.DatabaseName)

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		env.DatabaseUser, env.DatabasePassword, env.DatabaseHost, env.DatabasePort, env.DatabaseName, env.DatabaseSSLMode)
	logger.Info(ctx, connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil || db == nil {
		logger.Error(ctx, err)
		return
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil || driver == nil {
		logger.Error(ctx, err)
		return
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(
		"file://./db/migrations",
		"postgres", driver)
	if err != nil {
		logger.Error(ctx, err)
		return
	}
	if m == nil {
		logger.Error(ctx, "migrate instance is nil")
		return
	}
	defer m.Close()

	if steps == 0 {
		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			logger.Error(ctx, err)
			return
		}
	} else {
		err = m.Steps(steps)
		if err != nil && err != migrate.ErrNoChange {
			logger.Error(ctx, err)
			return
		}
	}

	logger.Info(ctx, "migrations completed successfully")
}

// ForceMigrationVersion resets the migration dirty status and sets the version you pass
func ForceMigrationVersion(version int) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Info(ctx, "fixing migrations...")

	env, err := config.LoadEnvironment()
	if err != nil {
		logger.Error(ctx, err)
		return
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		env.DatabaseUser, env.DatabasePassword, env.DatabaseHost, env.DatabasePort, env.DatabaseName, env.DatabaseSSLMode)
	logger.Info(ctx, connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil || db == nil {
		logger.Error(ctx, err)
		return
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil || driver == nil {
		logger.Error(ctx, err)
		return
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(
		"file://./db/migrations",
		"postgres", driver)
	if err != nil {
		logger.Error(ctx, err)
		return
	}
	if m == nil {
		logger.Error(ctx, "migrate instance is nil")
		return
	}
	defer m.Close()

	err = m.Force(version)
	if err != nil {
		logger.Error(ctx, err)
		return
	}

	logger.Info(ctx, "migrations fixed successfully")
}
