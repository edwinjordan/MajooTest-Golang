package seeders

import (
	"database/sql"
	"fmt"
)

func SeedPosts(db *sql.DB) error {
	// Insert posts
	result, err := db.Exec(`
		INSERT INTO posts (title, content, slug) VALUES
		('First Post', 'This is the content of the first post.', 'first-post'),
		('Second Post', 'This is the content of the second post.', 'second-post'),
		('Third Post', 'This is the content of the third post.', 'third-post')
		ON CONFLICT (slug) DO NOTHING;
	`)
	if err != nil {
		return fmt.Errorf("failed to insert posts: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	fmt.Printf("Inserted %d posts\n", rowsAffected)
	return nil
}

// Clear removes all posts from the posts table (useful for testing)
func ClearPosts(db *sql.DB) error {
	_, err := db.Exec(`TRUNCATE TABLE posts`)
	if err != nil {
		return fmt.Errorf("failed to clear posts: %w", err)
	}
	return nil
}
