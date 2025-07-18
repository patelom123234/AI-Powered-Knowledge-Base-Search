package kb

import (
	"encoding/json"
	"testing"
)

// TestArticleStruct tests the Article struct definition and JSON marshaling.
func TestArticleStruct(t *testing.T) {
	article := Article{
		ID:      "test-001",
		Title:   "Test Article",
		Content: "This is a test article content.",
	}

	if article.ID != "test-001" {
		t.Errorf("Expected ID 'test-001', got '%s'", article.ID)
	}

	if article.Title != "Test Article" {
		t.Errorf("Expected Title 'Test Article', got '%s'", article.Title)
	}

	if article.Content != "This is a test article content." {
		t.Errorf("Expected Content 'This is a test article content.', got '%s'", article.Content)
	}

	jsonData, err := json.Marshal(article)
	if err != nil {
		t.Errorf("Failed to marshal article to JSON: %v", err)
	}

	expectedJSON := `{"id":"test-001","title":"Test Article","content":"This is a test article content."}`
	if string(jsonData) != expectedJSON {
		t.Errorf("Expected JSON '%s', got '%s'", expectedJSON, string(jsonData))
	}

	var unmarshaledArticle Article
	err = json.Unmarshal(jsonData, &unmarshaledArticle)
	if err != nil {
		t.Errorf("Failed to unmarshal JSON to article: %v", err)
	}

	if unmarshaledArticle.ID != article.ID {
		t.Errorf("Unmarshaled ID doesn't match: expected '%s', got '%s'", article.ID, unmarshaledArticle.ID)
	}

	if unmarshaledArticle.Title != article.Title {
		t.Errorf("Unmarshaled Title doesn't match: expected '%s', got '%s'", article.Title, unmarshaledArticle.Title)
	}

	if unmarshaledArticle.Content != article.Content {
		t.Errorf("Unmarshaled Content doesn't match: expected '%s', got '%s'", article.Content, unmarshaledArticle.Content)
	}
}

// TestGetArticles tests the GetArticles function.
func TestGetArticles(t *testing.T) {
	articles := GetArticles()

	if articles == nil {
		t.Fatal("GetArticles() returned nil slice")
	}

	expectedCount := 3
	if len(articles) != expectedCount {
		t.Errorf("Expected %d articles, got %d", expectedCount, len(articles))
	}

	expectedArticles := []struct {
		ID      string
		Title   string
		Content string
	}{
		{
			ID:      "kb-001",
			Title:   "How to reset your password",
			Content: "To reset your password, go to the login page and click on the 'Forgot Password' link. You will receive an email with instructions on how to set a new password. Make sure to choose a strong password that you haven't used before.",
		},
		{
			ID:      "kb-002",
			Title:   "VPN Connection Issues",
			Content: "If you are having trouble connecting to the company VPN, first ensure you have the latest version of the VPN client installed. Second, check your internet connection to make sure it is stable. If the problem persists, try restarting your computer. Contact IT support if you are still unable to connect.",
		},
		{
			ID:      "kb-003",
			Title:   "Setting up a new printer",
			Content: "To set up a new printer, first connect it to the network via an ethernet cable or Wi-Fi. Then, go to your computer's system settings, find the 'Printers & Scanners' section, and click 'Add Printer'. Your computer should automatically detect the printer. If not, you may need to install drivers from the manufacturer's website.",
		},
	}

	for i, expected := range expectedArticles {
		if i >= len(articles) {
			t.Errorf("Article %d not found", i)
			continue
		}

		article := articles[i]
		if article.ID != expected.ID {
			t.Errorf("Article %d: Expected ID '%s', got '%s'", i, expected.ID, article.ID)
		}

		if article.Title != expected.Title {
			t.Errorf("Article %d: Expected Title '%s', got '%s'", i, expected.Title, article.Title)
		}

		if article.Content != expected.Content {
			t.Errorf("Article %d: Expected Content '%s', got '%s'", i, expected.Content, article.Content)
		}
	}
}

// TestGetArticlesConsistency tests that GetArticles returns consistent results.
func TestGetArticlesConsistency(t *testing.T) {
	articles1 := GetArticles()
	articles2 := GetArticles()

	if len(articles1) != len(articles2) {
		t.Errorf("GetArticles() returned different lengths: %d vs %d", len(articles1), len(articles2))
	}

	for i := range articles1 {
		if articles1[i].ID != articles2[i].ID {
			t.Errorf("Article %d ID inconsistent: '%s' vs '%s'", i, articles1[i].ID, articles2[i].ID)
		}
		if articles1[i].Title != articles2[i].Title {
			t.Errorf("Article %d Title inconsistent: '%s' vs '%s'", i, articles1[i].Title, articles2[i].Title)
		}
		if articles1[i].Content != articles2[i].Content {
			t.Errorf("Article %d Content inconsistent: '%s' vs '%s'", i, articles1[i].Content, articles2[i].Content)
		}
	}
}

// TestGetArticlesJSONMarshaling tests that all articles can be marshaled to JSON.
func TestGetArticlesJSONMarshaling(t *testing.T) {
	articles := GetArticles()

	for i, article := range articles {
		jsonData, err := json.Marshal(article)
		if err != nil {
			t.Errorf("Failed to marshal article %d to JSON: %v", i, err)
			continue
		}

		var unmarshaled map[string]interface{}
		err = json.Unmarshal(jsonData, &unmarshaled)
		if err != nil {
			t.Errorf("Failed to unmarshal article %d JSON: %v", i, err)
			continue
		}

		expectedFields := []string{"id", "title", "content"}
		for _, field := range expectedFields {
			if _, exists := unmarshaled[field]; !exists {
				t.Errorf("Article %d JSON missing field '%s'", i, field)
			}
		}

		if unmarshaled["id"] != article.ID {
			t.Errorf("Article %d JSON ID mismatch: expected '%s', got '%v'", i, article.ID, unmarshaled["id"])
		}
		if unmarshaled["title"] != article.Title {
			t.Errorf("Article %d JSON Title mismatch: expected '%s', got '%v'", i, article.Title, unmarshaled["title"])
		}
		if unmarshaled["content"] != article.Content {
			t.Errorf("Article %d JSON Content mismatch: expected '%s', got '%v'", i, article.Content, unmarshaled["content"])
		}
	}
}

// TestGetArticlesSliceProperties tests slice properties.
func TestGetArticlesSliceProperties(t *testing.T) {
	articles := GetArticles()

	if len(articles) == 0 {
		t.Error("GetArticles() returned empty slice")
	}

	count := 0
	for range articles {
		count++
	}
	if count != len(articles) {
		t.Errorf("Iteration count mismatch: expected %d, got %d", len(articles), count)
	}

	for i := range articles {
		if articles[i].ID == "" {
			t.Errorf("Article %d has empty ID", i)
		}
		if articles[i].Title == "" {
			t.Errorf("Article %d has empty Title", i)
		}
		if articles[i].Content == "" {
			t.Errorf("Article %d has empty Content", i)
		}
	}
}

// BenchmarkGetArticles benchmarks the GetArticles function.
func BenchmarkGetArticles(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetArticles()
	}
}
