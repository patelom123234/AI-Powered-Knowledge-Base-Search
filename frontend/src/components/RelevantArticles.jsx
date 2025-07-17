// In: frontend/src/components/RelevantArticles.jsx
import React from 'react';

const RelevantArticles = ({ articles }) => {
  return (
    <div>
      <h3>Relevant Articles</h3>
      {articles && articles.length > 0 ? (
        <ul>
          {articles.map((article) => (
            <li key={article.id}>{article.title}</li>
          ))}
        </ul>
      ) : (
        <p>No relevant articles found.</p>
      )}
    </div>
  );
};

export default RelevantArticles;