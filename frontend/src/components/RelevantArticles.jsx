// In: frontend/src/components/RelevantArticles.jsx
import React from 'react';
import './RelevantArticles.css';

const RelevantArticles = ({ articles }) => {
  if (!articles || articles.length === 0) {
    return (
      <div className="relevant-articles">
        <div className="articles-header">
          <div className="articles-icon">📚</div>
          <h3 className="articles-title">Relevant Articles</h3>
        </div>
        <div className="no-articles">
          <div className="no-articles-icon">🔍</div>
          <p className="no-articles-text">No relevant articles found for this query.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="relevant-articles">
      <div className="articles-header">
        <div className="articles-icon">📚</div>
        <h3 className="articles-title">Relevant Articles</h3>
        <div className="articles-count">
          {articles.length} {articles.length === 1 ? 'article' : 'articles'}
        </div>
      </div>
      
      <div className="articles-grid">
        {articles.map((article, index) => (
          <div 
            key={article.id} 
            className="article-card"
            style={{ animationDelay: `${index * 0.1}s` }}
          >
            <div className="article-header">
              <div className="article-number">#{index + 1}</div>
              <div className="article-id">{article.id}</div>
            </div>
            <h4 className="article-title">{article.title}</h4>
            <div className="article-footer">
              <div className="article-badge">
                <span className="badge-icon">📄</span>
                <span>Knowledge Base</span>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default RelevantArticles;