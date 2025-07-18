package ai

import (
	"ai-knowledge-base/internal/kb"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/generative-ai-go/genai"
)

// MockGenerativeAIModel implements GenerativeAIModel for testing.
type MockGenerativeAIModel struct {
	shouldError bool
	response    *genai.GenerateContentResponse
	errorMsg    string
}

func (m *MockGenerativeAIModel) GenerateContent(ctx context.Context, parts ...genai.Part) (*genai.GenerateContentResponse, error) {
	if m.shouldError {
		return nil, fmt.Errorf(m.errorMsg)
	}
	return m.response, nil
}

// mockModelFactory creates a mock model for testing.
func mockModelFactory(shouldError bool, response *genai.GenerateContentResponse, errorMsg string) ModelFactory {
	return func(ctx context.Context, apiKey string) (GenerativeAIModel, error) {
		return &MockGenerativeAIModel{
			shouldError: shouldError,
			response:    response,
			errorMsg:    errorMsg,
		}, nil
	}
}

// TestAIResponseStruct tests the AIResponse struct definition and JSON marshaling.
func TestAIResponseStruct(t *testing.T) {
	response := AIResponse{
		SummaryAnswer: "This is a test summary answer",
		RelevantArticles: []kb.Article{
			{
				ID:      "kb-001",
				Title:   "Test Article 1",
				Content: "Test content 1",
			},
			{
				ID:      "kb-002",
				Title:   "Test Article 2",
				Content: "Test content 2",
			},
		},
	}

	if response.SummaryAnswer != "This is a test summary answer" {
		t.Errorf("Expected SummaryAnswer 'This is a test summary answer', got '%s'", response.SummaryAnswer)
	}

	if len(response.RelevantArticles) != 2 {
		t.Errorf("Expected 2 relevant articles, got %d", len(response.RelevantArticles))
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Errorf("Failed to marshal AIResponse to JSON: %v", err)
	}

	var unmarshaled map[string]interface{}
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal JSON: %v", err)
	}

	if _, exists := unmarshaled["ai_summary_answer"]; !exists {
		t.Error("JSON missing field 'ai_summary_answer'")
	}

	if _, exists := unmarshaled["ai_relevant_articles"]; !exists {
		t.Error("JSON missing field 'ai_relevant_articles'")
	}

	var unmarshaledResponse AIResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	if err != nil {
		t.Errorf("Failed to unmarshal JSON to AIResponse: %v", err)
	}

	if unmarshaledResponse.SummaryAnswer != response.SummaryAnswer {
		t.Errorf("Unmarshaled SummaryAnswer doesn't match: expected '%s', got '%s'", response.SummaryAnswer, unmarshaledResponse.SummaryAnswer)
	}

	if len(unmarshaledResponse.RelevantArticles) != len(response.RelevantArticles) {
		t.Errorf("Unmarshaled RelevantArticles count doesn't match: expected %d, got %d", len(response.RelevantArticles), len(unmarshaledResponse.RelevantArticles))
	}
}

func TestBuildPrompt(t *testing.T) {
	articles := []kb.Article{
		{
			ID:      "kb-001",
			Title:   "How to reset password",
			Content: "Go to login page and click forgot password",
		},
		{
			ID:      "kb-002",
			Title:   "VPN Setup",
			Content: "Install VPN client and configure settings",
		},
	}

	userQuery := "How do I reset my password?"
	prompt := buildPrompt(userQuery, articles)

	if !strings.Contains(prompt, userQuery) {
		t.Error("Prompt does not contain user query")
	}

	if !strings.Contains(prompt, "kb-001") {
		t.Error("Prompt does not contain first article ID")
	}

	if !strings.Contains(prompt, "kb-002") {
		t.Error("Prompt does not contain second article ID")
	}

	if !strings.Contains(prompt, "How to reset password") {
		t.Error("Prompt does not contain first article title")
	}

	if !strings.Contains(prompt, "VPN Setup") {
		t.Error("Prompt does not contain second article title")
	}

	if !strings.Contains(prompt, "Go to login page and click forgot password") {
		t.Error("Prompt does not contain first article content")
	}

	if !strings.Contains(prompt, "Install VPN client and configure settings") {
		t.Error("Prompt does not contain second article content")
	}

	if !strings.Contains(prompt, "ai_summary_answer") {
		t.Error("Prompt does not contain expected JSON structure")
	}

	if !strings.Contains(prompt, "ai_relevant_articles") {
		t.Error("Prompt does not contain expected JSON structure")
	}

	// Test that prompt has proper structure
	if !strings.Contains(prompt, "START OF ARTICLES") {
		t.Error("Prompt does not contain article section markers")
	}

	if !strings.Contains(prompt, "END OF ARTICLES") {
		t.Error("Prompt does not contain article section markers")
	}
}

// TestBuildPromptWithEmptyArticles tests buildPrompt with empty articles slice.
func TestBuildPromptWithEmptyArticles(t *testing.T) {
	articles := []kb.Article{}
	userQuery := "Test query"
	prompt := buildPrompt(userQuery, articles)

	if !strings.Contains(prompt, userQuery) {
		t.Error("Prompt does not contain user query")
	}

	if !strings.Contains(prompt, "ai_summary_answer") {
		t.Error("Prompt does not contain expected JSON structure")
	}

	if !strings.Contains(prompt, "ai_relevant_articles") {
		t.Error("Prompt does not contain expected JSON structure")
	}
}

// TestCleanAIResponse tests the cleanAIResponse function.
func TestCleanAIResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple JSON",
			input:    `{"key": "value"}`,
			expected: `{"key": "value"}`,
		},
		{
			name:     "JSON with markdown",
			input:    "```json\n{\"key\": \"value\"}\n```",
			expected: `{"key": "value"}`,
		},
		{
			name:     "JSON with text before and after",
			input:    "Here is the response: {\"key\": \"value\"} That's it.",
			expected: `{"key": "value"}`,
		},
		{
			name:     "Multiple JSON blocks - takes first",
			input:    `{"first": "block"} {"second": "block"}`,
			expected: `{"first": "block"} {"second": "block"}`,
		},
		{
			name:     "No JSON braces",
			input:    "This is just text without JSON",
			expected: "This is just text without JSON",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Only opening brace",
			input:    "{",
			expected: "{",
		},
		{
			name:     "Only closing brace",
			input:    "}",
			expected: "}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanAIResponse(tt.input)
			if result != tt.expected {
				t.Errorf("cleanAIResponse(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestGetAIAnswerWithFactorySuccess tests successful AI response.
func TestGetAIAnswerWithFactorySuccess(t *testing.T) {
	mockResponse := &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{
			{
				Content: &genai.Content{
					Parts: []genai.Part{
						genai.Text(`{"ai_summary_answer": "To reset your password, go to the login page and click the forgot password link.", "ai_relevant_articles": [{"id": "kb-001", "title": "How to reset your password"}]}`),
					},
				},
			},
		},
	}

	factory := mockModelFactory(false, mockResponse, "")
	articles := []kb.Article{
		{
			ID:      "kb-001",
			Title:   "How to reset your password",
			Content: "Go to login page and click forgot password",
		},
	}

	originalKey := os.Getenv("GEMINI_API_KEY")
	os.Setenv("GEMINI_API_KEY", "test-key")
	defer os.Setenv("GEMINI_API_KEY", originalKey)

	response, err := getAIAnswerWithFactory(factory, "How do I reset my password?", articles)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if response == nil {
		t.Error("Expected response, got nil")
	}

	if response.SummaryAnswer != "To reset your password, go to the login page and click the forgot password link." {
		t.Errorf("Expected specific summary answer, got: %s", response.SummaryAnswer)
	}

	if len(response.RelevantArticles) != 1 {
		t.Errorf("Expected 1 relevant article, got %d", len(response.RelevantArticles))
	}

	if response.RelevantArticles[0].ID != "kb-001" {
		t.Errorf("Expected article ID 'kb-001', got '%s'", response.RelevantArticles[0].ID)
	}
}

// TestGetAIAnswerWithFactoryAPIError tests AI API error handling.
func TestGetAIAnswerWithFactoryAPIError(t *testing.T) {
	factory := mockModelFactory(true, nil, "API error occurred")
	articles := []kb.Article{
		{
			ID:      "kb-001",
			Title:   "Test Article",
			Content: "Test content",
		},
	}

	originalKey := os.Getenv("GEMINI_API_KEY")
	os.Setenv("GEMINI_API_KEY", "test-key")
	defer os.Setenv("GEMINI_API_KEY", originalKey)

	response, err := getAIAnswerWithFactory(factory, "test query", articles)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if !strings.Contains(err.Error(), "API error occurred") {
		t.Errorf("Expected error about API error, got: %v", err)
	}

	if response != nil {
		t.Error("Expected nil response when API error occurs")
	}
}

// TestGetAIAnswerWithFactoryEmptyResponse tests empty AI response handling.
func TestGetAIAnswerWithFactoryEmptyResponse(t *testing.T) {
	mockResponse := &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{},
	}

	factory := mockModelFactory(false, mockResponse, "")
	articles := []kb.Article{
		{
			ID:      "kb-001",
			Title:   "Test Article",
			Content: "Test content",
		},
	}

	originalKey := os.Getenv("GEMINI_API_KEY")
	os.Setenv("GEMINI_API_KEY", "test-key")
	defer os.Setenv("GEMINI_API_KEY", originalKey)

	response, err := getAIAnswerWithFactory(factory, "test query", articles)

	if err == nil {
		t.Error("Expected error for empty response, got nil")
	}

	if !strings.Contains(err.Error(), "received an empty response from AI") {
		t.Errorf("Expected error about empty response, got: %v", err)
	}

	if response != nil {
		t.Error("Expected nil response when AI response is empty")
	}
}

// TestGetAIAnswerWithFactoryInvalidJSON tests invalid JSON response handling.
func TestGetAIAnswerWithFactoryInvalidJSON(t *testing.T) {
	mockResponse := &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{
			{
				Content: &genai.Content{
					Parts: []genai.Part{
						genai.Text(`This is not valid JSON`),
					},
				},
			},
		},
	}

	factory := mockModelFactory(false, mockResponse, "")
	articles := []kb.Article{
		{
			ID:      "kb-001",
			Title:   "Test Article",
			Content: "Test content",
		},
	}

	originalKey := os.Getenv("GEMINI_API_KEY")
	os.Setenv("GEMINI_API_KEY", "test-key")
	defer os.Setenv("GEMINI_API_KEY", originalKey)

	response, err := getAIAnswerWithFactory(factory, "test query", articles)

	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}

	if !strings.Contains(err.Error(), "failed to parse AI response") {
		t.Errorf("Expected error about parsing AI response, got: %v", err)
	}

	if response != nil {
		t.Error("Expected nil response when JSON is invalid")
	}
}

// TestGetAIAnswerWithoutAPIKey tests GetAIAnswer when API key is not set.
func TestGetAIAnswerWithoutAPIKey(t *testing.T) {
	originalKey := os.Getenv("GEMINI_API_KEY")
	os.Unsetenv("GEMINI_API_KEY")
	defer os.Setenv("GEMINI_API_KEY", originalKey)

	articles := []kb.Article{
		{
			ID:      "kb-001",
			Title:   "Test Article",
			Content: "Test content",
		},
	}

	response, err := GetAIAnswer("test query", articles)

	// Should return error when API key is not set
	if err == nil {
		t.Error("Expected error when GEMINI_API_KEY is not set, got nil")
	}

	if !strings.Contains(err.Error(), "GEMINI_API_KEY environment variable not set") {
		t.Errorf("Expected error about missing API key, got: %v", err)
	}

	if response != nil {
		t.Error("Expected nil response when API key is not set")
	}
}

// TestGetAIAnswerWithEmptyArticles tests GetAIAnswer with empty articles slice.
func TestGetAIAnswerWithEmptyArticles(t *testing.T) {
	mockResponse := &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{
			{
				Content: &genai.Content{
					Parts: []genai.Part{
						genai.Text(`{"ai_summary_answer": "No relevant articles found.", "ai_relevant_articles": []}`),
					},
				},
			},
		},
	}

	factory := mockModelFactory(false, mockResponse, "")
	articles := []kb.Article{}

	originalKey := os.Getenv("GEMINI_API_KEY")
	os.Setenv("GEMINI_API_KEY", "test-key")
	defer os.Setenv("GEMINI_API_KEY", originalKey)

	response, err := getAIAnswerWithFactory(factory, "test query", articles)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if response == nil {
		t.Error("Expected response, got nil")
	}

	if response.SummaryAnswer != "No relevant articles found." {
		t.Errorf("Expected specific summary answer, got: %s", response.SummaryAnswer)
	}

	if len(response.RelevantArticles) != 0 {
		t.Errorf("Expected 0 relevant articles, got %d", len(response.RelevantArticles))
	}
}

// TestGetAIAnswerWithEmptyQuery tests GetAIAnswer with empty query.
func TestGetAIAnswerWithEmptyQuery(t *testing.T) {
	mockResponse := &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{
			{
				Content: &genai.Content{
					Parts: []genai.Part{
						genai.Text(`{"ai_summary_answer": "Please provide a query.", "ai_relevant_articles": []}`),
					},
				},
			},
		},
	}

	factory := mockModelFactory(false, mockResponse, "")
	articles := []kb.Article{
		{
			ID:      "kb-001",
			Title:   "Test Article",
			Content: "Test content",
		},
	}

	originalKey := os.Getenv("GEMINI_API_KEY")
	os.Setenv("GEMINI_API_KEY", "test-key")
	defer os.Setenv("GEMINI_API_KEY", originalKey)

	response, err := getAIAnswerWithFactory(factory, "", articles)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if response == nil {
		t.Error("Expected response, got nil")
	}

	if response.SummaryAnswer != "Please provide a query." {
		t.Errorf("Expected specific summary answer, got: %s", response.SummaryAnswer)
	}
}

// BenchmarkBuildPrompt benchmarks the buildPrompt function.
func BenchmarkBuildPrompt(b *testing.B) {
	articles := []kb.Article{
		{
			ID:      "kb-001",
			Title:   "Test Article 1",
			Content: "Test content 1",
		},
		{
			ID:      "kb-002",
			Title:   "Test Article 2",
			Content: "Test content 2",
		},
	}
	userQuery := "How do I reset my password?"

	for i := 0; i < b.N; i++ {
		buildPrompt(userQuery, articles)
	}
}

// BenchmarkCleanAIResponse benchmarks the cleanAIResponse function.
func BenchmarkCleanAIResponse(b *testing.B) {
	input := "Here is the response: {\"key\": \"value\"} That's it."
	for i := 0; i < b.N; i++ {
		cleanAIResponse(input)
	}
}
