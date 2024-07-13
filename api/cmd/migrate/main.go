package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"upworkapi/cmd/api/application"
	"upworkapi/pkg/db"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const migrationPath = "file://cmd/migrate/migrations"

var errUnknownCommand = errors.New("unknown command")

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func run(args []string, getenv func(string) string) error {
	cmd := args[0]
	dbName := getenv("DB_NAME")

	config := application.DatabaseConfig{
		Name:     getenv("DB_NAME"),
		Host:     getenv("DB_HOST"),
		Port:     getenv("DB_PORT"),
		Username: getenv("DB_USER"),
		Password: getenv("DB_PASSWORD"),
	}

	database, err := db.Connect(config)
	if err != nil {
		return err
	}
	defer database.Close()

	// Create migration instance
	driver, err := postgres.WithInstance(database, &postgres.Config{})
	if err != nil {
		return err
	}

	// Point to migration files
	m, err := migrate.NewWithDatabaseInstance(migrationPath, dbName, driver)
	if err != nil {
		return err
	}

	slog.Info("starting migration", "cmd", cmd)

	switch cmd {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	case "down":
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	default:
		slog.Error("unknown command", "cmd", cmd)
		return errUnknownCommand
	}

	slog.Info("migration complete")
	return nil
}

func main() {
	cmd := os.Args[len(os.Args)-1]
	args := []string{cmd}

	if err := run(args, os.Getenv); err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "%s\n", err)
		os.Exit(1)
	}
}
