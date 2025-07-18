// In: backend/internal/database/database.go

package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// SearchHistory defines the structure for our search history records in the database.
type SearchHistory struct {
	ID                 int64
	UserQuery          string
	AISummaryAnswer    string
	AIRelevantArticles string
	CreatedAt          time.Time
}

// InitDB initializes the SQLite database connection and creates the necessary tables.
func InitDB(filepath string) *sql.DB {
	// sql.Open creates a connection pool, but doesn't actually connect.
	// The connection is established lazily when it's first needed.
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Ping the database to verify the connection is alive.
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// SQL statement to create our search_history table if it doesn't already exist.
	// Using "IF NOT EXISTS" prevents errors on subsequent application startups.
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS search_history (
        "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "user_query" TEXT,
        "ai_summary_answer" TEXT,
        "ai_relevant_articles" TEXT,
        "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	log.Println("Database initialized successfully and table created.")
	return db
}

// SaveSearch saves a given search history record to the database.
// It uses prepared statements to prevent SQL injection vulnerabilities.
func SaveSearch(db *sql.DB, search SearchHistory) (int64, error) {
	// The '?' are placeholders for the actual values.
	stmt, err := db.Prepare("INSERT INTO search_history(user_query, ai_summary_answer, ai_relevant_articles) VALUES(?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	// Execute the prepared statement, passing in the values to use for the placeholders.
	res, err := stmt.Exec(search.UserQuery, search.AISummaryAnswer, search.AIRelevantArticles)
	if err != nil {
		return 0, err
	}

	// Get the ID of the newly inserted row.
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
