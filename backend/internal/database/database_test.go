package database

import (
	"database/sql"
	"os"
	"testing"
	"time"
)

// TestInitDB tests the database initialization function
func TestInitDB(t *testing.T) {
	tempFile := "test_db.sqlite"
	defer os.Remove(tempFile)

	db := InitDB(tempFile)
	if db == nil {
		t.Fatal("InitDB returned nil database")
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	// Test that the table was created by checking if it exists
	var tableName string
	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='search_history'").Scan(&tableName)
	if err != nil {
		t.Fatalf("search_history table was not created: %v", err)
	}
	if tableName != "search_history" {
		t.Fatalf("Expected table name 'search_history', got '%s'", tableName)
	}

	rows, err := db.Query("PRAGMA table_info(search_history)")
	if err != nil {
		t.Fatalf("Failed to get table info: %v", err)
	}
	defer rows.Close()

	expectedColumns := map[string]string{
		"id":                   "INTEGER",
		"user_query":           "TEXT",
		"ai_summary_answer":    "TEXT",
		"ai_relevant_articles": "TEXT",
		"created_at":           "TIMESTAMP",
	}

	columnCount := 0
	for rows.Next() {
		var cid, notnull, pk int
		var name, typ string
		var dflt_value sql.NullString
		err := rows.Scan(&cid, &name, &typ, &notnull, &dflt_value, &pk)
		if err != nil {
			t.Fatalf("Failed to scan table info: %v", err)
		}

		expectedType, exists := expectedColumns[name]
		if !exists {
			t.Errorf("Unexpected column: %s", name)
		} else if typ != expectedType {
			t.Errorf("Column %s has wrong type: expected %s, got %s", name, expectedType, typ)
		}
		columnCount++
	}

	if columnCount != len(expectedColumns) {
		t.Errorf("Expected %d columns, got %d", len(expectedColumns), columnCount)
	}
}

// TestInitDBWithInvalidPath tests database initialization with invalid path
func TestInitDBWithInvalidPath(t *testing.T) {
	t.Skip("Skipping test that requires log.Fatalf to panic")
}

// TestSaveSearch tests the SaveSearch function
func TestSaveSearch(t *testing.T) {
	tempFile := "test_save_search.sqlite"
	defer os.Remove(tempFile)

	db := InitDB(tempFile)
	defer db.Close()

	testSearch := SearchHistory{
		UserQuery:          "test query",
		AISummaryAnswer:    "test answer",
		AIRelevantArticles: `[{"id":"1","title":"Test Article"}]`,
		CreatedAt:          time.Now(),
	}

	id, err := SaveSearch(db, testSearch)
	if err != nil {
		t.Fatalf("SaveSearch failed: %v", err)
	}
	if id <= 0 {
		t.Errorf("Expected positive ID, got %d", id)
	}

	var savedSearch SearchHistory
	err = db.QueryRow("SELECT id, user_query, ai_summary_answer, ai_relevant_articles, created_at FROM search_history WHERE id = ?", id).
		Scan(&savedSearch.ID, &savedSearch.UserQuery, &savedSearch.AISummaryAnswer, &savedSearch.AIRelevantArticles, &savedSearch.CreatedAt)
	if err != nil {
		t.Fatalf("Failed to query saved search: %v", err)
	}

	if savedSearch.UserQuery != testSearch.UserQuery {
		t.Errorf("UserQuery mismatch: expected '%s', got '%s'", testSearch.UserQuery, savedSearch.UserQuery)
	}
	if savedSearch.AISummaryAnswer != testSearch.AISummaryAnswer {
		t.Errorf("AISummaryAnswer mismatch: expected '%s', got '%s'", testSearch.AISummaryAnswer, savedSearch.AISummaryAnswer)
	}
	if savedSearch.AIRelevantArticles != testSearch.AIRelevantArticles {
		t.Errorf("AIRelevantArticles mismatch: expected '%s', got '%s'", testSearch.AIRelevantArticles, savedSearch.AIRelevantArticles)
	}
	if savedSearch.ID != id {
		t.Errorf("ID mismatch: expected %d, got %d", id, savedSearch.ID)
	}
}

// TestSaveSearchWithEmptyData tests saving with empty fields
func TestSaveSearchWithEmptyData(t *testing.T) {
	tempFile := "test_empty_data.sqlite"
	defer os.Remove(tempFile)

	db := InitDB(tempFile)
	defer db.Close()

	emptySearch := SearchHistory{
		UserQuery:          "",
		AISummaryAnswer:    "",
		AIRelevantArticles: "",
		CreatedAt:          time.Now(),
	}

	id, err := SaveSearch(db, emptySearch)
	if err != nil {
		t.Fatalf("SaveSearch failed with empty data: %v", err)
	}
	if id <= 0 {
		t.Errorf("Expected positive ID, got %d", id)
	}

	var savedSearch SearchHistory
	err = db.QueryRow("SELECT user_query, ai_summary_answer, ai_relevant_articles FROM search_history WHERE id = ?", id).
		Scan(&savedSearch.UserQuery, &savedSearch.AISummaryAnswer, &savedSearch.AIRelevantArticles)
	if err != nil {
		t.Fatalf("Failed to query saved search: %v", err)
	}

	if savedSearch.UserQuery != "" {
		t.Errorf("Expected empty UserQuery, got '%s'", savedSearch.UserQuery)
	}
	if savedSearch.AISummaryAnswer != "" {
		t.Errorf("Expected empty AISummaryAnswer, got '%s'", savedSearch.AISummaryAnswer)
	}
	if savedSearch.AIRelevantArticles != "" {
		t.Errorf("Expected empty AIRelevantArticles, got '%s'", savedSearch.AIRelevantArticles)
	}
}

// TestSaveSearchWithSpecialCharacters tests saving data with special characters
func TestSaveSearchWithSpecialCharacters(t *testing.T) {
	tempFile := "test_special_chars.sqlite"
	defer os.Remove(tempFile)

	db := InitDB(tempFile)
	defer db.Close()

	specialSearch := SearchHistory{
		UserQuery:          "'; DROP TABLE search_history; --",
		AISummaryAnswer:    "Answer with 'quotes' and \"double quotes\"",
		AIRelevantArticles: `[{"id":"1","title":"Article with 'quotes'"}]`,
		CreatedAt:          time.Now(),
	}

	id, err := SaveSearch(db, specialSearch)
	if err != nil {
		t.Fatalf("SaveSearch failed with special characters: %v", err)
	}

	// Verify the data was saved correctly (should not cause SQL injection)
	var savedSearch SearchHistory
	err = db.QueryRow("SELECT user_query, ai_summary_answer, ai_relevant_articles FROM search_history WHERE id = ?", id).
		Scan(&savedSearch.UserQuery, &savedSearch.AISummaryAnswer, &savedSearch.AIRelevantArticles)
	if err != nil {
		t.Fatalf("Failed to query saved search: %v", err)
	}

	if savedSearch.UserQuery != specialSearch.UserQuery {
		t.Errorf("UserQuery mismatch: expected '%s', got '%s'", specialSearch.UserQuery, savedSearch.UserQuery)
	}
	if savedSearch.AISummaryAnswer != specialSearch.AISummaryAnswer {
		t.Errorf("AISummaryAnswer mismatch: expected '%s', got '%s'", specialSearch.AISummaryAnswer, savedSearch.AISummaryAnswer)
	}
	if savedSearch.AIRelevantArticles != specialSearch.AIRelevantArticles {
		t.Errorf("AIRelevantArticles mismatch: expected '%s', got '%s'", specialSearch.AIRelevantArticles, savedSearch.AIRelevantArticles)
	}
}

// TestSaveSearchMultipleRecords tests saving multiple records
func TestSaveSearchMultipleRecords(t *testing.T) {
	tempFile := "test_multiple_records.sqlite"
	defer os.Remove(tempFile)

	db := InitDB(tempFile)
	defer db.Close()

	searches := []SearchHistory{
		{
			UserQuery:          "first query",
			AISummaryAnswer:    "first answer",
			AIRelevantArticles: `[{"id":"1","title":"First Article"}]`,
			CreatedAt:          time.Now(),
		},
		{
			UserQuery:          "second query",
			AISummaryAnswer:    "second answer",
			AIRelevantArticles: `[{"id":"2","title":"Second Article"}]`,
			CreatedAt:          time.Now(),
		},
		{
			UserQuery:          "third query",
			AISummaryAnswer:    "third answer",
			AIRelevantArticles: `[{"id":"3","title":"Third Article"}]`,
			CreatedAt:          time.Now(),
		},
	}

	var ids []int64
	for _, search := range searches {
		id, err := SaveSearch(db, search)
		if err != nil {
			t.Fatalf("SaveSearch failed: %v", err)
		}
		ids = append(ids, id)
	}

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM search_history").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count records: %v", err)
	}
	if count != len(searches) {
		t.Errorf("Expected %d records, got %d", len(searches), count)
	}

	for i, id := range ids {
		var savedSearch SearchHistory
		err := db.QueryRow("SELECT user_query, ai_summary_answer, ai_relevant_articles FROM search_history WHERE id = ?", id).
			Scan(&savedSearch.UserQuery, &savedSearch.AISummaryAnswer, &savedSearch.AIRelevantArticles)
		if err != nil {
			t.Fatalf("Failed to query saved search %d: %v", i, err)
		}

		if savedSearch.UserQuery != searches[i].UserQuery {
			t.Errorf("Record %d UserQuery mismatch: expected '%s', got '%s'", i, searches[i].UserQuery, savedSearch.UserQuery)
		}
	}
}

// TestSaveSearchWithClosedDB tests saving to a closed database
func TestSaveSearchWithClosedDB(t *testing.T) {
	tempFile := "test_closed_db.sqlite"
	defer os.Remove(tempFile)

	db := InitDB(tempFile)
	db.Close()

	testSearch := SearchHistory{
		UserQuery:          "test query",
		AISummaryAnswer:    "test answer",
		AIRelevantArticles: `[{"id":"1","title":"Test Article"}]`,
		CreatedAt:          time.Now(),
	}

	_, err := SaveSearch(db, testSearch)
	if err == nil {
		t.Error("Expected error when saving to closed database, but got none")
	}
}

// BenchmarkSaveSearch benchmarks the SaveSearch function
func BenchmarkSaveSearch(b *testing.B) {
	tempFile := "benchmark_test.sqlite"
	defer os.Remove(tempFile)

	db := InitDB(tempFile)
	defer db.Close()

	testSearch := SearchHistory{
		UserQuery:          "benchmark query",
		AISummaryAnswer:    "benchmark answer",
		AIRelevantArticles: `[{"id":"1","title":"Benchmark Article"}]`,
		CreatedAt:          time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testSearch.UserQuery = "benchmark query " + string(rune(i))
		_, err := SaveSearch(db, testSearch)
		if err != nil {
			b.Fatalf("SaveSearch failed: %v", err)
		}
	}
}
