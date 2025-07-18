import React, { useState } from 'react'; 
import { postSearchQuery } from './services/api'; 
import SearchBar from './components/SearchBar';
import ResultsDisplay from './components/ResultsDisplay';
import RelevantArticles from './components/RelevantArticles';
import './App.css';

function App() {
  const [searchResult, setSearchResult] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSearch = async (query) => {
    setIsLoading(true);
    setError('');
    setSearchResult(null);

    try {
      const result = await postSearchQuery(query);
      setSearchResult(result);
    } catch (err) {
      setError('Sorry, something went wrong. Please try again.');
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="app">
      <div className="app-container">
        <header className="app-header">
          <div className="logo">
            <div className="logo-icon">🤖</div>
            <h1>AI Knowledge Base</h1>
          </div>
          <p className="app-subtitle">
            Get instant answers to your IT questions using our intelligent knowledge base
          </p>
        </header>

        <main className="app-main">
          <SearchBar onSearch={handleSearch} isLoading={isLoading} /> 
          
          {isLoading && (
            <div className="loading-container">
              <div className="loading-spinner"></div>
              <p className="loading-text">Searching for answers...</p>
            </div>
          )}

          {error && (
            <div className="error-container">
              <div className="error-icon">⚠️</div>
              <p className="error-text">{error}</p>
            </div>
          )}

          {searchResult && (
            <div className="results-container">
              <ResultsDisplay result={searchResult} />
              <RelevantArticles articles={searchResult.ai_relevant_articles} />
            </div>
          )}
        </main>

        <footer className="app-footer">
          <p>Powered by AI • Built with React & Go</p>
        </footer>
      </div>
    </div>
  );
}

export default App;