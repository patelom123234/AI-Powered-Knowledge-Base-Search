import React, { useState } from 'react'; 
import { postSearchQuery } from './services/api'; 
import SearchBar from './components/SearchBar';
import ResultsDisplay from './components/ResultsDisplay';
import RelevantArticles from './components/RelevantArticles';

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
    <div style={{ maxWidth: '800px', margin: 'auto' }}>
      <h1>AI-Powered Knowledge Base Search</h1>
      <SearchBar onSearch={handleSearch} isLoading={isLoading} /> 
      
      <hr style={{ margin: '2rem 0' }} />

      {isLoading && <p>Searching for answers...</p>}

      {error && <p style={{ color: 'red' }}>{error}</p>}

      {searchResult && (
        <>
          <ResultsDisplay result={searchResult} />
          <RelevantArticles articles={searchResult.ai_relevant_articles} />
        </>
      )}
    </div>
  );
}

export default App;