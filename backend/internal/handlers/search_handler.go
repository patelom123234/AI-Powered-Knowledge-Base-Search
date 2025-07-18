package handlers

import (
	"ai-knowledge-base/internal/ai"
	"ai-knowledge-base/internal/database"
	"ai-knowledge-base/internal/kb"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

// SearchRequest defines the structure of the incoming JSON request from the frontend.
type SearchRequest struct {
	Query string `json:"query"`
}

// SearchHandler is the main HTTP handler for the /api/search-query endpoint.
func SearchHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Decode the incoming JSON request body.
		var req SearchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if req.Query == "" {
			http.Error(w, "Query cannot be empty", http.StatusBadRequest)
			return
		}

		// 2. Get the simulated knowledge base articles.
		articles := kb.GetArticles()

		// 3. Call our (mocked) AI client to get a response.
		aiResponse, err := ai.GetAIAnswer(req.Query, articles)
		if err != nil {
			http.Error(w, "Failed to get response from AI service", http.StatusInternalServerError)
			return
		}

		// 4. Prepare the data to be saved to the database.
		// We marshal the relevant articles slice into a JSON string for storage.
		relevantArticlesJSON, _ := json.Marshal(aiResponse.RelevantArticles)

		searchRecord := database.SearchHistory{
			UserQuery:          req.Query,
			AISummaryAnswer:    aiResponse.SummaryAnswer,
			AIRelevantArticles: string(relevantArticlesJSON),
		}

		// 5. Save the interaction to the database.
		_, err = database.SaveSearch(db, searchRecord)
		if err != nil {
			log.Printf("Failed to save search to database: %v", err)
			// We don't return an error to the user here, as the primary function (getting an answer) succeeded.
			// Logging the error is sufficient for now.
		}

		// 6. Encode the AI response and send it back to the frontend.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(aiResponse)
	}
}
