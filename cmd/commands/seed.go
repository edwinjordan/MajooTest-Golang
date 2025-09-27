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
		if err := seeders.SeedTopics(db); err != nil {
			return fmt.Errorf("seeding Topics failed: %w", err)
		}
		// continue for other tables
	case "Topics":
		if err := seeders.SeedTopics(db); err != nil {
			return fmt.Errorf("seeding Topics failed: %w", err)
		}
		// continue for other tables
	default:
		return errors.New("unknown seed target: " + target)
	}

	return nil
}
