import React from 'react';

const ResultsDisplay = ({ result }) => {
  return (
    <div>
      <h2>AI Summary</h2>
      <p>{result.ai_summary_answer}</p>
    </div>
  );
};

export default ResultsDisplay;