package seeders

import (
	"database/sql"
	"fmt"

	"github.com/edwinjordan/MajooTest-Golang/utils"
)

// SeedUsers populates the users table with sample data
func SeedUsers(db *sql.DB) error {
	// Hash passwords
	password, err := utils.HashPassword("password123")
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Insert users with password hashes
	result, err := db.Exec(`
        INSERT INTO users (name, email, password) VALUES
        ('Alice Johnson', 'alice@example.com', $1),
        ('Bob Smith', 'bob@example.com', $1),
        ('Charlie Brown', 'charlie@example.com', $1),
        ('Diana Prince', 'diana@example.com', $1),
        ('John Doe', 'john@example.com', $1),
        ('Jane Smith', 'jane@example.com', $1)
        ON CONFLICT (email) DO NOTHING;
    `, password)

	if err != nil {
		return fmt.Errorf("failed to insert users: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	fmt.Printf("Inserted %d users\n", rowsAffected)
	return nil
}

// ClearUsers removes all users from the users table (useful for testing)
func ClearUsers(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM users")
	if err != nil {
		return fmt.Errorf("failed to clear users: %w", err)
	}

	fmt.Println("Cleared all users from the database")
	return nil
}
