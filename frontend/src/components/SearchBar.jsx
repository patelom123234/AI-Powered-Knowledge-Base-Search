// In: frontend/src/components/SearchBar.jsx
import React, { useState } from 'react';

const SearchBar = ({ onSearch, isLoading }) => {
  // State to hold the value of the input field
  const [query, setQuery] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault(); 
    if (!query.trim()) return; 
    onSearch(query);
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="text"
        placeholder="Ask a question about our IT systems..."
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        style={{ width: '100%', padding: '10px', fontSize: '16px', boxSizing: 'border-box' }}
        disabled={isLoading} 
      />
      <button type="submit" style={{ marginTop: '10px', padding: '10px 20px' }} disabled={isLoading}>
        {/* Show different text on the button when loading */}
        {isLoading ? 'Searching...' : 'Search'}
      </button>
    </form>
  );
};

export default SearchBar;