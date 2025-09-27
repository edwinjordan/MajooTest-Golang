package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/edwinjordan/MajooTest-Golang/internal/logging"
	"github.com/edwinjordan/MajooTest-Golang/seeders"
)

func runSeeder(db *sql.DB, target string) error {
	ctx := context.Background()
	logging.LogInfo(ctx, "Seeding target", slog.String("target", target))

	switch target {
	case "all":
		if err := seeders.SeedUsers(db); err != nil {
			return fmt.Errorf("seeding users failed: %w", err)
		}
		logging.LogInfo(ctx, "Successfully seeded all tables")
	case "users":
		if err := seeders.SeedUsers(db); err != nil {
			return fmt.Errorf("seeding users failed: %w", err)
		}
		logging.LogInfo(ctx, "Successfully seeded users table")
	case "clear":
		if err := seeders.ClearUsers(db); err != nil {
			return fmt.Errorf("clearing users failed: %w", err)
		}
		logging.LogInfo(ctx, "Successfully cleared users table")
	case "refresh":
		// Clear and then seed users
		if err := seeders.ClearUsers(db); err != nil {
			return fmt.Errorf("clearing users failed: %w", err)
		}
		if err := seeders.SeedUsers(db); err != nil {
			return fmt.Errorf("seeding users failed: %w", err)
		}
		logging.LogInfo(ctx, "Successfully refreshed users table")
	default:
		return errors.New("unknown seed target: " + target + ". Available targets: all, users, clear, refresh")
	}

	return nil
}
