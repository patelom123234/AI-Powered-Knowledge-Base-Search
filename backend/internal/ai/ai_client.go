package ai

import (
	"ai-knowledge-base/internal/kb"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type AIResponse struct {
	SummaryAnswer    string       `json:"ai_summary_answer"`
	RelevantArticles []kb.Article `json:"ai_relevant_articles"`
}

// GenerativeAIModel interface for dependency injection.
type GenerativeAIModel interface {
	GenerateContent(ctx context.Context, parts ...genai.Part) (*genai.GenerateContentResponse, error)
}

// ModelFactory function type for creating AI models.
type ModelFactory func(ctx context.Context, apiKey string) (GenerativeAIModel, error)

// GetAIAnswer is the REAL function that calls the Google Gemini API.
func GetAIAnswer(userQuery string, articles []kb.Article) (*AIResponse, error) {
	return getAIAnswerWithFactory(realModelFactory, userQuery, articles)
}

// getAIAnswerWithFactory is the internal testable function
func getAIAnswerWithFactory(factory ModelFactory, userQuery string, articles []kb.Article) (*AIResponse, error) {
	ctx := context.Background()

	// Get the API Key from the environment variable.
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	// Create a new client with your API key.
	model, err := factory(ctx, apiKey)
	if err != nil {
		log.Printf("Failed to create genai client: %v", err)
		return nil, err
	}

	prompt := buildPrompt(userQuery, articles)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Printf("Failed to generate content: %v", err)
		return nil, err
	}

	// The response from Gemini is inside resp.Candidates.
	// We need to parse this response to extract our JSON.
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("received an empty response from AI")
	}

	aiContent := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])

	// The AI's response might include markdown formatting for the JSON block (```json ... ```).
	// We need to clean this up before parsing.
	cleanedJSON := cleanAIResponse(aiContent)

	// Unmarshal the cleaned JSON string into our AIResponse struct.
	var aiResponse AIResponse
	err = json.Unmarshal([]byte(cleanedJSON), &aiResponse)
	if err != nil {
		log.Printf("Failed to unmarshal AI response. Raw response: %s", cleanedJSON)
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return &aiResponse, nil
}

// realModelFactory creates the real Gemini model.
func realModelFactory(ctx context.Context, apiKey string) (GenerativeAIModel, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	model := client.GenerativeModel("gemini-1.5-flash")
	return model, nil
}

func buildPrompt(userQuery string, articles []kb.Article) string {
	var articlesContext string
	for _, article := range articles {
		articlesContext += fmt.Sprintf("Article ID: %s\nTitle: %s\nContent: %s\n\n", article.ID, article.Title, article.Content)
	}

	return fmt.Sprintf(`
You are an expert IT support assistant for a corporate knowledge base.
Your task is to answer a user's question based ONLY on the provided knowledge base articles.

Here are the available articles:
--- START OF ARTICLES ---
%s
--- END OF ARTICLES ---

Here is the user's question: "%s"

Based on the articles, please perform the following two tasks:
1.  Provide a concise, one or two-sentence summary answer to the user's question. If the articles do not contain an answer, state that you could not find an answer.
2.  Identify the articles that are most relevant to the user's question.

Your entire response MUST be a single, valid JSON object with NO other text or explanation before or after it.
The JSON object must have the following structure:
{
  "ai_summary_answer": "Your concise summary answer here.",
  "ai_relevant_articles": [
    { "id": "The ID of the most relevant article", "title": "The title of the most relevant article" }
  ]
}
`, articlesContext, userQuery)
}

// cleanAIResponse removes everything before and after the JSON block.
func cleanAIResponse(rawResponse string) string {
	start := strings.Index(rawResponse, "{")
	end := strings.LastIndex(rawResponse, "}")
	if start == -1 || end == -1 {
		return rawResponse
	}
	return rawResponse[start : end+1]
}
