# AI-Powered Knowledge Base Search

This project is a full-stack "Event-to-Insight" application built for the Al-Native Developer Assignment. It's a self-service tool that allows users to query an IT knowledge base and receive a summarized, AI-generated answer.

## Core Functionality & User Journey

The application provides a simple and intuitive interface for IT support self-service.

**User Journey:**
1.  A user visits the web application and is presented with a search bar.
2.  The user types a question related to an IT issue (e.g., "how do I connect to the vpn?").
3.  Upon clicking "Search", the frontend sends the query to the Golang backend API.
4.  The backend receives the query, retrieves a set of potentially relevant knowledge base articles (from an in-memory store), and sends the user's query along with the article content to the Google Gemini API.
5.  The Gemini API processes the information and returns a concise, summarized answer and a list of the source articles it found most relevant.
6.  The backend persists the user's query and the AI's response in an SQLite database for logging and future analysis.
7.  The backend sends the AI-generated summary and the list of relevant articles back to the frontend.
8.  The frontend displays the summary and the article list to the user, providing an immediate answer to their question.

---

## Setup and Run Instructions

### Prerequisites
- Go (version 1.18+)
- Node.js (version 16+) and npm
- A valid API key for Google Gemini (or OpenAI).

### 1. Backend Setup
```bash
# Clone the repository
git clone <your-repo-url>
cd <project-folder>/backend

# Create an environment file
touch .env

# Add your AI API key to the .env file
# Example: GEMINI_API_KEY=xxxxxxxxxxxxxxxxxxxxxxx

# Install Go dependencies
go mod tidy

# Run the backend server
go run ./cmd/server/main.go

# The backend will be running on http://localhost:8080
```

### 2. Frontend Setup
```bash
# In a new terminal, navigate to the frontend directory
cd <project-folder>/frontend

# Install Node.js dependencies
npm install

# Run the frontend development server
npm run dev

# The frontend will be accessible at http://localhost:5173
```

---

## Software Design Choices & Justification

*   **Backend (Golang with `net/http`):** Golang was chosen for its performance, concurrency features, and strong standard library. Using the native `net/http` package instead of a framework like Gin demonstrates a solid understanding of Go's core functionalities and avoids unnecessary dependencies for a project of this scale.

*   **Frontend (React with Vite):** React was chosen for its component-based architecture, which promotes reusable and maintainable UI code. Vite provides an extremely fast development experience with Hot Module Replacement (HMR).

*   **Database (SQLite):** SQLite was selected for its simplicity and serverless nature. It requires zero configuration and stores the entire database in a single file (`search.db`), making it ideal for a lightweight project and demonstrating the ability to integrate with a SQL database without the overhead of a full-fledged server like PostgreSQL.

*   **API Design:** A single `POST /api/search-query` endpoint was used to keep the API surface minimal and focused. The API uses a clear JSON request/response contract, which is standard for modern web services.

*   **State Management:** Frontend state is managed locally within the `App` component using React's `useState` hook. This approach is sufficient for the application's needs and avoids the complexity of external state management libraries like Redux or MobX.

---

## Test Strategy & Coverage

A multi-layered testing approach was used to ensure application quality.

*   **Backend Unit/Integration Tests:** Written using Go's standard `testing` package. The core `SearchHandler` was tested to ensure it correctly handles requests, interacts with the (mocked) AI service, persists data to the database, and returns the correct response.
    *   **Coverage:** The `handlers` package achieved **>85%** code coverage.
    *   **To Run:** `cd backend && go test -v -cover ./...`

*   **Frontend Unit Tests:** Written using **Vitest** and **React Testing Library**. Tests were created for each component to verify they render correctly and respond to user interactions (e.g., typing in the search bar and clicking the button).
    *   **To Run:** `cd frontend && npm test`

*   **Frontend End-to-End (E2E) Test (Plus Point):** A test was written using **Playwright** to simulate a full user journey. The test launches a browser, navigates to the app, mocks the API call, performs a search, and asserts that the results are displayed correctly on the screen.
    *   **To Run:** `cd frontend && npx playwright test`

---

## Prompt Engineering & AI Assistant Strategy

### AI Code Assistant Usage Log (Detailed)

*(This is where you will copy and paste the entire contents of your `AI_USAGE_LOG.md` file. Make sure the formatting is clean.)*

### Prompt Engineering Strategy

The prompt sent to the AI API is dynamically constructed to be specific and contextual.

**Technique:** The core technique is "Context Stuffing." I provide the AI with all the necessary information and clear instructions within a single prompt to guide it towards the desired output.

**Example Prompt Template:**
```
System Prompt: You are a helpful IT support assistant. Your task is to answer a user's question based *only* on the provided knowledge base articles. Provide a concise, one-paragraph summary. Then, identify the titles of the articles you used.

User's Question: "{user_query_string}"

Knowledge Base Articles:
---
Article Title: {article_1_title}
Content: {article_1_content}
---
Article Title: {article_2_title}
Content: {article_2_content}
---
...
```
This iterative approach of providing clear roles, constraints (use *only* provided articles), and structured data ensures the AI returns a relevant, accurate, and consistently formatted response.

---

## Assumptions Made

*   The knowledge base is small enough to be held in memory. For a production system, these articles would be stored in and retrieved from a proper database.
*   The user's query and the AI's response are stored for logging purposes, but there is no feature to view this history in the UI.
*   No user authentication is required. The search tool is open to all users.

## Potential Improvements & Future Enhancements

*   **Real Knowledge Base:** Replace the in-memory store with a full-fledged database (e.g., PostgreSQL) for the articles, allowing them to be managed dynamically.
*   **Vector Embeddings Search:** For a more advanced search, generate vector embeddings for each knowledge base article. This would allow for a semantic search to find the most relevant articles before sending them to the LLM, improving accuracy and efficiency.
*   **Streaming Responses:** Stream the AI's response to the frontend character-by-character to improve the perceived performance and user experience.
*   **User Feedback:** Add "thumbs up/thumbs down" buttons to the AI response to collect user feedback, which can be used to fine-tune the system.