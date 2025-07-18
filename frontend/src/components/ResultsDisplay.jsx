import React from 'react';
import './ResultsDisplay.css';

const ResultsDisplay = ({ result }) => {
  return (
    <div className="results-display">
      <div className="results-header">
        <div className="results-icon">💡</div>
        <h2 className="results-title">AI Summary</h2>
      </div>
      
      <div className="results-content">
        <div className="summary-card">
          <div className="summary-text">
            {result.ai_summary_answer}
          </div>
          <div className="summary-footer">
            <div className="ai-badge">
              <span className="ai-icon">🤖</span>
              <span>AI Generated</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ResultsDisplay;