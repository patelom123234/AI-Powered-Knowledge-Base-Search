import React, { useState } from 'react';
import './SearchBar.css';

const SearchBar = ({ onSearch, isLoading }) => {
  const [query, setQuery] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault(); 
    if (!query.trim()) return; 
    onSearch(query);
  };

  const handleKeyDown = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      if (!query.trim()) return;
      onSearch(query);
    }
  };

  return (
    <div className="search-container">
      <form onSubmit={handleSubmit} className="search-form">
        <div className="search-input-container">
          <div className="search-icon">🔍</div>
          <input
            type="text"
            placeholder="Ask anything about our IT systems..."
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            onKeyDown={handleKeyDown}
            className="search-input"
            disabled={isLoading} 
          />
          <button 
            type="submit" 
            className={`search-button ${isLoading ? 'loading' : ''}`}
            disabled={isLoading || !query.trim()}
          >
            {isLoading ? (
              <>
                <div className="button-spinner"></div>
                <span>Searching...</span>
              </>
            ) : (
              <>
                <span>Search</span>
                <div className="button-arrow">→</div>
              </>
            )}
          </button>
        </div>
        
        <div className="search-suggestions">
          <p className="suggestions-text">Try asking about:</p>
          <div className="suggestion-tags">
            <button 
              type="button" 
              className="suggestion-tag"
              onClick={() => onSearch("How do I reset my password?")}
              disabled={isLoading}
            >
              Password Reset
            </button>
            <button 
              type="button" 
              className="suggestion-tag"
              onClick={() => onSearch("VPN connection issues")}
              disabled={isLoading}
            >
              VPN Setup
            </button>
            <button 
              type="button" 
              className="suggestion-tag"
              onClick={() => onSearch("Setting up a new printer")}
              disabled={isLoading}
            >
              Printer Setup
            </button>
          </div>
        </div>
      </form>
    </div>
  );
};

export default SearchBar;