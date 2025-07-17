// In: backend/internal/handlers/search_handler_test.go

package handlers

import (
	"ai-knowledge-base/internal/database" // <-- Replace with your module name
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TestMain is a special function that runs once before all tests in this package.
// We use it to set up and tear down our test database.
func TestMain(m *testing.M) {
	// Setup: Create a temporary test database.
	// We'll run the tests, and then tear it down regardless of pass or fail.
	setup()
	code := m.Run() // Run the tests
	teardown()
	os.Exit(code)
}

const testDBFile = "./test_search.db"

func setup() {
	// Clean up any previous test database file before starting.
	teardown()
}

func teardown() {
	// Remove the test database file.
	os.Remove(testDBFile)
}

// TestSearchHandler is our main test function.
func TestSearchHandler(t *testing.T) {
	// --- ARRANGE ---
	// 1. Initialize a clean, temporary database for this test.
	db := database.InitDB(testDBFile)
	defer db.Close()

	// 2. Create our handler, passing in the test database.
	handler := SearchHandler(db)

	// 3. Create the request body (the JSON we want to send).
	requestBody := SearchRequest{
		Query: "how to reset password?",
	}
	bodyBytes, _ := json.Marshal(requestBody)

	// 4. Create a new HTTP request object for our test.
	// We use `bytes.NewReader` to create a stream from our JSON bytes.
	req, err := http.NewRequest("POST", "/api/search-query", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 5. Create a ResponseRecorder.
	// This is a special tool from httptest that acts like a fake web browser,
	// capturing the HTTP response that our handler writes.
	rr := httptest.NewRecorder()

	// --- ACT ---
	// Call the handler's ServeHTTP method directly, passing our fake request and recorder.
	handler.ServeHTTP(rr, req)

	// --- ASSERT ---
	// 1. Check the status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// 2. Check the response body.
	// We expect the mocked response from our `ai_client`.
	expectedSummary := "To reset your password, please navigate to the login page and click the 'Forgot Password' link. This is a mocked response."

	// Decode the JSON response we received into a struct.
	var responseBody map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&responseBody); err != nil {
		t.Fatalf("could not decode response body: %v", err)
	}

	// Check if the summary in the response matches what we expect.
	if summary := responseBody["ai_summary_answer"]; summary != expectedSummary {
		t.Errorf("handler returned unexpected body: got %v want %v", summary, expectedSummary)
	}

	// 3. (Optional but good) Check if the data was saved to the database.
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM search_history WHERE user_query = ?", requestBody.Query).Scan(&count)
	if err != nil {
		t.Fatalf("could not query database: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 record to be saved in db, but found %d", count)
	}
}

// In: backend/internal/handlers/search_handler_test.go
// (Add these functions after the TestSearchHandler function)

// TestSearchHandler_EmptyQuery tests the case where the user sends an empty query.
func TestSearchHandler_EmptyQuery(t *testing.T) {
	// --- ARRANGE ---
	db := database.InitDB(testDBFile) // Use the test DB
	defer db.Close()

	handler := SearchHandler(db)

	// Create a request with an empty query string
	requestBody := SearchRequest{Query: " "} // A query with only a space
	bodyBytes, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/search-query", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// --- ACT ---
	handler.ServeHTTP(rr, req)

	// --- ASSERT ---
	// We expect a 400 Bad Request status code.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

// TestSearchHandler_InvalidJSON tests the case where the frontend sends malformed JSON.
func TestSearchHandler_InvalidJSON(t *testing.T) {
	// --- ARRANGE ---
	db := database.InitDB(testDBFile)
	defer db.Close()

	handler := SearchHandler(db)

	// Create a request with invalid JSON
	invalidJSON := []byte(`{"query": "test"`) // Missing closing brace

	req, _ := http.NewRequest("POST", "/api/search-query", bytes.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// --- ACT ---
	handler.ServeHTTP(rr, req)

	// --- ASSERT ---
	// We expect a 400 Bad Request status code.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
